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

//TODO:
// 	 To change song groupings - playlists, orders, albums, artists, etc.  We should call the backend and update the
//  	song list to match.
//	 So, the source of truth should just be the index, so we can look up the song in the list and get what we
//		want off it, including easily know where it is to fetch more out of the list.

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

func (a *App) UpdateNavQueues() {
	if len(navQueuePrevious) >= (NAV_QUEUE_PREVIOUS_LIMIT) {
		index := len(navQueuePrevious) - (NAV_QUEUE_PREVIOUS_LIMIT - 1)
		// Take all the things in the back of the list from the limit so that we leave room to add the new one.
		navQueuePrevious = navQueuePrevious[index:]
	}
	if len(navQueueNext) >= (NAV_QUEUE_NEXT_LIMIT) {
		index := len(navQueueNext) - (NAV_QUEUE_NEXT_LIMIT - 1)
		navQueueNext = navQueueNext[index:]
	}

	navQueuePrevious = append(navQueuePrevious, a.SelectedSongIndex)

	//All this is to pre-populate the navQueueNext with the next batch of songs to play, either based on the current
	//index if we don't have a previous list, or append to the previous list with the songs that follow it in the
	//current song list.
	currentIndex := a.SelectedSongIndex
	if len(navQueueNext) != 0 {
		index := len(navQueueNext) - 1
		currentIndex = navQueueNext[index]
	}
	for len(navQueueNext) < NAV_QUEUE_NEXT_LIMIT {
		// If we're about to index passed the end of the songList, then reset to 0 and start over
		if (currentIndex + 1) > len(a.songs.songList) {
			currentIndex = 0
		}
		nextIndex := currentIndex + 1
		navQueueNext = append(navQueueNext, nextIndex)
		currentIndex++
	}
	return
}

func (a *App) playSong() {
	currentSongId := <-playing
	songId := a.songs.songList[a.SelectedSongIndex].Id
	songName := a.songs.songList[a.SelectedSongIndex].Name

	if currentSongId == a.SelectedSongIndex {
		playing <- currentSongId
		log.Printf("CurrentSongAlreadyPlaying - Index: %d, Id: %d, Name: %s", a.SelectedSongIndex, songId, songName)
	} else {
		playing <- a.SelectedSongIndex
		go func() {
			log.Printf("PlayingSong - Index: %d, Id: %d, Name: %s", a.SelectedSongIndex, songId, songName)
			streamer, format, err := a.getSong(songId)
			if err != nil {
				log.Printf("ClickSongFailure - %+v", err)
				return
			}
			defer streamer.Close()

			go a.UpdateNavQueues()

			resampled := beep.Resample(4, a.SampleRate, format.SampleRate, streamer)
			//Clear any existing songs that might be going on
			speaker.Clear()
			//Create this to make the app wait for this song to finish before the callback fires
			done := make(chan bool)
			speaker.Play(beep.Seq(resampled, beep.Callback(func() {
				// TODO - Set up shuffle, and better song list navigation (can't just increment the song id)
				//	TODO - set up a channel to handle queue waiting & and reverse for recently played
				a.SelectedSongIndex = navQueueNext[0]
				done <- true
				go a.playSong()
			})))
			<-done
		}()
	}
}

func (a *App) clickStop() {
	log.Println("StoppingCurrentSong")
	<-playing
	playing <- 0
	speaker.Clear()
}

//TODO - NEXT - make these work with the navQueueNext & navQueuePrevious
//TODO - NEXT - also make tabs for viewing "Up Next" & "Recently Played"
func (a *App) clickNext() {

}

func (a *App) clickPrevious() {

}

func (a *App) clickSong(index int) {
	a.SelectedSongIndex = index
	// If we actually clicked on a new song, we can wipe the up next list because it's not relevant any more
	//considering our new index/position.
	navQueueNext = []int{}
	a.playSong()
}



func (a *App) SongsList(gtx layout.Context) layout.Dimensions {
	if !a.songs.populated {
		if !a.songs.initSongsInProgress {
			return material.Button(th, &a.songs.reload, "Retry").Layout(gtx)
		}
		return material.Button(th, &a.songs.loadingButton, "Loading... Click me if I take too long").Layout(gtx)
	}

	listDimensions := a.displayList.Layout(gtx, len(a.songs.songList), func(gtx layout.Context, index int) layout.Dimensions {
		song := &a.songs.songList[index]
		if song.line.Clicked() {
			a.clickSong(index)
		}

		// Song line clickable, wrapping the line meta labels, etc., wrapped in margins and borders.
		line := material.Clickable(gtx, &song.line, func(gtx layout.Context) layout.Dimensions {
			// Build the rest of the song line info in here.
			return buildSongLine(gtx, song)
		})
		return songLineMargins(gtx, line)
	})
	return listDimensions
}

func (a *App) MediaToolBar(gtx layout.Context) layout.Dimensions {
	margins := layout.Inset{
		Top:    unit.Dp(5),
		Right:  unit.Dp(5),
		Bottom: unit.Dp(10),
		Left:   unit.Dp(5),
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			// TODO - make these icons
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &a.previous, "Previous").Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// TODO - make these icons
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &a.play, "Play").Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// TODO - make these icons
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &a.stop, "Stop").Layout(gtx)
			})
		}),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			// TODO - make these icons
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &a.next, "Next").Layout(gtx)
			})
		}),
	)
}

func (s *Songs) initSongs() {
	if s.initSongsInProgress || s.populated {
		return
	}
	log.Println("getSongsStart")
	s.initSongsInProgress = true
	defer func() {
		s.initSongsInProgress = false
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

func buildSongLine(gtx layout.Context, s *Song) layout.Dimensions {
	lineDimenstions := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(
			layout.Spacer{Width: unit.Dp(10)}.Layout,
		),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return songFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), strconv.Itoa(s.Id)).Layout(gtx))
		}),
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return songFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), s.Name).Layout(gtx))
			},
		),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return songFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), s.ArtistName).Layout(gtx))
		}),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return songFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), s.AlbumName).Layout(gtx))
		}),
		layout.Flexed(0.15, func(gtx layout.Context) layout.Dimensions {
			return songFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), strconv.Itoa(s.TrackNumber)).Layout(gtx))
		}),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return songFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), strconv.Itoa(s.PlayCount)).Layout(gtx))
		}),
		layout.Rigid(
			// The height of the spacer is 25 Device independent pixels
			layout.Spacer{Width: unit.Dp(10)}.Layout,
		),
	)
	return lineDimenstions
}

//TODO - figure out something more elegant than this for the borders and margins
// I'm going to need to wrap all kinds of stuff in borders/margins, and I'll need something more elegant
// FIXME - think about reworking the grid at at some point, they're evenly spaced but also look a little wonky.
func SongsHeader(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(
			layout.Spacer{Width: unit.Dp(10)}.Layout,
		),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return headerFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), "ID").Layout(gtx))
		}),
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return headerFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), "Title").Layout(gtx))
			},
		),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return headerFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), "Artist").Layout(gtx))
		}),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			return headerFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), "Album").Layout(gtx))
		}),
		layout.Flexed(0.15, func(gtx layout.Context) layout.Dimensions {
			return headerFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), "Track Number").Layout(gtx))
		}),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return headerFieldsMargins(gtx, material.Label(th, unit.Dp(float32(20)), "Play Count").Layout(gtx))
		}),
		layout.Rigid(
			layout.Spacer{Width: unit.Dp(10)}.Layout,
		),
	)
}
