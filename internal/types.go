package internal

type Config struct {
	StaticDirectory string
	DataDirectory   string
	InputDirectory  string
	OutputDirectory string
	DownloadPath    string
	Port            int
}

type ProcessedImage struct {
	Name     string `json:"name"`
	Filepath string `json:"filepath"`
}
