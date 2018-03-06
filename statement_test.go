package influx_help

import "testing"

func TestSelectStr(t *testing.T) {
	db := NewDB("addr", "user", "pswd")
	db.ShowSQL = true
	session := db.NewSession("testlog", "us")
	stmt := Statement{}
	stmt.Init(session)
	stmt.Select("id, content")
	stmt.Table("testlog")
	stmt.And("id = ?", 5).And("price = ?", 3.4).Or("name = 'ddd'")
	condSQL, condArgs, err := stmt.condSQL()
	t.Log(err, condArgs)
	t.Log(stmt.columnSQL())
	query := stmt.selectSQLNoArgs(stmt.columnSQL(), condSQL)
	t.Log(query)
	t.Log(stmt.selectSQL(query, condArgs))
	stmt2 := Statement{}
	stmt2.Init(session)
	stmt2.Table("oplog").Select("content").Where("id = ?", 3).GroupBy("action, id")
	condSQL, condArgs, err = stmt2.condSQL()
	query = stmt2.selectSQLNoArgs(stmt2.columnSQL(), condSQL)
	t.Log(query)
	t.Log(stmt2.selectSQL(query, condArgs))

	stmt.Query()
	stmt2.Query()

}