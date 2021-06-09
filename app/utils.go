package app

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
