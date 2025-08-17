package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sql_operation/db"
)

type UserHandler struct {
	DB *db.Database
}

type RegisterRequest struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type RegisterResponse struct {
	UserId  int64  `json:"userid"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type LoginRequest struct {
	UserId   int64  `json:"userid"`
	PassWord string `json:"password"`
}

type LoginResponse struct {
	UserName string `json:"username,omitempty"`
	Message  string `json:"message"`
	Success  bool   `json:"success"`
}

type ResetNameRequest struct {
	UserId  int64  `json:"userid"`
	NewName string `json:"newname"`
}

type ResetNameResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ResetPasswordRequest struct {
	UserId      int64  `json:"userid"`
	NewPassword string `json:"newpassword"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (h *UserHandler) RegisterHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "请求体格式错误", http.StatusBadRequest)
	}

	userid, ok, err := h.DB.USER_Register(req.UserName, req.PassWord)
	if err != nil {
		// 如果 err 不是 nil，说明是数据库查询等系统内部错误
		// 我们应该记录这个错误，并返回一个通用的服务器错误信息
		log.Printf("数据库查询失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json") // 设置响应头为 JSON
	var res RegisterResponse
	if ok {
		w.WriteHeader(http.StatusOK)
		res = RegisterResponse{
			UserId:  userid,
			Message: "用户注册成功",
			Success: true,
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		res = RegisterResponse{
			Message: "用户注册失败",
			Success: false,
		}
	}
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) LoignHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "请求体格式错误", http.StatusBadRequest)
	}

	username, ok, err := h.DB.USER_Login(req.UserId, req.PassWord)
	if err != nil {
		// 如果 err 不是 nil，说明是数据库查询等系统内部错误
		// 我们应该记录这个错误，并返回一个通用的服务器错误信息
		log.Printf("数据库查询失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json") // 设置响应头为 JSON
	var res LoginResponse
	if ok {
		w.WriteHeader(http.StatusOK)
		res = LoginResponse{
			UserName: username,
			Message:  "登陆认证成功",
			Success:  true,
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		res = LoginResponse{
			Message: "登陆认证失败",
			Success: false,
		}
	}
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) ResetNameHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var req ResetNameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "请求体格式错误", http.StatusBadRequest)
	}

	ok, err := h.DB.USER_ResetName(req.UserId, req.NewName)
	if err != nil {
		// 如果 err 不是 nil，说明是数据库查询等系统内部错误
		// 我们应该记录这个错误，并返回一个通用的服务器错误信息
		log.Printf("数据库查询失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json") // 设置响应头为 JSON
	var res ResetNameResponse
	if ok {
		w.WriteHeader(http.StatusOK)
		res = ResetNameResponse{
			Message: "昵称修改成功",
			Success: true,
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		res = ResetNameResponse{
			Message: "昵称修改失败",
			Success: false,
		}
	}
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) ResetPasswordHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var req ResetPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "请求体格式错误", http.StatusBadRequest)
	}

	ok, err := h.DB.USER_ResetPassword(req.UserId, req.NewPassword)
	if err != nil {
		// 如果 err 不是 nil，说明是数据库查询等系统内部错误
		// 我们应该记录这个错误，并返回一个通用的服务器错误信息
		log.Printf("数据库查询失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json") // 设置响应头为 JSON
	var res ResetPasswordResponse
	if ok {
		w.WriteHeader(http.StatusOK)
		res = ResetPasswordResponse{
			Message: "密码修改成功",
			Success: true,
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		res = ResetPasswordResponse{
			Message: "密码修改失败",
			Success: false,
		}
	}
	json.NewEncoder(w).Encode(res)
}
