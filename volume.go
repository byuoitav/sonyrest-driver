package sonyrest

import (
	"context"
	"errors"
	"time"

	"encoding/json"

	"github.com/byuoitav/common/log"
)

func (t *TV) GetVolumeByBlock(ctx context.Context, block string) (int, error) {
	log.L.Infof("Getting volume for %v", t.Address)
	parentResponse, err := t.getAudioInformation(ctx)
	if err != nil {
		return 0, err
	}
	log.L.Infof("%v", parentResponse)

	var output int
	for _, outerResult := range parentResponse.Result {

		for _, result := range outerResult {

			if result.Target == "speaker" {

				output = result.Volume
			}
		}
	}
	log.L.Infof("Done")

	return output, nil
}

func (t *TV) SetVolumeByBlock(ctx context.Context, block string, volume int) error {

	if volume > 100 || volume < 0 {
		return errors.New("Error: volume must be a value from 0 to 100!")
	}

	log.L.Debugf("Setting volume for %s to %v...", t.Address, volume)

	params := make(map[string]interface{})
	params["target"] = "speaker"
	params["volume"] = volume

	err := t.BuildAndSendPayload(ctx, t.Address, "audio", "setAudioVolume", params)
	if err != nil {
		return err
	}

	//do the same for the headphone
	params = make(map[string]interface{})
	params["target"] = "headphone"
	params["volume"] = volume

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

func (t *TV) GetMutedByBlock(ctx context.Context, block string) (bool, error) {
	log.L.Infof("Getting mute status for %v", t.Address)
	parentResponse, err := t.getAudioInformation(ctx)
	if err != nil {
		return false, err
	}
	var output bool
	for _, outerResult := range parentResponse.Result {
		for _, result := range outerResult {
			if result.Target == "speaker" {
				log.L.Infof("local mute: %v", result.Mute)
				output = result.Mute
			}
		}
	}

	log.L.Infof("Done")

	return output, nil
}

func (t *TV) SetMutedByBlock(ctx context.Context, block string, muted bool) error {
	params := make(map[string]interface{})
	params["status"] = muted

	err := t.BuildAndSendPayload(ctx, t.Address, "audio", "setAudioMute", params)
	if err != nil {
		return err
	}
	//we need to validate that it was actually muted
	postStatus, err := t.GetMutedByBlock(ctx, block)
	if err != nil {
		return err
	}

	if postStatus == muted {
		return nil
	}

	//wait for a short time
	time.Sleep(10 * time.Millisecond)

	return nil
}
