package main

import (
	"log"

	"github.com/ashmitsharp/net-pack/internal/capture"
	"github.com/ashmitsharp/net-pack/internal/ui"
)

func main() {
	// Initialize packet capture
	capturer, err := capture.NewCapture()
	if err != nil {
		log.Fatalf("Failed to initialize packet capturer: %v", err)
	}

	// Initialize UI
	gui, err := ui.NewGui()
	if err != nil {
		log.Fatalf("Failed to initialize UI: %v", err)
	}
	defer gui.Close()

	// Start capturing packets
	packetChan := make(chan capture.Packet)
	go capturer.StartCapture(packetChan)

	// Run the UI
	if err := gui.Run(packetChan); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
