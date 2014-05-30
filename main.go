package main

import (
	"flag"

	"net/http"
)

var (
	flagListen = flag.String("listen", ":9999", "address to listen on")
)

func main() {
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(".")))
	panic(http.ListenAndServe(*flagListen, nil))
}
