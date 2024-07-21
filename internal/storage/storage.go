package storage

import (
	"sync"
	"time"

	"github.com/ashmitsharp/net-pack/internal/capture"
)

type Storage struct {
	mu             sync.RWMutex
	packets        []capture.Packet
	protocolCount  map[string]int
	timeSeriesData map[time.Time]int
}

func NewStorage() *Storage {
	return &Storage{
		packets:        make([]capture.Packet, 0),
		protocolCount:  make(map[string]int),
		timeSeriesData: make(map[time.Time]int),
	}
}

func (s *Storage) AddPacket(packet capture.Packet) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.packets = append(s.packets, packet)
	s.protocolCount[packet.Protocol]++

	roundedTime := packet.TimeStp.Round(time.Second)
	s.timeSeriesData[roundedTime]++
}

func (s *Storage) GetPackets() []capture.Packet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	packetsCopy := make([]capture.Packet, len(s.packets))
	copy(packetsCopy, s.packets)
	return packetsCopy
}

func (s *Storage) GetPacket(index int) *capture.Packet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if index >= 0 && index < len(s.packets) {
		packetCopy := s.packets[index]
		return &packetCopy
	}
	return nil
}

func (s *Storage) GetProtocolDistribution() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	distributionCopy := make(map[string]int)
	for k, v := range s.protocolCount {
		distributionCopy[k] = v
	}
	return distributionCopy
}

func (s *Storage) GetTimeSeriesData() map[time.Time]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	timeSeriesCopy := make(map[time.Time]int)
	for k, v := range s.timeSeriesData {
		timeSeriesCopy[k] = v
	}
	return timeSeriesCopy
}

func (s *Storage) FilterPackets(filter func(capture.Packet) bool) []capture.Packet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []capture.Packet
	for _, packet := range s.packets {
		if filter(packet) {
			filtered = append(filtered, packet)
		}
	}
	return filtered
}
