package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

func resizeCompressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 << 20) //limit to 10 mb
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}
	//get file from that we will find image
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	targetWidthStr := r.FormValue("width")
	targetHeightStr := r.FormValue("height")
	qualityStr := r.FormValue("quality")
	format := r.FormValue("format")
	//set default format
	if format == "" {
		format = "jpeg"
	}
	//get height , width and quality
	targetWidth, err := strconv.ParseUint(targetWidthStr, 10, 32) //base 10 decimal which is a 32 bit integer
	if err != nil {
		http.Error(w, "Invalid width value", http.StatusBadRequest)
		return
	}
	targetHeight, err := strconv.ParseUint(targetHeightStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid height value", http.StatusBadRequest)
		return
	}
	quality, err := strconv.ParseUint(qualityStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid quality value", http.StatusBadRequest)
		return
	}
	image, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusInternalServerError)
		return
	}
	resizedImage := resize.Resize(uint(targetWidth), uint(targetHeight), image, resize.Lanczos3)
	var outputBytes bytes.Buffer
	switch strings.ToLower((format)) {
	case "jpeg":
		err = jpeg.Encode(&outputBytes, resizedImage, &jpeg.Options{Quality: int(quality)})
		if err != nil {
			http.Error(w, "Failed to encode image", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
	case "png":
		err = png.Encode(&outputBytes, resizedImage)
		if err != nil {
			http.Error(w, "Failed to encode image", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
	default:
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}
	w.Write(outputBytes.Bytes())
}

func main() {
	http.HandleFunc("/resize", resizeCompressHandler)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
