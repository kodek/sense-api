package sense_api

import (
	"net/url"

	"fmt"

	"encoding/json"

	"github.com/golang/glog"
)

type RealtimeResponse struct {
	Type    string `json:"type"`
	Payload struct {
		Hz       float64   `json:"hz"`
		Channels []float64 `json:"channels"`
		Devices  []struct {
			W    float64 `json:"w"`
			Name string  `json:"name"`
			Icon string  `json:"icon"`
			ID   string  `json:"id"`
			Tags struct {
				DefaultUserDeviceType       string `json:"DefaultUserDeviceType"`
				TimelineAllowed             string `json:"TimelineAllowed"`
				UserDeviceType              string `json:"UserDeviceType"`
				UserEditable                string `json:"UserEditable"`
				UserDeviceTypeDisplayString string `json:"UserDeviceTypeDisplayString"`
				DeviceListAllowed           string `json:"DeviceListAllowed"`
				//UserEditable                string `json:"user_editable"`
			} `json:"tags"`
		} `json:"devices"`
		W     float64 `json:"w"`
		Stats struct {
			Mrcv float64 `json:"mrcv"`
			Brcv float64 `json:"brcv"`
			Msnd float64 `json:"msnd"`
		} `json:"_stats"`
		Epoch   int64         `json:"epoch"`
		Deltas  []interface{} `json:"deltas"`
		Voltage []float64     `json:"voltage"`
		Frame   int           `json:"frame"`
		SolarW  float64       `json:"solar_w"`
	} `json:"payload"`
}

const WSS_URL_FORMAT = "wss://clientrt.sense.com/monitors/%d/realtimefeed?access_token=%s"

func (c *ClientImpl) Realtime(done <-chan struct{}) (<-chan RealtimeResponse, error) {
	u, err := url.Parse(fmt.Sprintf(WSS_URL_FORMAT, c.monitorId, c.accessToken))
	if err != nil {
		return nil, err
	}

	recv, err := wssConnect(done, *u)
	if err != nil {
		return nil, err
	}

	recvParsed := make(chan RealtimeResponse)
	go func() {
		defer close(recvParsed)
		for {
			select {
			case msg := <-recv:
				var r RealtimeResponse
				err := json.Unmarshal(msg, &r)
				if err != nil {
					glog.Error("Cannot parse to JSON: ", err, string(msg))
					return
				}
				// Send realtime_update messages. Ignore all other message types
				if r.Type == "realtime_update" {
					recvParsed <- r
				} else {
					glog.Info("Ignored message (not type 'realtime_update'): ", string(msg))
				}
			case <-done:
				return
			}

		}
	}()
	return recvParsed, nil
}
