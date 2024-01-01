package main

import (
	"backend/internal"
	"log"
	"os"
	"path/filepath"
)

// TODO: Get config from external
var staticDirectory = "static"
var dataDirectory = "data"
var inputDirectory = filepath.Join(dataDirectory, "input")
var outputDirectory = filepath.Join(dataDirectory, "output")
var downloadPath = "/download/"
var port = 8080

func prepareDirectories() error {
	os.Mkdir(dataDirectory, 0775)
	os.Mkdir(inputDirectory, 0775)
	os.Mkdir(outputDirectory, 0775)
	return nil
}

func main() {
	prepareDirectories()
	config := internal.Config{
		StaticDirectory: staticDirectory,
		DataDirectory:   dataDirectory,
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
		DownloadPath:    downloadPath,
		Port:            port,
	}

	processor, err := internal.NewProcessor(config)
	if err != nil {
		log.Fatal(err)
	}
	go processor.Start()

	server := internal.NewServer(config, processor)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
