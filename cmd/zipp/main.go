package main

import (
	"path/filepath"
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strings"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./static")
		return
	}

	if r.Method == http.MethodPost {

		//check attachment size
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			w.Write([]byte("your file is too big! Need to be lesser than 10mb"))
			fmt.Println(err)
			return
		}

		//get attachment
		_, fileHeader, err := r.FormFile("file-container")
		fileType := filepath.Ext(fileHeader.Filename)
		fmt.Println(fileType)

		//JPEG to PNG
		if strings.ToLower(fileType) == ".jpeg" {
			if err != nil {
				w.Write([]byte("Its look like file sended wrong"))
				fmt.Println(err)
				return
			}

			//geting image content
			f, err := fileHeader.Open()
			if err != nil {
				w.Write([]byte("Something is not ok"))
				fmt.Println(err)
				return
			}
			defer f.Close()

			//decoding jpeg encoded image to language recognizeable image.Image
			img, err := jpeg.Decode(f)
			if err != nil {
				w.Write([]byte("It is probably not jpeg"))
				fmt.Println(err, "jpeg decode error!")
				return
			}

			//writing to response result of converting jpeg to png
			var b bytes.Buffer
			if err := png.Encode(&b, img); err != nil {
				w.Write([]byte("Something is not ok"))
				fmt.Println(err)
				return
			}

			w.Header().Add("Content-Disposition", "attachment")
			io.Copy(w, &b)
		}

		//PNG to JPEG
		if strings.ToLower(fileType) == ".png" {
			if err != nil {
				w.Write([]byte("Its look like file sended wrong way"))
				fmt.Println(err)
				return
			}

			//geting image content
			f, err := fileHeader.Open()
			if err != nil {
				w.Write([]byte("Something is not ok"))
				fmt.Println(err)
				return
			}
			defer f.Close()

			//decoding jpeg encoded image to language recognizeable image.Image
			img, err := png.Decode(f)
			if err != nil {
				w.Write([]byte("It is probably not png"))
				fmt.Println(err, "jpeg decode error!")
				return
			}

			//writing to response result of converting png to jpeg
			var b bytes.Buffer
			if err := jpeg.Encode(&b, img, nil); err != nil {
				w.Write([]byte("Something is not ok"))
				fmt.Println(err)
				return
			}
			w.Header().Add("Content-Disposition", "attachment")
			io.Copy(w, &b)

		}

	}
	w.Write([]byte("Mehod not allowed!!"))
	return

}

func main() {
	http.HandleFunc("/", rootHandler)

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
