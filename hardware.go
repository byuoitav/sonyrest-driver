package sonyrest

import (
	"context"
	"encoding/json"
	"net"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
)

func (t *TV) GetInfo(ctx context.Context) (interface{}, error) {
	return nil, nil
}

// GetHardwareInfo returns the hardware information for the device
func (t *TV) GetHardwareInfo(ctx context.Context) (structs.HardwareInfo, *nerr.E) {
	var toReturn structs.HardwareInfo

	// get the hostname
	addr, e := net.LookupAddr(t.Address)
	if e != nil {
		toReturn.Hostname = t.Address
	} else {
		toReturn.Hostname = strings.Trim(addr[0], ".")
	}

	// get Sony TV system information
	systemInfo, err := t.getSystemInfo(ctx)
	if err != nil {
		err.Addf("Could not get system info from %s", t.Address)
		return toReturn, err
	}

	toReturn.ModelName = systemInfo.Model
	toReturn.SerialNumber = systemInfo.Serial
	toReturn.FirmwareVersion = systemInfo.Generation

	// get Sony TV network settings
	networkInfo, err := t.getNetworkInfo(ctx)
	if err != nil {
		err.Addf("Could not get network info from %s", t.Address)
		return toReturn, err
	}

	toReturn.NetworkInfo = structs.NetworkInfo{
		IPAddress:  networkInfo.IPv4,
		MACAddress: networkInfo.HardwareAddress,
		Gateway:    networkInfo.Gateway,
		DNS:        networkInfo.DNS,
	}

	log.L.Info(toReturn)

	// get power status
	powerStatus, e := t.GetPower(context.TODO())
	if e != nil {
		err = nerr.Translate(e).Addf("Could not get power status from %s")
		return toReturn, err
	}

	toReturn.PowerStatus = powerStatus

	return toReturn, nil
}

func (t *TV) getSystemInfo(ctx context.Context) (SonySystemInformation, *nerr.E) {
	var system SonyTVSystemResponse

	payload := SonyTVRequest{
		Params: []map[string]interface{}{},
		Method: "getSystemInformation", Version: "1.0",
		ID: 1,
	}

	response, err := t.PostHTTPWithContext(ctx, "system", payload)
	if err != nil {
		return SonySystemInformation{}, nerr.Translate(err)
	}

	err = json.Unmarshal(response, &system)
	if err != nil {
		return SonySystemInformation{}, nerr.Translate(err)
	}

	return system.Result[0], nil
}

func (t *TV) getNetworkInfo(ctx context.Context) (SonyTVNetworkInformation, *nerr.E) {
	var network SonyNetworkResponse

	payload := SonyTVRequest{
		ID:      2,
		Method:  "getNetworkSettings",
		Version: "1.0",
		Params: []map[string]interface{}{
			map[string]interface{}{
				"netif": "eth0",
			},
		},
	}

	response, err := t.PostHTTPWithContext(ctx, "system", payload)
	if err != nil {
		return SonyTVNetworkInformation{}, nerr.Translate(err)
	}

	err = json.Unmarshal(response, &network)
	if err != nil {
		return SonyTVNetworkInformation{}, nerr.Translate(err)
	}

	return network.Result[0][0], nil
}
