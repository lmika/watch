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

	Message              string
	VisibleSavedSnapshot int

	errCol func(s string) string
}

// Init initialises the UI context
func (uc *UICtx) Init() {
	uc.Gui.SetManagerFunc(uc.LayoutFn)

	uc.VisibleSavedSnapshot = -1
	uc.Sampler.OnSampled = uc.sampleUpdated

	uc.errCol = ansi.ColorFunc("red")
}

// LayoutFn is the update function for gocui
func (uc *UICtx) LayoutFn(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	titleView(uc, g, maxX, maxY)
	outputView(uc, g, maxX, maxY)

	return nil
}

func (uc *UICtx) sampleUpdated() {
	uc.Message = ""
	uc.Update()
}

// Update requests an update
func (uc *UICtx) Update() {
	uc.Gui.Update(uc.LayoutFn)
}

// format the title
func (uc *UICtx) formatTitle() string {
	if uc.Message != "" {
		return uc.Message
	}

	if uc.VisibleSavedSnapshot >= 0 {
		return fmt.Sprintf("Snapshot %v of %v", uc.VisibleSavedSnapshot+1, uc.Sampler.NSavedSnapshot())
	}

	if lastSnapshot := uc.Sampler.LastSnapshot(); lastSnapshot != nil {
		return lastSnapshot.Command
	}

	return uc.Sampler.Command
}

// format the date
func (uc *UICtx) formatDate() string {
	if uc.VisibleSavedSnapshot >= 0 {
		snapshot := uc.Sampler.SavedSnapshot(uc.VisibleSavedSnapshot)
		return snapshot.Started.Format("2006-01-02 15:04:05")
	}

	if lastSnapshot := uc.Sampler.LastSnapshot(); lastSnapshot != nil {
		return fmt.Sprintf("%s (every %.1f sec)",
			lastSnapshot.Started.Format("2006-01-02 15:04:05"), uc.Sampler.Interval.Seconds())
	}

	return ""
}

func (uc *UICtx) CycleSnapshots(d int) {
	n := uc.Sampler.NSavedSnapshot()
	if d < 0 {
		if uc.VisibleSavedSnapshot < 0 {
			uc.VisibleSavedSnapshot = n - 1
		} else {
			uc.VisibleSavedSnapshot -= 1
		}
	} else if d > 0 {
		if uc.VisibleSavedSnapshot < n-1 {
			uc.VisibleSavedSnapshot += 1
		} else {
			uc.VisibleSavedSnapshot = -1
		}
	}
	uc.Update()
}
