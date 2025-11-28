package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
	"todo/actors"
	"todo/logger"
	"todo/web"

	"github.com/google/uuid"
)

/*
go run api.go
http://127.0.0.1:8080/about
http://127.0.0.1:8080/list/<TodoName>
*/
func main() {
	ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())

	mux := http.NewServeMux()

	mux.HandleFunc("/createlist", web.CreateListHandler)
	mux.HandleFunc("/getlist", web.GetListHandler)
	mux.HandleFunc("/additem", web.AddItemHandler)
	mux.HandleFunc("/deleteitem", web.DeleteItemHandler)
	mux.HandleFunc("/updateitemdescription", web.UpdateItemDescriptionHandler)
	mux.HandleFunc("/updateitemstatus", web.UpdateItemStatusHandler)
	mux.HandleFunc("/list/{listname}", web.DynamicListHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(strings.ToLower(r.URL.Path), ".html") {
			http.ServeFile(w, r, "./web/static/"+path.Base(r.URL.Path))
		} else {
			http.ServeFile(w, r, "./web/static/"+path.Base(r.URL.Path)+".html")
		}
	})

	handler := web.AddLogLayer(mux)
	handler = web.AddTraceIDLayer(handler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// goroutine to run the actor processing
	go func() {
		ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())
		logger.GetCtxLogger(ctx).Info("Starting Actors Thread")
		actors.GetActor().ProcessMessages(ctx)
		logger.GetCtxLogger(ctx).Info("Actor thread stopped")
	}()

	// goroutine to run the server
	go func() {
		ctx := context.WithValue(context.Background(), logger.TraceIdKey{}, uuid.NewString())
		logger.GetCtxLogger(ctx).Info("Starting HTTP Server Thread")

		logger.GetCtxLogger(ctx).Info("Server listening on :8080")
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.GetCtxLogger(ctx).Error("HTTP server Close error: " + err.Error())
		}
		logger.GetCtxLogger(ctx).Info("Server stopped listening on :8080")
		actors.GetActor().Messages <- actors.Message{
			Request: actors.Request{Operation: "quit"},
			Ctx:     ctx,
			Quit:    true,
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.GetCtxLogger(ctx).Error("HTTP shutdown error: " + err.Error())
		err = srv.Close()
		if err != nil {
			logger.GetCtxLogger(ctx).Error("Server shutdown error: " + err.Error())
		}
	}
	logger.GetCtxLogger(ctx).Info("Shutdown complete")
}
