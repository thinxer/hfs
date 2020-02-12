package main

import (
	"flag"
	"io"
	"log"
	"mime/multipart"
	"net"
	"os"
	"path"

	"net/http"
)

var (
	bind   = flag.String("bind", ":", "Address to listen on")
	upload = flag.Bool("upload", false, "Enable uploading")
)

func main() {
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(".")))
	if *upload {
		log.Println("Uploading enabled")
		http.Handle("/_/upload", uploadHandler)
	}
	lis, err := net.Listen("tcp", *bind)
	check(err)
	log.Println("Listening at", lis.Addr().String())
	addr := lis.Addr().(*net.TCPAddr)
	if addr.IP.IsUnspecified() {
		for _, ip := range localIP() {
			a := net.TCPAddr{
				IP:   ip,
				Port: addr.Port,
				Zone: addr.Zone,
			}
			log.Printf("Visit http://%s", a.String())
		}
	} else {
		log.Printf("Visit http://%s", lis.Addr().String())
	}
	check(http.Serve(lis, nil))
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

var uploadHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		mr, err := r.MultipartReader()
		check(err)
		handlePart := func(p *multipart.Part) {
			defer p.Close()
			name := path.Base(p.FileName())
			f, err := os.Create(name)
			check(err)
			defer f.Close()
			_, err = io.Copy(f, p)
			check(err)
		}
		for {
			p, err := mr.NextPart()
			if err != nil {
				break
			}
			handlePart(p)
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		w.Write([]byte(tmplUpload))
	}
})
