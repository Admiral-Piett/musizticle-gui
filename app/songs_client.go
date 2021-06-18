package app

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Song struct {
	ID             int    `json:"Id"`
	Name           string `json:"Name"`
	ArtistID       int    `json:"ArtistId"`
	ArtistName     string `json:"ArtistName,omitempty"`
	AlbumID        int    `json:"AlbumId"`
	AlbumName      string `json:"AlbumName,omitempty"`
	TrackNumber    int    `json:"TrackNumber"`
	PlayCount      int    `json:"PlayCount"`
	FilePath       string `json:"FilePath"`
	CreatedAt      string `json:"CreatedAt"`
	LastModifiedAt string `json:"LastModifiedAt"`
}

type Songs []Song

func getSongs() (Songs, error) {
	songs := Songs{}
	url := fmt.Sprintf("%s/songs", hostApi)
	resp, err := http.Get(url)
	if err != nil {
		return songs, err
	}
	err = json.NewDecoder(resp.Body).Decode(&songs)
	if err != nil {
		return songs, err
	}
	return songs, nil
}

func (g *Gui) ReturnSongs() *fyne.Container {
	songs, err := getSongs()
	if err != nil {
		g.Logger.WithFields(logrus.Fields{LogFields.ErrorMessage: err}).Error("GetArtistsFailure")
		box := container.New(layout.NewCenterLayout(), widget.NewButton("Retry?", func() {
			g.Logger.Info("Retrying")
			// TODO - Add retry
		}))
		return box
	}

	count := 1
	listWidget := widget.NewList(
		func() int {
			return len(songs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("name\tartist\talbum\ttrackNumber\tplayCount")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("%d \t %s \t %s \t %s \t %d \t %d", songs[i].ID, songs[i].Name, songs[i].ArtistName, songs[i].AlbumName, songs[i].TrackNumber, songs[i].PlayCount))
			count ++
		})
	listWidget.OnSelected = func(id widget.ListItemID) {
		g.SelectedSongId = id + 1
		g.Logger.Info(fmt.Sprintf("SelectedSong - %d", g.SelectedSongId))
	}

	content := container.NewMax(listWidget)
	return content
}

func (g *Gui) SetUpSpeaker(){
	g.Logger.Info("SettingUpSpeakerStart")
	_, format, err := g.getSong(g.SelectedSongId)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{LogFields.ErrorMessage: err}).Error("SettingUpSpeakerFailure")
		panic(err)
	}
	g.SampleRate = format.SampleRate
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	g.Logger.Info("SettingUpSpeakerFinish")
}

func (g *Gui) getSong(songId int) (beep.StreamSeekCloser, beep.Format, error) {
	url := fmt.Sprintf("%s/songs/%d", hostApi, songId)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, beep.Format{}, err
	}
	request.Header.Set("Content-Type", "multipart/form-data;")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, beep.Format{}, err
	}
	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		return nil, beep.Format{}, err
	}
	return streamer, format, err
}

func (g *Gui) playSong() {
	currentSongId := <- playing
	if currentSongId == g.SelectedSongId {
		playing <- currentSongId
		g.Logger.Info(fmt.Sprintf("CurrentSongAlreadyPlaying - %d", currentSongId))
	} else {
		go func() {
			g.Logger.Info(fmt.Sprintf("Playing Song - %d", g.SelectedSongId))
			streamer, format, err := g.getSong(g.SelectedSongId)
			if err != nil {
				g.Logger.WithFields(logrus.Fields{LogFields.ErrorMessage: err}).Error("PlaySongFailure")
				return
			}
			defer streamer.Close()

			playing <- g.SelectedSongId

			resampled := beep.Resample(4, g.SampleRate, format.SampleRate, streamer)
			//Clear any existing songs that might be going on
			speaker.Clear()
			//Create this to make the app wait for this song to finish before the call back fires
			done := make(chan bool)
			speaker.Play(beep.Seq(resampled, beep.Callback(func() {
				// TODO - wire up channels to indicate playing or not to not have duplicate playing contexts over lapping
				//	TODO - set up a channel to handle queue waiting & and reverse for recently played
				g.NextSongId++
				g.SelectedSongId++
				done <- true
				go g.playSong()
			})))
			<-done
		}()
	}
}

func (g *Gui) clickStop() {
	g.Logger.Info("StoppingCurrentSong")
	<- playing
	playing <- 0
	speaker.Clear()
}