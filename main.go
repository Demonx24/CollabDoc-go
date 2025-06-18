package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type CallbackRequest struct {
	Status int `json:"status"`
	// 你可以根据需求解析其他字段
}

type CallbackResponse struct {
	Error int `json:"error"`
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	var req CallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 这里可以添加保存文件等业务逻辑
	// 例如根据req.Status判断是否成功

	resp := CallbackResponse{Error: 0}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func main() {
	// 1. 静态资源（.docx 文件）服务
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("./file"))))

	// 2. 回调接口（只打印信息）
	http.HandleFunc("/onlyoffice/callback", callbackHandler)

	// 3. 编辑器页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/editor.html")
		if err != nil {
			http.Error(w, "Template error", 500)
			return
		}

		config := map[string]interface{}{
			"document": map[string]interface{}{
				"fileType": "docx",
				"key":      "example-doc-001",
				"title":    "示例文档.docx",
				"url":      "http://localhost:8080/file/example.docx",
			},
			"editorConfig": map[string]interface{}{
				"callbackUrl": "http://host.docker.internal:8080/onlyoffice/callback",
				"user": map[string]string{
					"id":   "u1",
					"name": "张三",
				},
			},
		}

		tmpl.Execute(w, config)
	})

	log.Println("Go server running at http://localhost:8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}
