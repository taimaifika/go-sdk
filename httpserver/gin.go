package httpserver

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/taimaifika/go-sdk/httpserver/middleware"
	"github.com/taimaifika/go-sdk/logger"
	"github.com/taimaifika/go-sdk/plugin/otel"
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
	version   string

	logger   logger.Logger
	svr      *myHttpServer
	router   *gin.Engine
	mu       *sync.Mutex
	handlers []func(*gin.Engine)
}

func New(name string) *ginService {
	return &ginService{
		name:     name,
		mu:       &sync.Mutex{},
		handlers: []func(*gin.Engine){},
	}
}

func (gs *ginService) Name() string {
	return gs.name + "-gin"
}

func (gs *ginService) Version() string {
	return gs.version
}

func (gs *ginService) InitFlags() {
	prefix := "gin"
	flag.IntVar(&gs.Config.Port, prefix+"Port", defaultPort, "gin server Port. If 0 => get a random Port")
	flag.StringVar(&gs.BindAddr, prefix+"addr", "", "gin server bind address")
	flag.StringVar(&ginMode, "gin-mode", "", "gin mode")
	flag.BoolVar(&ginNoLogger, "gin-no-logger", false, "disable default gin logger middleware")
}

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

		// recovery middleware
		gs.router.Use(gin.Recovery())

		// panic logger middleware
		gs.router.Use(middleware.PanicLogger())
		// otel middleware
		gs.router.Use(otelgin.Middleware(gs.name))
	}

	gs.svr = &myHttpServer{
		Server: http.Server{
			Handler: gs.router,
		},
	}

	return nil
}

func formatBindAddr(s string, p int) string {
	if strings.Contains(s, ":") && !strings.Contains(s, "[") {
		s = "[" + s + "]"
	}
	return fmt.Sprintf("%s:%d", s, p)
}

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

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Set up OpenTelemetry.
	otelShutdown, err := otel.SetupOTelSDK(ctx, gs.Name(), gs.Version())
	if err != nil {
		return err
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Start the server
	err = gs.svr.Serve(lis)

	if err != nil && err == http.ErrServerClosed {
		return nil
	}
	return err
}

func getPort(lis net.Listener) int {
	addr := lis.Addr()
	tcp, _ := net.ResolveTCPAddr(addr.Network(), addr.String())
	return tcp.Port
}

func (gs *ginService) Port() int {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	return gs.Config.Port
}

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

func (gs *ginService) URI() string {
	return formatBindAddr(gs.BindAddr, gs.Config.Port)
}

func (gs *ginService) AddHandler(hdl func(*gin.Engine)) {
	gs.isEnabled = true
	gs.handlers = append(gs.handlers, hdl)
}

func (gs *ginService) Reload(config Config) error {
	gs.Config = config
	<-gs.Stop()
	return gs.Run()
}

func (gs *ginService) GetConfig() Config {
	return gs.Config
}

func (gs *ginService) IsRunning() bool {
	return gs.svr != nil
}
