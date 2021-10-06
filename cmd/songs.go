package main

import (
	"encoding/json"
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"io"
	"log"
	"net/http"
)

// TODO - environmentalize
var HOST = "http://localhost:9000/api"


func (s *Songs) clickSong(songId int) {
	s.selected = songId
	if s.inProgress {
		return
	}
	log.Println("clickSongStart")
	s.inProgress = true
	defer func() {
		s.inProgress = false
		displayChange <- true
	}()
	url := fmt.Sprintf("%s/songs/%d", HOST, songId)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("clickSongFailure")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// TODO - NEXT - start here, add streaming
	//  Consider adding an App Struct and carrying app status here (selected song, currently playing a song, beep SampleRate etc.)
	fmt.Println(body)
	log.Println("clickSongComplete")
}

func (s *Songs) ShowSongs(gtx layout.Context) layout.Dimensions {
	go s.getSongs()

	if !s.populated {
		text := "Loading... Click me if I take too long"
		//FIXME - this should be able to swap, but the concurrency isn't right
		//if !s.inProgress {
		//	text = "Retry"
		//}
		// If we're not populated yet, and we're no longer in progress that means we failed to get the songs so expose
		//a retry button
		return material.Button(th, &s.loadingButton, text).Layout(gtx)
	}

	listDimensions := displayList.Layout(gtx, len(s.songList), func(gtx layout.Context, index int) layout.Dimensions {
		song := &s.songList[index]
		if song.line.Clicked() {
			s.clickSong(song.Id)
		}
		dims := material.Clickable(gtx, &song.line, func(gtx layout.Context) layout.Dimensions {
			return material.Label(th, unit.Dp(float32(20)), song.Name).Layout(gtx)
		})
		return dims
	})
	return listDimensions
}

func (s *Songs) getSongs() {
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
