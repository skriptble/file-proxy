package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/skriptble/froxy"
)

var proxy froxy.ProxyBuilder

func main() {
	local := froxy.Dir(".")
	proxy = froxy.NewProxy()
	proxy.AddFileSource(local, "local")
	href, err := url.Parse("http://i.imgur.com/")
	if err != nil {
		log.Println(err)
	}
	rmt := froxy.NewRemote(*href)
	proxy.AddFileSource(rmt, "remote")

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	http.ListenAndServe(":8080", mux)
}

func handler(w http.ResponseWriter, req *http.Request) {
	paths := strings.SplitN(req.URL.Path, "/", 3)
	if len(paths) < 3 {
		// We don't have enough pieces of the url. Return a 404 Not Found.
		return
	}
	source := paths[1]
	name := paths[2]
	file, err := proxy.RetrieveFile(name, source)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	io.Copy(w, file)
}
