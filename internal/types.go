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

type ByName []ProcessedImage

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
