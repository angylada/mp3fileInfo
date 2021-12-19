# mp3fileInfo
Enriches ID3-data (artist and title) based on the filename of all mp3-files  in either a specific directory if 
given via command line arguments. If the parameter is being omitted, the current working directory will be used.
The filenames need to be in the format "ARTIST - TITLE.mp3" since " - " is the separator.

## How to build

### Linux
`go build cmd/mp3fileInfo.go`

### Windows 
Built on linux for execution on Windows:

`GOOS=windows GOARCH=amd64 go build cmd/mp3fileInfo.go`