package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var staticDirectory = "static"
var dataDirectory = "data"
var inputDirectory = filepath.Join(dataDirectory, "input")
var outputDirectory = filepath.Join(dataDirectory, "output")
var downloadPath = "/download/"
var unprocessedFiles = make(map[string]string, 0)
var processedFiles = make(map[string]string, 0)
var port = 8080

func processFile(filename string) error {
	log.Printf("INFO: processing file %s", filename)

	fPath := filepath.Join(inputDirectory, filename)
	cmd := exec.Command(
		"python3",
		"./Real-ESRGAN/inference_realesrgan.py",
		"-n",
		"RealESRGAN_x4plus",
		"-i",
		fPath,
		"-o",
		outputDirectory,
	)

	var stdout, errout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &errout

	err := cmd.Run()
	fmt.Println(stdout.String())
	fmt.Println(errout.String())
	if err != nil {
		return err
	}

	err = os.Remove(fPath)
	if err != nil {
		return err
	}
	delete(unprocessedFiles, filename)

	err = readProcessedFiles()
	if err != nil {
		return err
	}

	return nil
}

func startFileProcessor() {
	for {
		for _, imagePath := range unprocessedFiles {
			err := processFile(imagePath)
			if err != nil {
				log.Printf("ERROR: %s", err)
			}
		}
		time.Sleep(time.Second * 10)
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
	http.Handle(downloadPath, http.StripPrefix(downloadPath, http.FileServer(http.Dir(outputDirectory))))
	http.HandleFunc("/api/upscale", HandleUpscale)

	log.Printf("Open a web browser and go to http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

type ProcessedImage struct {
	Name     string `json:"name"`
	Filepath string `json:"filepath"`
}

func HandleUpscale(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		list := make([]ProcessedImage, 0)
		for _, filename := range processedFiles {
			fPath := filepath.Join(downloadPath, filename)
			list = append(list, ProcessedImage{Name: filename, Filepath: fPath})
		}
		json.NewEncoder(w).Encode(list)

	} else if r.Method == "POST" {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}

		var buf bytes.Buffer

		for _, formFile := range r.MultipartForm.File["file"] {
			log.Printf("INFO: received image %s", formFile.Filename)
			file, err := formFile.Open()
			if err != nil {
				log.Printf("ERROR: %s", err)
				continue
			}
			defer file.Close()

			_, err = io.Copy(&buf, file)
			if err != nil {
				log.Printf("ERROR: %s", err)
				continue
			}

			f, err := os.Create(filepath.Join(inputDirectory, formFile.Filename))
			if err != nil {
				log.Printf("ERROR: %s", err)
				continue
			}
			defer f.Close()
			f.Write(buf.Bytes())

			buf.Reset()

			unprocessedFiles[formFile.Filename] = formFile.Filename
		}

		w.WriteHeader(200)
	} else {
		w.WriteHeader(405)
	}
}
