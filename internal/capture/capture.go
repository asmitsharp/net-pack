package capture

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type Packet struct {
	TimeStp  time.Time
	SrcIP    string
	DestIP   string
	Protocol string
	Length   int
}

type Capture struct {
	handle *pcap.Handle
}

func NewCapture() (*Capture, error) {
	// finding all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	device := devices[0].Name

	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	return &Capture{
		handle: handle,
	}, nil
}

func (c *Capture) StartCapture(packetChn chan<- Packet) {
	packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
	for packet := range packetSource.Packets() {

		networkLayer := packet.NetworkLayer()
		if networkLayer == nil {
			continue
		}

		transportLayer := packet.TransportLayer()
		if transportLayer == nil {
			continue
		}

		srcIp := networkLayer.NetworkFlow().Src().String()
		dstIp := networkLayer.NetworkFlow().Dst().String()
		protocol := transportLayer.LayerType().String()

		packetChn <- Packet{
			TimeStp:  packet.Metadata().Timestamp,
			SrcIP:    srcIp,
			DestIP:   dstIp,
			Protocol: protocol,
			Length:   packet.Metadata().Length,
		}
	}
}
