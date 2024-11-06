package fileupload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type HandlerGroup struct {
	FilePath string
}

func (g *HandlerGroup) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", g.upload)
	return mux
}

func (h *HandlerGroup) upload(rw http.ResponseWriter, req *http.Request) {
	const NameFormKey = "name"
	name := req.FormValue(NameFormKey)
	if len(name) == 0 {
		http.Error(rw, "empty name", http.StatusBadRequest)
		return
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		http.Error(rw, fmt.Sprint("read file: ", err), http.StatusBadRequest)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("Failed to close file: ", err)
		}
	}()

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll(h.FilePath, 0o755)
	if err != nil {
		http.Error(rw, fmt.Sprint("create directory: ", err), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(filepath.Join(h.FilePath, name))
	if err != nil {
		http.Error(rw, fmt.Sprint("create file: ", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("Failed to close file:", err)
		}
	}()

	err = dst.Chmod(0o744)
	if err != nil {
		http.Error(rw, fmt.Sprint("set permissions: ", err), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(rw, fmt.Sprint("upload file: ", err), http.StatusInternalServerError)
		return
	}
}
