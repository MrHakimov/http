package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MrHakimov/http/pkg/banners"
)

// Server class for main data
type Server struct {
	mux        *http.ServeMux
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

func (s *Server) Init() {
	log.Println("Init method")
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerById)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveById)
}

func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	items, err := s.bannersSvc.All(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonResponse(writer, data)
}

func (s *Server) handleGetBannerById(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonResponse(writer, data)
}

func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {
	log.Println("requestURI: ", request.RequestURI)

	idParam := request.PostFormValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	banner := &banners.Banner{
		ID:      id,
		Title:   request.FormValue("title"),
		Content: request.FormValue("content"),
		Button:  request.FormValue("button"),
		Link:    request.FormValue("link"),
	}
	image, header, err := request.FormFile("image")
	if err == nil {
		var name = strings.Split(header.Filename, ".")
		//banner.Image = name[len(name)-1]
		banner.Image = name[1]
	}

	item, err := s.bannersSvc.Save(request.Context(), banner, image)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonResponse(writer, data)
}

func (s *Server) handleRemoveById(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonResponse(writer, data)
}

func jsonResponse(writer http.ResponseWriter, data []byte) {
	writer.Header().Set("Content-Type", "application/json")
	_, err := writer.Write(data)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}
