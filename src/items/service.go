package items

import (
	"database/sql"
	"strings"
)

type Service interface {
	InsertItem(post *ItemPost) (int, error)
	GetItem(itemId int) (*ItemGet, error)
	GetItems(limit int) (*[]ItemGet, error)
	UpdateItem(item *ItemPatch) (int, error)
	DeleteItem(itemId int) (int, error)
	InsertCategory(category *ItemCategory) (int64, error)
	DeleteCategory(category *ItemCategory) error
	InsertBrand(brand *Brand) (int64, error)
	InsertSize(size *Size) (int, error)
	DeleteSize(size *Size) error
	InsertLocation(location *Location) (int, error)
	DeleteLocation(location *Location) error
	InsertDiscount(dis *Discount) (int, error)
	DeleteDiscount(dis *Discount) error
	InsertItemDiscount(itemdis *ItemDiscount) error
}

type Rdbms interface {
	ExecuteQuery(query string, values ...interface{}) (sql.Result, error)
	GetItem(query string, id int) (*ItemGet, error)
	GetItems(query string, limit int) (*[]ItemGet, error)
}

type service struct {
	mysql Rdbms
}

func NewPostsService(db Rdbms) Service {
	return &service{db}
}

func (s *service) InsertItem(item *ItemPost) (int, error) {
	query := "INSERT INTO items(user_id, category_id, brand_id, created_at, price, discounted_price, description) VALUES (?, ?, ?, FROM_UNIXTIME(?), ?, ?, ?)"
	res, err := s.mysql.ExecuteQuery(query, item.UserId, item.CategoryId, item.BrandId, item.CreatedAt, item.Price, item.DiscountedPrice, item.Description)
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
	query := "SELECT id, user_id, category_id, brand_id, UNIX_TIMESTAMP(created_at), price, discounted_price, description, UNIX_TIMESTAMP(modified_at) FROM items WHERE id = (?) AND deleted_at IS NULL;"
	item, err := s.mysql.GetItem(query, itemId)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *service) GetItems(limit int) (*[]ItemGet, error) {
	query := "SELECT id, user_id, category_id, brand_id, UNIX_TIMESTAMP(created_at), price, discounted_price, description, UNIX_TIMESTAMP(modified_at) FROM items WHERE deleted_at IS NULL LIMIT ?;"
	items, err := s.mysql.GetItems(query, limit)
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
	res, err := s.mysql.ExecuteQuery(query, params...)
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
	res, err := s.mysql.ExecuteQuery(query, itemId)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}

func (s *service) InsertCategory(category *ItemCategory) (int64, error) {
	query := "call shop.add_subcategory(?, ?, ?);"
	res, err := s.mysql.ExecuteQuery(query, category.Name, category.ParentName, category.UserId)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func (s *service) DeleteCategory(category *ItemCategory) error {
	query := "call shop.delete_category(?);"
	_, err := s.mysql.ExecuteQuery(query, category.Name)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) InsertBrand(brand *Brand) (int64, error) {
	query := "INSERT INTO brands(name, user_id) VALUES(?, ?);"
	res, err := s.mysql.ExecuteQuery(query, brand.Name, brand.UserId)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func (s *service) InsertSize(size *Size) (int, error) {
	query := "INSERT INTO sizes(name, user_id) VALUES(?, ?);"
	res, err := s.mysql.ExecuteQuery(query, size.Name, size.UserId)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastId), nil
}

func (s *service) DeleteSize(size *Size) error {
	query := "DELETE FROM sizes WHERE name = ?;"
	_, err := s.mysql.ExecuteQuery(query, size.Name)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) InsertLocation(location *Location) (int, error) {
	query := "INSERT INTO locations(address, user_id) VALUES(?, ?);"
	res, err := s.mysql.ExecuteQuery(query, location.Address, location.UserId)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastId), nil
}

func (s *service) DeleteLocation(location *Location) error {
	query := "DELETE FROM locations WHERE id = (?);"
	_, err := s.mysql.ExecuteQuery(query, location.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) InsertDiscount(dis *Discount) (int, error) {
	query := "INSERT INTO discounts(code, amount, expires_at, user_id) VALUES(?, ?, FROM_UNIXTIME(?), ?);"
	res, err := s.mysql.ExecuteQuery(query, dis.Code, dis.Amount, dis.ExpiresAt, dis.UserId)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastId), nil
}

func (s *service) DeleteDiscount(dis *Discount) error {
	query := "DELETE FROM discounts WHERE id = (?);"
	_, err := s.mysql.ExecuteQuery(query, dis.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) InsertItemDiscount(itemdis *ItemDiscount) error {
	query := "INSERT INTO items_discounts(item_id, discount_id, valid_at) VALUES(?, ?, FROM_UNIXTIME(?));"
	_, err := s.mysql.ExecuteQuery(query, itemdis.ItemId, itemdis.DiscountId, itemdis.ValidAt)
	if err != nil {
		return err
	}
	return nil
}
