package db

// Product 结构体映射 products 表
type Product struct {
	ID          int
	Name        string
	TotalStock  uint
	Description string
}

// CreateProduct 创建一个新商品
func (d *Database) CreateProduct(name string, totalStock uint, description string) (int64, error) {
	query := "INSERT INTO products (name, total_stock, description) VALUES (?, ?, ?)"
	result, err := d.Conn.Exec(query, name, totalStock, description)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetProductByID 通过 ID 获取商品信息
func (d *Database) GetProductByID(productID int) (*Product, error) {
	var p Product
	query := "SELECT id, name, total_stock, description FROM products WHERE id = ?"
	row := d.Conn.QueryRow(query, productID)
	err := row.Scan(&p.ID, &p.Name, &p.TotalStock, &p.Description)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
