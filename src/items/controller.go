package items

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fnmzgdt/e_shop/src/responses"
	"github.com/fnmzgdt/e_shop/src/utils"
)

const MAX_UPLOAD_SIZE int64 = 10 << 20
const MAX_UPLOAD_SIZE_SINGLE_IMAGE int64 = 2 << 20

func getItem(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		itemId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/items/items/"))
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		item, err := s.GetItem(itemId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				responses.JSONError(w, "Item not found", http.StatusNotFound)
				return
			}
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		images, err := s.GetImages(itemId)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		item.Images = *images
		responses.JSONResponse(w, "Success.", []ItemGet{*item}, http.StatusOK)
		return
	}
}

func getItems(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const limit = 10
		if err := r.ParseForm(); err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, err := json.Marshal(r.Form)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		filter := new(Filter)
		if err = json.Unmarshal(data, filter); err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		items, err := s.GetItems(limit, filter)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responses.JSONResponse(w, "Success.", *items, http.StatusOK)
		return
	}
}

func postItem(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			if err.Error() == "http: request body too large" {
				responses.JSONError(w, fmt.Sprintf("Maximum request body size is %vMB", MAX_UPLOAD_SIZE/(1000*1000)), http.StatusRequestEntityTooLarge)
				return
			}
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		userId, _ := strconv.Atoi(r.Header.Get("userId"))
		//checking the images before inserting item in db
		files := r.MultipartForm.File["imageFile"]
		if len(files) > 20 {
			responses.JSONError(w, "You can upload up to 20 images at once.", http.StatusBadRequest)
			return
		}
		if statusCode, err := utils.InspectImages(w, r, files, MAX_UPLOAD_SIZE_SINGLE_IMAGE); err != nil {
			responses.JSONError(w, err.Error(), statusCode)
			return
		}
		data := r.MultipartForm.Value["data"][0]
		item := NewItemPost(userId)
		if err := json.Unmarshal([]byte(data), &item); err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := item.checkFields(); err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		lastId, err := s.InsertItem(&item)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		item.Id = lastId
		imageArray := make([]Image, 0)
		for i := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			imageUrl, err := s.UploadItemImage(r.Context(), file)
			if err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			imageId, err := s.InsertItemImage(imageUrl)
			if err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = s.InsertItemImageJunction(item.Id, imageId, i)
			if err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			imageArray = append(imageArray, Image{Url: imageUrl})
		}
		item.Images = imageArray
		responses.JSONResponse(w, "Successful entry.", []ItemPost{item}, http.StatusCreated)
		return
	}
}

func updateItem(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		itemId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/items/items/"))
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		item := NewItemPatch(itemId)
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := item.checkFields(); err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		rowsAffected, err := s.UpdateItem(&item)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			responses.JSONResponse(w, fmt.Sprintf("Successfully updated %d rows", rowsAffected), nil, http.StatusNoContent)
			return
		}
		responses.JSONResponse(w, fmt.Sprintf("Successfully updated %d rows", rowsAffected), item, http.StatusOK)
		return
	}
}

func deleteItem(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		itemId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/items/items/"))
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		rowsAffected, err := s.DeleteItem(itemId)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			responses.JSONResponse(w, fmt.Sprintf("Successfully updated %d rows", rowsAffected), nil, http.StatusNoContent)
			return
		}
		//delete images related to the item
		responses.JSONResponse(w, fmt.Sprintf("Successfully updated %d rows", rowsAffected), nil, http.StatusOK)
		return
	}
}
