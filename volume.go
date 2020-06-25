package sonyrest

import (
	"context"
	"errors"
	"strconv"
	"time"

	"encoding/json"

	"github.com/byuoitav/common/log"
)

func (t *TV) GetVolumes(ctx context.Context, blocks []string) (map[string]int, error) {
	log.L.Infof("Getting volume for %v", t.Address)
	toReturn := make(map[string]int)
	parentResponse, err := t.getAudioInformation(ctx)
	if err != nil {
		return toReturn, err
	}
	log.L.Infof("%v", parentResponse)

	for _, outerResult := range parentResponse.Result {

		for _, result := range outerResult {

			if result.Target == "speaker" {

				toReturn[""] = result.Volume
			}
		}
	}
	log.L.Infof("Done")

	return toReturn, nil
}

func (t *TV) SetVolume(ctx context.Context, block string, volume int) error {

	if volume > 100 || volume < 0 {
		return errors.New("Error: volume must be a value from 0 to 100!")
	}

	log.L.Debugf("Setting volume for %s to %v...", t.Address, volume)
	params := make(map[string]interface{})
	params["target"] = "speaker"
	params["volume"] = strconv.Itoa(volume)

	err := t.BuildAndSendPayload(ctx, t.Address, "audio", "setAudioVolume", params)
	if err != nil {
		return err
	}

	//do the same for the headphone
	params = make(map[string]interface{})
	params["target"] = "headphone"
	params["volume"] = strconv.Itoa(volume)

	err = t.BuildAndSendPayload(ctx, t.Address, "audio", "setAudioVolume", params)
	if err != nil {
		return err
	}

	log.L.Debugf("Done.")
	return nil
}

func (t *TV) getAudioInformation(ctx context.Context) (SonyAudioResponse, error) {
	payload := SonyTVRequest{
		Params:  []map[string]interface{}{},
		Method:  "getVolumeInformation",
		Version: "1.0",
		ID:      1,
	}

	log.L.Infof("%+v", payload)

	resp, err := t.PostHTTPWithContext(ctx, "audio", payload)

	parentResponse := SonyAudioResponse{}

	log.L.Infof("%s", resp)

	err = json.Unmarshal(resp, &parentResponse)
	return parentResponse, err

}

func (t *TV) GetMutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	toReturn := make(map[string]bool)
	log.L.Infof("Getting mute status for %v", t.Address)
	parentResponse, err := t.getAudioInformation(ctx)
	if err != nil {
		return toReturn, err
	}

	for _, outerResult := range parentResponse.Result {
		for _, result := range outerResult {
			if result.Target == "speaker" {
				log.L.Infof("local mute: %v", result.Mute)
				toReturn[""] = result.Mute
			}
		}
	}

	log.L.Infof("Done")

	return toReturn, nil
}

func (t *TV) SetMute(ctx context.Context, block string, mute bool) error {
	params := make(map[string]interface{})
	params["status"] = mute

	err := t.BuildAndSendPayload(ctx, t.Address, "audio", "setAudioMute", params)
	if err != nil {
		return err
	}
	//we need to validate that it was actually muted
	blocks := []string{block}
	postStatus, err := t.GetMutes(ctx, blocks)
	if err != nil {
		return err
	}

	if postStatus[""] == mute {
		return nil
	}

	//wait for a short time
	time.Sleep(10 * time.Millisecond)

	return nil
}
