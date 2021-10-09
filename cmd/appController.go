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


func (a *App) TabsBar(gtx layout.Context) layout.Dimensions {
	margins := layout.Inset{
		Top:    unit.Dp(5),
		Right:  unit.Dp(2),
		Bottom: unit.Dp(5),
		Left:   unit.Dp(2),
	}

	weightMap := make(map[string]float32)
	for _, tab := range TAB_LIST {
		weightMap[tab] = 0.2
		if a.selectedTab == tab {
			weightMap[tab] = 1
		}
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(
			layout.Spacer{Width: unit.Dp(10)}.Layout,
		),
		layout.Flexed(weightMap[HOME_TAB], func(gtx layout.Context) layout.Dimensions {
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &a.homeTab, "Home").Layout(gtx)
			})
		}),
		layout.Flexed(weightMap[NEXT_TAB], func(gtx layout.Context) layout.Dimensions {
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &a.nextTab, "Up Next").Layout(gtx)
			})
		}),
		layout.Flexed(weightMap[PREVIOUS_TAB], func(gtx layout.Context) layout.Dimensions {
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &a.previousTab, "Recently Played").Layout(gtx)
			})
		}),
		layout.Rigid(
			layout.Spacer{Width: unit.Dp(10)}.Layout,
		),
	)
}

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

func (a *App) reloadNavQueues() {

}

//FIXME - HERE - Instead of managing these all manually, could I maintain a list of indexes on the entire songsList, and
// index pieces of it to populate these?
func (a *App) UpdateNavQueues() {
	//Every time we come through here, if we don't have room for another song, we slice the array at a placement of
	//one less than the limit from the front of the array.  Meaning if the limit is 20, we get the last 19 to be added.
	if len(a.navQueuePrevious) >= (NAV_QUEUE_PREVIOUS_LIMIT) {
		index := len(a.navQueuePrevious) - (NAV_QUEUE_PREVIOUS_LIMIT - 1)
		// Take all the things in the back of the list from the limit so that we leave room to add the new one.
		a.navQueuePrevious = a.navQueuePrevious[index:]
	}
	if len(a.navQueueNext) >= (NAV_QUEUE_NEXT_LIMIT) {
		index := len(a.navQueueNext) - (NAV_QUEUE_NEXT_LIMIT - 1)
		a.navQueueNext = a.navQueueNext[index:]
	}

	a.navQueuePrevious = append(a.navQueuePrevious, a.SelectedSongIndex)

	a.populateNavQueues()
	log.Printf("navQueueNext - %+v", a.navQueueNext)
	log.Printf("navQueuePrevious - %+v", a.navQueuePrevious)
	return
}

func (a *App) populateNavQueues() {
	// This is to pre-populate the navQueuePrevious index, if we're starting fresh or it gets wiped on a new song list.
	// TODO - make sure to wipe this on a new songList population?
	currentIndex := a.SelectedSongIndex
	if len(a.navQueuePrevious) != 0 {
		index := len(a.navQueuePrevious) - 1
		currentIndex = a.navQueuePrevious[index]
	}
	for len(a.navQueuePrevious) < NAV_QUEUE_NEXT_LIMIT {
		// If we're about to index passed the end of the songList, then reset to start at the back of the song list
		// to start over at the last index (we subtract 1 before adding which should make up for the 0 index diff)
		if (currentIndex - 1) < 0 {
			currentIndex = len(a.songs.songList)
		}
		previousIndex := currentIndex - 1
		a.navQueuePrevious = append([]int{previousIndex}, a.navQueuePrevious...)
		currentIndex--
	}
	//All this is to pre-populate the navQueueNext with the next batch of songs to play, either based on the current
	//index if we don't have a previous list, or append to the previous list with the songs that follow it in the
	//current song list.
	currentIndex = a.SelectedSongIndex
	if len(a.navQueueNext) != 0 {
		index := len(a.navQueueNext) - 1
		currentIndex = a.navQueueNext[index]
	}
	for len(a.navQueueNext) < NAV_QUEUE_NEXT_LIMIT {
		// If we're about to index passed the end of the songList, then reset to 0 and start over
		if (currentIndex + 1) > len(a.songs.songList) {
			currentIndex = 0
		}
		nextIndex := currentIndex + 1
		a.navQueueNext = append(a.navQueueNext, nextIndex)
		currentIndex++
	}
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
				a.SelectedSongIndex = a.navQueueNext[0]
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

func (a *App) clickNext() {
	if len(a.navQueueNext) == 0 {
		log.Println("Nothing in next (wtf are you doing?) - skipping")
		return
	}
	// QUESTION - should I have navQueuePrevious preform like this too?
	// I don't have to pop this off the queue here, because UpdateNavQueues handles this queue's values automatically.
	a.SelectedSongIndex = a.navQueueNext[0]
	a.playSong()
}

func (a *App) clickPrevious() {
	// FIXME - should we just start from the back of the current songList?  If so, how would we populate the previous
	//  list so that it preforms like the navQueueNext?
	//If we don't have any previous songs, just return
	if len(a.navQueuePrevious) == 0 {
		log.Println("Nothing previously played - skipping")
		return
	}
	// Get the last song off the back of this queue and strip it off since that should be what's currently playing (we
	//populate the queues before actually executing the song in case you cut it off halfway through), so that we can
	//click previous more than once and actually move through the queue.
	currentSongIndex := len(a.navQueuePrevious) - 1
	a.navQueuePrevious = a.navQueuePrevious[:currentSongIndex]
	// FIXME - this is too tedious for every click - upgrade with a little more book keeping for performance down
	//  the road.
	// Reset this here because this will get populated completely from where we are currently as the song plays
	a.navQueueNext = []int{}

	previousSongIndex := len(a.navQueuePrevious) - 1
	a.SelectedSongIndex = a.navQueuePrevious[previousSongIndex]
	a.navQueuePrevious = a.navQueuePrevious[:previousSongIndex]
	a.playSong()
}

func (a *App) clickSong(index int) {
	a.SelectedSongIndex = index
	// If we actually clicked on a new song, we can wipe the up next list because it's not relevant any more
	//considering our new index/position.
	a.navQueueNext = []int{}
	a.playSong()
}



func (a *App) SongsList(gtx layout.Context, songsList []Song) layout.Dimensions {
	if !a.songs.populated {
		if !a.songs.initSongsInProgress {
			return material.Button(th, &a.songs.reload, "Retry").Layout(gtx)
		}
		return material.Button(th, &a.songs.loadingButton, "Loading... Click me if I take too long").Layout(gtx)
	}

	listDimensions := a.displayList.Layout(gtx, len(songsList), func(gtx layout.Context, index int) layout.Dimensions {
		song := &songsList[index]
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

//TODO - figure out something more elegant than this for the borders and margins
// I'm going to need to wrap all kinds of stuff in borders/margins, and I'll need something more elegant
// FIXME - think about reworking the grid at at some point, they're evenly spaced but also look a little wonky.
func (a *App) SongsHeader(gtx layout.Context) layout.Dimensions {
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
	go func() {
		for index, _ := range s.songList {
			s.songIndexes = append(s.songIndexes, index)
		}
	}()
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
