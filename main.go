package main

import (
	"time"

	"github.com/coip/rebus/ibsvc"
)

func main() {

	// internal api log is zap log, you could use GetLogger to get the logger
	// besides, you could use SetAPILogger to set you own log option
	// or you can just use the other logger
	log := ibsvc.GetLogger().Sugar()
	defer log.Sync()

	// implement your own IbWrapper to handle the msg delivered via tws or gateway
	// Wrapper{} below is a default implement which just log the msg
	ic := ibsvc.NewIbClient(&ibsvc.Wrapper{})

	// tcp connect with tws or gateway
	// fail if tws or gateway had not yet set the trust IP
	if err := ic.Connect("localhost", 7496, 0); err != nil {
		log.Panic("Connect failed:", err)
	}

	// handshake with tws or gateway, send handshake protocol to tell tws or gateway the version of client
	// and receive the server version and connection time from tws or gateway.
	// fail if someone else had already connected to tws or gateway with same clientID
	if err := ic.HandShake(); err != nil {
		log.Panic("HandShake failed:", err)
	}

	// make some request, msg would be delivered via wrapper.
	// req will not send to TWS or Gateway until ic.Run()
	// you could just call ic.Run() before these
	ic.ReqCurrentTime()

	// start to send req and receive msg from tws or gateway after this
	ic.Run()
	for {
		select {
		case <-time.After(time.Second * 2):

			ic.Disconnect()

		}
	}
}
