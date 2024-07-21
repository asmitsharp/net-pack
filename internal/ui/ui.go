package ui

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/ashmitsharp/net-pack/internal/capture"
	"github.com/ashmitsharp/net-pack/internal/storage"
	"github.com/jroimartin/gocui"
)

type GUI struct {
	gui     *gocui.Gui
	storage *storage.Storage
	mu      sync.Mutex
}

func NewGui(storage *storage.Storage) (*GUI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	return &GUI{gui: g, storage: storage}, nil
}

func (g *GUI) layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	if v, err := gui.SetView("packets", 0, 0, maxX/2-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Packet List"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		g.gui.SetCurrentView("packets")
	}

	if v, err := gui.SetView("details", maxX/2, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Packet Details"
		v.Wrap = true
	}

	if v, err := gui.SetView("protocols", 0, maxY/2, maxX/2-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Protocol Distribution"
	}

	if v, err := gui.SetView("timeseries", maxX/2, maxY/2, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Packet Rate Over Time"
	}

	if v, err := gui.SetView("status", 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Status"
		fmt.Fprintln(v, "Press Ctrl+C to quit, F to filter")
	}

	return nil
}

func (g *GUI) updatePacketList() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.gui.Update(func(gui *gocui.Gui) error {
		v, err := gui.View("packets")
		if err != nil {
			return err
		}
		v.Clear()

		packets := g.storage.GetPackets()
		for i, packet := range packets {
			fmt.Fprintf(v, "%d: %s | %s:%d -> %s:%d | %s | %d bytes\n",
				i,
				packet.TimeStp.Format("15:04:05"),
				packet.SrcIP, packet.SrcPort,
				packet.DestIP, packet.DstPort,
				packet.Protocol,
				packet.Length)
		}
		return nil
	})
}

func (g *GUI) updatePacketDetails(packet capture.Packet) {
	g.gui.Update(func(gui *gocui.Gui) error {
		v, err := gui.View("details")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintf(v, "Time: %s\n", packet.TimeStp)
		fmt.Fprintf(v, "Source MAC: %s\n", packet.SrcMAC)
		fmt.Fprintf(v, "Destination MAC: %s\n", packet.DstMAC)
		fmt.Fprintf(v, "Source IP: %s\n", packet.SrcIP)
		fmt.Fprintf(v, "Destination IP: %s\n", packet.DestIP)
		fmt.Fprintf(v, "Protocol: %s\n", packet.Protocol)
		fmt.Fprintf(v, "Source Port: %d\n", packet.SrcPort)
		fmt.Fprintf(v, "Destination Port: %d\n", packet.DstPort)
		fmt.Fprintf(v, "Length: %d bytes\n", packet.Length)
		fmt.Fprintf(v, "Payload: %s\n", string(packet.Payload))
		return nil
	})
}

func (g *GUI) updateProtocolDistribution() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.gui.Update(func(gui *gocui.Gui) error {
		v, err := gui.View("protocols")
		if err != nil {
			return err
		}
		v.Clear()
		distribution := g.storage.GetProtocolDistribution()
		for protocol, count := range distribution {
			fmt.Fprintf(v, "%s: %d\n", protocol, count)
		}
		return nil
	})
}

func (g *GUI) updateTimeSeriesData() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.gui.Update(func(gui *gocui.Gui) error {
		v, err := gui.View("timeseries")
		if err != nil {
			return err
		}
		v.Clear()
		timeSeriesData := g.storage.GetTimeSeriesData()
		for timestamp, count := range timeSeriesData {
			fmt.Fprintf(v, "%s: %d\n", timestamp.Format("15:04:05"), count)
		}
		return nil
	})
}

func (g *GUI) filterPrompt(gui *gocui.Gui, v *gocui.View) error {
	maxX, maxY := gui.Size()
	if v, err := gui.SetView("filter", maxX/4, maxY/4, 3*maxX/4, maxY/4+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = "Enter filter (e.g., 'ip 192.168.1.1' or 'port 80')"
		gui.SetCurrentView("filter")
	}
	return nil
}

func (g *GUI) applyFilter(gui *gocui.Gui, v *gocui.View) error {
	filterStr := strings.TrimSpace(v.Buffer())
	gui.DeleteView("filter")
	gui.SetCurrentView("packets")

	filter := func(p capture.Packet) bool {
		if strings.Contains(filterStr, "ip") && strings.Contains(filterStr, p.SrcIP) {
			return true
		}
		if strings.Contains(filterStr, "port") {
			portStr := strings.Split(filterStr, " ")[1]
			return fmt.Sprintf("%d", p.SrcPort) == portStr || fmt.Sprintf("%d", p.DstPort) == portStr
		}
		return false
	}

	filteredPackets := g.storage.FilterPackets(filter)
	g.updateFilteredPackets(filteredPackets)
	return nil
}

func (g *GUI) updateFilteredPackets(packets []capture.Packet) {
	g.gui.Update(func(gui *gocui.Gui) error {
		v, err := gui.View("packets")
		if err != nil {
			return err
		}
		v.Clear()
		for i, packet := range packets {
			fmt.Fprintf(v, "%d: %s | %s:%d -> %s:%d | %s | %d bytes\n",
				i,
				packet.TimeStp.Format("15:04:05"),
				packet.SrcIP, packet.SrcPort,
				packet.DestIP, packet.DstPort,
				packet.Protocol,
				packet.Length)
		}
		return nil
	})
}

func (g *GUI) Close() {
	g.gui.Close()
}

func (g *GUI) setKeybindings() error {
	if err := g.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, g.quit); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("packets", gocui.KeyArrowUp, gocui.ModNone, g.cursorUp); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("packets", gocui.KeyArrowDown, gocui.ModNone, g.cursorDown); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("", gocui.KeyF1, gocui.ModNone, g.filterPrompt); err != nil {
		return err
	}

	if err := g.gui.SetKeybinding("filter", gocui.KeyEnter, gocui.ModNone, g.applyFilter); err != nil {
		return err
	}

	return nil
}

func (g *GUI) Run(packetChan <-chan capture.Packet) error {
	g.gui.SetManagerFunc(g.layout)

	if err := g.setKeybindings(); err != nil {
		return err
	}

	go func() {
		for packet := range packetChan {
			g.storage.AddPacket(packet)
			g.updatePacketList()
			g.updateProtocolDistribution()
			g.updateTimeSeriesData()
		}
	}()

	if err := g.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (g *GUI) cursorUp(gui *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
		g.updateSelectedPacket(v)
	}
	return nil
}

func (g *GUI) cursorDown(gui *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
		g.updateSelectedPacket(v)
	}
	return nil
}

func (g *GUI) updateSelectedPacket(v *gocui.View) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		return
	}

	if len(l) == 0 {
		return
	}

	parts := strings.Split(l, ":")
	if len(parts) < 2 {
		return
	}

	packetIndex, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return
	}

	packet := g.storage.GetPacket(packetIndex)
	if packet != nil {
		g.updatePacketDetails(*packet)
	}
}

func (g *GUI) quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
