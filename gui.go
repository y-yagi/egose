package main

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/dghubble/go-twitter/twitter"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/y-yagi/gocui"
)

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, openLink); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("tweets", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Tweets"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		for _, tweet := range tweets {
			fmt.Fprintln(v, buildLine(tweet))
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func buildLine(tweet twitter.Tweet) string {
	re := regexp.MustCompile(`\r?\n`)
	return "[" + runewidth.Truncate(tweet.User.Name, 30, "...") + "] " + re.ReplaceAllString(tweet.Text, " ")
}

func openLink(g *gocui.Gui, v *gocui.View) error {
	browser := "google-chrome"

	if v == nil {
		v = g.Views()[0]
	}

	_, cy := v.Cursor()
	tweet := tweets[cy]
	return exec.Command(browser, tweetURL(tweet)).Run()
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		v = g.Views()[0]
	}

	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		v = g.Views()[0]
	}

	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}
