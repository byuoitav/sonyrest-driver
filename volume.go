package sonyrest

import (
	"context"
	"encoding/json"

	"github.com/byuoitav/common/log"

	"github.com/byuoitav/common/status"
)

func (t *TV) GetVolume(ctx context.Context) (status.Volume, error) {
	log.L.Infof("Getting volume for %v", t.Address)
	parentResponse, err := t.getAudioInformation(ctx)
	if err != nil {
		return status.Volume{}, err
	}
	log.L.Infof("%v", parentResponse)

	var output status.Volume
	for _, outerResult := range parentResponse.Result {

		for _, result := range outerResult {

			if result.Target == "speaker" {

				output.Volume = result.Volume
			}
		}
	}
	log.L.Infof("Done")

	return output, nil
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

func (t *TV) GetMute(ctx context.Context) (status.Mute, error) {
	log.L.Infof("Getting mute status for %v", t.Address)
	parentResponse, err := t.getAudioInformation(ctx)
	if err != nil {
		return status.Mute{}, err
	}
	var output status.Mute
	for _, outerResult := range parentResponse.Result {
		for _, result := range outerResult {
			if result.Target == "speaker" {
				log.L.Infof("local mute: %v", result.Mute)
				output.Muted = result.Mute
			}
		}
	}

	log.L.Infof("Done")

	return output, nil
}
