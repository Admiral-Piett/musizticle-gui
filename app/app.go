package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"os"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/sirupsen/logrus"
)

var minHeight float32 = float32(1200)
var minWidth float32 = float32(2000)

type Gui struct {
	Logger *logrus.Logger
	Data	map[string]interface{}
	NowPlaying *fyne.Container
	SelectedSong int
}

func NewApp() {
	logger := logrus.New()
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	logger.Info("Starting Sound Control GUI...")

	myApp := app.New()
	myWindow := myApp.NewWindow("Media Master")
	myWindow.Resize(fyne.Size{minWidth, minHeight})

	g := &Gui{
		Logger: logger,
		Data: map[string]interface{}{},
		SelectedSong: 0,
		NowPlaying: container.New(layout.NewCenterLayout(), widget.NewLabel("Welcome")),
	}

	tabs := container.NewAppTabs(
		//container.NewTabItemWithIcon("", theme.HomeIcon(), g.NowPlaying),
		container.NewTabItemWithIcon("Songs", theme.HomeIcon(), container.NewScroll(g.ReturnSongs())),
		container.NewTabItem("Artist", g.ReturnArtists()),
		container.NewTabItem("Albums", g.ReturnAlbums()),
		//container.NewTabItem("Songs", container.NewScroll(g.ReturnSongs())),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
