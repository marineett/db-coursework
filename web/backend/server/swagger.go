package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

func SwaggerUIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("server/files/index.html")
		if err != nil {
			http.Error(w, "Swagger UI not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func SwaggerSpecHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("server/files/ccc_fixed.yaml")
		if err != nil {
			http.Error(w, "OpenAPI spec not found", http.StatusInternalServerError)
			return
		}

		// Проверяем, запрашивается ли JSON версия по расширению файла или заголовку Accept
		urlPath := r.URL.Path
		acceptHeader := r.Header.Get("Accept")
		shouldReturnJSON := strings.HasSuffix(urlPath, ".json") || strings.Contains(acceptHeader, "application/json")

		if shouldReturnJSON {
			// Конвертируем YAML в JSON
			var yamlData interface{}
			if err := yaml.Unmarshal(data, &yamlData); err != nil {
				log.Printf("ERROR: Failed to parse YAML: %v", err)
				http.Error(w, "Failed to parse YAML", http.StatusInternalServerError)
				return
			}
			jsonData, err := json.Marshal(yamlData)
			if err != nil {
				log.Printf("ERROR: Failed to convert to JSON: %v", err)
				http.Error(w, "Failed to convert to JSON", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
			return
		}

		// Отдаем YAML
		if strings.Contains(acceptHeader, "application/yaml") || strings.Contains(acceptHeader, "application/x-yaml") {
			w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func DocumentationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		readmePath := filepath.Join("server", "files", "README.md")
		data, err := os.ReadFile(readmePath)
		if err != nil {
			http.Error(w, "Documentation not found", http.StatusInternalServerError)
			return
		}

		md := data

		imagePathRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
		md = imagePathRegex.ReplaceAllFunc(md, func(match []byte) []byte {
			matches := imagePathRegex.FindSubmatch(match)
			if len(matches) >= 3 {
				altText := string(matches[1])
				imagePath := string(matches[2])

				filename := filepath.Base(imagePath)
				apiPath := "/api/v2/static/" + filename
				return []byte("![" + altText + "](" + apiPath + ")")
			}
			return match
		})

		extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.SuperSubscript | parser.Tables
		p := parser.NewWithExtensions(extensions)
		doc := p.Parse(md)

		htmlFlags := html.CommonFlags | html.HrefTargetBlank
		opts := html.RendererOptions{Flags: htmlFlags}
		renderer := html.NewRenderer(opts)
		htmlBytes := markdown.Render(doc, renderer)

		htmlDoc := createGitHubStyledHTML(string(htmlBytes))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlDoc))
	}
}

func createGitHubStyledHTML(bodyContent string) string {
	return `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Документация проекта</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
            line-height: 1.6;
            color: #24292e;
            background-color: #ffffff;
            max-width: 1012px;
            margin: 0 auto;
            padding: 32px;
        }
        h1, h2, h3, h4, h5, h6 {
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }
        h1 {
            font-size: 2em;
            border-bottom: 1px solid #eaecef;
            padding-bottom: 0.3em;
        }
        h2 {
            font-size: 1.5em;
            border-bottom: 1px solid #eaecef;
            padding-bottom: 0.3em;
        }
        h3 {
            font-size: 1.25em;
        }
        p {
            margin-bottom: 16px;
        }
        ul, ol {
            margin-bottom: 16px;
            padding-left: 2em;
        }
        li {
            margin-bottom: 0.25em;
        }
        code {
            padding: 0.2em 0.4em;
            margin: 0;
            font-size: 85%;
            background-color: rgba(27, 31, 35, 0.05);
            border-radius: 3px;
            font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace;
        }
        pre {
            padding: 16px;
            overflow: auto;
            font-size: 85%;
            line-height: 1.45;
            background-color: #f6f8fa;
            border-radius: 6px;
            margin-bottom: 16px;
        }
        pre code {
            display: inline;
            max-width: auto;
            padding: 0;
            margin: 0;
            overflow: visible;
            line-height: inherit;
            word-wrap: normal;
            background-color: transparent;
            border: 0;
        }
        blockquote {
            padding: 0 1em;
            color: #6a737d;
            border-left: 0.25em solid #dfe2e5;
            margin: 0 0 16px 0;
        }
        table {
            border-spacing: 0;
            border-collapse: collapse;
            margin-bottom: 16px;
            width: 100%;
        }
        table th, table td {
            padding: 6px 13px;
            border: 1px solid #dfe2e5;
        }
        table th {
            font-weight: 600;
            background-color: #f6f8fa;
        }
        img {
            max-width: 100%;
            box-sizing: content-box;
            background-color: #fff;
            border-style: none;
        }
        a {
            color: #0366d6;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
` + bodyContent + `
</body>
</html>`
}

func ApiV2Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("server/files/ccc_fixed.yaml")
		if err != nil {
			http.Error(w, "OpenAPI spec not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func StaticFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := vars["filename"]
		if filename == "" {
			filename = "BPMN.png"
		}

		filePath := filepath.Join("server", "files", "static", filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			filePath = filepath.Join("server", "files", "static", "png", filename)
			var err2 error
			data, err2 = os.ReadFile(filePath)
			if err2 != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
		}

		var ctype string
		if strings.HasSuffix(filename, ".css") {
			ctype = "text/css; charset=utf-8"
		} else if strings.HasSuffix(filename, ".js") {
			ctype = "application/javascript; charset=utf-8"
		} else if strings.HasSuffix(filename, ".json") {
			ctype = "application/json; charset=utf-8"
		} else if strings.HasSuffix(filename, ".png") {
			ctype = "image/png"
		} else if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
			ctype = "image/jpeg"
		} else if strings.HasSuffix(filename, ".svg") {
			ctype = "image/svg+xml"
		} else if strings.HasSuffix(filename, ".ico") {
			ctype = "image/x-icon"
		} else {
			ctype = http.DetectContentType(data)
		}

		w.Header().Set("Content-Type", ctype)
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func SwaggerUIStaticHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := vars["filename"]

		// Если filename не указан в vars, пытаемся извлечь из URL
		if filename == "" {
			// Извлекаем имя файла из пути URL
			// Например: /api/v2/static/swagger-ui/swagger-ui.css -> swagger-ui.css
			path := r.URL.Path
			prefix := "/api/v2/static/swagger-ui/"
			if strings.HasPrefix(path, prefix) {
				filename = strings.TrimPrefix(path, prefix)
				// Убираем начальный слэш, если есть
				filename = strings.TrimPrefix(filename, "/")
			}
		}

		if filename == "" {
			// Если filename все еще не указан, возвращаем ошибку
			http.Error(w, "Filename required", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join("server", "files", "swagger-ui", filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("ERROR: Failed to read swagger-ui file %s: %v", filePath, err)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Определяем Content-Type на основе расширения файла
		var ctype string
		if strings.HasSuffix(filename, ".js") {
			ctype = "application/javascript; charset=utf-8"
		} else if strings.HasSuffix(filename, ".css") {
			ctype = "text/css; charset=utf-8"
		} else if strings.HasSuffix(filename, ".json") {
			ctype = "application/json; charset=utf-8"
		} else if strings.HasSuffix(filename, ".png") {
			ctype = "image/png"
		} else if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
			ctype = "image/jpeg"
		} else if strings.HasSuffix(filename, ".svg") {
			ctype = "image/svg+xml"
		} else {
			ctype = http.DetectContentType(data)
		}

		w.Header().Set("Content-Type", ctype)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func ReservedStaticFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := vars["filename"]
		if filename == "" {
			// Если filename не указан, возвращаем список файлов или ошибку
			http.Error(w, "Filename required", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join("server", "files", "reserved", filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		ctype := http.DetectContentType(data)
		w.Header().Set("Content-Type", ctype)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// FaviconHandler обрабатывает запросы к favicon.ico
func FaviconHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Пробуем найти favicon.ico в static папке
		filePath := filepath.Join("server", "files", "static", "favicon.ico")
		data, err := os.ReadFile(filePath)
		if err != nil {
			// Если favicon не найден, возвращаем 204 No Content (стандартная практика)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func ManagementHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем список PNG файлов из static/png
		pngFiles := []string{}
		pngDir := filepath.Join("server", "files", "static", "png")
		if entries, err := os.ReadDir(pngDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
					pngFiles = append(pngFiles, entry.Name())
				}
			}
		}

		var html strings.Builder
		html.WriteString("<!DOCTYPE html>\n<html>\n<head><meta charset=\"UTF-8\"><title>API Management</title></head>\n<body>\n")

		// Ручки из setup.go (строки 22-36)
		endpoints := []struct {
			path string
			desc string
		}{
			{"/api/v2/management", "API Management"},
			{"/api/v2/openapi.yaml", "OpenAPI YAML"},
			{"/api/v2/openapi.json", "OpenAPI JSON"},
			{"/api/v2/static/openapi.yaml", "OpenAPI YAML (альтернативный путь)"},
			{"/api/v2/static/documentation", "README документация"},
			{"/api/v2/legacy/", "Legacy архив"},
		}

		html.WriteString("<ul>\n")
		for _, ep := range endpoints {
			html.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a> - %s</li>\n", ep.path, ep.path, ep.desc))
		}
		html.WriteString("</ul>\n")

		// Ручки с картинками
		if len(pngFiles) > 0 {
			html.WriteString("<h2>Картинки</h2>\n<ul>\n")
			for _, png := range pngFiles {
				imgPath := "/api/v2/static/" + png
				html.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s</a> <img src=\"%s\" alt=\"%s\"></li>\n", imgPath, imgPath, imgPath, png))
			}
			html.WriteString("</ul>\n")
		}

		html.WriteString("</body>\n</html>")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html.String()))
	}
}

// StatusHandler возвращает HTML страницу со статусом сервера
func StatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Читаем HTML файл
		htmlPath := filepath.Join("server", "files", "status.html")
		htmlData, err := os.ReadFile(htmlPath)
		if err != nil {
			http.Error(w, "Status page not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(htmlData)
	}
}

// StatusDataHandler возвращает данные статуса в формате JSON
// Получает статус от Nginx через HTTP запрос
func StatusDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Пытаемся получить статус от Nginx
		// Если Nginx доступен через переменную окружения или по умолчанию
		nginxStatusURL := os.Getenv("NGINX_STATUS_URL")
		if nginxStatusURL == "" {
			// Пробуем получить статус локально или через прокси
			nginxStatusURL = "http://nginx/status"
		}

		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get(nginxStatusURL)
		if err != nil {
			// Если не удалось получить статус от Nginx, возвращаем ошибку
			http.Error(w, "Failed to fetch status from Nginx", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Nginx status endpoint returned error", http.StatusServiceUnavailable)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read status response", http.StatusInternalServerError)
			return
		}

		// Парсим статус Nginx и конвертируем в JSON
		statusData := parseNginxStatus(string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(statusData)
	}
}

// parseNginxStatus парсит текст статуса Nginx и возвращает структуру
func parseNginxStatus(text string) map[string]interface{} {
	result := make(map[string]interface{})
	lines := strings.Split(text, "\n")

	// Компилируем регулярные выражения один раз для производительности
	numbersRe := regexp.MustCompile(`^\s*\d+\s+\d+\s+\d+`)
	readingRe := regexp.MustCompile(`Reading:\s*(\d+)\s+Writing:\s*(\d+)\s+Waiting:\s*(\d+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Active connections: 1
		if strings.HasPrefix(line, "Active connections:") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				result["active_connections"] = parts[2]
			}
		}

		// server accepts handled requests
		//  123 456 789
		if strings.Contains(line, "server accepts handled requests") {
			// Следующая строка содержит числа
			continue
		}

		// Проверяем строку с числами (принятые, обработанные, запросы)
		if numbersRe.MatchString(line) {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				result["accepted"] = parts[0]
				result["handled"] = parts[1]
				result["requests"] = parts[2]
			}
		}

		// Reading: X Writing: Y Waiting: Z
		if strings.HasPrefix(line, "Reading:") {
			matches := readingRe.FindStringSubmatch(line)
			if len(matches) == 4 {
				result["reading"] = matches[1]
				result["writing"] = matches[2]
				result["waiting"] = matches[3]
			}
		}
	}

	return result
}

// // WebAdminHandler возвращает HTML страницу веб-админки
// func WebAdminHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		htmlPath := filepath.Join("server", "files", "web_admin.html")
// 		htmlData, err := os.ReadFile(htmlPath)
// 		if err != nil {
// 			http.Error(w, "Web admin page not found", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(htmlData)
// 	}
// }

// // WebAdminTablesHandler возвращает список таблиц в БД
// func WebAdminTablesHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		query := `
// 			SELECT table_name
// 			FROM information_schema.tables
// 			WHERE table_schema = 'public'
// 			ORDER BY table_name
// 		`

// 		rows, err := db.Query(query)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Error querying tables: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		defer rows.Close()

// 		var tables []string
// 		for rows.Next() {
// 			var tableName string
// 			if err := rows.Scan(&tableName); err != nil {
// 				http.Error(w, fmt.Sprintf("Error scanning table name: %v", err), http.StatusInternalServerError)
// 				return
// 			}
// 			tables = append(tables, tableName)
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"tables": tables,
// 		})
// 	}
// }

// WebAdminTableHandler возвращает данные из таблицы с пагинацией
// func WebAdminTableHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		tableName := vars["tableName"]

// 		if tableName == "" {
// 			http.Error(w, "Table name required", http.StatusBadRequest)
// 			return
// 		}

// 		// Безопасность: проверяем, что имя таблицы содержит только допустимые символы
// 		if !regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(tableName) {
// 			http.Error(w, "Invalid table name", http.StatusBadRequest)
// 			return
// 		}

// 		// Получаем параметры пагинации
// 		limit := r.URL.Query().Get("limit")
// 		offset := r.URL.Query().Get("offset")
// 		if limit == "" {
// 			limit = "50"
// 		}
// 		if offset == "" {
// 			offset = "0"
// 		}

// 		// Получаем колонки таблицы
// 		columnsQuery := `
// 			SELECT column_name
// 			FROM information_schema.columns
// 			WHERE table_name = $1
// 			ORDER BY ordinal_position
// 		`

// 		rows, err := db.Query(columnsQuery, tableName)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Error querying columns: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		defer rows.Close()

// 		var columns []string
// 		for rows.Next() {
// 			var colName string
// 			if err := rows.Scan(&colName); err != nil {
// 				http.Error(w, fmt.Sprintf("Error scanning column: %v", err), http.StatusInternalServerError)
// 				return
// 			}
// 			columns = append(columns, colName)
// 		}

// 		if len(columns) == 0 {
// 			http.Error(w, "Table not found or has no columns", http.StatusNotFound)
// 			return
// 		}

// 		// Получаем общее количество строк
// 		countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, tableName)
// 		var total int
// 		err = db.QueryRow(countQuery).Scan(&total)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Error counting rows: %v", err), http.StatusInternalServerError)
// 			return
// 		}

// 		// Получаем данные с пагинацией
// 		// Строим SELECT запрос безопасно
// 		selectQuery := fmt.Sprintf(`SELECT %s FROM %s LIMIT $1 OFFSET $2`,
// 			strings.Join(columns, ", "), tableName)

// 		dataRows, err := db.Query(selectQuery, limit, offset)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Error querying data: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		defer dataRows.Close()

// 		// Преобразуем строки в массив объектов
// 		var resultRows []map[string]interface{}
// 		for dataRows.Next() {
// 			// Создаем slice для значений
// 			values := make([]interface{}, len(columns))
// 			valuePtrs := make([]interface{}, len(columns))
// 			for i := range values {
// 				valuePtrs[i] = &values[i]
// 			}

// 			if err := dataRows.Scan(valuePtrs...); err != nil {
// 				http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
// 				return
// 			}

// 			row := make(map[string]interface{})
// 			for i, col := range columns {
// 				val := values[i]
// 				// Преобразуем []byte в string для JSON
// 				if b, ok := val.([]byte); ok {
// 					val = string(b)
// 				}
// 				row[col] = val
// 			}
// 			resultRows = append(resultRows, row)
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"columns": columns,
// 			"rows":    resultRows,
// 			"total":   total,
// 		})
// 	}
// }

// // WebAdminQueryHandler выполняет SQL запрос (только SELECT)
// func WebAdminQueryHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != "POST" {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		var req struct {
// 			Query string `json:"query"`
// 		}

// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			http.Error(w, fmt.Sprintf("Error decoding request: %v", err), http.StatusBadRequest)
// 			return
// 		}

// 		query := strings.TrimSpace(req.Query)
// 		if query == "" {
// 			http.Error(w, "Query is required", http.StatusBadRequest)
// 			return
// 		}

// 		// Безопасность: разрешаем только SELECT запросы
// 		upperQuery := strings.ToUpper(strings.TrimSpace(query))
// 		if !strings.HasPrefix(upperQuery, "SELECT") {
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(map[string]interface{}{
// 				"error": "Only SELECT queries are allowed",
// 			})
// 			return
// 		}

// 		// Выполняем запрос
// 		rows, err := db.Query(query)
// 		if err != nil {
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(map[string]interface{}{
// 				"error": err.Error(),
// 			})
// 			return
// 		}
// 		defer rows.Close()

// 		// Получаем колонки
// 		columns, err := rows.Columns()
// 		if err != nil {
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(map[string]interface{}{
// 				"error": fmt.Sprintf("Error getting columns: %v", err),
// 			})
// 			return
// 		}

// 		// Читаем данные
// 		var resultRows []map[string]interface{}
// 		for rows.Next() {
// 			values := make([]interface{}, len(columns))
// 			valuePtrs := make([]interface{}, len(columns))
// 			for i := range values {
// 				valuePtrs[i] = &values[i]
// 			}

// 			if err := rows.Scan(valuePtrs...); err != nil {
// 				w.Header().Set("Content-Type", "application/json")
// 				json.NewEncoder(w).Encode(map[string]interface{}{
// 					"error": fmt.Sprintf("Error scanning row: %v", err),
// 				})
// 				return
// 			}

// 			row := make(map[string]interface{})
// 			for i, col := range columns {
// 				val := values[i]
// 				if b, ok := val.([]byte); ok {
// 					val = string(b)
// 				}
// 				row[col] = val
// 			}
// 			resultRows = append(resultRows, row)
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"columns": columns,
// 			"rows":    resultRows,
// 		})
// 	}
// }
