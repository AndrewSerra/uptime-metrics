package utils

import (
	"errors"
)

type UniqueChannel struct {
	Id string
	C  chan bool
}

type ChannelCoordinator struct {
	channels []*UniqueChannel
}

func (cc *ChannelCoordinator) Add(id string) (*UniqueChannel, error) {

	for _, channel := range cc.channels {
		if channel.Id == id {
			return nil, errors.New("Id exists.")
		}
	}

	uc := &UniqueChannel{
		Id: id,
		C:  make(chan bool),
	}

	cc.channels = append(cc.channels, uc)

	return uc, nil
}

func (cc *ChannelCoordinator) Remove(id string) error {
	targetIndex := -1

	for index, channel := range cc.channels {
		if channel.Id == id {
			targetIndex = index
			// send signal to end ticker
			channel.C <- false
			break
		}
	}

	if targetIndex != -1 {
		cc.channels = append(cc.channels[:targetIndex], cc.channels[targetIndex+1:]...)
		return nil
	} else {
		return errors.New("Channel with id does not exist.")
	}
}
