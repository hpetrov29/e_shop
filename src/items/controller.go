package items

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fnmzgdt/e_shop/src/responses"
)

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
		responses.JSONResponse(w, "Success.", []ItemGet{*item}, http.StatusOK)
		return
	}
}

func getItems(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := s.GetItems(10)
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
		userId, _ := strconv.Atoi(r.Header.Get("userId"))
		item := NewItemPost(userId)
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
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
		responses.JSONResponse(w, fmt.Sprintf("Successfully updated %d rows", rowsAffected), nil, http.StatusOK)
		return
	}
}

func postCategory(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// check for special claims
		userId := r.Header.Get("userId")
		category := newItemCategory(userId)
		_ = json.NewDecoder(r.Body).Decode(&category)
		if err := category.checkFields(); err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		lastId, err := s.InsertCategory(&category)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		category.Id = int(lastId)
		responses.JSONResponse(w, "Successful entry.", []ItemCategory{category}, 200)
		return
	}
}

func deleteCategory(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("userId")
		category := newItemCategory(userId)
		_ = json.NewDecoder(r.Body).Decode(&category)
		if err := category.checkName(); err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		//deletes the category and all its subcategories if any
		if err := s.DeleteCategory(&category); err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responses.JSONResponse(w, fmt.Sprintf("Successfully deleted category %s", category.Name), nil, 200)
		return
	}
}

func postBrand(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("userId")
		brand := createBrand(userId)
		_ = json.NewDecoder(r.Body).Decode(&brand)
		if err := brand.checkFields(); err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		lastId, err := s.InsertBrand(&brand)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		brand.Id = int(lastId)
		fmt.Println(brand)
		responses.JSONResponse(w, "Successful entry", []Brand{brand}, 200)
		return
	}
}

func deleteBrand(s Service) {

}

func postSizes(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("userId")
		var sizes []Size
		_ = json.NewDecoder(r.Body).Decode(&sizes)

		if len(sizes) == 0 {
			responses.JSONError(w, "Include at least one size", http.StatusBadRequest)
			return
		}
		for i := 0; i < len(sizes); i++ {
			sizes[i].setUserId(userId)
			if err := sizes[i].checkFields(); err != nil {
				responses.JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		for i := 0; i < len(sizes); i++ {
			lastId, err := s.InsertSize(&sizes[i])
			if err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sizes[i].setId(lastId)
		}
		responses.JSONResponse(w, fmt.Sprintf("Successfully inserted %d sizes.", len(sizes)), sizes, 200)
		return
	}
}

func deleteSizes(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//check claims
		var sizes []Size
		_ = json.NewDecoder(r.Body).Decode(&sizes)
		if len(sizes) == 0 {
			responses.JSONError(w, "Empty request body", http.StatusBadRequest)
			return
		}
		for i := 0; i < len(sizes); i++ {
			if err := sizes[i].checkName(); err != nil {
				responses.JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		for i := 0; i < len(sizes); i++ {
			if err := s.DeleteSize(&sizes[i]); err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		responses.JSONResponse(w, fmt.Sprintf("Successfully deleted sizes %v", sizes), nil, 200)
		return
	}
}

func postLocations(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("userId")
		var locations []Location
		_ = json.NewDecoder(r.Body).Decode(&locations)
		if len(locations) == 0 {
			responses.JSONError(w, "Empty request body", http.StatusBadRequest)
			return
		}
		for i := 0; i < len(locations); i++ {
			locations[i].setUserId(userId)
			if err := locations[i].checkFields(); err != nil {
				responses.JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		for i := 0; i < len(locations); i++ {
			lastId, err := s.InsertLocation(&locations[i])
			if err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			locations[i].setId(lastId)
		}
		responses.JSONResponse(w, "Successful entry", locations, 200)
		return
	}
}

//ON DELETE RESTRICT
func deleteLocations(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var locations []Location
		_ = json.NewDecoder(r.Body).Decode(&locations)
		if len(locations) == 0 {
			responses.JSONError(w, "Empty request body", http.StatusBadRequest)
			return
		}
		for i := 0; i < len(locations); i++ {
			if err := locations[i].checkId(); err != nil {
				responses.JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		for i := 0; i < len(locations); i++ {
			if err := s.DeleteLocation(&locations[i]); err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		responses.JSONResponse(w, fmt.Sprintf("Successfully deleted locations %v", locations), nil, 200)
		return
	}
}

func postDiscounts(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("userId")
		var discounts []Discount
		_ = json.NewDecoder(r.Body).Decode(&discounts)
		if len(discounts) == 0 {
			responses.JSONError(w, "Empty request body", http.StatusBadRequest)
			return
		}
		for i := 0; i < len(discounts); i++ {
			discounts[i].setUserId(userId)
			if err := discounts[i].checkFields(); err != nil {
				responses.JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		for i := 0; i < len(discounts); i++ {
			lastId, err := s.InsertDiscount(&discounts[i])
			if err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			discounts[i].setId(lastId)
		}
		responses.JSONResponse(w, "Successful entry", discounts, 200)
		return
	}
}

//deleting discount deletes all items_discount pairs
func deleteDiscounts(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var discounts []Discount
		_ = json.NewDecoder(r.Body).Decode(&discounts)
		if len(discounts) == 0 {
			responses.JSONError(w, "Empty request body", http.StatusBadRequest)
			return
		}
		for i := 0; i < len(discounts); i++ {
			if err := discounts[i].checkId(); err != nil {
				responses.JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		for i := 0; i < len(discounts); i++ {
			if err := s.DeleteDiscount(&discounts[i]); err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		responses.JSONResponse(w, fmt.Sprintf("Successfully deleted locations %v", discounts), nil, 200)
		return
	}
}

func applyDiscounts(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var itemdiscounts []ItemDiscount
		_ = json.NewDecoder(r.Body).Decode(&itemdiscounts)
		if len(itemdiscounts) == 0 {
			responses.JSONError(w, "Empty request body", http.StatusBadRequest)
			return
		}
		for i := 0; i < len(itemdiscounts); i++ {
			if err := itemdiscounts[i].checkFields(); err != nil {
				responses.JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		for i := 0; i < len(itemdiscounts); i++ {
			if err := s.InsertItemDiscount(&itemdiscounts[i]); err != nil {
				responses.JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		responses.JSONResponse(w, fmt.Sprintf("Successfully applied discounts %v", itemdiscounts), itemdiscounts, 200)
		return
	}
}

func ceaseDiscounts(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func createInventories(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func deleteInventories(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
