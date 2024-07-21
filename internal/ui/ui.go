package ui

import (
	"fmt"

	"github.com/ashmitsharp/net-pack/internal/capture"
	"github.com/jroimartin/gocui"
)

type GUI struct {
	gui *gocui.Gui
}

func NewGui() (*GUI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	return &GUI{
		gui: g,
	}, nil
}

func (g *GUI) Close() {
	g.gui.Close()
}

func (g *GUI) Run(packetChn <-chan capture.Packet) error {
	g.gui.SetManagerFunc(g.layout)

	if err := g.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, g.quit); err != nil {
		return err
	}

	go g.updatePackets(packetChn)

	if err := g.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (g *GUI) layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	if v, err := gui.SetView("packets", 0, 0, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Packet List"
		v.Wrap = false
		v.Autoscroll = true
	}

	if v, err := gui.SetView("status", 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Status"
		fmt.Fprintln(v, "Press Ctrl+C to quit")
	}

	return nil
}

func (g *GUI) updatePackets(packetChn <-chan capture.Packet) {
	for packet := range packetChn {
		g.gui.Update(func(gui *gocui.Gui) error {
			v, err := gui.View("packets")
			if err != nil {
				return err
			}
			fmt.Fprintf(v, "%s | %s -> %s | %s | %d bytes\n",
				packet.TimeStp.Format("15:04:05"),
				packet.SrcIP,
				packet.DestIP,
				packet.Protocol,
				packet.Length)
			return nil
		})
	}
}

func (g *GUI) quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
