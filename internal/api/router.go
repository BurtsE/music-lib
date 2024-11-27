package api

import (
	"MusicLibrary/internal/config"
	"MusicLibrary/internal/database"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	timeout     = 10 * time.Second
	maxFileSize = 128
)

// Mock для симуляции похода в апи с деталями песни
type getDetailsFunc func(song songConstructor) (details detailsConstructor, err error)

var getDetails getDetailsFunc

type Server struct {
	server             *http.Server
	storage            *database.Postgres
	externalStorageURL string
	client             http.Client
	port               int
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

func NewServer(cfg *config.Config, storage *database.Postgres) *Server {
	srv := &Server{
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", cfg.Host.Address, cfg.Host.Port),
		},
		storage:            storage,
		externalStorageURL: cfg.HelperApi,
		client:             http.Client{Timeout: timeout},
	}
	srv.server.Handler = srv.catcher(http.MaxBytesHandler(http.DefaultServeMux, maxFileSize))

	getDetails = srv.mockGetDetails //Для тестов: srv.mockGetDetails

	http.HandleFunc("GET /api/ping/", srv.ping)

	http.Handle("GET /api/songs/list/", paginate(http.HandlerFunc(srv.list)))

	http.HandleFunc("GET /api/songs/text/{song_id}", srv.text)
	http.HandleFunc("POST /api/songs/upload", srv.uploadSong)
	http.HandleFunc("PATCH /api/songs/update/{song_id}", srv.updateSong)
	http.HandleFunc("DELETE /api/songs/delete/{song_id}", srv.deleteSong)

	return srv
}

func (s *Server) catcher(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("recovered from: %v", rec)
			}
		}()
		log.Println("new request: ", r.URL.RequestURI())
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	pageNumber := 0
	ok := false
	group := r.FormValue("group")
	name := r.FormValue("name")
	if pageNumber, ok = r.Context().Value(PageIDKey).(int); !ok {
		log.Println("page number not found in context")
	}
	songs, err := s.storage.GetSongs(group, name, pageNumber)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "unable to get songs", http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(songs)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(res)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (s *Server) text(w http.ResponseWriter, r *http.Request) {
	songIDStr := r.PathValue("song_id")
	songID, err := strconv.ParseUint(songIDStr, 10, 64)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "invalid song id", http.StatusBadRequest)
		return
	}
	song, err := s.storage.GetSong(uint(songID))
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "could not get song", http.StatusInternalServerError)
		return
	}
	data, err := splitText(song.Text)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "could not get song", http.StatusInternalServerError)
		return
	}
}

func (s *Server) deleteSong(w http.ResponseWriter, r *http.Request) {
	songIDStr := r.PathValue("song_id")
	songID, err := strconv.ParseUint(songIDStr, 10, 64)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "invalid song id", http.StatusBadRequest)
		return
	}
	err = s.storage.DeleteSong(uint(songID))
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "could not delete song", http.StatusInternalServerError)
		return
	}
	log.Printf("song with id: %s deleted", songIDStr)
}

func (s *Server) uploadSong(w http.ResponseWriter, r *http.Request) {
	var (
		song songConstructor
	)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "could not get request body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &song)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "could not parse request body", http.StatusInternalServerError)
		return
	}
	details, err := getDetails(song)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "could not get song details", http.StatusServiceUnavailable)
		return
	}
	err = s.storage.CreateSong(songFromConstructors(song, details))
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "could not create song", http.StatusInternalServerError)
		return
	}
	log.Printf("song %s by %s added", song.Name, song.Group)
}

func (s *Server) updateSong(w http.ResponseWriter, r *http.Request) {
	songIDStr := r.PathValue("song_id")
	songID, err := strconv.ParseUint(songIDStr, 10, 64)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "invalid song id", http.StatusBadRequest)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "could not get request body", http.StatusInternalServerError)
		return
	}
	songParams := songParams{}
	err = json.Unmarshal(data, &songParams)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "could not parse request body", http.StatusInternalServerError)
		return
	}
	song := fillSongParams(songParams)
	song.ID = uint(songID)
	err = s.storage.UpdateSong(song)
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "could not update song", http.StatusInternalServerError)
		return
	}
	log.Printf("song with id %d changed: %v", songID, song)
	_, err = w.Write([]byte("successful update"))
	if err != nil {
		log.Println(r.URL.RequestURI(), ":", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) mockGetDetails(song songConstructor) (details detailsConstructor, err error) {
	body := []byte(`{
			"releaseDate": "16.07.2006",
			"text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			"link": " https://www.youtube.com/watch?v=Xsp3_a-PMTw"
		}`)
	err = json.Unmarshal(body, &details)
	if err != nil {
		return detailsConstructor{}, err
	}
	return
}
func (s *Server) getDetailsFromApi(song songConstructor) (details detailsConstructor, err error) {
	req, err := http.NewRequest("GET", s.externalStorageURL, nil)
	if err != nil {
		return detailsConstructor{}, err
	}
	q := req.URL.Query()
	q.Add("group", song.Group)
	q.Add("song", song.Name)
	req.URL.RawQuery = q.Encode()

	resp, err := s.client.Do(req)
	if err != nil {
		return detailsConstructor{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return detailsConstructor{}, errors.New("service request finished with status code: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return detailsConstructor{}, err
	}
	err = json.Unmarshal(body, &details)
	if err != nil {
		return detailsConstructor{}, err
	}
	return
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
