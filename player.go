package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

const (
	filePrefix = "/f/"
)

var (
	addr = flag.String("http", ":8080", "http listen address")
	root = flag.String("root", "/Volumes/media/Music/", "music root")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", Index)
	log.Printf("About to listen on 8080. Go to https://localhost:8080/")
	http.HandleFunc(filePrefix, File)
	// http.ListenAndServe(*addr, nil)
	err := http.ListenAndServeTLS(*addr, "moderation-cert.pem", "moderation-key.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	log.Println("Request:", r)
	log.Println("Response:", w)
	http.ServeFile(w, r, "./index.html")
}

func File(w http.ResponseWriter, r *http.Request) {
	fn := *root + r.URL.Path[len(filePrefix):]
	fi, err := os.Stat(fn)
	if err != nil {
		http.Error(w, err.String(), http.StatusNotFound)
		return
	}
	if fi.IsDirectory() {
		serveDirectory(fn, w, r)
		return
	}
	log.Println("Request:", r)
	log.Println("Response:", w)
	http.ServeFile(w, r, fn)
}

func serveDirectory(fn string, w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err, ok := recover().(os.Error); ok {
			http.Error(w, err.String(), http.StatusInternalServerError)
		}
	}()
	d, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	files, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}
	j := json.NewEncoder(w)
	if err := j.Encode(files); err != nil {
		panic(err)
	}
}
