package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/faiface/beep"
)

// TODO - Explore other font collections
//var th = material.NewTheme(gofont.Collection())
var th = CreateTheme(gofont.Collection())

// TODO - environmentalize
var HOST = "http://localhost:9000/api"

var playing = make(chan int, 1)
var displayChange = make(chan bool)

var NAV_QUEUE_PREVIOUS_LIMIT = 20
var NAV_QUEUE_NEXT_LIMIT = 20
var navQueuePrevious = make(chan int)
var navQueueNext = make(chan int)

type App struct {
	displayList    *layout.List
	songs          Songs
	window         *app.Window
	SelectedSongId int
	SelectedSongIndex int
	NextSongId     int
	SampleRate     beep.SampleRate
	play           widget.Clickable
	stop           widget.Clickable
	next           widget.Clickable
	previous       widget.Clickable
}

type Songs struct {
	songList            []Song
	selected            int
	populated           bool
	initSongsInProgress bool
	reload              widget.Clickable
	loadingButton       widget.Clickable
}

type Song struct {
	line widget.Clickable
	Id   int
	// TODO - rename to Title
	Name           string
	ArtistId       int
	ArtistName     string
	AlbumId        int
	AlbumName      string
	TrackNumber    int
	PlayCount      int
	FilePath       string
	CreatedAt      string
	LastModifiedAt string
}
