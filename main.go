package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var staticDirectory = "static"
var dataDirectory = "data"
var inputDirectory = filepath.Join(dataDirectory, "input")
var outputDirectory = filepath.Join(dataDirectory, "output")
var unprocessedFiles = make(map[string]string, 0)
var processedFiles = make(map[string]string, 0)
var port = 8080

func processFile(filename string) error {
	log.Printf("INFO: processing file %s", filename)

	err := os.Rename(filepath.Join(inputDirectory, filename), filepath.Join(outputDirectory, filename))
	if err != nil {
		return err
	}

	delete(unprocessedFiles, filename)
	processedFiles[filename] = filename

	return nil
}

func startFileProcessor() {
	for {
		time.Sleep(time.Second * 10)
		for _, imagePath := range unprocessedFiles {
			processFile(imagePath)
		}
	}
}

func prepareDirectories() error {
	os.Mkdir(dataDirectory, 0775)
	os.Mkdir(inputDirectory, 0775)
	os.Mkdir(outputDirectory, 0775)

	err := readProcessedFiles()
	if err != nil {
		return err
	}
	err = readUnprocessedFiles()
	if err != nil {
		return err
	}
	return nil
}

func readProcessedFiles() error {
	files, err := os.ReadDir(outputDirectory)
	if err != nil {
		return err
	}
	for _, file := range files {
		processedFiles[file.Name()] = file.Name()
	}
	return nil
}
func readUnprocessedFiles() error {
	files, err := os.ReadDir(inputDirectory)
	if err != nil {
		return err
	}
	unprocessedFiles = make(map[string]string, 0)
	for _, file := range files {
		unprocessedFiles[file.Name()] = file.Name()
	}
	return nil
}

func main() {
	prepareDirectories()
	go startFileProcessor()

	http.Handle("/", http.FileServer(http.Dir(staticDirectory)))
	http.Handle("/upscaled", http.FileServer(http.Dir(outputDirectory)))
	http.HandleFunc("/api/upscale", HandleUpscale)

	log.Printf("Open a web browser and go to http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func HandleUpscale(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(processedFiles)
		err := readProcessedFiles()
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}
	} else if r.Method == "POST" {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}

		var buf bytes.Buffer

		file, header, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}
		defer file.Close()
		log.Printf("INFO: received image %s", header.Filename)

		_, err = io.Copy(&buf, file)
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}

		f, err := os.Create(filepath.Join(inputDirectory, header.Filename))
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}
		defer f.Close()
		f.Write(buf.Bytes())
		buf.Reset()

		unprocessedFiles[header.Filename] = header.Filename

		w.WriteHeader(200)
	} else {
		w.WriteHeader(405)
	}
}
