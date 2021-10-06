package main

import (
	"gioui.org/layout"
	"gioui.org/widget"
)
import "gioui.org/app"

var displayChange = make(chan bool)

type App struct {
	displayList *layout.List
	songs Songs
	window *app.Window
}

type Songs struct {
	songList []Song
	selected int
	populated bool
	inProgress    bool
	reload widget.Clickable
	loadingButton widget.Clickable
}

type Song struct {
	line  widget.Clickable
	Id int
	Name string
	ArtistId int
	ArtistName string
	AlbumId int
	AlbumName string
	TrackNumber int
	PlayCount int
	FilePath string
	CreatedAt string
	LastModifiedAt string
}
