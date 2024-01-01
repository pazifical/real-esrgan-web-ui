package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	mux       *http.ServeMux
	port      int
	processor *Processor
}

func NewServer(config Config, processor *Processor) Server {
	server := Server{processor: processor, port: config.Port}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(config.StaticDirectory)))
	mux.Handle(config.DownloadPath, http.StripPrefix(config.DownloadPath, http.FileServer(http.Dir(config.OutputDirectory))))
	mux.HandleFunc("/api/upscale", server.handleUpscale)
	server.mux = mux
	return server
}

func (s *Server) Start() error {
	log.Printf("Open a web browser and go to http://localhost:%d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func (s *Server) handleUpscale(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGet(w, r)
	case "POST":
		s.handlePost(w, r)
	case "DELETE":
		s.handleDelete(w, r)
	default:
		w.WriteHeader(405)
	}
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list := s.processor.GetProcessedImages()
	json.NewEncoder(w).Encode(list)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	err := s.processor.DeleteAllProcessed()
	if err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(422)
		return
	}

	upscaleFactor, err := strconv.Atoi(r.FormValue("factor"))
	if err != nil {
		log.Printf("ERROR: %s", err)
		upscaleFactor = 2
	}
	s.processor.UpscaleFactor = upscaleFactor

	for _, formFile := range r.MultipartForm.File["file"] {
		log.Printf("INFO: received image %s", formFile.Filename)
		err := s.processor.AddImage(formFile)
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
	}
	w.WriteHeader(200)
}
