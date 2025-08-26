package db

import (
	"context"
	"database/sql"
	"fmt"
)

type UserPoints struct {
	UserID    int64  `json:"userId"`
	Points    uint   `json:"points"` // 使用 uint (无符号整数) 确保积分不为负
	UpdatedAt string `json:"updatedAt"`
}

func (d *Database) InitPoints(userid int64) error {
	query := "INSERT INTO user_points (user_id,points) VALUES (?,0)"
	_, err := d.Conn.Exec(query, userid)
	return err
}

func (d *Database) ShowPoints(userid int64) (uint, error) {
	var points uint
	query := "SELECT points FROM user_points WHERE user_id = ?"
	row := d.Conn.QueryRow(query, userid)
	err := row.Scan(&points)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user with id %d not found in points table", userid)
		}
		return 0, err
	}
	return points, nil
}

func (d *Database) IncreasePoints(userid int64, amount uint) error {
	tx, err := d.Conn.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := "UPDATE user_points SET points = points + ? WHERE user_id = ?"
	result, err := tx.Exec(query, amount, userid)
	if err != nil {
		return fmt.Errorf("failed to increase points in transaction: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found, no points added", userid)

	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (d *Database) DecreasePoints(userid int64, amount uint) error {
	tx, err := d.Conn.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentPoints uint

	querycheck := "SELECT points FROM user_points WHERE user_id = ? FOR UPDATE"
	err = tx.QueryRow(querycheck, userid).Scan(&currentPoints)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with id %d not found", userid)
		}
		return fmt.Errorf("failed to select points for update: %w", err)
	}
	if currentPoints < amount {
		return fmt.Errorf("insufficient points : current %d, tried to decrease %d", currentPoints, amount)
	}
	query := "UPDATE user_points SET points = points - ? WHERE user_id = ?"
	result, err := tx.Exec(query, amount, userid)
	if err != nil {
		return fmt.Errorf("failed to decrease points in transaction: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found, no points added", userid)

	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
