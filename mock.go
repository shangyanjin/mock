package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 获取命令行参数或使用默认端口 3000
	port := "3000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	// 遍历目录并找到所有的 JSON 文件
	var files []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			normalizedPath := strings.ReplaceAll(path, "\\", "/")
			files = append(files, normalizedPath)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error finding JSON files:", err)
		return
	}

	// 显示所有 JSON 文件及其路由
	for _, file := range files {
		route := "/" + strings.TrimPrefix(file, "./")
		route = strings.TrimSuffix(route, ".json")
		fmt.Printf("Route: %s -> File: ./%s\n", route, file)
	}

	// 通用路由处理
	r.Any("/*path", func(c *gin.Context) {
		path := c.Param("path")
		localFilePath := "./" + path + ".json"

		if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		content, err := ioutil.ReadFile(localFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
			return
		}

		c.Data(http.StatusOK, "application/json", content)
	})

	// 启动服务器，端口由命令行参数指定，若无则默认为 3000
	fmt.Printf("Starting server on port %s\n", port)
	r.Run(":" + port)
}
