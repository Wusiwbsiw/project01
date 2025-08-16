##      数据库管理


推荐创建非root用户管理单一数据库
CREATE USER 'login_project01'@'localhost' IDENTIFIED BY 'Aa!123456';
GRANT ALL PRIVILEGES ON PROJECT01.* TO 'login_project01'@'localhost';

登陆账户
mysql -u login_project01 -p

显示数据库
SHOW DATABASES;

创建数据库
CREATE DATABASE PROJECT01;

使用数据库
USE PROJECT01;

删除数据库
DROP DATABASE PROJECT01;

新建表格

CREATE TABLE IF NOT EXISTS user_login_check (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
) AUTO_INCREMENT=100000;


CREATE TABLE IF NOT EXISTS user_chat_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    message_type ENUM('text', 'file_path', 'image_url', 'audio_url') NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    is_read TINYINT(1) NOT NULL DEFAULT 0,
    INDEX idx_chat_partners (sender_id, receiver_id, created_at)
);

显示表格
SHOW TABLES

表格新增列
ALTER TABLE table_name ADD COLUMN new_column_name data_type [constraints];

表格删除列
ALTER TABLE table_name DROP COLUMN column_name;

表格新增行
INSERT INTO table_name (column1, column2, column3) VALUES (value1, value2, value3);

表格删除行
DELETE FROM table_name WHERE condition;

##      常用查询语句

登陆查询
SELECT * FROM user_login_check WHERE id is "" AND WHERE password is ""

创建用户
"INSERT INTO user_login_check (username, password) VALUES (?, ?)"


聊天记录查询
(SELECT * FROM user_chat_history WHERE sender_id = ? AND receiver_id = ?)
UNION ALL
(SELECT * FROM user_chat_history WHERE sender_id = ? AND receiver_id = ?)
ORDER BY created_at ASC;

聊天记录插入
"INSERT INTO user_chat_history (sender_id, receiver_id,message_type,content,created_at,is_read) VALUES (?,?,?,?,?,?)"

