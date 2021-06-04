package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var artists = []string{"Hans Zimmer", "John Williams", "The Impressionist"}
var albums = []string{"Dark Kight Rises", "Last Samurai"}
var songs = []string{"The Fire Rises", "Bane's Speech"}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TabContainer Widget")
	myWindow.Resize(fyne.Size{2000, 1200})

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.HomeIcon(), widget.NewLabel("Home tab")),
		container.NewTabItem("Artist", getArtists()),
		container.NewTabItem("Albums", getAlbums()),
		container.NewTabItem("Songs", getSongs()),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

func getArtists() *widget.List{
	list := widget.NewList(
		func() int {
			return len(artists)
		},
		// Placeholder??
		func() fyne.CanvasObject {
			return widget.NewLabel("artists")
		},
		// TODO - Fetch artists here
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(artists[i])
		})
	return list
}

func getAlbums() *widget.List{
	list := widget.NewList(
		func() int {
			return len(albums)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("artists")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(albums[i])
		})
	return list
}

func getSongs() *widget.List{
	list := widget.NewList(
		func() int {
			return len(songs)
		},
		// Placeholder??
		func() fyne.CanvasObject {
			return widget.NewLabel("artists")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(songs[i])
		})
	return list
}
