package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var uploadDir string

func uploadGet(w http.ResponseWriter) {
	body := `
		<html>
			<body>
				<p>Upload file</p>
				<form action="/upload" method="post" enctype="multipart/form-data">
					<input type="file" name="file"/>
					<input type="Submit" value="Upload" name="submit"/>
				</form>
			</body>
		</html>
	`
	fmt.Fprintln(w, body)
}

func uploadPost(req *http.Request) {
	file, handler, err := req.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	filename := filepath.Base(handler.Filename)
	filePath := fmt.Sprintf("%s%s%s", uploadDir, string(filepath.Separator), filename)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0444)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

}

func uploadHanlder(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		//fmt.Fprintln(w, "GET upload")
		uploadGet(w)
	case http.MethodPost:
		uploadPost(req)
		fmt.Fprintln(w, "POST upload")

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	var port int
	var root string

	currentDir, _ := os.Getwd()
	flag.IntVar(&port, "port", 8090, "port to listen to...")
	flag.StringVar(&root, "root", currentDir, "server root directory ")
	flag.StringVar(&uploadDir, "uploadDir", currentDir, "directory where file will be uploaded")
	flag.Parse()
	fmt.Printf("(+) Starting web server on port %d....\n", port)
	fmt.Printf("(+) Server root directory: %s\n", root)
	fmt.Printf("(+) Upload directory: %s\n", uploadDir)
	http.HandleFunc("/upload", uploadHanlder)
	http.Handle("/", http.FileServer(http.Dir(root)))
	adress := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(adress, logRequest(http.DefaultServeMux))
	fmt.Println(err)
}
