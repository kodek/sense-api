package sense_api

import (
	"net/url"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

func wssConnect(u url.URL) (<-chan []byte, chan<- struct{}, error) {
	glog.Infof("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	done := make(chan struct{})
	sendChan := make(chan []byte)

	// Close the channel when the user is done.
	go func() {
		defer c.Close()
		<-done
	}()

	// When received, send new messages back to the user.
	go func() {
		for {
			defer c.Close()
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

	return sendChan, done, err
}
