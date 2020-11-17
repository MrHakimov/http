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

//ServeHTTP calls s.mux.ServeHTTP
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)

}

// Init wHmunannsupyet cepsep (permctpupyet sce Handler's)
func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	//	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)

	s.mux.HandleFunc("/banners.removeById", s.handleRemoveByID)
	s.mux.HandleFunc("/banners.getById", s.handleGetPostByID)

}

func (s *Server) handleGetPostByID(writer http.ResponseWriter, request *http.Request) {

	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)

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

	log.Print(request)

	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)

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

	log.Print(item)
}

func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {
	//log.Print(request)
	//log.Print(request.Header)
	//log.Print(request.Body)

	log.Print(request.RequestURI)
	log.Print(request.Method)
	log.Print(request.Header)
	log.Print(request.Header.Get("Content-Type"))

	log.Print(request.FormValue("id"))
	log.Print(request.FormValue("title"))
	log.Print(request.FormValue("content"))
	log.Print(request.FormValue("button"))
	log.Print(request.FormValue("link"))

	//log.Print(request.PostFormValue("tags"))

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print(err)
	}
	log.Printf("%s", body)

	err = request.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		log.Print(err)
	}

	log.Print(request.Form)
	log.Print(request.PostForm)
	log.Print(request.FormFile("image"))
	fileA, fileHeader, _ := request.FormFile("image")
	idParam := request.FormValue("id")
	fileNameInBanner := ""

	if fileA != nil {
		fileName := fileHeader.Filename
		log.Print(fileName)
		extenIndex := strings.Index(fileName, ".")
		fileExtension := fileName[extenIndex:]

		//fileA.Read()
		//content := make([]byte, 0)

		fileNameInBanner = idParam + fileExtension
	}
	//os.Create()
	//wd, err := os.Getwd()

	// wd, err := os.Executable()
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// fileN, _ := strconv.ParseInt(idParam, 10, 64)
	// fileN = fileN + 1
	// wd = wd + "/" + idParam + fileExtension

	// fileNameNew := item.Image
	//wdd := "web/banners" + "/" + strconv.FormatInt(fileN, 10) + fileExtension
	// wdd := "c:/projects/http/web/banners" + "/" + strconv.FormatInt(fileN, 10) + fileExtension
	// log.Print(wdd)
	// err = ioutil.WriteFile(wdd, content, 0600)
	// if err != nil {
	// 	log.Print(err)

	// }
	// fileN = fileN + 1
	// wdd1 := "web/banners" + "/" + strconv.FormatInt(fileN, 10) + fileExtension
	// err = ioutil.WriteFile(wdd1, content, 0600)
	// if err != nil {
	// 	log.Print(err)

	// }

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)

	}

	titleParam := request.FormValue("title")
	contentParam := request.FormValue("content")
	buttonParam := request.FormValue("button")
	linkParam := request.FormValue("link")

	banner := banners.Banner{
		ID: id,

		Title: titleParam,

		Content: contentParam,

		Button: buttonParam,

		Link: linkParam,

		Image: fileNameInBanner,
	}

	//banner := s.bannersSvc.Initial(request)

	//idParam := request.URL.Query().Get("id")

	//id, err := strconv.ParseInt(idParam, 10, 64)

	// if err != nil {
	// 	log.Print(err)
	// 	http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 	return
	// }

	item, err := s.bannersSvc.Save(request.Context(), &banner, fileA)

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

func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request) {
	log.Print(request)
}

func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	log.Print(request)
	log.Print(request.Header)
	log.Print(request.Body)
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
