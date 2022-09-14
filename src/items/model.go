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

type ItemCategory struct {
	Name       string `json:"name,omitempty"`
	ParentName string `json:"parentName,omitempty"`
	UserId     string `json:"userId,omitempty"`
	Id         int    `json:"id,omitempty"`
}

func newItemCategory(userId string) ItemCategory {
	return ItemCategory{UserId: userId}
}

func (c ItemCategory) checkFields() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("Category Name field can't be empty.")
	}
	if strings.TrimSpace(c.ParentName) == "" {
		return errors.New("Parent Category Name field can't be empty.")
	}
	if strings.TrimSpace(c.UserId) == "" {
		return errors.New("User Id field can't be empty.")
	}
	return nil
}

func (c ItemCategory) checkName() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("Category Name field can't be empty.")
	}
	return nil
}

type Brand struct {
	Name   string `json:"name,omitempty"`
	UserId string `json:"userId,omitempty"`
	Id     int    `json:"id,omitempty"`
}

func createBrand(userId string) Brand {
	return Brand{UserId: userId}
}

func (b Brand) checkFields() error {
	if strings.TrimSpace(b.Name) == "" {
		return errors.New("Brand Name field can't be empty.")
	}
	if strings.TrimSpace(b.UserId) == "" {
		return errors.New("User Id field can't be empty.")
	}
	return nil
}

type Size struct {
	Name   string `json:"name,omitempty"`
	UserId string `json:"userId,omitempty"`
	Id     int    `json:"id,omitempty"`
}

type Sizes struct {
	SizesArr []Size `json:"sizes,omitempty"`
}

func (s Size) checkFields() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("Size Name field can't be empty.")
	}
	if strings.TrimSpace(s.UserId) == "" {
		return errors.New("User Id field can't be empty.")
	}
	return nil
}

func (s *Size) setUserId(userId string) {
	s.UserId = userId
}
func (s *Size) setId(id int) {
	s.Id = id
}

func (s Size) checkName() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("Size Name field can't be empty.")
	}
	return nil
}

type Location struct {
	Id      int    `json:"id,omitempty"`
	UserId  string `json:"userId,omitempty"`
	Address string `json:"address,omitempty"`
}

type Locations struct {
	LocationsArr []Location `json:"locations,omitempty"`
}

func (loc Location) checkFields() error {
	if strings.TrimSpace(loc.Address) == "" {
		return errors.New("Address field can't be empty.")
	}
	if strings.TrimSpace(loc.UserId) == "" {
		return errors.New("User Id field can't be empty.")
	}
	return nil
}

func (loc Location) checkId() error {
	if loc.Id == 0 {
		return errors.New("Id field can't be empty.")
	}
	return nil
}

func (loc *Location) setUserId(userId string) {
	loc.UserId = userId
}
func (loc *Location) setId(id int) {
	loc.Id = id
}

type Discount struct {
	Id        int    `json:"id,omitempty"`
	Code      string `json:"code,omitempty"`
	Amount    string `json:"amount,omitempty"`
	ExpiresAt int    `json:"expiresAt,omitempty"`
	UserId    string `json:"userId,omitempty"`
}

type Discounts struct {
	DiscountsArr []Discount `json:"discounts,omitempty"`
}

func (dis Discount) checkFields() error {
	if strings.TrimSpace(dis.Code) == "" {
		return errors.New("Code field can't be empty.")
	}
	if strings.TrimSpace(dis.Amount) == "" {
		return errors.New("Amount field can't be empty.")
	}
	if dis.ExpiresAt == 0 {
		return errors.New("ExpiresAt field can't be empty.")
	}
	if strings.TrimSpace(dis.UserId) == "" {
		return errors.New("UserId field can't be empty.")
	}
	return nil
}

func (dis *Discount) setUserId(userId string) {
	dis.UserId = userId
}
func (dis *Discount) setId(id int) {
	dis.Id = id
}

func (dis Discount) checkId() error {
	if dis.Id == 0 {
		return errors.New("Id field can't be empty.")
	}
	return nil
}

type ItemDiscount struct {
	ItemId     string `json:"itemId,omitempty"`
	DiscountId string `json:"discountId,omitempty"`
	ValidAt    int    `json:"validAt,omitempty"`
}

type ItemDiscounts struct {
	ItemDiscountsArr []ItemDiscount `json:"itemdiscounts,omitempty"`
}

func (i *ItemDiscount) checkFields() error {
	if strings.TrimSpace(i.ItemId) == "" {
		return errors.New("ItemId field can't be empty.")
	}
	if strings.TrimSpace(i.DiscountId) == "" {
		return errors.New("DiscountId field can't be empty.")
	}
	if i.ValidAt == 0 {
		return errors.New("ValidAt field can't be empty.")
	}
	return nil
}
