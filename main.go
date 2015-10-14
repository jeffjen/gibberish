package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/jeffjen/gibberish/server"

	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type GibberishResponse struct {
	Next   int64    `json:"next" desc:"the cursor index for which next should be"`
	Mumble []string `json:"mumble" desc: "a slice of gibberish string"`
}

func Gibberish(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 403)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	const (
		MAX_GIBBERISH_ACCOUNT = 50

		GIBBERISH_ACCOUNT_SIZE = 10
	)

	next, err := strconv.ParseInt(r.Form.Get("next"), 10, 32)
	if err != nil {
		next = 0
	}

	log.WithFields(log.Fields{"next": next}).Info("Gibberish")

	if next+GIBBERISH_ACCOUNT_SIZE == MAX_GIBBERISH_ACCOUNT {
		next = 0 // reset iterator
	} else {
		next += GIBBERISH_ACCOUNT_SIZE
	}
	gibber := &GibberishResponse{Next: next}

	for iter := 0; iter < GIBBERISH_ACCOUNT_SIZE; iter++ {
		h := sha1.New()
		io.WriteString(h, time.Now().String())
		gibber.Mumble = append(gibber.Mumble, fmt.Sprintf("%x", h.Sum(nil)))
	}

	json.NewEncoder(w).Encode(gibber)
}

func init() {
	handler := server.GetServeMux()

	// start spewing out gibberish
	handler.HandleFunc("/gibberish", Gibberish)

	// load html page
	handler.Handle("/", http.FileServer(http.Dir("/var/public")))
}

func main() {
	srv := server.GetServer()

	// This is where I will serve the site
	srv.Addr = "0.0.0.0:9090"

	log.WithFields(log.Fields{"netloc": srv.Addr}).Info("service begin")

	log.Fatal(srv.ListenAndServe())
}
