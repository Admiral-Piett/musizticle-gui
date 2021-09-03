package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"math"
)

type LogFieldStruct = struct {
	AlbumId string
	ArtistId string
	ErrorMessage string
	FilePath     string
	SongID       string
	RequestBody  string
	Size string
	StackContext string
}

var LogFields = LogFieldStruct{
	AlbumId: "album_id",
	ArtistId: "artist_id",
	ErrorMessage: "error_message",
	FilePath:     "file_path",
	SongID:       "song_id",
	RequestBody:  "request_body",
	Size: "size",
	StackContext: "stack_context",
}

var hostApi = "http://localhost:9000/api"

type SongGrid struct {
	columns  int
	vertical bool
}

func NewSongGrid(columns int, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(&SongGrid{columns, false}, objects...)
}

func (s *SongGrid) MinSize(objects []fyne.CanvasObject) fyne.Size{
	w, h := minWidth, minHeight
	for _, o := range objects {
		objSize := o.MinSize()
		if objSize.Width > w {
			w = objSize.Width
		}
		if objSize.Height > h {
			h = objSize.Height
		}
	}
	return fyne.NewSize(w, h)
}

func (s *SongGrid) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	padWidth := theme.Padding()
	padHeight := theme.Padding()

	//TODO - for some reason this isn't working, maybe swap to a list version?  I don't know
	row, col := 0, 0
	i := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		size := child.MinSize()
		cellWidth := float64(size.Width+padWidth)
		//All the cells will be the same height for now
		cellHeight := float64(size.Height+padHeight)

		x1 := getLeading(cellWidth, col)
		y1 := getLeading(cellHeight, row)
		x2 := getTrailing(cellWidth, col)
		y2 := getTrailing(cellHeight, row)

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize(x2-x1, y2-y1))

		if !s.vertical {
			if (i+1)%s.columns == 0 {
				row++
				col = 0
			} else {
				col++
			}
		} else {
			if (i+1)%s.columns == 0 {
				col++
				row = 0
			} else {
				row++
			}
		}
		i++
	}
}

// NOTE: copied from fyne.gridlayout for now
// Get the leading (top or left) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getLeading(size float64, offset int) float32 {
	ret := (size + float64(theme.Padding())) * float64(offset)

	return float32(math.Round(ret))
}

// Get the trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getTrailing(size float64, offset int) float32 {
	return getLeading(size, offset+1) - theme.Padding()
}
