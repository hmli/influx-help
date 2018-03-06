package influx_help

import "testing"

// 本地局域网的一台测试服务器
func Test_lserver(t *testing.T) {
	address := "http://192.168.15.95:8086"
	username := ""
	password  := ""
	db := NewDB(address, username, password)
	db.ShowSQL = true
	sess := db.NewSession("testlog", "us")
	res, err := sess.Table("temperature").Where("action = ?", "view").Limit(3).Find()
	t.Log(err)
	t.Log(res)

}
