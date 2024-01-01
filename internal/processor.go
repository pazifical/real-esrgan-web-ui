package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Processor struct {
	processedFiles   map[string]string
	unprocessedFiles map[string]string
	downloadPath     string
	inputDirectory   string
	outputDirectory  string
}

func NewProcessor(config Config) (Processor, error) {
	processor := Processor{
		processedFiles:   make(map[string]string, 0),
		unprocessedFiles: make(map[string]string, 0),
		downloadPath:     config.DownloadPath,
		inputDirectory:   config.InputDirectory,
		outputDirectory:  config.OutputDirectory,
	}

	err := processor.readProcessedFiles()
	if err != nil {
		return processor, err
	}

	err = processor.readUnprocessedFiles()
	if err != nil {
		return processor, err
	}

	return processor, nil
}

func (p *Processor) GetProcessedImages() []ProcessedImage {
	list := make([]ProcessedImage, 0)
	for _, filename := range p.processedFiles {
		fPath := filepath.Join(p.downloadPath, filename)
		list = append(list, ProcessedImage{Name: filename, Filepath: fPath})
	}
	return list
}

func (p *Processor) AddImage(formFile *multipart.FileHeader) error {
	var buf bytes.Buffer
	defer buf.Reset()

	file, err := formFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(&buf, file)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(p.inputDirectory, formFile.Filename))
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(buf.Bytes())

	buf.Reset()

	p.unprocessedFiles[formFile.Filename] = formFile.Filename
	return nil
}

func (p *Processor) processFile(filename string) error {
	log.Printf("INFO: processing file %s", filename)

	fPath := filepath.Join(p.inputDirectory, filename)
	cmd := exec.Command(
		"python3",
		"./Real-ESRGAN/inference_realesrgan.py",
		"-n",
		"RealESRGAN_x4plus",
		"-i",
		fPath,
		"-o",
		p.outputDirectory,
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
	delete(p.unprocessedFiles, filename)

	err = p.readProcessedFiles()
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) Start() {
	for {
		for _, imagePath := range p.unprocessedFiles {
			err := p.processFile(imagePath)
			if err != nil {
				log.Printf("ERROR: %s", err)
			}
		}
		time.Sleep(time.Second * 10)
	}
}

func (p *Processor) readProcessedFiles() error {
	files, err := os.ReadDir(p.outputDirectory)
	if err != nil {
		return err
	}
	for _, file := range files {
		p.processedFiles[file.Name()] = file.Name()
	}
	return nil
}

func (p *Processor) readUnprocessedFiles() error {
	files, err := os.ReadDir(p.inputDirectory)
	if err != nil {
		return err
	}
	p.unprocessedFiles = make(map[string]string, 0)
	for _, file := range files {
		p.unprocessedFiles[file.Name()] = file.Name()
	}
	return nil
}
