package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"DistributedCalculator/http/orchestrator"
	"DistributedCalculator/pkg/ExpParser"
	"DistributedCalculator/pkg/database"
	"DistributedCalculator/pkg/env"
	"DistributedCalculator/pkg/expression"
	"DistributedCalculator/pkg/logger"

	"go.uber.org/zap"
)
var Log = logger.SetupLogger()

type expTempl struct {
	Items []string
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Log.Info("HTTP request", zap.String("method", r.Method), zap.String("path", r.URL.Path))
		f, err := os.Open("frontend/main.html")
		if err != nil {
			Log.Error("Reading", zap.String("ERROR", "reading html error"))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		body, _ := io.ReadAll(f)
		fmt.Fprint(w, string(body))
	}
	if r.Method == http.MethodPost {
		Log.Info("HTTP request", zap.String("method", r.Method), zap.String("path", r.URL.Path))

		body, _ := io.ReadAll(r.Body)
		expr := expression.NewExpression(string(body))

		database.WriteExp(*expr, Log)
		node, err := ExpParser.ParseExpr(expr.Name)
		if err != nil {
			expr.Status = 3
			database.UpdateExpr(*expr)
			return
		}
		expr.Node = *node
		b, _ := json.Marshal(expr)
		rb := bytes.NewReader(b)
		http.Post("http://localhost:8081", "application/json", rb)
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Log.Info("HTTP request", zap.String("method", r.Method), zap.String("path", r.URL.Path))
		tmpl := template.Must(template.ParseFiles("frontend/expressions.html"))
		exprs, err := database.GetAllExpressions(Log)
		if err != nil {
			Log.Error("Error", zap.String("Unexpected", "unexpected error"))
		}
		line := expTempl{}

		arr := make([]string, 0)
		for _, i := range exprs {
			arr = append(arr, i.ForTemplate())
		}
		line.Items = arr
		tmpl.Execute(w, line)

	}
}

func main() {
	database.DeleteAll()
	orchestrator.StartServer(Log)
	env.SetupEnv()
	filepath.Abs("/")
	mux := http.NewServeMux()
	mux.HandleFunc("/expressions", resultHandler)
	mux.HandleFunc("/", calculateHandler)
	Log.Info("ListenAndServe", zap.String("port", "8080"))
	fmt.Println("Server is running on http://localhost:8080	")
	http.ListenAndServe(":8080", mux)
}
