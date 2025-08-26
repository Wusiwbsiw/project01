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


# 实际使用中可能会用到:

CREATE TABLE IF NOT EXISTS user_points (
    -- 用户 ID，既是主键也是外键
    user_id INT NOT NULL PRIMARY KEY,

    -- 积分，使用无符号整数，并设置默认值为 0
    points INT UNSIGNED NOT NULL DEFAULT 0,

    -- 最后更新时间，每次记录更新时自动更新为当前时间
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- 设置外键约束，关联到 login 表的 id 字段
    -- ON DELETE CASCADE 表示如果 login 表中的用户被删除了，他对应的积分记录也会被自动删除
    CONSTRAINT fk_user
    FOREIGN KEY (user_id) REFERENCES user_login_check(id)
    ON DELETE CASCADE
);

-- 创建 products 表，用于存储秒杀商品信息
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    total_stock INT UNSIGNED NOT NULL, -- 商品总库存
    description TEXT
)

-- 创建 orders 表，用于存储秒杀成功的订单
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 建议添加外键以保证数据完整性
    FOREIGN KEY (user_id) REFERENCES user_login_check(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    
    -- 一个用户对一个商品只能下一单，创建联合唯一索引来防止重复下单
    UNIQUE KEY idx_user_product (user_id, product_id)
)
