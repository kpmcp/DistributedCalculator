package database

import (
	"DistributedCalculator/pkg/expression"
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var mu sync.Mutex

func WriteExp(exp expression.Expression, logger *zap.Logger) {
	mu.Lock()
	defer mu.Unlock()
	db, err := sql.Open("sqlite3", "pkg/database/data.sql")
	if err != nil {
		logger.Error("DB ERROR", zap.String("ERROR", err.Error()))
	}
	db.Exec("INSERT INTO Expressions (Id, Name, Status, Result) values ($1, $2, $3, $4)", exp.Id, exp.Name, exp.Status, exp.Result)

}

func ReadExp(id int, logger *zap.Logger) *expression.Expression {
	mu.Lock()
	defer mu.Unlock()
	db, err := sql.Open("sqlite3", "pkg/database/data.sql")
	if err != nil {
		logger.Error("DB ERROR", zap.String("ERROR", err.Error()))
	}
	defer db.Close()
	row := db.QueryRow("select * from Expressions where Id = $1", id)
	exp := expression.NewExpression("")
	err = row.Scan(&exp.Id, &exp.Name, &exp.Status, &exp.Result)
	if err != nil {
		return nil
	}
	return exp
}


func DeleteAll() {
	db, _ := sql.Open("sqlite3", "pkg/database/data.sql")
	db.Exec("DELETE From Expressions")
	defer db.Close()
}

func UpdateExpr(expr expression.Expression) {
	mu.Lock()
	defer mu.Unlock()
	db, _ := sql.Open("sqlite3", "pkg/database/data.sql")
	if expr.Status == 0 {
		db.Exec("update Expressions set status = $1 where id = $2", expr.Status, expr.Id)
		db.Exec("update Expressions set result = $1 where id = $2", expr.Result, expr.Id)
	} else {
		db.Exec("update Expressions set status = $1 where id = $2", expr.Status, expr.Id)
	}

	defer db.Close()

}

func GetAllExpressions(logger *zap.Logger) ([]*expression.Expression, error) {
	all := make([]*expression.Expression, 0)
	db, err := sql.Open("sqlite3", "pkg/database/data.sql")
	if err != nil {
		logger.Error("DB ERROR", zap.String("ERROR", err.Error()))
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Expressions")
	if err != nil {
		logger.Error("DB ERROR", zap.String("ERROR", err.Error()))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		expr := expression.NewExpression("")
		err = rows.Scan(&expr.Id, &expr.Name, &expr.Status, &expr.Result)
		if err != nil {
			logger.Error("DB ERROR", zap.String("ERROR", err.Error()))
			return nil, err
		}
		if expr.Name != "added" {
			all = append(all, expr)
		}
	}
	return all, nil
}
