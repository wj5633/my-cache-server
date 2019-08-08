package http

import (
	"net/http"
	"wj5633/my-cache-server/chapter1/server/cache"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen(addr string) {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	if addr == "" {
		addr = ":12345"
	}
	http.ListenAndServe(addr, nil)
}

func New(c cache.Cache) *Server {
	return &Server{c}
}
