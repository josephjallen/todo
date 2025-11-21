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
	"todo/logger"
	"todo/web"
)

/*
go run api.go
*/
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/createlist", web.CreateListHandler)
	mux.HandleFunc("/getlist", web.GetListHandler)
	mux.HandleFunc("/additem", web.AddItemHandler)
	mux.HandleFunc("/deleteitem", web.DeleteItemHandler)
	mux.HandleFunc("/updateitemdescription", web.UpdateItemDescriptionHandler)
	mux.HandleFunc("/updateitemstatus", web.UpdateItemStatusHandler)
	/* http://127.0.0.1:8080/list/weeklytodo */
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

	// Parallel goroutine to run the server
	go func() {
		logger.InfoLog(nil, "Server listening on :8080")
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorLog(nil, "HTTP server Close error: "+err.Error())
		}
		logger.InfoLog(nil, "Server stopped listening on :8080")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.ErrorLog(nil, "HTTP shutdown error: "+err.Error())
		err = srv.Close()
		if err != nil {
			logger.ErrorLog(nil, "Server shutdown error: "+err.Error())
		}
	}
	logger.InfoLog(nil, "Shutdown complete")
}
