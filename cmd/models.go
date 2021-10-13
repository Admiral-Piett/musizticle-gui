package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/faiface/beep"
)

// TODO - Explore other font collections
var th = CreateTheme(gofont.Collection())

// TODO - SECRETIFY - environmentalize
var HOST = "http://localhost:9000/api"

var playing = make(chan int, 1)
var displayChange = make(chan bool)

// TODO - SECRETIFY - environmentalize
var NAV_QUEUE_PREVIOUS_LIMIT = 20
var NAV_QUEUE_NEXT_LIMIT = 20

var HOME_TAB = "home"
var NEXT_TAB = "nextBtn"
var PREVIOUS_TAB = "previousBtn"
var TAB_LIST = []string{HOME_TAB, NEXT_TAB, PREVIOUS_TAB}

type App struct {
	// --- Display
	displayList *layout.List
	songs       Songs
	window      *app.Window
	// Speaker Execution
	SampleRate     beep.SampleRate
	speakerControl *beep.Ctrl
	// --- Song Execution
	songList         []*Song
	selectedSong     *Song
	navQueueNext     []*Song
	navQueuePrevious []*Song
	playBtn          widget.Clickable
	stopBtn          widget.Clickable
	nextBtn          widget.Clickable
	previousBtn widget.Clickable
	// --- Tabs
	selectedTab string
	homeTab     widget.Clickable
	nextTab     widget.Clickable
	previousTab widget.Clickable
}

type Songs struct {
	populated           bool
	initSongsInProgress bool
	reload              widget.Clickable
	loadingButton       widget.Clickable
}

type Song struct {
	line          widget.Clickable
	songListIndex int
	Id            int
]	Title    string
	ArtistId int
	ArtistName     string
	AlbumId        int
	AlbumName      string
	TrackNumber    int
	PlayCount      int
	FilePath       string
	CreatedAt      string
	LastModifiedAt string
}
