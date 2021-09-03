package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/faiface/beep"
	"os"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/sirupsen/logrus"
)

var minHeight float32 = float32(1200)
var minWidth float32 = float32(2000)

type Gui struct {
	Logger         *logrus.Logger
	Data           map[string]interface{}
	NowPlaying     *fyne.Container
	SelectedSongId int
	NextSongId int
	SampleRate     beep.SampleRate
}

var playing chan int = make(chan int, 1)

func NewApp() {
	logger := logrus.New()
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	logger.Info("Starting Sound Control GUI...")

	myApp := app.NewWithID("Musisticles")
	myWindow := myApp.NewWindow("Media Master")
	myWindow.Resize(fyne.Size{minWidth, minHeight})

	g := &Gui{
		Logger:         logger,
		Data:           map[string]interface{}{},
		//TODO wire up a display of currently playing song
		SelectedSongId: 1,
		NextSongId: 2,
		NowPlaying:     container.New(layout.NewCenterLayout(), widget.NewLabel("Welcome")),
	}

	playing <-0
	g.SetUpSpeaker()

	play := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		g.playSong()
	})
	stop := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
		g.clickStop()
	})
	//controller := container.NewHSplit(play, stop)
	controller := container.NewAdaptiveGrid(2, play, stop)
	//content := container.NewMax(container.New(layout.NewGridLayout(1), controller, g.ReturnSongs()))
	content := container.New(layout.NewBorderLayout(controller, container.NewScroll(g.ReturnSongs()), nil, nil, ), controller, container.NewScroll(g.ReturnSongs()))

	tabs := container.NewAppTabs(
		container.NewTabItem("Songs", content),
		container.NewTabItem("Playlists", container.New(layout.NewCenterLayout(), widget.NewLabel("Playlists"))),
		container.NewTabItem("Now Playing", container.New(layout.NewCenterLayout(), widget.NewLabel("Now Playing"))),
		//container.NewTabItem("Albums", g.ReturnAlbums()),
		//container.NewTabItem("Songs", container.NewScroll(g.ReturnSongs())),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	myWindow.SetContent(tabs)
	//myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
