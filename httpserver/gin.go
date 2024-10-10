package httpserver

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/taimaifika/go-sdk/httpserver/middleware"
	"github.com/taimaifika/go-sdk/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var (
	ginMode     string
	ginNoLogger bool
	defaultPort = 3000
)

type Config struct {
	Port         int    `json:"http_port"`
	BindAddr     string `json:"http_bind_addr"`
	GinNoDefault bool   `json:"http_no_default"`
}

type GinService interface {
	// block until ready
	Port() int
	isGinService()
}

type ginService struct {
	Config
	isEnabled bool
	name      string

	logger   logger.Logger
	svr      *myHttpServer
	router   *gin.Engine
	mu       *sync.Mutex
	handlers []func(*gin.Engine)
}

// New creates a new GinService.
func New(name string) *ginService {
	return &ginService{
		name:     name,
		mu:       &sync.Mutex{},
		handlers: []func(*gin.Engine){},
	}
}

// isGinService is a marker function to indicate that the service is a GinService.
func (gs *ginService) Name() string {
	return gs.name + "-gin"
}

// InitFlags initializes the flags.
func (gs *ginService) InitFlags() {
	prefix := "gin"
	flag.IntVar(&gs.Config.Port, prefix+"-port", defaultPort, "gin server Port. If 0 => get a random Port")
	flag.StringVar(&gs.BindAddr, prefix+"-addr", "", "gin server bind address")
	flag.StringVar(&ginMode, prefix+"-mode", "", "gin mode: debug, release, default is debug")

	// Logger
	flag.BoolVar(&ginNoLogger, prefix+"-no-logger", false, "disable default gin logger middleware, default is false")
}

// Configure configures the service.
func (gs *ginService) Configure() error {
	gs.logger = logger.GetCurrent().GetLogger("gin")

	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	gs.logger.Debug("init gin engine...")
	gs.router = gin.New()
	if !gs.GinNoDefault {
		if !ginNoLogger {
			gs.router.Use(gin.Logger())
		}

		// // recovery middleware (default)
		// gs.router.Use(gin.Recovery())

		// recovery middleware (custom)
		gs.router.Use(middleware.PanicLogger())

		// otelgin middleware
		gs.router.Use(otelgin.Middleware(gs.name))
	}

	gs.svr = &myHttpServer{
		Server: http.Server{
			Handler: gs.router,
		},
	}

	return nil
}

// formatBindAddr formats the bind address.
func formatBindAddr(s string, p int) string {
	if strings.Contains(s, ":") && !strings.Contains(s, "[") {
		s = "[" + s + "]"
	}
	return fmt.Sprintf("%s:%d", s, p)
}

// Run starts the service.
func (gs *ginService) Run() error {
	if !gs.isEnabled {
		return nil
	}

	if err := gs.Configure(); err != nil {
		return err
	}

	for _, hdl := range gs.handlers {
		hdl(gs.router)
	}

	addr := formatBindAddr(gs.BindAddr, gs.Config.Port)
	gs.logger.Debugf("start listen tcp %s...", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		gs.logger.Fatalf("failed to listen: %v", err)
	}

	gs.Config.Port = getPort(lis)

	gs.logger.Infof("listen on %s...", lis.Addr().String())

	// Start the server
	err = gs.svr.Serve(lis)

	if err != nil && err == http.ErrServerClosed {
		return nil
	}
	return err
}

// getPort returns the Port of the listener.
func getPort(lis net.Listener) int {
	addr := lis.Addr()
	tcp, _ := net.ResolveTCPAddr(addr.Network(), addr.String())
	return tcp.Port
}

// Port returns the Port of the service.
func (gs *ginService) Port() int {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	return gs.Config.Port
}

// Stop stops the service.
func (gs *ginService) Stop() <-chan bool {
	c := make(chan bool)

	go func() {
		if gs.svr != nil {
			_ = gs.svr.Shutdown(context.Background())
		}
		c <- true
	}()
	return c
}

// URI returns the URI of the service.
func (gs *ginService) URI() string {
	return formatBindAddr(gs.BindAddr, gs.Config.Port)
}

// AddHandler adds a handler to the gin service.
func (gs *ginService) AddHandler(hdl func(*gin.Engine)) {
	gs.isEnabled = true
	gs.handlers = append(gs.handlers, hdl)
}

// Reload reloads the service with the new config.
func (gs *ginService) Reload(config Config) error {
	gs.Config = config
	<-gs.Stop()
	return gs.Run()
}

// GetConfig returns the value of Config.
func (gs *ginService) GetConfig() Config {
	return gs.Config
}

// IsEnabled returns the value of isEnabled.
func (gs *ginService) IsEnabled() bool {
	return gs.isEnabled
}

// IsRunning returns true if the service is running.
func (gs *ginService) IsRunning() bool {
	return gs.svr != nil
}
