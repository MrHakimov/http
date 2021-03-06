package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"strconv"
	"strings"

	"github.com/Iftikhor99/httpForms/pkg/banners"
)

// Server class for main data
type Server struct {
	mux *http.ServeMux

	bannersSvc *banners.Service
}

// NewServer creates new server
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}

// ServeHTTP calls s.mux.ServeHTTP
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init initializes all supported operations
func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getById", s.handleGetPostByID)
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)

	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveByID)
}

func (s *Server) handleGetPostByID(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleRemoveByID(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {
	_, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print(err)
	}

	err = request.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		log.Print(err)
	}

	fileContent, fileHeader, _ := request.FormFile("image")
	idParam := request.FormValue("id")
	fileNameInBanner := ""

	if fileContent != nil {
		fileExtension := fileHeader.Filename[strings.Index(fileHeader.Filename, "."):]
		fileNameInBanner = idParam + fileExtension
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}

	_title := request.FormValue("title")
	_content := request.FormValue("content")
	_button := request.FormValue("button")
	_ling := request.FormValue("link")

	banner := banners.Banner{
		ID:      id,
		Title:   _title,
		Content: _content,
		Button:  _button,
		Link:    _ling,
		Image:   fileNameInBanner,
	}

	item, err := s.bannersSvc.Save(request.Context(), &banner, fileContent)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request) {
	log.Print(request)
}

func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	item, err := s.bannersSvc.All(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		log.Print(err)
	}

	log.Print(item)
}
