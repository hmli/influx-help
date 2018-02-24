package influx_help

import "testing"

func TestSelectStr(t *testing.T) {
	stmt := Statement{}
	stmt.Init("testlog")
	stmt.Select("id, content")
	stmt.Table("testlog")
	stmt.And("id = ?", 5).And("price = ?", 3.4).Or("name = 'ddd'")
	condSQL, condArgs, err := stmt.condSQL()
	t.Log(err)
	t.Log(condArgs)
	t.Log(stmt.columnSQL())
	query := stmt.selectSQLNoArgs(stmt.columnSQL(), condSQL)
	t.Log(query)
	t.Log(stmt.selectSQL(query, condArgs))
}