package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/ivankuchin/connectopia.org/internal/pkg/api_handlers"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("HTTP request: " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func Run() {
	r := mux.NewRouter()
	// r.Use(loggingMiddleware)
	// test it like this: curl.exe -X POST http://localhost/api/v1/register -d '{\"domain\":\"www.bestbounty.ru\"}' -v
	r.HandleFunc("/api/v1/register", api_handlers.RegisterDomain).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/domains", api_handlers.GetDomains).Methods(http.MethodGet)
	// r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./web/docs/"))))
	// r.PathPrefix("/docs/").Handler(http.FileServer(http.Dir("./web/")))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	srv := &http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(80),
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 5,
		Handler:      r,
	}

	c := make(chan os.Signal, 1)

	// Run our server in a goroutine so that it doesn't block.
	go func(c chan os.Signal) {
		log.Printf("listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Print(err)
			c <- os.Interrupt
		}
	}(c)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	log.Print("shutthing down http-server")
}
