package influx_help

import (
	"github.com/influxdata/influxdb/client/v2"
	"testing"
)

var u = "http://192.168.15.95:8086"



func TestGroup(t *testing.T) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     u,
	})
	if err != nil {
		t.Fatal(err)
	}
	q := "SELECT content FROM oplog group by action,id"
	t.Log(q)
	res, err := queryDB(c, q)
	if err != nil {
		t.Fatal(err)
	}
	data, err := GroupValues(res)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", data)
	if len(data) == 0 {
		t.Log("no data")
		t.FailNow()
	}
	// Tag Decode
	for k := range data {
		tags := TagDecode(k)
		t.Log(tags)
		if len(tags) == 0 {
			t.FailNow()
		}
	}
}

func TestNormalAndMap(t *testing.T) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     u,
	})
	if err != nil {
		t.Fatal(err)
	}
	q := "select * from oplog where id='1'"
	t.Log(q)
	res, err := queryDB(c, q)
	if err != nil {
		t.Fatal(err)
	}
	// NormalValues
	rows, err := NormalValues(res)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", rows)
	// MapValues
	data ,err := MapValues(rows)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", data)
	if len(data) != len(rows.Values) {
		t.FailNow()
	}
	// ColToMap
	if len(rows.Values) > 0 {
		data, err := ColToMap(rows.Columns, rows.Values[0])
		if err != nil {
			t.Fatal(err)
		}
		t.Log(data)
		if len(data) == 0 {
			t.FailNow()
		}
	}
	// Col
	col := Col(rows, "content")
	t.Log(col)
	if len(col) == 0 {
		t.FailNow()
	}
	// FirstValue
	v := FirstValue(rows)
	t.Log(v)
	if len(v) == 0 {
		t.FailNow()
	}
}


func queryDB(clnt client.Client, cmd string) (res *client.Response, err error) {
	q := client.Query{
		Command:  cmd,
		Database: "testlog",
		Precision: "ns",
	}
	res, err = clnt.Query(q)
	if err != nil {
		return res, err
	}
	if res.Error() != nil {
		return res, res.Error()
	}
	return res, nil
}