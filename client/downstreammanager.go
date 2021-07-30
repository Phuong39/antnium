package client

import (
	"github.com/dobin/antnium/model"
)

type DownstreamManager struct {
	downstreamClient   DownstreamClient
	downstreamLocaltcp DownstreamLocaltcp
}

func MakeDownstreamManager() DownstreamManager {
	downstreamClient := MakeDownstreamClient()
	downstreamLocaltcp := MakeDownstreamLocaltcp()

	downstreamManager := DownstreamManager{
		downstreamClient,
		downstreamLocaltcp,
	}
	return downstreamManager
}

func (dm *DownstreamManager) start() {
	go dm.downstreamClient.start()
	go dm.downstreamLocaltcp.start()
}

func (dm *DownstreamManager) GetFor(p model.Packet) chan model.Packet {
	channelId, ok := p.Arguments["channelId"]
	if !ok {
		return dm.downstreamClient.channel
	}

	if channelId == "client" {
		return dm.downstreamClient.channel
	} else if channelId == "net#1" {
		return dm.downstreamLocaltcp.channel
	} else {
		return dm.downstreamClient.channel
	}
}

func (dm *DownstreamManager) GetList() []string {
	ret := make([]string, 0)

	ret = append(ret, "client")
	ret = append(ret, "net#1")

	return ret
}