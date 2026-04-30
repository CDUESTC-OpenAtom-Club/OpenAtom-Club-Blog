package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 用户结构体（用于注册/登录参数）
type User struct {
	Username string `json:"username" binding:"required,min=3,max=20"` // 用户名：必填，3-20位
	Password string `json:"password" binding:"required,min=6"`        // 密码：必填，至少6位
}

// 模拟数据库（实际项目用MySQL+GORM）
var userDB = make(map[string]string) // key:用户名, value:加密后的密码

func main() {
	r := gin.Default()

	// 路由：注册和登录
	r.POST("/register", registerHandler) // 注册接口
	r.POST("/login", loginHandler)       // 登录接口

	r.Run(":8080") // 启动服务
}

// 注册接口逻辑
func registerHandler(c *gin.Context) {
	var req User
	// 校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误：" + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	if _, exists := userDB[req.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{
			"code": 409,
			"msg":  "用户名已被注册",
		})
		return
	}

	// 密码加密（bcrypt）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	// 存入“数据库”
	userDB[req.Username] = string(hashedPassword)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

// 登录接口逻辑
func loginHandler(c *gin.Context) {
	var req User
	// 校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误：" + err.Error(),
		})
		return
	}

	// 检查用户是否存在
	hashedPassword, exists := userDB[req.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 登录成功（实际项目会返回token，这里简化）
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{"username": req.Username},
	})
}
