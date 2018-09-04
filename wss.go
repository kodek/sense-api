package sense_api

import (
	"net/url"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

func wssConnect(done <-chan struct{}, u url.URL) (<-chan []byte, error) {
	glog.Infof("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	sendChan := make(chan []byte)

	// When received, send new messages back to the user.
	go func() {
		for {
			defer c.Close()
			defer close(sendChan)
			_, message, err := c.ReadMessage()
			if err != nil {
				glog.Error("Read error:", err)
				return
			}
			select {
			case sendChan <- message:
			case <-done:
				return
			default:
			}
		}
	}()
	return sendChan, err
}
