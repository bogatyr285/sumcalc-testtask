package http

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type HTTPServer struct {
	h            http.Handler
	readTimeout  int
	writeTimeout int
	addr         string
	log          *zap.Logger
}

func NewHTTPServer(address string, h http.Handler, log *zap.Logger) (*HTTPServer, error) {
	return &HTTPServer{
		h:    h,
		addr: address,
		log:  log.Named("http-server"),
	}, nil
}

func (s *HTTPServer) WithReadTimeout(readTimeout int) *HTTPServer {
	s.readTimeout = readTimeout
	return s
}

func (s *HTTPServer) WithWriteTimeout(writeTimeout int) *HTTPServer {
	s.writeTimeout = writeTimeout
	return s
}

func (s *HTTPServer) Start() (func() error, error) {
	hserver := &http.Server{
		Addr:         s.addr,
		Handler:      s.h,
		WriteTimeout: time.Duration(s.writeTimeout),
		ReadTimeout:  time.Duration(s.writeTimeout),
	}

	go func() {
		if err := hserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("http listening err: %v", zap.Error(err))
		}
	}()

	return func() error {
		s.log.Info("shutting down http")
		return hserver.Close()
	}, nil
}
