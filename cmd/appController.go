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

//TODO: GENERAL
// 	 - To change song groupings - playlists, orders, albums, artists, etc.  We should call the backend and update the
//  	song list to match.
//   - Shuffle

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
	return
}

func (a *App) populateNextNavQueue() {
	//All this is to pre-populate the navQueueNext with the nextBtn batch of songs to playBtn, either based on the current
	//index if we don't have a previousBtn list, or append to the previousBtn list with the songs that follow it in the
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

func (a *App) clickNext() {
	if len(a.navQueueNext) == 0 {
		log.Println("Nothing in nextBtn (wtf are you doing?) - skipping")
		return
	}
	// I don't have to pop this off the queue here, because UpdateNavQueues handles this queue's values automatically.
	a.selectedSong = a.navQueueNext[0]
	a.playSong()
}

func (a *App) clickPrevious() {
	//If we don't have any previousBtn songs, just return
	if len(a.navQueuePrevious) == 0 {
		log.Println("Nothing previously played - skipping")
		return
	}
	// Get the last song off the back of this queue and strip it off since that should be what's currently paused (we
	//populate the queues before actually executing the song in case you cut it off halfway through), so that we can
	//click previousBtn more than once and actually move through the queue.
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

func (a *App) clickStop() {
	log.Println("StoppingCurrentSong")
	//a.stop <-true
	<-playing
	playing <- 0
	speaker.Clear()
}

// FIXME - OPTIONAL - we could potentially do something like https://github.com/pfortin-urbn/stalk/blob/master/collectors/aws/aws_collector.go#L106-L125 \
//  here but would require some revamping.
func (a *App) clickPlay() {
	if a.speakerControl == nil {
		log.Println("Playing me")
		a.playSong()
	} else if a.speakerControl.Paused {
		log.Println("Un-Pausing me")
		speaker.Lock()
		a.speakerControl.Paused = false
		speaker.Unlock()
	} else {
		log.Println("Pausing me")
		speaker.Lock()
		a.speakerControl.Paused = true
		speaker.Unlock()
	}
}

func (a *App) clickSong(song *Song) {
	a.selectedSong = song
	// If we actually clicked on a new song, we can wipe the up nextBtn list because it's not relevant any more
	//considering our new index/position.
	a.navQueueNext = []*Song{}
	a.playSong()
}

func (a *App) playSong() {
	if a.selectedSong == nil {
		a.selectedSong = a.songList[0]
	}
	currentSongId := <-playing
	songId := a.selectedSong.Id
	songName := a.selectedSong.Title

	if currentSongId == songId {
		playing <- currentSongId
		log.Printf("CurrentSongAlreadyPlaying - Index: %d, Id: %d, Title: %s", a.selectedSong.songListIndex, songId, songName)
		return
	}

	go func() {
		log.Printf("PlayingSong - Index: %d, Id: %d, Title: %s", a.selectedSong.songListIndex, songId, songName)
		streamer, format, err := a.getSong(songId)
		//streamer, _, err := a.getSong(songId)
		if err != nil {
			log.Printf("ClickSongFailure - %+v", err)
			return
		}
		defer streamer.Close()

		playing <- songId
		go a.UpdateNavQueues()

		a.speakerControl = &beep.Ctrl{Streamer: streamer, Paused: false}
		resampled := beep.Resample(4, a.SampleRate, format.SampleRate, a.speakerControl)
		//Clear any existing songs that might be going on
		speaker.Clear()
		//Create this to make the app wait for this song to finish before the callback fires
		done := make(chan bool)
		speaker.Play(beep.Seq(resampled, beep.Callback(func() {
			a.selectedSong = a.navQueueNext[0]
			done <- true
			go a.playSong()
		})))
		<-done
	}()
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
