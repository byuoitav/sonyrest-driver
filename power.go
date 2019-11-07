package sonyrest

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
)

func (t *TV) GetPower(ctx context.Context) (string, error) {
	var output string

	payload := SonyTVRequest{
		Params: []map[string]interface{}{},
		Method: "getPowerStatus", Version: "1.0",
		ID: 1,
	}

	response, err := t.PostHTTPWithContext(ctx, "system", payload)
	if err != nil {
		return "", err
	}

	powerStatus := string(response)
	if strings.Contains(powerStatus, "active") {
		output = "on"
	} else if strings.Contains(powerStatus, "standby") {
		output = "standby"
	} else {
		return "", errors.New("Error getting power status")
	}

	return output, nil
}

func (t *TV) SetPower(ctx context.Context, power string) error {
	var status bool
	if power == "standby" {
		status = false
	} else {
		status = true
	}
	params := make(map[string]interface{})
	params["status"] = status

	payload := SonyTVRequest{
		Params:  []map[string]interface{}{params},
		Method:  "setPowerStatus",
		Version: "1.0",
		ID:      1,
	}

	log.L.Infof("Setting power to %v", status)

	_, err := t.PostHTTPWithContext(ctx, "system", payload)
	if err != nil {
		return err
	}

	// wait for the display to turn on
	ticker := time.NewTicker(256 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.New("context timed out while waiting for display to turn on")
		case <-ticker.C:
			power, err := t.GetPower(ctx)
			if err != nil {
				return err
			}

			log.L.Infof("Waiting for display power to change to %v, current status %s", status, power)

			switch {
			case status && power == "on":
				return nil
			case !status && power == "standby":
				return nil
			}
		}
	}
}
