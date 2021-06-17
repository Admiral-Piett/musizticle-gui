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

	//var table [][]string
	//row := []string{"Track", "Artist", "Album", "Track Number", "Play Count"}
	//table = append(table, row)
	//for _, song := range(songs) {
	//	row := []string{song.Name, song.ArtistName, song.AlbumName, fmt.Sprintf("%d", song.TrackNumber), fmt.Sprintf("%d", song.PlayCount)}
	//	table = append(table, row)
	//}
	//
	//tableWidget := widget.NewTable(
	//	func() (int, int) {
	//		return len(table), len(table[0])
	//	},
	//	func() fyne.CanvasObject {
	//		return widget.NewLabel("wide content")
	//	},
	//	func(i widget.TableCellID, o fyne.CanvasObject) {
	//		o.(*widget.Label).SetText(table[i.Row][i.Col])
	//	})
	//
	//tableWidget.OnSelected = func(id widget.TableCellID) {
	//	g.Logger.Info(fmt.Sprintf("Clicked - %v+", id))
	//}


	//grid := container.New(layout.NewGridLayout(6))
	//for _, song := range(songs){
	//	grid.Add(widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
	//		g.playSong()
	//	}))
	//	grid.Add(canvas.NewText(song.Name, color.White))
	//	grid.Add(canvas.NewText(song.ArtistName, color.White))
	//	grid.Add(canvas.NewText(song.AlbumName, color.White))
	//	grid.Add(canvas.NewText(fmt.Sprintf("%d", song.TrackNumber), color.White))
	//	grid.Add(canvas.NewText(fmt.Sprintf("%d", song.PlayCount), color.White))
	//}
	//return grid`

	listWidget := widget.NewList(
		func() int {
			return len(songs)
		},
		func() fyne.CanvasObject {
			//name := canvas.NewText("Name", color.White)
			//artist := canvas.NewText("Artist", color.White)
			//album := canvas.NewText("Album", color.White)
			//trackNumber := canvas.NewText("Track Number", color.White)
			//playCount := canvas.NewText("Play Count", color.White)
			//return container.New(layout.NewHBoxLayout(), name, artist, album, trackNumber, playCount)
			//return container.NewHBox(name, artist, album, trackNumber, playCount)
			//return widget.NewTextGridFromString("name\tartist\talbum\ttrackNumber\tplayCount")
			return widget.NewLabel("name\tartist\talbum\ttrackNumber\tplayCount")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			//name := canvas.NewText(songs[i].Name, color.White)
			//artist := canvas.NewText(songs[i].ArtistName, color.White)
			//album := canvas.NewText(songs[i].AlbumName, color.White)
			//trackNumber := canvas.NewText(fmt.Sprintf("%d", songs[i].TrackNumber), color.White)
			//playCount := canvas.NewText(fmt.Sprintf("%d", songs[i].PlayCount), color.White)
			//o.(*fyne.CanvasObject).SetCo layout.NewHBoxLayout(), name, artist, album, trackNumber, playCount)
			o.(*widget.Label).SetText(fmt.Sprintf("%s \t %s \t %s \t %d \t %d", songs[i].Name, songs[i].ArtistName, songs[i].AlbumName, songs[i].TrackNumber, songs[i].PlayCount))
		})

	listWidget.OnSelected = func(id widget.ListItemID) {
		g.Logger.Info(fmt.Sprintf("Clicked - %s", id))
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
	g.Logger.Info(fmt.Sprintf("Playing Song - %d", g.SelectedSongId))
	done := make(chan bool)
	go func() {
		streamer, format, err := g.getSong(g.SelectedSongId)
		if err != nil {
			g.Logger.WithFields(logrus.Fields{LogFields.ErrorMessage: err}).Error("PlaySongFailure")
			return
		}
		defer streamer.Close()
		resampled := beep.Resample(4, g.SampleRate, format.SampleRate, streamer)
		speaker.Play(beep.Seq(resampled, beep.Callback(func() {
			// TODO - wire up channels to indicate playing or not to not have duplicate playing contexts over lapping
			//	TODO - set up a channel to handle queue waiting & and reverse for recently played
			g.NextSongId ++
			g.SelectedSongId ++
			done <- true
			go g.playSong()
		})))
		<-done
	}()
}

func (g *Gui) stopSong() {
	g.Logger.Info("Playing Song - %d", g.SelectedSongId)
	speaker.Clear()
}