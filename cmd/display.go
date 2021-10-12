package main

import "gioui.org/layout"

//FIXME - (with FIXME on UpdateNavQueues()) Instead of managing these all independently could I merge into one
// tabDisplay method that populates a song list?
// 10/9/2021- meh...I'm not sure I want to do this.  I hate that this code isn't shared but the extra complication would
//   be a headache.
func (a *App) homeTabDisplay() []layout.FlexChild {
	displayArray := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.TabsBar(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.SongsHeader(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.SongsList(gtx, a.songList)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.MediaToolBar(gtx)
		}),
	}
	return displayArray
}

func (a *App) nextTabDisplay() []layout.FlexChild {
	// Populate this list with the songs currently in the navQueueNext
	songsList := []*Song{}
	for _, index := range a.navQueueNext {
		songsList = append(songsList, a.songList[index])
	}
	displayArray := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.TabsBar(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.SongsHeader(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.SongsList(gtx, songsList)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.MediaToolBar(gtx)
		}),
	}
	return displayArray
}

func (a *App) previousTabDisplay() []layout.FlexChild {
	// Populate this list with the songs currently in the navQueueNext
	songsList := []*Song{}
	for _, index := range a.navQueuePrevious {
		songsList = append(songsList, a.songList[index])
	}
	displayArray := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.TabsBar(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.SongsHeader(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.SongsList(gtx, songsList)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.MediaToolBar(gtx)
		}),
	}
	return displayArray
}
