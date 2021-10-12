package main

import (
	"encoding/json"
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io"
	"log"
	"net/http"
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

	a.navQueuePrevious = append(a.navQueuePrevious, a.selectedSong)

	a.populateNextNavQueue()
	//log.Printf("navQueueNext - %+v", a.navQueueNext)
	//log.Printf("navQueuePrevious - %+v", a.navQueuePrevious)
	return
}

func (a *App) populateNextNavQueue() {
	//All this is to pre-populate the navQueueNext with the next batch of songs to play, either based on the current
	//index if we don't have a previous list, or append to the previous list with the songs that follow it in the
	//current song list.
	currentIndex := a.selectedSong.songListIndex
	if len(a.navQueueNext) != 0 {
		currentIndex = len(a.navQueueNext) - 1
	}
	for len(a.navQueueNext) < NAV_QUEUE_NEXT_LIMIT {
		// If we're about to index passed the end of the songList, then reset to 0 and start over
		nextIndex := currentIndex + 1
		if nextIndex > len(a.songList) {
			currentIndex = 0
		}
		a.navQueueNext = append(a.navQueueNext, a.songList[nextIndex])
		currentIndex++
	}
}


func (a *App) playSong() {
	if a.selectedSong == nil {
		a.selectedSong = a.songList[0]
	}
	currentSongId := <-playing
	songId := a.selectedSong.Id
	songName := a.selectedSong.Name

	if currentSongId == songId {
		playing <- currentSongId
		log.Printf("CurrentSongAlreadyPlaying - Index: %d, Id: %d, Name: %s", a.selectedSong.songListIndex, songId, songName)
	} else {
		playing <- songId
		go func() {
			log.Printf("PlayingSong - Index: %d, Id: %d, Name: %s", a.selectedSong.songListIndex, songId, songName)
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
				a.selectedSong = a.navQueueNextSongs[0]
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
	a.selectedSong = a.navQueueNext[0]
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
	// Get the last song off the back of this queue and strip it off since that should be what's currently paused (we
	//populate the queues before actually executing the song in case you cut it off halfway through), so that we can
	//click previous more than once and actually move through the queue.
	currentSongIndex := len(a.navQueuePrevious) - 1
	a.navQueuePrevious = a.navQueuePrevious[:currentSongIndex]
	// FIXME - this is too tedious for every click - upgrade with a little more book keeping for performance down
	//  the road.
	// Reset this here because this will get populated completely from where we are currently as the song plays
	a.navQueueNext = []*Song{}

	previousSongIndex := len(a.navQueuePrevious) - 1
	a.selectedSong = a.navQueuePrevious[previousSongIndex]
	a.navQueuePrevious = a.navQueuePrevious[:previousSongIndex]
	a.playSong()
}

func (a *App) clickPlay() {
	if a.paused {
		log.Println("Unpause me")
		a.paused = false
		a.playSong()
	} else {
		a.paused = true
		speaker.Lock()
		a.speakerControl.Paused = a.paused
		speaker.Unlock()
	}
}

func (a *App) clickSong(song *Song) {
	a.selectedSong = song
	// If we actually clicked on a new song, we can wipe the up next list because it's not relevant any more
	//considering our new index/position.
	a.navQueueNext = []*Song{}
	a.playSong()
}

func (a *App) SongsList(gtx layout.Context, songsList []*Song) layout.Dimensions {
	if !a.songs.populated {
		if !a.songs.initSongsInProgress {
			return material.Button(th, &a.songs.reload, "Retry").Layout(gtx)
		}
		return material.Button(th, &a.songs.loadingButton, "Loading... Click me if I take too long").Layout(gtx)
	}

	listDimensions := a.displayList.Layout(gtx, len(songsList), func(gtx layout.Context, index int) layout.Dimensions {
		song := songsList[index]
		if song.line.Clicked() {
			a.clickSong(song)
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

func (a *App) initSongs() {
	if a.songs.initSongsInProgress || a.songs.populated {
		return
	}
	log.Println("getSongsStart")
	a.songs.initSongsInProgress = true
	defer func() {
		a.songs.initSongsInProgress = false
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
	json.Unmarshal(body, &a.songList)
  	// Make the songs aware of where they are in the "master" songList
	go func() {
		for index, song := range a.songList {
			 song.songListIndex = index
		}
	}()
	a.songs.populated = true
	log.Println("getSongsComplete")
}
