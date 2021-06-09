package app

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"net/http"
)

var hostApi = "http://localhost:9000/api"

//var artists = []string{"Hans Zimmer", "John Williams", "The Impressionist"}
var albums = []string{"Dark Kight Rises", "Last Samurai"}
var songs = []string{"The Fire Rises", "Bane's Speech"}

type Artist struct {
	ID             int    `json:"Id"`
	Name           string `json:"Name"`
	CreatedAt      string `json:"CreatedAt"`
	LastModifiedAt string `json:"LastModifiedAt"`
}

type Artists []Artist

func getArtists() (Artists, error) {
	artists := Artists{}
	url := fmt.Sprintf("%s/artists", hostApi)
	resp, err := http.Get(url)
	if err != nil {
		return artists, err
	}
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return artists, err
	}
	fmt.Println(artists)
	return artists, nil
}

func (g *Gui) ReturnArtists() *fyne.Container {
	artists, err := getArtists()
	if err != nil {
		g.Logger.WithFields(logrus.Fields{LogFields.ErrorMessage: err}).Error("GetArtistsFailure")
		box := container.New(layout.NewCenterLayout(), widget.NewButton("Retry?", func() {
			g.Logger.Info("Retrying")
		}))
		return box
	}
	g.Data["artists"] = artists
	list := widget.NewList(
		func() int {
			return len(artists)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("artists")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(artists[i].Name)
		})
	list.OnSelected = g.clickedArtist
	box := container.New(layout.NewMaxLayout(), list)
	box.Resize(fyne.NewSize(minWidth, minHeight))
	return box
}

func (g *Gui) clickedArtist(i int) {
	fmt.Printf("CLICKED - %d", i)
	//Covert data from interface into a list of structs that I can access
	artistList := g.Data["artists"].(Artists)
	fmt.Printf("DATA - %v+", artistList[i])
	g.getSongsByArtist(artistList[i])
}

//TODO - probably break this up before this client gets nuts?
type Song struct {
	ID             int    `json:"Id"`
	Name           string `json:"Name"`
	ArtistID       int    `json:"ArtistId"`
	AlbumID        int    `json:"AlbumId"`
	TrackNumber    int    `json:"TrackNumber"`
	PlayCount      int    `json:"PlayCount"`
	FilePath       string `json:"FilePath"`
	CreatedAt      string `json:"CreatedAt"`
	LastModifiedAt string `json:"LastModifiedAt"`
}
type Songs []Song

func (g *Gui) getSongsByArtist(a Artist) {
	fmt.Println(a)
	songs := Songs{}
	url := fmt.Sprintf("%s/songs/artists/%d", hostApi, a.ID)
	resp, err := http.Get(url)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{LogFields.ErrorMessage: err}).Error("GetSongsByArtistFailure")
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&songs)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{LogFields.ErrorMessage: err}).Error("GetSongsByArtistFailure")
		return
	}
	fmt.Println(songs)
	//TODO - update, refresh, and show now playing screen
	return

}

func (g *Gui) ReturnAlbums() *widget.List{
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

func (g *Gui) ReturnSongs() *widget.List{
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


