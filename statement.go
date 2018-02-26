package influx_help

import (
	"strings"
	"github.com/go-xorm/builder"
	"fmt"
	"bytes"
	"github.com/influxdata/influxdb/client/v2"
	"time"
	"errors"
	"strconv"
)
var ErrSkip = errors.New("skip fast-path; continue as if unimplemented")
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
	Database string
}

func (stmt *Statement) Init(database string) {
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
	stmt.Database = database
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

//func (stmt *Statement) In(column string, args ...interface{}) *Statement {
//	in := builder.In(column, args...)
//	stmt.cond = stmt.cond.And(in)
//	return stmt
//}

//func (stmt *Statement) NotIn(column string, args ...interface{}) *Statement {
//	notIn := builder.NotIn(column, args...)
//	stmt.cond = stmt.cond.And(notIn)
//	return stmt
//}

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

func (stmt *Statement) GroupBy(keys string) *Statement {
	stmt.GroupByStr = keys
	return stmt
}

func (stmt *Statement) Query() (query *client.Query, err error) {
	condSQL, condArgs, err := stmt.condSQL()
	if err != nil {
		return nil, err
	}
	columnSQL := stmt.columnSQL()
	queryStr := stmt.selectSQLNoArgs(columnSQL, condSQL)
	selectSQL, err := stmt.selectSQL(queryStr, condArgs)
	if err != nil {
		return nil, err
	}
	q := client.Query{Command: selectSQL, Database: stmt.Database}
	return &q, nil
}

func (stmt *Statement) condSQL() (sql string, args []interface{}, err error) {
	// TODO 先使用 xorm 自带的 builder，后面还需要修改，要补上一些 ' 等。
	sql, args, err = builder.ToSQL(stmt.cond)
	return
}

func (stmt *Statement) columnSQL() string{
	if stmt.selectStr != "" {
		return stmt.selectStr
	}
	return stmt.ColumnStr
}

// selectSQL 拼接出select SQL语句，但参数仍以？代替
func (stmt *Statement) selectSQLNoArgs(columnStr, condSQL string) string {
	var buf bytes.Buffer
	if len(condSQL) > 0 {
		fmt.Fprintf(&buf, " WHERE %v", condSQL)
	}
	var whereStr = buf.String()
	var fromStr = " FROM " + stmt.tableName
	a := fmt.Sprintf("SELECT %v%v%v", columnStr, fromStr, whereStr)
	if len(stmt.OrderStr) > 0 {
		a +=  " ORDER BY " + stmt.OrderStr
	}
	if len(stmt.GroupByStr) > 0 {
		a += " GROUP BY " + stmt.GroupByStr
	}
	if stmt.Start > 0 {
		a = fmt.Sprintf("%s LIMIT %v OFFSET %v", a, stmt.LimitN, stmt.Start)
	} else if stmt.LimitN > 0 {
		a = fmt.Sprintf("%s LIMIT %v", a, stmt.LimitN)
	}
	return a
}

func (stmt *Statement) selectSQL(query string, args []interface{}) (a string, err error) {
	// TODO 目前是用MySQL的格式拼接的，需测试逐步改到influx格式
	if strings.Count(query, "?") != len(args) {
		return "", ErrSkip
	}
	argPos := 0
	buf := make([]byte, 0) // TODO buf pool
	for i:=0; i<len(query); i++ {
		q := strings.IndexByte(query[i:], '?')
		if q == -1 {
			buf = append(buf, query[i:]...)
			break
		}
		buf = append(buf, query[i:i+q]...)
		i += q

		arg := args[argPos]
		argPos++

		if arg == nil {
			buf = append(buf, "NULL"...)
			continue
		}
		switch v := arg.(type) {
		case int64:
			buf = strconv.AppendInt(buf, v, 10)
		case int:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case int32:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case int16:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case int8:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case float64:
			buf = strconv.AppendFloat(buf, v, 'g', -1, 64)
		case float32:
			buf = strconv.AppendFloat(buf, float64(v), 'g', -1, 64)
		case bool:
			if v {
				buf = append(buf, '1')
			} else {
				buf = append(buf, '0')
			}
		case time.Time:
			if v.IsZero() {
				buf = append(buf, "'0000-00-00'"...)
			} else {
				buf = append(buf, []byte(v.Format("'2006-01-02 15:04:05'"))...)
			}
		case []byte:
			buf = append(buf, '\'')
			buf = append(buf, v...)
			buf = append(buf, '\'')
		case string:
			buf = append(buf, '\'')
			buf = append(buf, []byte(v)...)
			buf = append(buf, '\'')
		default:
			return "", ErrSkip
		}
	}
	return string(buf), nil

}

// TODO 自动识别 tag/field, 并给 tag 里的值加''  (schema 随时可能会变，暂没有好的方法

//var QuoteStr = "'"
//
//func quote(sql string) string {
//	return QuoteStr + sql + QuoteStr
//}