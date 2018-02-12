package influx_help

import (
	"strings"
	"github.com/go-xorm/builder"
)

type Statement struct {
	Start           int
	LimitN          int
	OrderStr        string
	GroupByStr      string
	HavingStr       string
	ColumnStr       string
	selectStr       string
	//columnMap       map[string]bool
	//useAllCols  bool
	tableName       string
	RawSQL          string // RawSQL 和 RawParams 优先级最高。 如果 RawSQL 不为空，就只使用RawSQL； 如果为空，才使用其它字段拼接
	RawParams       []interface{}
	//UseCascade      bool
	//Charset         string
	UseAutoTime     bool
	//noAutoCondition bool
	//IsDistinct      bool
	TableAlias      string
	cond            builder.Cond
}

func (stmt *Statement) Init() {
	stmt.Start = 0
	stmt.LimitN = 0
	stmt.OrderStr = ""
	stmt.GroupByStr = ""
	stmt.HavingStr = ""
	stmt.ColumnStr = ""
	stmt.tableName = ""
	stmt.RawSQL = ""
	stmt.RawParams = make([]interface{}, 0)
	stmt.cond = builder.NewCond()
}

func (stmt *Statement) SQL(query string, args ...interface{}) *Statement {
	stmt.RawSQL = query
	stmt.RawParams = args
	return stmt
}

func (stmt *Statement) Table(name string) *Statement {
	stmt.tableName = name
	return stmt
}

func (stmt *Statement) Cols(columns ...string) *Statement {
	stmt.ColumnStr = strings.Join(columns, ", ")
	return stmt
}

func (stmt *Statement) Select(str string) *Statement {
	stmt.selectStr = str
	return stmt
}

func (stmt *Statement) Where(query string, args ...interface{}) *Statement {
	return stmt.And(query, args...)
}

func (stmt *Statement) And(query string, args ...interface{}) *Statement {
			cond := builder.Expr(query, args...)
			stmt.cond = stmt.cond.And(cond)
			return stmt
}

func (stmt *Statement) Or(query string, args ...interface{}) *Statement {
	cond := builder.Expr(query, args...)
	stmt.cond = stmt.cond.Or(cond)
	return stmt
}

func (stmt *Statement) In(column string, args ...interface{}) *Statement {
	in := builder.In(column, args...)
	stmt.cond = stmt.cond.And(in)
	return stmt
}

func (stmt *Statement) NotIn(column string, args ...interface{}) *Statement {
	notIn := builder.NotIn(column, args...)
	stmt.cond = stmt.cond.And(notIn)
	return stmt
}

func (stmt *Statement) Limit(limit int, start ...int) *Statement {
	stmt.LimitN = limit
	if len(start) > 0 {
		stmt.Start = start[0]
	}
	return stmt
}

func (stmt *Statement) OrderBy(order string) *Statement {
	if len(stmt.OrderStr) > 0 {
		stmt.OrderStr += ", "
	}
	stmt.OrderStr += order
	return stmt
}

func (stmt *Statement) Query() (sql string) {
	return
}

// TODO test Cols/Select, Cond, Limit
// TODO 1. Desc,Asc  2. GroupBy  3. Make SQL

var QuoteStr = "'"

func quote(sql string) string {
	return QuoteStr + sql + QuoteStr
}