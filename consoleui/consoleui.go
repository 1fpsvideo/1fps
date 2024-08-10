package consoleui

import (
	"fmt"
	"os"
	"time"

	"github.com/rivo/tview"
)

type ScreenSize int
type Quality int

const (
	Small ScreenSize = iota
	Medium
	Large
)

const (
	Low Quality = iota
	Normal
	High
)

type ConsoleUI struct {
	app         *tview.Application
	topPanel    *tview.Form
	bottomPanel *tview.TextView
	url         string
	ScreenSize  ScreenSize
	Quality     Quality
}

func (ui *ConsoleUI) init() {
	ui.app = tview.NewApplication()

	// Create the top panel for static information and settings
	ui.topPanel = tview.NewForm().
		AddTextView("Your URL:", ui.url, 0, 1, false, false).
		AddTextView("Instructions:", "Use Tab to jump between options.\nUse Up/Down arrows to select, Enter to confirm.\nSettings will be applied on the fly.", 0, 3, false, false).
		AddDropDown("Screen size:", []string{"Small", "Medium", "Large"}, int(ui.ScreenSize), func(text string, index int) {
			ui.ScreenSize = ScreenSize(index)
			go func() { ui.WriteBottom("Screen size set to %s", text) }()
		}).
		AddDropDown("Quality:", []string{"Low", "Normal", "High"}, int(ui.Quality), func(text string, index int) {
			ui.Quality = Quality(index)
			go func() { ui.WriteBottom("Quality set to %s", text) }()
		})
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
	if ui.app == nil || ui.bottomPanel == nil {
		return
	}

	currentTime := time.Now().Format("15:04:05")
	fmt.Fprintf(ui.bottomPanel, "[yellow]%s[white] "+format+"\n", append([]interface{}{currentTime}, a...)...)
	ui.app.Draw()
	ui.bottomPanel.ScrollToEnd()
}

// Start initializes and runs the ConsoleUI
func Start() *ConsoleUI {
	ui := &ConsoleUI{
		ScreenSize: Medium,
		Quality:    Normal,
	}
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
