package consoleui

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
)

type ConsoleUI struct {
	app         *tview.Application
	topPanel    *tview.Form
	bottomPanel *tview.TextView
	url         string
}

func (ui *ConsoleUI) init() {
	ui.app = tview.NewApplication()

	// Create the top panel for static information and settings
	ui.topPanel = tview.NewForm().
		AddTextView("Your URL:", ui.url, 0, 1, false, false).
		AddTextView("Instructions", "Use Tab to jump between options.\nUse Up/Down arrows to select.\nSettings will be applied on the fly.", 0, 3, false, false).
		AddDropDown("Screen size:", []string{"Small", "Medium", "Large"}, 1, nil).
		AddDropDown("Quality:", []string{"Low", "Normal", "High"}, 1, nil)
	ui.topPanel.SetBorder(true).SetTitle("1FPS.video Settings").SetTitleAlign(tview.AlignLeft)

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
		AddItem(ui.topPanel, 0, 1, true).
		AddItem(ui.bottomPanel, 0, 2, false)

	// Set the root of the application to our flex layout
	ui.app.SetRoot(flex, true)

	// Clear bottom panel
	ui.bottomPanel.Clear()
}

func (ui *ConsoleUI) SetUrl(url string) {
	ui.url = url
	if ui.topPanel != nil {
		ui.topPanel.GetFormItem(0).(*tview.TextView).SetText(url)
	}
}

// WriteBottom writes a formatted message to the bottom panel
func (ui *ConsoleUI) WriteBottom(format string, a ...interface{}) {
	fmt.Fprintf(ui.bottomPanel, format, a...)
	ui.app.Draw()
	ui.bottomPanel.ScrollToEnd()
}

// Start initializes and runs the ConsoleUI
func Start() *ConsoleUI {
	ui := &ConsoleUI{}
	ui.init()

	// Run the application in a goroutine
	go func() {
		if err := ui.app.Run(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	return ui
}
