package main

import (
	"flag"
	"fmt"
	"github.com/lmika/gopkgs/fp/slices"
	"log"
	"mvdan.cc/sh/v3/syntax"
	"os"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [FLAGS] CMD ARGS...\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	waitSec := flag.Int("s", 2, "Wait time between updates")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	quotedCli, err := slices.MapWithError(flag.Args(), func(s string) (string, error) {
		return syntax.Quote(s, syntax.LangBash)
	})
	if err != nil {
		log.Fatal(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// Setup the sampler
	sampler := &Sampler{
		Command:  strings.Join(quotedCli, " "),
		Interval: time.Duration(*waitSec) * time.Second,
	}
	sampler.Init()

	// Setup the UI
	uiCtx := &UICtx{
		Sampler: sampler,
		Gui:     g,
	}
	uiCtx.Init()

	// Start sampling
	go sampler.Start()

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
