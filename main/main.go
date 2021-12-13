package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Server struct {
	resMap map[string]string
}

func (s *Server) Download(rw http.ResponseWriter, r *http.Request) {
	idx := strings.Index(r.RequestURI, "/")
	fileFullPath := r.RequestURI[idx+1:]
	println(fileFullPath)
	file, err := ioutil.ReadFile(s.resMap[fileFullPath])
	if err != nil {
		return
	}
	reader := bytes.NewReader(file)

	fileName := path.Base(s.resMap[fileFullPath])
	fileName = url.QueryEscape(fileName)
	rw.Header().Add("Content-Type", "application/octet-stream")
	rw.Header().Add("content-disposition", "attachment; filename=\""+fileName+"\"")
	_, err = io.Copy(rw, reader)
	if err != nil {
		return
	}
}

func (s *Server) MakeALink(fileName, realPath string) []byte {
	return []byte("<tr><a href = " + realPath + ">" + fileName + "</a></tr><br>")
}
func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	println(r.RequestURI)
	if r.RequestURI == "/favicon.ico" {
		return
	}
	if r.RequestURI == "" || r.RequestURI == "/" {
		// os.Open("C:\Users\yj\Desktop")
		rw.Write([]byte("<html> <head>server Response</head><body><br>"))
		for realPath, fileName := range s.resMap {
			rw.Write(s.MakeALink(fileName, realPath))
		}
		rw.Write([]byte("</body></html>"))
	} else {
		s.Download(rw, r)
	}
}

func main() {
	flag.Parse()
	s := &Server{
		resMap: make(map[string]string),
	}
	for _, filePath := range flag.Args() {
		encoded := base64.StdEncoding.EncodeToString([]byte(filePath))
		http.Handle("/"+encoded, s)
		s.resMap[encoded] = filePath
	}
	err := http.ListenAndServe(":9090", s)
	if err != nil {
		return
	}
}
