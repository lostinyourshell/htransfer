package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func coucou(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "coucou\n")
	fmt.Fprint(w, req.Method+"\n")
}

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
	f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0444)
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
	fmt.Println("Starting web server....")
	http.HandleFunc("/coucou", coucou)
	http.HandleFunc("/upload", uploadHanlder)
	http.Handle("/", http.FileServer(http.Dir("/home/seb")))
	http.ListenAndServe(":8090", logRequest(http.DefaultServeMux))
}
