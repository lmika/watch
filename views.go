package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/willf/pad"
)

func titleView(uictx *UICtx, g *gocui.Gui, maxX, maxY int) error {
	v, err := g.SetView("title", -1, -1, maxX, 1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		v.BgColor = gocui.ColorDefault | gocui.AttrReverse
		v.FgColor = gocui.ColorDefault | gocui.AttrReverse
	}

	titleString := uictx.formatTitle()
	dateString := uictx.formatDate()

	if maxTitleLen := maxX - len(dateString) - 3; len(titleString) > maxTitleLen {
		dateString = ""
	}

	v.Clear()

	fmt.Fprint(v, titleString)
	if dateString != "" {
		fmt.Fprintln(v, pad.Left(dateString, maxX-len(titleString), " "))
	}
	return nil
}

func outputView(uictx *UICtx, g *gocui.Gui, maxX, maxY int) error {
	v, err := g.SetView("output", -1, 1, maxX, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
	}

	// Maximum number of lines to display until they are clipped
	var rows []Line
	var lastLine string

	var snapshot *Snapshot
	if uictx.VisibleSavedSnapshot >= 0 {
		snapshot = uictx.Sampler.SavedSnapshot(uictx.VisibleSavedSnapshot)
	} else {
		snapshot = uictx.Sampler.LastSnapshot()
	}

	if snapshot != nil {
		rows = snapshot.Lines
		if snapshot.Err != nil {
			lastLine = uictx.errCol(fmt.Sprintf("error: %v", snapshot.Err))
		}
	}

	maxRows := maxY - 2
	if lastLine != "" {
		maxRows--
	}

	v.Clear()
	for row, line := range rows {
		if row >= maxRows {
			break
		}
		fmt.Fprintln(v, line.Line)
	}

	fmt.Fprintln(v, lastLine)

	return nil
}
