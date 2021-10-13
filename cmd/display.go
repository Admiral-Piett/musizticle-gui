package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"strconv"
)

//FIXME - (with FIXME on UpdateNavQueues()) Instead of managing these all independently could I merge into one
// tabDisplay method that populates a song list?
// 10/9/2021- meh...I'm not sure I want to do this.  I hate that this code isn't shared but the extra complication would
//   be a headache.
func (a *App) tabDisplay(songList []*Song) []layout.FlexChild {
	displayArray := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.TabsBar(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.SongsHeader(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.SongsList(gtx, songList)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.MediaToolBar(gtx)
		}),
	}
	return displayArray
}

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
				return material.Button(th, &a.previousBtn, "Previous").Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// TODO - make these icons
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				text := "Play"
				if a.speakerControl != nil && !a.speakerControl.Paused {
					text = "Pause"
				}
				return material.Button(th, &a.playBtn, text).Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// TODO - make these icons
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &a.stopBtn, "Stop").Layout(gtx)
			})
		}),
		layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
			// TODO - make these icons
			return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &a.nextBtn, "Next").Layout(gtx)
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
