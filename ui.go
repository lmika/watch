package main

import (
	"fmt"

	"github.com/mgutz/ansi"

	"github.com/jroimartin/gocui"
)

// UICtx Context manages the UI
type UICtx struct {
	Sampler *Sampler
	Gui     *gocui.Gui

	errCol func(s string) string
}

// Init initialises the UI context
func (uc *UICtx) Init() {
	uc.Gui.SetManagerFunc(uc.LayoutFn)

	uc.Sampler.OnSampled = uc.Update

	uc.errCol = ansi.ColorFunc("red")
}

// LayoutFn is the update function for gocui
func (uc *UICtx) LayoutFn(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	titleView(uc, g, maxX, maxY)
	outputView(uc, g, maxX, maxY)

	return nil
}

// Update requests an update
func (uc *UICtx) Update() {
	uc.Gui.Update(uc.LayoutFn)
}

// format the title
func (uc *UICtx) formatTitle() string {
	if lastSnapshot := uc.Sampler.LastSnapshot(); lastSnapshot != nil {
		return lastSnapshot.Command
	}

	return uc.Sampler.Command
}

// format the date
func (uc *UICtx) formatDate() string {
	if lastSnapshot := uc.Sampler.LastSnapshot(); lastSnapshot != nil {
		return fmt.Sprintf("%s (every %.1f sec)",
			lastSnapshot.Started.Format("2006-01-02 15:04:05"), uc.Sampler.Interval.Seconds())
	}

	return ""
}
