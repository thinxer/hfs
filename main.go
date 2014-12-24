package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path"

	"net/http"
)

var (
	flagListen = flag.String("listen", ":9999", "address to listen on")
)

func main() {
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/_/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			mr, err := r.MultipartReader()
			check(err)
			for {
				p, err := mr.NextPart()
				if err != nil {
					break
				}
				name := path.Base(p.FileName())
				f, err := os.Create(name)
				check(err)
				_, err = io.Copy(f, p)
				check(err)
				f.Close()
				p.Close()
			}
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			w.Write([]byte(tmplUpload))
		}
	})
	log.Println("Listing at", *flagListen)
	panic(http.ListenAndServe(*flagListen, nil))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const tmplUpload = `
<!DOCTYPE html>
<html>
<body>
<form enctype="multipart/form-data" method="POST">
    <input type="file" name="upload">
    <button type="submit">upload</button>
</form>
`
