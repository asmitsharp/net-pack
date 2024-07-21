package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashmitsharp/net-pack/internal/capture"
	"github.com/ashmitsharp/net-pack/internal/storage"
	"github.com/ashmitsharp/net-pack/internal/ui"
)

func main() {
	// Set up logging
	logFile, err := os.Create("network_analyzer.log")
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Getting all the Interfaces
	interfaces, err := capture.ListInterfaces()
	if err != nil {
		log.Fatalf("Failed to list all Interfaces: %v", err)
	}

	// List all Interfaces
	fmt.Println("Available interfaces:")
	for i, iface := range interfaces {
		fmt.Printf("%d: %s (%s)\n", i, iface.Name, iface.Description)
	}

	// User input for Interface selection
	var selectedIndex int
	fmt.Print("Select an interface by number: ")
	_, err = fmt.Scanf("%d", &selectedIndex)
	if err != nil || selectedIndex < 0 || selectedIndex >= len(interfaces) {
		log.Fatalf("Invalid selection: %v", err)
	}

	// Initialize packet capture
	capturer, err := capture.NewCapture(interfaces[selectedIndex].Name)
	if err != nil {
		log.Fatalf("Failed to initialize packet capturer: %v", err)
	}
	defer capturer.Close()

	// Getting Filter from user
	var filter string
	fmt.Print("Enter a BPF filter (e.g., 'tcp port 80' for HTTP traffic, or press Enter for all traffic): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		filter = scanner.Text()
	}

	if filter != "" {
		err = capturer.SetFilter(filter)
		if err != nil {
			log.Fatalf("Failed to set filter: %v", err)
		}
	}

	// Initialize storage
	store := storage.NewStorage()

	// Initialize UI
	gui, err := ui.NewGui(store)
	if err != nil {
		log.Fatalf("Failed to initialize UI: %v", err)
	}

	// Set up channel for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start capturing packets
	packetChan := make(chan capture.Packet)
	go func() {
		capturer.StartCapture(packetChan)
		close(packetChan)
	}()

	// Run the UI in a separate goroutine
	errChan := make(chan error)
	go func() {
		errChan <- gui.Run(packetChan)
	}()

	// Wait for either an error from the UI or a signal to shut down
	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("Application error: %v", err)
		}
	case <-sigChan:
		log.Println("Received interrupt, shutting down...")
	}

	// Perform cleanup
	gui.Close()
	log.Println("Application stopped")
}
