package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"disCom/internal/database"
	"disCom/internal/env"
	constants "disCom/internal/env"
	"disCom/internal/expression"
	"disCom/internal/logger"
	"disCom/internal/orchestrator"
	"disCom/internal/parser"

	"go.uber.org/zap"
)

type expTempl struct {
	Items []string
}

type jsonSet struct {
	Plus  string `json:"plus"`
	Minus string `json:"minus"`
	Mul   string `json:"mul"`
	Div   string `json:"div"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		f, err := os.Open("frontend/main.html")
		if err != nil {
			zap.String("ERROR", "reading html error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		body, err := io.ReadAll(f)
		fmt.Fprintf(w, string(body))
	}
	if r.Method == http.MethodPost {
		body, _ := io.ReadAll(r.Body)
		expr := expression.NewExpression(string(body))

		database.WriteExp(*expr)
		node, err := parser.ParseExpr(expr.Name)
		if err != nil {
			expr.Status = 3
			database.UpdateExpr(*expr)
			return
		}
		expr.Node = *node
		b, _ := json.Marshal(expr)
		rb := bytes.NewReader(b)
		http.Post("http://localhost:8081", "application/json", rb)
		//req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))\
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("frontend/expressions.html"))
		exprs, err := database.GetAllExpressions()
		if err != nil {
			zap.String("ERROR", "unexpected error")
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

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("frontend/settings.html"))
		f := struct {
			Plus  int
			Minus int
			Mul   int
			Div   int
		}{
			Plus:  env.Plus,
			Minus: env.Minus,
			Mul:   env.Mul,
			Div:   env.Div,
		}
		tmpl.Execute(w, f)
	} else if r.Method == http.MethodPost {
		var set jsonSet
		var plus, minus, mul, div int
		body, err := io.ReadAll(r.Body)
		if err != nil {
			zap.String("ERROR", "reading error")
		}
		err = json.Unmarshal(body, &set)
		if err != nil {
			zap.String("ERROR", "reading jsondb error")
		}
		plus, _ = strconv.Atoi(set.Plus)
		constants.Plus = plus
		
		minus, _ = strconv.Atoi(set.Minus)
		constants.Minus = minus

		mul, _ = strconv.Atoi(set.Mul)
		constants.Mul = mul

		div, _ = strconv.Atoi(set.Div)
		constants.Div = div

		env.Save()

	}
}


func main() {
	logger.SetupLogger()
	database.DeleteAll()
	orchestrator.StartServer()
	env.Init()
	filepath.Abs("/")
	mux := http.NewServeMux()
	mux.HandleFunc("/expressions", resultHandler)
	mux.HandleFunc("/", calculateHandler)
	fmt.Println("Server is running on http://localhost:8080	")
	http.ListenAndServe(":8080", mux)
}