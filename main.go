package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
)

// 環境変数からポート番号を取得するための構造体
type Env struct {
	Port uint16 `envconfig:"PORT" default:"8000"`
}

// ハンドラの実装
func newHandler() http.Handler {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		// 何かAPIを足したい場合はここに足す
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		})
	})

	// シングルページアプリケーションを配布するハンドラをNotFoundに設定
	router.NotFound(NotFoundHandler)
	return router
}

func main() {
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't parse environment variables: %s\n", err.Error())
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	server := &http.Server{
		Addr:    ":" + strconv.FormatUint(uint64(env.Port), 10),
		Handler: newHandler(),
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()
	fmt.Fprintf(os.Stderr, "start receiving at :%d\n", env.Port)
	fmt.Fprintln(os.Stderr, server.ListenAndServe())
}
