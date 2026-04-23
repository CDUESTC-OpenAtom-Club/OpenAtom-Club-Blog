package main

import "github.com/gin-gonic/gin"

func main() {
    // 创建Gin引擎
    r := gin.Default()
    // 定义一个GET接口，访问根路径返回"Hello Gin!"
    r.GET("/", func(c *gin.Context) {
        c.String(200, "Hello Gin!")
    })
    // 启动服务，监听8080端口
    r.Run(":8080")
}