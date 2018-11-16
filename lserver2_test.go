package influx_help

import (
	"testing"
)

// 本地局域网的一台测试服务器
func Test_lserver(t *testing.T) {
	address := "http://192.168.12.137:8086" // a local server
	username := ""
	password  := ""
	db := NewDB(address, username, password)
	db.ShowSQL = true
	sess := db.NewSession("testlog", "us")
	sess.Table("temperature").InsertOne(map[string]string{
		"action": "view",
	}, map[string]interface{}{
		"qty": 1,
	})
	res, err := sess.Table("temperature").Where("action = ?", "view").Limit(3).Find()
	t.Log(err, res)
	row, err := NormalRow(res)
	t.Log(err)
	for _, values := range row.Values {
		t.Log(values)
	}
	res2, err := sess.Table("temperature").GroupBy("action").Find()
	t.Log(res2, err)
	row2, err := GroupRows(res2)
	t.Log(err)
	for k, r := range row2{
		for _, values := range r.Values {
			t.Log(k, "->", values)
		}
	}
}
