package items

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/fnmzgdt/e_shop/src/utils"
)

type Service interface {
	InsertItem(post *ItemPost) (int, error)
	GetItem(itemId int) (*ItemGet, error)
	GetImages(itemId int) (*[]Image, error)
	GetItems(limit int, filter *Filter) (*[]ItemGet, error)
	UpdateItem(item *ItemPatch) (int, error)
	DeleteItem(itemId int) (int, error)
	UploadItemImage(ctx context.Context, imageFile multipart.File) (string, error)
	InsertItemImage(imageUrl string) (int, error)
	InsertItemImageJunction(itemId, imageId, display_order int) (int, error)
}

type MySQL interface {
	ExecuteQuery(query string, values ...interface{}) (sql.Result, error)
	GetItem(query string, id int) (*ItemGet, error)
	GetItems(query string, values ...interface{}) (*[]ItemGet, error)
	GetImages(query string, itemId int) (*[]Image, error)
}

type GoogleImageBucket interface {
	UploadImage(ctx context.Context, objName string, imageFile multipart.File) (string, error)
}

type service struct {
	sqldb       MySQL
	imageBucket GoogleImageBucket
}

func NewPostsService(sqldb MySQL, cloudstorage GoogleImageBucket) Service {
	return &service{sqldb, cloudstorage}
}

func (s *service) InsertItem(item *ItemPost) (int, error) {
	query := "INSERT INTO items(user_id, category_id, brand_id, created_at, price, discounted_price, description) VALUES (?, ?, ?, FROM_UNIXTIME(?), ?, NULLIF(?, 0), ?)"
	res, err := s.sqldb.ExecuteQuery(query, item.UserId, item.CategoryId, item.BrandId, item.CreatedAt, item.Price, item.DiscountedPrice, item.Description)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *service) GetItem(itemId int) (*ItemGet, error) {
	query := "SELECT i.id, i.user_id, c.name as gategory_name, b.name as brand_name, UNIX_TIMESTAMP(created_at), price, discounted_price, description, UNIX_TIMESTAMP(modified_at) FROM items i JOIN categories c ON i.category_id = c.id JOIN brands b ON i.brand_id = b.id WHERE i.id = (?) AND i.deleted_at IS NULL;"
	item, err := s.sqldb.GetItem(query, itemId)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *service) GetImages(itemId int) (*[]Image, error) {
	query := "SELECT url FROM items_images JOIN images ON items_images.image_id = images.id WHERE items_images.item_id = (?) ORDER BY display_order ASC;"
	images, err := s.sqldb.GetImages(query, itemId)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (s *service) GetItems(limit int, filter *Filter) (*[]ItemGet, error) {
	query := "SELECT items.id, items.user_id, categories.name, brands.name, UNIX_TIMESTAMP(created_at), price, discounted_price, description, UNIX_TIMESTAMP(modified_at), images.url FROM items JOIN brands ON items.brand_id = brands.id JOIN categories ON items.category_id = categories.id LEFT JOIN items_images ON items.id = items_images.item_id LEFT JOIN images ON items_images.image_id = images.id WHERE deleted_at IS NULL AND items_images.display_order = 0"
	var arguments []interface{}
	if len(filter.Brand) > 0 {
		query += " AND brand_id IN ("
		brandQueryParams := strings.Split(filter.Brand[0], ",")
		for i := range brandQueryParams {
			if brandId, err := strconv.Atoi(brandQueryParams[i]); err == nil {
				query += "?, "
				arguments = append(arguments, brandId)
			}
		}
		if len(arguments) == 0 {
			query = strings.TrimSuffix(query, " AND brand_id IN (")
		} else {
			query = strings.TrimSuffix(query, ", ")
			query += ")"
		}
	}
	if len(filter.Prices) == 1 {
		lowPricehighPrice := strings.Split(filter.Prices[0], "-")
		lowPrice := lowPricehighPrice[0]
		highPrice := lowPricehighPrice[1]
		query += " AND (CASE WHEN discounted_price IS NULL THEN ((price > ?) AND (price < ?)) ELSE ((discounted_price > ?) AND (discounted_price < ?)) END)"
		arguments = append(arguments, lowPrice, highPrice, lowPrice, highPrice)
	}
	query += " LIMIT ?;"
	arguments = append(arguments, limit)
	fmt.Println(query)
	items, err := s.sqldb.GetItems(query, arguments...)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *service) UpdateItem(item *ItemPatch) (int, error) {
	var params []interface{}
	query := "UPDATE items SET"
	if item.CategoryId != 0 {
		query += " category_id = (?),"
		params = append(params, item.CategoryId)
	}
	if item.BrandId != 0 {
		query += " brand_id = (?),"
		params = append(params, item.BrandId)
	}
	if item.Price != 0 {
		query += " price = (?),"
		params = append(params, item.Price)
	}
	if item.Discount { //when discounted is true and discountedPrice = 0 / null : set it to null
		query += " discounted_price = NULLIF(?, 0),"
		params = append(params, item.DiscountedPrice)
	}
	if item.Description != "" {
		query += " description = (?),"
		params = append(params, item.Description)
	}
	if item.ModifiedAt != 0 {
		query += " modified_at = FROM_UNIXTIME(?),"
		params = append(params, item.ModifiedAt)
	}
	if item.ChangeDeleted {
		query += " deleted_at = FROM_UNIXTIME(NULLIF(?, 0)),"
		params = append(params, item.DeletedAt)
	}
	query = strings.TrimSuffix(query, ",")
	query += " WHERE id = (?);"
	params = append(params, item.Id)
	res, err := s.sqldb.ExecuteQuery(query, params...)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}

func (s *service) DeleteItem(itemId int) (int, error) {
	query := "DELETE FROM items WHERE id = (?);"
	res, err := s.sqldb.ExecuteQuery(query, itemId)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}

func (s *service) UploadItemImage(ctx context.Context, imageFile multipart.File) (string, error) {
	objName, err := utils.ObjNameFromUrl("")
	if err != nil {
		return "", err
	}
	imageUrl, err := s.imageBucket.UploadImage(ctx, objName, imageFile)
	if err != nil {
		return "", err
	}
	return imageUrl, nil
}

func (s *service) InsertItemImage(imageUrl string) (int, error) {
	query := "INSERT INTO images(url) VALUES (?);"
	res, err := s.sqldb.ExecuteQuery(query, imageUrl)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastId), nil
}

func (s *service) InsertItemImageJunction(itemId, imageId, display_order int) (int, error) { //mysql
	query := "INSERT INTO items_images(item_id, image_id, display_order) VALUES (?, ?, ?);"
	res, err := s.sqldb.ExecuteQuery(query, itemId, imageId, display_order)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}
