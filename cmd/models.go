package main

import "gioui.org/widget"

var displayChange = make(chan bool)

type Songs struct {
	songList []Song
	selected int
	populated bool
	inProgress    bool
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
