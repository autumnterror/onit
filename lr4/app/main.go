package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = "Нода не определена"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
 	<meta charset="UTF-8">
    <title>Балансировка Nginx</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            text-align: center;
            margin-top: 100px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
        }
        .node-card {
            background: rgba(255,255,255,0.2);
            border-radius: 20px;
            padding: 40px;
            display: inline-block;
            backdrop-filter: blur(10px);
            box-shadow: 0 8px 32px rgba(0,0,0,0.1);
        }
        h1 {
            font-size: 48px;
            margin-bottom: 10px;
        }
        .node-id {
            font-size: 72px;
            font-weight: bold;
            color: #ffd700;
            margin: 20px 0;
        }
        .ip {
            font-size: 18px;
            opacity: 0.9;
        }
    </style>
</head>
<body>
    <div class="node-card">
        <h1>Балансировка нагрузки</h1>
        <div class="node-id">%s</div>
        <p>Ответ от сервера</p>
        <div class="ip">Ваш IP: %s</div>
    </div>
</body>
</html>
`, nodeID, r.RemoteAddr)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK "+nodeID)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Node %s starting on port %s", nodeID, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
