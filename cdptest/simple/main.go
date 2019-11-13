package main

import (
	"fmt"
	"github.com/pingcap/parser"
	. "github.com/pingcap/parser/ast"
	. "github.com/pingcap/parser/format"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"reflect"
	"strings"
)

type visitor struct {
}

func (visitor) Enter(n Node) (node Node, skipChildren bool) {
	fmt.Printf("%#v\n", n)
	switch c := n.(type) {
	case *ColumnDef:
		fmt.Printf("%#v\n", c.Tp)
	}
	return n, false
}

func (visitor) Leave(n Node) (node Node, ok bool) {
	return n, true
}

// For test only.
func CleanNodeText(node Node) {
	var cleaner nodeTextCleaner
	node.Accept(&cleaner)
}

// nodeTextCleaner clean the text of a node and it's child node.
// For test only.
type nodeTextCleaner struct {
}

// Enter implements Visitor interface.
func (checker *nodeTextCleaner) Enter(in Node) (out Node, skipChildren bool) {
	in.SetText("")
	switch node := in.(type) {
	case *CreateTableStmt:
		for _, opt := range node.Options {
			switch opt.Tp {
			case TableOptionCharset:
				opt.StrValue = strings.ToUpper(opt.StrValue)
			case TableOptionCollate:
				opt.StrValue = strings.ToUpper(opt.StrValue)
			}
		}
		for _, col := range node.Cols {
			col.Tp.Charset = strings.ToUpper(col.Tp.Charset)
			col.Tp.Collate = strings.ToUpper(col.Tp.Collate)
		}
	case *Constraint:
		if node.Option != nil {
			if node.Option.KeyBlockSize == 0x0 && node.Option.Tp == 0 && node.Option.Comment == "" {
				node.Option = nil
			}
		}
	case *AggregateFuncExpr:
		node.F = strings.ToLower(node.F)
	case *AlterTableSpec:
		for _, opt := range node.Options {
			opt.StrValue = strings.ToLower(opt.StrValue)
		}
	}
	return in, false
}

// Leave implements Visitor interface.
func (checker *nodeTextCleaner) Leave(in Node) (out Node, ok bool) {
	return in, true
}
func main() {

	//sql1 := "CREATE TABLE silver_deduct_repay (update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近操作时间')"

	//sql2 := "CREATE TABLE `t1` (`accout_id` INT(11) DEFAULT '0',`summoner_id` INT(11) DEFAULT '0',`union_name` VARBINARY(52) NOT NULL,`union_id` INT(11) DEFAULT '0',PRIMARY KEY(`union_name`)) ENGINE = MyISAM DEFAULT CHARACTER SET = BINARY;"

	//sql3 := `CREATE DATABASE "";`

	//sql3 := "GRANT ALL PRIVILEGES ON `utcs_fansou`.* TO 'utcs_fansou'@'%'"
	sql3 := "select * from t where t.a >= '2019-01-01 00:00:00' and t.a < '2019-01-01 23:59:59' and user = 'ficl';"

	parser := parser.New()

	stmt, err := parser.ParseOneStmt(sql3, "", "")
	if err != nil {
		fmt.Println(err)
	}

	stmt1 := stmt.(*SelectStmt).Where
	a := reflect.TypeOf(stmt1).String()
	fmt.Println(a)
	var sb1 strings.Builder

	_ = stmt1.Restore(NewRestoreCtx(DefaultRestoreFlags, &sb1))

	fmt.Println(sb1.String())
	//fmt.Println("after parser")
	//fmt.Printf("%#v\n", stmt)

	//switch n := stmt.(type) {
	//case GrantStmt:
	//	for _, priv := range n.Privs {
	//		for _, col := range priv.Cols {
	//			fmt.Println(col.Schema.O, col.Table.O)
	//		}
	//	}
	//case ast.SelectStmt:

	//}
	//fmt.Printf("%#v\n", s)
	//
	//var sb strings.Builder
	//err = stmt.Restore(NewRestoreCtx(DefaultRestoreFlags, &sb))
	//fmt.Println("sql1", sb.String())
	//
	//fmt.Println("=========")
	//var sbFrom strings.Builder
	//err = stmt.(*SelectStmt).Fields.Restore(NewRestoreCtx(DefaultRestoreFlags, &sbFrom))
	//fmt.Println(sbFrom.String())
	//fmt.Println("=========")
	//
	//stmt1, err := parser.ParseOneStmt(sql2, "", "")
	//var sb1 strings.Builder
	//err = stmt1.Restore(NewRestoreCtx(DefaultRestoreFlags, &sb1))
	//fmt.Println("sql2", sb1.String())
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//CleanNodeText(stmt)
	//CleanNodeText(stmt1)
	//
	//stmt.Accept(visitor{})
	//fmt.Println()
	//fmt.Println()
	//fmt.Println()
	//fmt.Println()
	//stmt1.Accept(visitor{})
	//
	//result := reflect.DeepEqual(stmt, stmt1)
	//fmt.Println(result)

}
