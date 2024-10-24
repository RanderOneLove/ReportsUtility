package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

const W = 300
const H = 300

type Config struct {
	Replacements map[string]string `json:"replacements"`
}

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())

	w := a.NewWindow("ReportsUtility")
	wIcon, _ := fyne.LoadResourceFromPath("assets/icon.png")
	w.SetIcon(wIcon)
	w.Resize(fyne.NewSize(W, H))

	colEmojiTab := createColEmojiTab()
	imgPath := "assets/important.jpg"
	img := canvas.NewImageFromFile(imgPath)
	img.FillMode = canvas.ImageFillContain
	labelUpdated := widget.NewLabel("There's nothing here, go away")

	updatedTab := container.NewMax(
		img,
		labelUpdated,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem(
			"ColEmoji",
			colEmojiTab,
		),
		container.NewTabItem(
			"ClickMe",
			updatedTab,
		),
	)

	w.SetContent(tabs)
	w.Show()
	a.Run()

}

func createColEmojiTab() *fyne.Container {
	entry := widget.NewEntry()
	entry.MultiLine = true
	labelContainer := container.NewHBox()
	button := widget.NewButton("Updated Text", func() {
		input := "&xFFFFFF" + entry.Text
		labelContainer.Objects = []fyne.CanvasObject{}
		labelUpdater(input, labelContainer)
		labelContainer.Refresh()
	})

	go func() {
		for {
			time.AfterFunc(2*time.Minute, func() {
				labelContainer.Objects = []fyne.CanvasObject{}
				labelContainer.Refresh()
				entry.SetText("")
			})
			time.Sleep(2 * time.Minute)
		}
	}()

	return container.NewVBox(
		entry,
		button,
		labelContainer,
	)
}

func loadConfig(filename string) (Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func labelUpdater(input string, container *fyne.Container) {
	config, err := loadConfig("cfg/config.json")
	if err != nil {
		fmt.Println("Какое-то говно, извинись:", err)
		return
	}
	for key, value := range config.Replacements {
		input = strings.ReplaceAll(input, key, value)
	}
	re := regexp.MustCompile(`&x([0-9A-Fa-f]{6})([^&]*)|([^&]+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if len(match) == 4 {
			if match[1] != "" {
				hexColor := match[1]
				text := match[2]
				r, g, b := hexToRGB(hexColor)
				col := color.RGBA{R: r, G: g, B: b, A: 255}
				coloredLabel := canvas.NewText(text, col)
				container.Add(coloredLabel)
			}
		}
	}
}

func hexToRGB(hex string) (uint8, uint8, uint8) {
	var r, g, b uint8
	fmt.Sscanf(hex, "%2X%2X%2X", &r, &g, &b)
	return r, g, b
}
