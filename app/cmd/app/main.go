package main

import (
	"auth-service/app/cmd/internal/domain/user/storage"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"auth-service/app/cmd/internal/config"
	"auth-service/app/cmd/package/logger"

	"github.com/gorilla/mux"
)

func main() {
	l := logger.NewLogger()
	cfg := config.NewConfig(l)

	userRepo := storage.NewUserRepository(nil, l)
	fmt.Println(userRepo)

	router := mux.NewRouter()

	// register handlers
	postR := router.Methods(http.MethodPost).Subrouter()
	getR := router.Methods(http.MethodGet).Subrouter()
	putR := router.Methods(http.MethodPut).Subrouter()
	mailR := router.Path("verify").Methods(http.MethodPost).Subrouter()

	postR.HandleFunc("/123", HomeHandler)
	getR.HandleFunc("/123", HomeHandler)
	putR.HandleFunc("/123", HomeHandler)
	mailR.HandleFunc("/123", HomeHandler)

	// create server
	svrAddr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	svr := &http.Server{
		Addr:    svrAddr,
		Handler: router,
	}

	// start server
	go func() {
		err := svr.ListenAndServe()
		if err != nil {
			l.Error("failed t0 listen and serve server, err: %v", err)
			os.Exit(1)
		}
	}()

	// look for Interrupt or Kill for shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	l.Error("shutting down the server, signal %v", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	if err := svr.Shutdown(ctx); err != nil {
		l.Error("failed to shutdown server, err: %v", err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

}
