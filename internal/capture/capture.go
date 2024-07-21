package capture

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Packet struct {
	TimeStp  time.Time
	SrcMAC   string
	DstMAC   string
	SrcIP    string
	DestIP   string
	SrcPort  uint16
	DstPort  uint16
	Protocol string
	Length   int
	Payload  []byte
}

type Capture struct {
	handle *pcap.Handle
	filter string
}

func (c *Capture) SetFilter(filter string) error {
	fmt.Println(filter)
	err := c.handle.SetBPFFilter(filter)
	if err != nil {
		return fmt.Errorf("failed to set BPF filter: %v", err)
	}
	c.filter = filter
	return nil
}

func ListInterfaces() ([]pcap.Interface, error) {
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		return nil, fmt.Errorf("failed to list interfaces: %v", err)
	}
	return interfaces, nil
}

func NewCapture(deviceName string) (*Capture, error) {
	handle, err := pcap.OpenLive(deviceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("failed to open device %s: %v", deviceName, err)
	}
	return &Capture{
		handle: handle,
	}, nil
}

func (c *Capture) StartCapture(packetChan chan<- Packet) {
	packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
	for {
		packet, err := packetSource.NextPacket()
		if err != nil {
			log.Printf("Error capturing packet: %v", err)
			if err == pcap.NextErrorTimeoutExpired {
				continue // This error is normal, just try again
			}
			break // For any other error, break the loop
		}
		parsedPacket := c.parsePacket(packet)
		packetChan <- parsedPacket
	}
	log.Println("Packet capture stopped")
}

func (c *Capture) parsePacket(packet gopacket.Packet) Packet {
	parsed := Packet{
		TimeStp: packet.Metadata().Timestamp,
		Length:  packet.Metadata().Length,
	}

	// Ethernet Layer
	if ethernetLayer := packet.Layer(layers.LayerTypeEthernet); ethernetLayer != nil {
		ethernet, _ := ethernetLayer.(*layers.Ethernet)
		parsed.SrcMAC = ethernet.SrcMAC.String()
		parsed.DstMAC = ethernet.DstMAC.String()
	}

	// IP Layer
	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		parsed.SrcIP = ip.SrcIP.String()
		parsed.DestIP = ip.DstIP.String()
	}

	// TCP layer
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		parsed.SrcPort = uint16(tcp.SrcPort)
		parsed.DstPort = uint16(tcp.DstPort)
		parsed.Protocol = "TCP"
	}

	// UDP layer
	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		parsed.SrcPort = uint16(udp.SrcPort)
		parsed.DstPort = uint16(udp.DstPort)
		parsed.Protocol = "UDP"
	}

	// Application layer
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		parsed.Payload = applicationLayer.Payload()
	}

	return parsed
}

func (c *Capture) Close() {
	if c.handle != nil {
		c.handle.Close()
	}
}
