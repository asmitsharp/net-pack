# Packet Capture and Inspection Tool

## Overview

This command-line application, implemented in Go, allows users to capture and inspect network packets in real-time. The tool offers various features for monitoring network traffic, including packet capture, detailed inspection, and filtering.

## Features

### Program Launch

- Launch the program from the terminal.
- Displays a welcome message and available network interfaces.

### Interface Selection

- Select a network interface to monitor.
- Starts capturing packets on the chosen interface.

### Main Interface

The main interface is divided into sections:

- **Real-time Packet List**: Displays packets as they are captured.
- **Detailed View of Selected Packet**: Shows detailed information about the selected packet.
- **Protocol Distribution**: Visual representation of the distribution of different protocols.
- **Time Series Data**: Graphical representation of packet data over time.

### Packet Capture and Display

- New packets appear in the real-time list as they are captured.
- Each line shows:
  - Timestamp
  - Source IP
  - Destination IP
  - Protocol
  - Length

### Packet Inspection

- Navigate the packet list using arrow keys.
- The selected packet's details are displayed in the middle section.
- Details include all protocol layers:
  - Ethernet
  - IP
  - TCP/UDP
  - Application

### Filtering

- Apply filters to update the packet list to show only matching packets.
- Recalculate statistics based on the current filter.

## Future Enhancements

- **Save Capture**: Save the current capture to a file (`s` command).
- **Load Capture**: Load a previously saved capture file (`l` command).
- **Pause/Resume Capture**: Pause or resume live capture (`p` command).

## Usage

1. Launch the program from the terminal.
2. Select the network interface you want to monitor.
3. Use the real-time packet list, detailed packet view, protocol distribution, and time series data sections to monitor and analyze network traffic.
4. Navigate and inspect packets using arrow keys.
5. Apply filters to narrow down the packet list and view updated statistics.

## Requirements

- Go programming language (version 1.XX or later)

## Installation

Clone the repository:

````bash
git clone https://github.com/ashmitsharp/net-pack.git

```bash
cd net-pack

```bash
cd build -o net-pack

```bash
./net-pack
````
