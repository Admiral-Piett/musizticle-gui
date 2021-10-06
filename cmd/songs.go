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
	"time"
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
	fmt.Println(body)
	log.Println("clickSongComplete")
}

func (a *App) ShowSongs(gtx layout.Context) layout.Dimensions {

	if !a.songs.populated {
		if !a.songs.inProgress {
			return material.Button(th, &a.songs.reload, "Retry").Layout(gtx)
		}
		return material.Button(th, &a.songs.loadingButton, "Loading... Click me if I take too long").Layout(gtx)
	}

	listDimensions := a.displayList.Layout(gtx, len(a.songs.songList), func(gtx layout.Context, index int) layout.Dimensions {
		song := &a.songs.songList[index]
		if song.line.Clicked() {
			a.songs.clickSong(song.Id)
		}
		dims := material.Clickable(gtx, &song.line, func(gtx layout.Context) layout.Dimensions {
			return material.Label(th, unit.Dp(float32(20)), song.Name).Layout(gtx)
		})
		return dims
	})
	return listDimensions
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
	time.Sleep(5 * time.Second)
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
