package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func InspectImages(w http.ResponseWriter, r *http.Request, files []*multipart.FileHeader, maxSize int64) (int, error) {
	for i := range files {
		if files[i].Size > maxSize {
			return http.StatusUnprocessableEntity, fmt.Errorf("The uploaded image is too big: %s. Please use an image less than 2MB in size", files[i].Filename)
		}
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			return http.StatusInternalServerError, err
		}
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			return http.StatusUnprocessableEntity, fmt.Errorf("The uploaded file format is not allowed: %s. Please upload a JPEG or PNG image", files[i].Filename)
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}
	return 0, nil
}
