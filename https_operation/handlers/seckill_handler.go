package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"https_operation/auth"
	"https_operation/redis"
	"log"
	"net/http"
	"os"
	"sql_operation/db"
	"time"
)

var (
	initStockLua string
	doSeckillLua string
)

const orderQueueKey = "order:queue"

// SeckillHandler 结构体，持有数据库和 Redis 客户端的连接
type SeckillHandler struct {
	DB    *db.Database
	Redis *redis.AppRedisClient
}

// InitSeckillRequest 结构体，用于解析初始化请求的 JSON
type InitSeckillRequest struct {
	ProductID int `json:"productId"`
}

type DoSeckillRequest struct {
	ProductID int `json:"productId"`
	// UserID 将从 JWT Token 中获取，而不是由客户端传递
}

func init() {
	initBytes, err := os.ReadFile("scripts/init_stock.lua")
	if err != nil {
		panic(err)
	}
	initStockLua = string(initBytes)

	doBytes, err := os.ReadFile("scripts/do_seckill.lua")
	if err != nil {
		panic(err)
	}
	doSeckillLua = string(doBytes)
	fmt.Println("秒杀 Lua 脚本已成功加载到内存。")
}

// InitSeckillHandler 是处理初始化秒杀商品库存的 HTTP 处理器
func (h *SeckillHandler) InitSeckillHandler(w http.ResponseWriter, r *http.Request) {
	// 仅允许 POST 方法
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	// 1. 解析请求体，获取 productID
	var req InitSeckillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求体格式错误: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.ProductID <= 0 {
		http.Error(w, "无效的 productId", http.StatusBadRequest)
		return
	}

	// 2. 从 MySQL 中查询该商品的真实总库存
	product, err := h.DB.GetProductByID(req.ProductID)
	if err != nil {
		// 如果在数据库中找不到该商品，返回 404 Not Found
		http.Error(w, fmt.Sprintf("商品 ID %d 未找到", req.ProductID), http.StatusNotFound)
		return
	}

	// 3. 将库存信息写入 Redis
	// 我们定义两个 Redis Key:
	// - 一个用于存储库存数量 (String 类型)
	// - 一个用于存储已成功秒杀的用户ID (Set 类型)
	stockKey := fmt.Sprintf("product:stock:%d", req.ProductID)
	userSetKey := fmt.Sprintf("product:users:%d", req.ProductID)
	ctx := context.Background()

	// 创建 Lua 脚本对象
	result, err := h.Redis.Client.Eval(ctx, initStockLua, []string{stockKey, userSetKey}, product.TotalStock).Result()
	if err != nil {
		log.Printf("执行 Redis Lua 脚本失败: %v", err)
		http.Error(w, "服务器内部错误，无法初始化库存", http.StatusInternalServerError)
		return
	}

	// --- 根据 Lua 脚本的返回值进行响应 ---
	w.Header().Set("Content-Type", "application/json")
	resultStr, ok := result.(string)
	if !ok {
		log.Printf("Lua 脚本返回了非预期的类型: %T", result)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}

	switch resultStr {
	case "already_initialized":
		w.WriteHeader(http.StatusOK) // 也可以返回 200 OK，表示操作已知悉
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("商品 %d 已被初始化过，无需重复操作", req.ProductID),
		})
		fmt.Printf("秒杀商品 %d 初始化请求被跳过 (已存在)\n", req.ProductID)

	case "ok":
		w.WriteHeader(http.StatusCreated) // 201 Created 更符合语义
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"message":   fmt.Sprintf("商品 %d 初始化成功", req.ProductID),
			"productId": req.ProductID,
			"stock":     product.TotalStock,
		})
		fmt.Printf("秒杀商品 %d 已成功加载到 Redis，库存: %d\n", req.ProductID, product.TotalStock)
	}
}

// DoSeckillHandler 是处理用户执行秒杀操作的 HTTP 处理器
// 我们先创建一个空的函数框架，下一步再来实现它
// DoSeckillHandler 是处理用户执行秒杀操作的 HTTP 处理器
func (h *SeckillHandler) DoSeckillHandler(w http.ResponseWriter, r *http.Request) {
	// --- 1. 认证和授权 (这是简化版，完整的应该用中间件) ---
	// 客户端应在 Header 中提供 Token: "Authorization: Bearer <your_token>"
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "请求未授权", http.StatusUnauthorized)
		return
	}
	// 去掉 "Bearer " 前缀
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "无效的 Token", http.StatusUnauthorized)
		return
	}
	// 从 Token 中获取 userID
	userID := claims.UserID

	// --- 2. 解析请求体 ---
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}
	var req DoSeckillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求体格式错误: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.ProductID <= 0 {
		http.Error(w, "无效的 productId", http.StatusBadRequest)
		return
	}

	// --- 3. 执行 Redis Lua 脚本进行原子操作 ---
	stockKey := fmt.Sprintf("product:stock:%d", req.ProductID)
	userSetKey := fmt.Sprintf("product:users:%d", req.ProductID)
	ctx := context.Background()

	// 使用 Eval 执行脚本 (v8 兼容)
	result, err := h.Redis.Client.Eval(ctx, doSeckillLua, []string{stockKey, userSetKey}, userID).Result()
	if err != nil {
		log.Printf("执行秒杀 Lua 脚本失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}

	// 将 result (interface{}) 转换为 int64
	resultCode, ok := result.(int64)
	if !ok {
		log.Printf("Lua 脚本返回了非预期的类型: %T", result)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// --- 4. 根据 Lua 返回码处理结果 ---
	var response map[string]interface{}

	switch resultCode {
	case 0:
		// --- 5. 秒杀成功，将订单信息推入异步队列 ---
		// 我们将 userID 和 productID 序列化为 JSON 字符串再推入队列
		orderData := map[string]interface{}{
			"userId":    userID,
			"productId": req.ProductID,
			"timestamp": time.Now().Unix(),
		}
		orderJSON, _ := json.Marshal(orderData)

		// 使用 RPUSH 将订单推入列表尾部
		if err := h.Redis.Client.RPush(ctx, orderQueueKey, orderJSON).Err(); err != nil {
			log.Printf("将订单推入队列失败: %v", err)
			// 这是一个严重问题，用户以为抢到了但订单可能丢失
			// 生产环境需要有重试或补偿机制
			http.Error(w, "服务器内部错误，无法创建订单", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		response = map[string]interface{}{
			"success": true,
			"message": "秒杀成功！订单正在处理中。",
		}

	case 1:
		w.WriteHeader(http.StatusOK) // 业务上是正常的，所以返回 200
		response = map[string]interface{}{
			"success": false,
			"message": "商品已售罄！",
		}
	case 2:
		w.WriteHeader(http.StatusOK)
		response = map[string]interface{}{
			"success": false,
			"message": "您已抢购过此商品，请勿重复抢购！",
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		response = map[string]interface{}{
			"success": false,
			"message": "未知的服务器错误",
		}
	}

	json.NewEncoder(w).Encode(response)
}
