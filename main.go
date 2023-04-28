package main

import (
	"flag"
	stdlog "log"
	"net/http"
	"os"

	"github.com/dombott/updog/github"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

type updog struct {
	log    logr.Logger
	client *github.Client
}

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	owner := flag.String("owner", "", "github repo owner")
	repo := flag.String("repo", "", "github repo name")
	token := os.Getenv("GH_TOKEN")
	flag.Parse()

	log := stdr.NewWithOptions(stdlog.New(os.Stderr, "", stdlog.LstdFlags), stdr.Options{LogCaller: stdr.All})
	log = log.WithName("updog")

	updog := &updog{
		log:    log,
		client: github.NewClient(token, *owner, *repo),
	}

	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/webhook", updog.webhook)
	log.Error(http.ListenAndServe(*addr, nil), "failed to serve http")
}
