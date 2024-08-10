package consoleui

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
)

type ConsoleUI struct {
	app         *tview.Application
	topPanel    *tview.TextView
	bottomPanel *tview.TextView
}

func (ui *ConsoleUI) init() {
	ui.app = tview.NewApplication()

	// Create the top panel for static information
	ui.topPanel = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(false).
		SetWrap(false)

	// Create the bottom panel for event logs
	ui.bottomPanel = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(false).
		SetWrap(true).
		SetMaxLines(1000).
		SetScrollable(true)

	// Create a flex layout to divide the screen
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.topPanel, 5, 0, false).
		AddItem(ui.bottomPanel, 0, 1, true)

	// Set the root of the application to our flex layout
	ui.app.SetRoot(flex, true)

	// Clear top panel
	ui.topPanel.Clear()
}

// WriteTop writes a formatted message to the top panel

func (ui *ConsoleUI) WriteTop(format string, a ...interface{}) {
	fmt.Fprintf(ui.topPanel, format, a...)
	ui.app.Draw()
}

// WriteBottom writes a formatted message to the bottom panel

func (ui *ConsoleUI) WriteBottom(format string, a ...interface{}) {
	fmt.Fprintf(ui.bottomPanel, format, a...)
	ui.app.Draw()
	ui.bottomPanel.ScrollToEnd()
}

// Similar to New method, but has slightly more logic in the "constructor", so I renamed it to Start to reflect that
// it's not just a dumb method that initializes the UI - so you cannot create hundreds of these objects.

func Start() *ConsoleUI {
	ui := &ConsoleUI{}
	ui.init()

	// tview library app.Run() call is blocking. And it will block the app until you hit Ctrl+C if you don't run it
	// as goroutine.

	go func() {
		if err := ui.app.Run(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	return ui
}
