package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Conn *sql.DB
}

type ChatMessage struct {
	ID         		int64     `json:"id"`          // 消息的唯一ID
	SENDER_ID   	int       `json:"senderId"`    // 发送者ID
	RECEIVER_ID 	int       `json:"receiverId"`  // 接收者ID
	MESSAGE_TYPE   	string    `json:"type"`        // 消息类型，例如 "text", "file"
	CONTENT    		string    `json:"content"`     // 消息内容 (文本或文件路径/URL)
	CREATED_AT  	time.Time `json:"createdAt"`   // 发送时间
	IS_READ     	bool      `json:"isRead"`      // 是否已读
}

func InitDatabase(user, password, host, dbName string) (*Database, error) {
	// 初始化数据库连接池
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		user, password, host, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	fmt.Println("数据库连接池初始化成功！")
	return &Database{Conn: db}, nil
}

func (d *Database) Close() error {
	if d.Conn != nil {
		return d.Conn.Close()
	}
	return nil
}

func (d* Database) USER_Register(user_name string, password string)(int64,error){
	query := "INSERT INTO user_login_check (username, password) VALUES (?, ?)"
	result,err := d.Exec(query,user_name,password)
	if err != nil {
		return 0,err
	}

	id,err := result.LastInsertId()
	if err != nil{
		return 0,err
	}
	return id,nil
}

func (d *Database) USER_Login(user_id int64, user_password string) (ispermitted bool, error) {
	var storedPassword string
	query := "SELECT password FROM user_login_check WHERE id = ?"
	row := d.QueryRow(query,user_id)
	err := row.Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows{
			return false,nil
		}else{
			return false,err
		}
	}
	if storedPassword == user_password {
		return true,nil
	}else{
		return false,nil
	}
}

func (d* Database) USER_Logout(user_id int64)(ispermitted bool, error){
	// 可能用于管理用户状态
	
}

func (d* Database) USER_Reset_password(user_id int64, user_new_password string)(ispermitted bool,error){
	query := "UPDATE user_login_check SET password = ? WHERE id = ?"
	result, err := d.Exec(query, user_new_password, id)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (d* Database) USER_Reset_name(user_id int64, user_new_name string)(ispermitted bool,error){
	query := "UPDATE user_login_check SET name = ? WHERE id = ?"
	result, err := d.Exec(query, user_new_name, id)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (d* Database) CHAT_HISTORY_Push(user0_id int64,user1_id int64,message_type string,content string)(ispermitted bool,error){
	query := "INSERT INTO user_chat_history (sender_id, receiver_id,message_type,content) VALUES (?,?,?,?)"
	result,err := d.Exec(query,user0_id,user1_id,message_type,content)
		if err != nil {
		return 0, err
	}
	lastID, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

func (d* Database) CHAT_HISTORY_Pull(user0_id int64,user1_id int64,limit int64,offset int64)([]ChatMessage,error){
	query := "(SELECT id, sender_id, receiver_id, message_type, content, created_at, is_read 
		 FROM chat_messages 
		 WHERE sender_id = ? AND receiver_id = ?) UNION ALL (SELECT id, sender_id, receiver_id, message_type, content, created_at, is_read 
		 FROM chat_messages 
		 WHERE sender_id = ? AND receiver_id = ?) ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows,err := d.Query(query,user0_id,user1_id,user1_id,user0_id,limit,offset)
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	var message []ChatMessage

	for rows.Next(){
		var msg ChatMessage
		err := rows.Scan(
			&msg.ID,
			&msg.SENDER_ID,
			&msg.RECEIVER_ID,
			&msg.MESSAGE_TYPE,
			&msg.CONTENT,
			&msg.CREATED_AT,
			&msg.IS_READ
		)
		if err != nil {
			return nil,err
		}
		message = append(message,msg)
	}

	if err = rows.Err();err != nil {
		return nil,err
	}

	for i,j := 0,len(message)-1;i < j;i,j = i+1,j-1 {
		message[i],message[j] = message[j],message[i]
	}

	return message,nil
}

