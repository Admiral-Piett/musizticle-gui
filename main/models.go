package main

import (
    "gioui.org/app"
    "gioui.org/font/gofont"
    "gioui.org/layout"
    "gioui.org/widget"
    "github.com/faiface/beep"
    "time"
)

// TODO - Explore other font collections
var th = CreateTheme(gofont.Collection())

// TODO - SECRETIFY - environmentalize
var HOST = "http://localhost:9000/api"

var currentSongId int
var displayChange = make(chan bool)

var progressIncrementer chan float32
var progress float32

var loginRequired bool
var loginUsername widget.Editor
var loginPassword widget.Editor
var loginButton widget.Clickable
var authToken string

// TODO - SECRETIFY - environmentalize
var authExpirationTime time.Time

var backgroundLoginInProgress bool
var backgroundLoginRetryCount int

// TODO - SECRETIFY - environmentalize
var NAV_QUEUE_PREVIOUS_LIMIT = 20
var NAV_QUEUE_NEXT_LIMIT = 20

var HOME_TAB = "home"
var NEXT_TAB = "nextBtn"
var PREVIOUS_TAB = "previousBtn"
var TAB_LIST = []string{HOME_TAB, NEXT_TAB, PREVIOUS_TAB}

type App struct {
    // --- Display
    displayList *layout.List
    songs       Songs
    window      *app.Window
    // Speaker Execution
    SampleRate     beep.SampleRate
    speakerControl *beep.Ctrl
    // --- Song Execution
    playing          bool
    songList         []*Song
    selectedSong     *Song
    navQueueNext     []*Song
    navQueuePrevious []*Song
    playBtn          widget.Clickable
    stopBtn          widget.Clickable
    nextBtn          widget.Clickable
    previousBtn      widget.Clickable
    // --- Tabs
    selectedTab string
    homeTab     widget.Clickable
    nextTab     widget.Clickable
    previousTab widget.Clickable
}

type Songs struct {
    populated           bool
    initSongsInProgress bool
    reload              widget.Clickable
    loadingButton       widget.Clickable
}

type Song struct {
    line           widget.Clickable
    songListIndex  int
    Id             int
    Title          string
    ArtistId       int
    ArtistName     string
    AlbumId        int
    AlbumName      string
    TrackNumber    int
    PlayCount      int
    FilePath       string
    Duration       int
    CreatedAt      string
    LastModifiedAt string
}

type AuthRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type AuthResponse struct {
    AuthToken      string `json:"authToken"`
    ReauthToken    string `json:"reauthToken"`
    ExpirationTime string `json:"expirationTime"`
}
