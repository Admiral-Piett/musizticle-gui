package main

import (
	"encoding/json"
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"net/http"
	"os"
	"time"
)

//TODO: GENERAL
// 	- Playlists
//	- Find songs by Artists/Albums (already done in the back end)
//	- Repeat (Single and list - list might just work?)
//	- Search - by song, artist, album, playlist, etc.
//  - Shuffle
//	- Rig up to spotify?

func (a *App) getSong(songId int) (beep.StreamSeekCloser, beep.Format, error) {
	url := fmt.Sprintf("%s/songs/%d", HOST, songId)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, beep.Format{}, err
	}
	request.Header.Set("Content-Type", "multipart/form-data;")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
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
	// Get the index of the next song in the list
	nextIndex := currentIndex + 1
	for len(a.navQueueNext) < NAV_QUEUE_NEXT_LIMIT {
		// If we're about to index passed the end of the songList, then reset to 0 and start over
		if nextIndex > (len(a.songList) - 1) {
			nextIndex = 0
		}
		a.navQueueNext = append(a.navQueueNext, a.songList[nextIndex])
		nextIndex++
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
	currentSongId = 0
	a.playing = false
	progress = 0
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
		a.playing = true
		speaker.Lock()
		a.speakerControl.Paused = false
		speaker.Unlock()
	} else {
		log.Println("Pausing me")
		a.playing = false
		speaker.Lock()
		a.speakerControl.Paused = true
		speaker.Unlock()
	}
}

func (a *App) clickSong(song *Song) {
	a.selectedSong = song
	progress = float32(song.Duration)
	// If we actually clicked on a new song, we can wipe the up nextBtn list because it's not relevant any more
	//considering our new index/position.
	a.navQueueNext = []*Song{}
	a.playSong()
}

//TODO - HERE - check this out
func (a *App) playSong() {
	if a.selectedSong == nil {
		a.selectedSong = a.songList[0]
	}
	songId := a.selectedSong.Id
	songName := a.selectedSong.Title

	if currentSongId == songId {
		log.Printf("CurrentSongAlreadyPlaying - Index: %d, Id: %d, Title: %s", a.selectedSong.songListIndex, songId, songName)
		return
	}

	resetProgress(a.selectedSong.Duration)
	go func() {
		log.Printf("PlayingSong - Index: %d, Id: %d, Title: %s", a.selectedSong.songListIndex, songId, songName)
		// TODO - pull from a cache
		streamer, format, err := a.getSong(songId)
		//streamer, _, err := a.getSong(songId)
		if err != nil {
			log.Printf("ClickSongFailure - %+v", err)
			return
		}
		defer streamer.Close()

		currentSongId = songId
		go a.UpdateNavQueues()

		a.playing = true
		progress = 0

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
	err := Get(url, &a.songList,true)
	if err != nil {
		log.Println("Failed to get songs")
		return
	}
	go func() {
		for index, song := range a.songList {
			song.songListIndex = index
		}
	}()
	a.songs.populated = true
	log.Println("getSongsComplete")
}

func (a *App) startUp() {
	loginRequired = true
	// Attempt to login with any stored config we have
	login()
	if loginRequired {
		return
	}

	a.initSongs()
	//Put an invalid song id on the currentSongId queue to start with
	currentSongId = -1
	a.SetUpSpeaker()
}

func login() {
	log.Println("loginStart")
	url := fmt.Sprintf("%s/auth", HOST)

	requestBody := AuthRequest{}
	// If we don't have anything in the input fields (should only happen at start up), then try to read from the file.
	//If we can't do that either just return and let them open the login menu.
	if loginUsername.Text() == "" || loginPassword.Text() == "" {
		authCreds, err := os.ReadFile("config")
		if err != nil {
			return
		}
		err = json.Unmarshal(authCreds, &requestBody)
		if err != nil || requestBody.Username == "" || requestBody.Password == "" {
			return
		}
	} else {
		requestBody.Username = loginUsername.Text()
		requestBody.Password = loginPassword.Text()
	}

	responseBody := &AuthResponse{}

	err := Post(url, requestBody, responseBody, false)
	if err != nil {
		log.Printf("LoginFailure: %s\n", err)
	}

	expiration, err := time.Parse(time.RFC3339,responseBody.ExpirationTime)
	if err != nil {
		log.Printf("ExpirationParsingFailure: %s\n", err)
	}
	authToken = responseBody.AuthToken
	// Subtract 5 minutes from the expiration time so we are always ahead of when it actually expires
	authExpirationTime = expiration.Add(-5*time.Minute)

	// Make sure these are set since we may have just got them from the file.
	if loginUsername.Text() == "" || loginPassword.Text() == "" {
		loginUsername.SetText(requestBody.Username)
		loginPassword.SetText(requestBody.Password)
	}
	loginRequired = false

	// Convert current credentials to bytes and save off to a file
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("SaveToFileFailure: %s\n", err)
	}
	_ = os.WriteFile("config", jsonBody, 0777)
	log.Println("loginComplete")

	// Reset this to 0, since we just got in successfully, and if we haven't already started the background
	//process, do so.
	backgroundLoginRetryCount = 0
	if !backgroundLoginInProgress {
		go backgroundLogin()
	}
}

// FIXME - there's probably a more elegant way of handling this, but for now, we're here.
func backgroundLogin() {
	backgroundLoginInProgress = true
	for {
		switch time.Now().After(authExpirationTime) {
		case true:
			// Wait if this is set.  The user is going to get prompted to fill in their password anyway,
			//	so we don't also have to try.
			if loginRequired {
				continue
			}
			log.Printf("BackgroundLoginRefreshStart - Attempt: %d\n", backgroundLoginRetryCount)
			backgroundLoginRetryCount += 1
			// This doesn't respond with an error, but it wil internally reset the `authExpirationTime` so if we come
			//	around again here we'll know.
			login()

			if backgroundLoginRetryCount >= 5 {
				log.Println("Error - retry limit exceeded")
				loginRequired = true
				return
			}
		}
		time.Sleep(10*time.Second)
	}
}
