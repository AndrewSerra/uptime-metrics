package utils

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const listenAddr = "0.0.0.0"

// addr parameter is an ip address
func Ping(addr string) (*time.Duration, error) {
	// Open connection to listen to icmp replies
	conn, err := icmp.ListenPacket("ip4:icmp", listenAddr)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	// Create message of type Echo
	message := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  0,
			Data: []byte(""),
		},
	}

	messageBin, err := message.Marshal(nil)

	if err != nil {
		return nil, err
	}

	dest, err := net.ResolveIPAddr("ip4", addr)

	if err != nil {
		return nil, err
	}

	// Start sending process
	startTime := time.Now()

	dataSize, err := conn.WriteTo(messageBin, dest)

	if err != nil {
		return nil, err
	} else if len(messageBin) != dataSize {
		return nil, fmt.Errorf("Message size does not match. Received %d, actual %d", dataSize, len(messageBin))
	}

	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	if err != nil {
		return nil, err
	}

	replyMessage := make([]byte, len(messageBin))
	_, _, err = conn.ReadFrom(replyMessage)

	if err != nil {
		return nil, err
	}

	// Transmission complete - get end time
	pingDuration := time.Since(startTime)

	return &pingDuration, nil
}
