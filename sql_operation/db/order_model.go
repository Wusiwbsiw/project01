package db

import "time"

// Order 结构体映射 orders 表
type Order struct {
	ID        int64
	UserID    int
	ProductID int
	CreatedAt time.Time
}

// CreateOrderInBatch 批量创建订单 (用于后台任务)
// 我们接收一个 Order 对象的切片，在一个事务中批量插入
func (d *Database) CreateOrderInBatch(orders []*Order) error {
	tx, err := d.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 准备一个插入语句
	stmt, err := tx.Prepare("INSERT INTO orders (user_id, product_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// 遍历所有待插入的订单，在事务中执行
	for _, order := range orders {
		_, err := stmt.Exec(order.UserID, order.ProductID)
		if err != nil {
			// 如果任何一个插入失败，整个事务都会回滚
			return err
		}
	}

	// 所有都成功后，提交事务
	return tx.Commit()
}
