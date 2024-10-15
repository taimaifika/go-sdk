package httpserver

import (
	"net"
	"net/http"
	"time"
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept implements the Accept method in the Listener interface
func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	_ = tc.SetKeepAlive(true)
	_ = tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// myHttpServer is a wrapper around http.Server that adds TCP keep-alive support
type myHttpServer struct {
	http.Server
}

// Serve accepts incoming connections on the Listener l, creating a new
func (srv *myHttpServer) Serve(lis net.Listener) error {
	return srv.Server.Serve(tcpKeepAliveListener{lis.(*net.TCPListener)})
}

// ServeTLS accepts incoming connections on the Listener l, creating a new Server with a TLS configuration
func (srv *myHttpServer) ServeTLS(lis net.Listener, certFile, keyFile string) error {
	return srv.Server.ServeTLS(tcpKeepAliveListener{lis.(*net.TCPListener)}, certFile, keyFile)
}
