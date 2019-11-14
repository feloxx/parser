package main

import (
	"fmt"
	"github.com/pingcap/parser"
	. "github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

func whereIterator(stmt ExprNode) {
	switch w := stmt.(type) {
	case *BinaryOperationExpr:
		switch l := w.L.(type) {
		case *BinaryOperationExpr:
			if col, ok := l.L.(*ColumnNameExpr); ok {
				if val, ok := l.R.(ValueExpr); ok {
					if col.Name.Name.L == "query_time" || col.Name.Name.L == "user" {
						fmt.Println(col.Name.Name.L)
						fmt.Println(val.GetString())
						fmt.Println(l.Op.String())
					}
				}
			} else {
				whereIterator(l)
			}
		}

		switch r := w.R.(type) {
		case *BinaryOperationExpr:
			if col, ok := r.L.(*ColumnNameExpr); ok {
				if val, ok := r.R.(ValueExpr); ok {
					if col.Name.Name.L == "query_time" || col.Name.Name.L == "user" {
						fmt.Println(col.Name.Name.L)
						fmt.Println(val.GetString())
						fmt.Println(r.Op.String())
					}
				}
			} else {
				whereIterator(r)
			}
		}
	}
}

func main() {
	sql := "select * from t where t.query_time >= '2019-01-01 00:00:00' and t.query_time < '2019-01-01 23:59:59' and user = 'ficl' and t.aa = '1';"

	stmt, err := parser.New().ParseOneStmt(sql, "", "")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", stmt)

	whereIterator(stmt.(*SelectStmt).Where)
}
