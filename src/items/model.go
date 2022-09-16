package items

import (
	"errors"
	"strings"
	"time"
)

type ItemGet struct {
	Id              int    `json:"id,omitempty"`
	UserId          int    `json:"userId,omitempty"`
	CategoryId      int    `json:"categoryId,omitempty"`
	BrandId         int    `json:"brandId,omitempty"`
	CreatedAt       int    `json:"createdAt,omitempty"`
	Price           int    `json:"price,omitempty"`
	DiscountedPrice int    `json:"discountedPrice,omitempty"`
	Description     string `json:"description,omitempty"`
	ModifiedAt      int    `json:"modifiedAt,omitempty"`
	DeletedAt       int    `json:"deletedAt,omitempty"`
}

type ItemPost struct {
	Id              int    `json:"id,omitempty"`
	UserId          int    `json:"userId,omitempty"`
	CategoryId      int    `json:"categoryId,omitempty"`
	BrandId         int    `json:"brandId,omitempty"`
	CreatedAt       int    `json:"createdAt,omitempty"`
	Price           int    `json:"price,omitempty"`
	DiscountedPrice int    `json:"discountedPrice,omitempty"`
	Description     string `json:"description,omitempty"`
}

func NewItemPost(userId int) ItemPost {
	now := int(time.Now().Unix())
	return ItemPost{UserId: userId, CreatedAt: now}
}

func (item ItemPost) checkFields() error {
	if item.UserId == 0 {
		return errors.New("UserId field can't be empty.")
	}
	if item.CategoryId == 0 {
		return errors.New("CategoryId field can't be empty.")
	}
	if item.BrandId == 0 {
		return errors.New("BrandId field can't be empty.")
	}
	if item.CreatedAt == 0 {
		return errors.New("CreatedAt field can't be empty.")
	}
	if item.Price == 0 {
		return errors.New("Price field can't be empty.")
	}
	if strings.TrimSpace(item.Description) == "" {
		return errors.New("Description field can't be empty.")
	}
	return nil
}

type ItemPatch struct {
	Id              int    `json:"id,omitempty"`
	CategoryId      int    `json:"categoryId,omitempty"`
	BrandId         int    `json:"brandId,omitempty"`
	Price           int    `json:"price,omitempty"`
	Discount        bool   `json:"discount,omitempty"`
	DiscountedPrice int    `json:"discountedPrice,omitempty"`
	Description     string `json:"description,omitempty"`
	ModifiedAt      int    `json:"modifiedAt,omitempty"`
	ChangeDeleted   bool   `json:"changeDeleted,omitempty"`
	DeletedAt       int    `json:"deletedAt,omitempty"`
}

func NewItemPatch(itemId int) ItemPatch {
	now := int(time.Now().Unix())
	return ItemPatch{Id: itemId, ModifiedAt: now}
}

func (item ItemPatch) checkFields() error {
	if item.Id == 0 {
		return errors.New("Wrong URL Path.")
	}
	if item.ModifiedAt == 0 {
		return errors.New("ModifiedAt field can't be empty.")
	}
	if item.CategoryId == 0 && item.BrandId == 0 && item.Price == 0 && item.DeletedAt == 0 && strings.TrimSpace(item.Description) == "" && item.DiscountedPrice == 0 && item.Discount == false {
		return errors.New("Include fields to be updated.")
	}
	return nil
}

type Filter struct {
	Brand []string `json:"brand,omitempty"`
	Prices []string  `json:"prices,omitempty"`
}
