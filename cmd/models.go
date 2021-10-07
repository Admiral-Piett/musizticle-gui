package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/faiface/beep"
)

// TODO - Explore other font collections
var th = material.NewTheme(gofont.Collection())
// TODO - environmentalize
var HOST = "http://localhost:9000/api"

var playing = make(chan int, 1)
var displayChange = make(chan bool)

type App struct {
	displayList    *layout.List
	songs          Songs
	window         *app.Window
	SelectedSongId int
	NextSongId     int
	SampleRate     beep.SampleRate
}

type Songs struct {
	songList      []Song
	selected      int
	populated     bool
	inProgress    bool
	reload        widget.Clickable
	loadingButton widget.Clickable
}

type Song struct {
	line           widget.Clickable
	Id             int
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
