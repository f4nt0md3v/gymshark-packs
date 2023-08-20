package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/f4nt0md3v/gymshark-packs/config"
	"github.com/f4nt0md3v/gymshark-packs/server/api/controllers"
	"github.com/f4nt0md3v/gymshark-packs/server/api/services"
)

const compressLevel = 5

type Server struct {
	basePath       string
	address        string
	allOrigins     bool
	packer         services.Packer
	logger         *zap.Logger
	idleConnClosed chan struct{}
	masterCtx      context.Context
}

func NewServer(ctx context.Context, config *config.Configuration, opts ...ServerResource) *Server {
	srv := &Server{
		basePath:       config.BasePath,
		address:        config.ListenAddr,
		idleConnClosed: make(chan struct{}), // a way to identify unclosed connections
		masterCtx:      ctx,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

type ServerResource func(srv *Server)

func WithPackerService(p services.Packer) ServerResource {
	return func(srv *Server) {
		srv.packer = p
	}
}

func WithLogger(logger *zap.Logger) ServerResource {
	return func(srv *Server) {
		srv.logger = logger
	}
}

// setupRouter initialize HTTP router.
// Function used to connect middlewares and map resources.
func (srv *Server) setupRouter(
	packController *controllers.PackController,
) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.NoCache)                                      // no-cache
	r.Use(middleware.Logger)                                       // logs beginning and end of each request with time to process
	r.Use(middleware.Recoverer)                                    // recovers from panics and logs stack trace
	r.Use(middleware.RealIP)                                       // sets RemoteAddr for each request with headers X-Forwarded-For or X-Real-IP
	r.Use(middleware.Compress(compressLevel, []string{"gzip"}...)) // uses default compression
	// mount controllers
	r.Mount(srv.basePath+"/pack", packController.Routes())

	return r
}

// Run runs HTTP server using APIServer{}.
func (srv *Server) Run(
	packController *controllers.PackController,
) error {
	const (
		readTimeout              = 5 * time.Second
		writeTimeout             = 30 * time.Second
		defaultReadHeaderTimeout = 2 * time.Second
	)

	s := &http.Server{
		Addr: srv.address,
		Handler: srv.setupRouter(
			packController,
		),
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	go srv.gracefulShutdown(s)
	srv.logger.Info("serving http", zap.String("address", srv.address))
	if err := s.ListenAndServe(); err != nil {
		srv.logger.Error(err.Error())
		return err
	}

	return nil
}

// GracefulShutdown handles unfinished calls before shutting down gracefully.
func (srv *Server) gracefulShutdown(httpSrv *http.Server) {
	<-srv.masterCtx.Done()

	if err := httpSrv.Shutdown(context.Background()); err != nil {
		srv.logger.Error("http server shutdown error", zap.Error(err))
	}

	srv.logger.Info("http server has processed all idle connections")
	close(srv.idleConnClosed)
}

// Wait awaiting to complete all connections.
func (srv *Server) Wait() {
	<-srv.idleConnClosed
	srv.logger.Info("http server has processed all idle connections")
}
