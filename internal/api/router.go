package api

import "net/http"

type Server struct {
	srv *http.Server
}

func NewServer() *Server {
	srv := &Server{
		srv: &http.Server{
			Handler: http.DefaultServeMux,
		},
	}
	http.HandleFunc("GET /api/", srv.index)
	http.HandleFunc("GET /api/info/{song_name}", srv.info)
	http.HandleFunc("DELETE /api/delete", srv.deleteSong)
	http.HandleFunc("POST /api/upload", srv.uploadSong)
	return srv
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) info(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) deleteSong(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) uploadSong(w http.ResponseWriter, r *http.Request) {

}
