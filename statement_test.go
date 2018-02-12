package influx_help

import "testing"

func TestSelectStr(t *testing.T) {
	stmt := Statement{}
	stmt.Init()
	stmt.Select("id, content")
	stmt.Table("testlog")
	t.Log(stmt.selectStr)
}