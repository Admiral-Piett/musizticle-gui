package main

import (
	"encoding/json"
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io"
	"log"
	"net/http"
	"strconv"
)


func (a *App) getSong(songId int) (beep.StreamSeekCloser, beep.Format, error) {
	url := fmt.Sprintf("%s/songs/%d", HOST, songId)
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


func (a *App) clickSong() {
	currentSongId := <- playing
	if currentSongId == a.SelectedSongId {
		playing <- currentSongId
		log.Printf("CurrentSongAlreadyPlaying - %d", currentSongId)
	} else {
		go func() {
			log.Printf("PlayingSong - %d", a.SelectedSongId)
			streamer, format, err := a.getSong(a.SelectedSongId)
			if err != nil {
				log.Printf("ClickSongFailure - %+v", err)
				return
			}
			defer streamer.Close()

			playing <- a.SelectedSongId

			resampled := beep.Resample(4, a.SampleRate, format.SampleRate, streamer)
			//Clear any existing songs that might be going on
			speaker.Clear()
			//Create this to make the app wait for this song to finish before the callback fires
			done := make(chan bool)
			speaker.Play(beep.Seq(resampled, beep.Callback(func() {
				//	TODO - set up a channel to handle queue waiting & and reverse for recently played
				a.SelectedSongId++
				a.NextSongId = a.SelectedSongId + 1
				done <- true
				go a.clickSong()
			})))
			<-done
		}()
	}
}

func buildSongLine(gtx layout.Context, s *Song) layout.Dimensions {
	lineDimenstions := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return songFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), strconv.Itoa(s.Id)).Layout(gtx))
		}),
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return songFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), s.Name).Layout(gtx))
			},
		),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return songFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), s.ArtistName).Layout(gtx))
		}),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return songFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), s.AlbumName).Layout(gtx))
		}),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return songFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), strconv.Itoa(s.TrackNumber)).Layout(gtx))
		}),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return songFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), strconv.Itoa(s.PlayCount)).Layout(gtx))
		}),
	)
	return lineDimenstions
}

func (a *App) SongsList(gtx layout.Context) layout.Dimensions {
	if !a.songs.populated {
		if !a.songs.inProgress {
			return material.Button(th, &a.songs.reload, "Retry").Layout(gtx)
		}
		return material.Button(th, &a.songs.loadingButton, "Loading... Click me if I take too long").Layout(gtx)
	}

	listDimensions := a.displayList.Layout(gtx, len(a.songs.songList), func(gtx layout.Context, index int) layout.Dimensions {
		song := &a.songs.songList[index]
		if song.line.Clicked() {
			a.SelectedSongId = song.Id
			a.clickSong()
		}

		// Song line clickable, wrapping the line meta labels, etc., wrapped in margins and borders.
		line := material.Clickable(gtx, &song.line, func(gtx layout.Context) layout.Dimensions {
			// Build the rest of the song line info in here.
			return buildSongLine(gtx, song)
		})
		return songLineWrapMarginsBorder(gtx, line)
	})
	return listDimensions
}

//TODO - NEXT - figure out something more elegant than this for the borders and margins
// I'm going to need to wrap all kinds of stuff in borders/margins, and I'll need something more elegant
// FIXME - think about reworking the grid at at some point, they're evenly spaced but also look a little wonky.
func SongsHeader(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return headerFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), "ID").Layout(gtx))
		}),
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return headerFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), "Title").Layout(gtx))
			},
		),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return headerFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), "Artist").Layout(gtx))
		}),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return headerFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), "Album").Layout(gtx))
		}),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return headerFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), "Track Number").Layout(gtx))
		}),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return headerFieldsMarginsBorder(gtx, material.Label(th, unit.Dp(float32(20)), "Play Count").Layout(gtx))
		}),
	)
}

func (s *Songs) initSongs() {
	if s.inProgress || s.populated {
		return
	}
	log.Println("getSongsStart")
	s.inProgress = true
	defer func() {
		s.inProgress = false
		displayChange <- true
	}()
	url := fmt.Sprintf("%s/songs", HOST)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to get songs")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	json.Unmarshal(body, &s.songList)
	s.populated = true
	log.Println("getSongsComplete")
}
