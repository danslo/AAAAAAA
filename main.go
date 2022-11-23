package main

import (
	"github.com/charmbracelet/wish"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// SSH server begin
	s, err := wish.NewServer(
		ssh.PasswordAuth(passHandler),
		ssh.PublicKeyAuth(pkHandler),
		wish.WithAddress(":2222"),
		wish.WithMiddleware(
			lm.Middleware(),
			loginBubbleteaMiddleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	s.ConnectionFailedCallback = func(conn net.Conn, err error) {
		log.Println("Connection failed:", err)
	}
	done := make(chan os.Signal, 0)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", host, port)
	go func() {
		err = s.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	// start web based terminal
	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
