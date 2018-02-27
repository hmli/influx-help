package influx_help

import (
	"errors"
	"sort"
	"strings"
	"github.com/influxdata/influxdb/models"
	"github.com/influxdata/influxdb/client/v2"
)


type InfluxHelp struct {}

var (
	ErrNoResult = errors.New("err: No result")
	ErrNoSeries = errors.New("err: No series")
	ErrInconsistent = errors.New("err: Amount of column and values are inconsistent")
)

// !! 只考虑一次请求一条sql语句的情况

// models.Row: 类似于 sql.Rows
// NormalValues 普通的(非group by 语句)， 这时只有一个 series 中有值
func NormalValues(res *client.Response) (data *models.Row, err error) {
	if len(res.Results) == 0 {
		return nil, ErrNoResult
	}
	series := res.Results[0].Series
	if len(series) == 0 {
		return nil, ErrNoSeries
	}
	return &series[0], nil
}

// GroupValues GroupBy 语句返回的数据， 这时会有多个series, 每一个series使用相应的tags区分。
// 因此返回值是将所有tags encode 成一个字符串作Key的map
func GroupValues(res *client.Response) (data map[string]*models.Row, err error) {
	if len(res.Results) == 0 {
		return nil, ErrNoResult
	}
	series := res.Results[0].Series
	if len(series) == 0 {
		return nil, ErrNoSeries
	}
	data = make(map[string]*models.Row)
	for _, row := range series {
		key := TagEncode(row.Tags)
		data[key] = &row
	}
	return data, nil
}

// MapValues 根据row中的Column和Value,  将 Column []string 和 Values [][]interface{} 组装成 []map[string]interface{}
func MapValues(row *models.Row) (data []map[string]interface{}, err error) {
	if len(row.Values) == 0{
		return data, nil
	}
	if len(row.Columns) != len(row.Values[0]) {
		return nil, ErrInconsistent
	}
	for _, value := range row.Values {
		colValue := make(map[string]interface{})
		for j, col := range row.Columns {
			colValue[col] = value[j]
		}
		data = append(data, colValue)
	}
	return data, nil
}

func ColToMap(columns []string, value []interface{}) (data map[string]interface{}, err error) {
	if len(columns) != len(value) {
		return nil, ErrInconsistent
	}
	data = make(map[string]interface{})
	for i, colName := range columns {
		data[colName] = value[i]
	}
	return
}

// Col 根据字段名找到其中的一列数据
func Col(row *models.Row, name string) (col []interface{}) {
	for i, col := range row.Columns {
		if col == name && i < len(row.Values) {
			return row.Values[i]
		}
	}
	return
}

func FirstValue(row *models.Row) (value []interface{}){
	if len(row.Values) == 0{
		return value
	}
	return row.Values[0]
}

// TagEncode 将返回的 map 格式 tags 转化成正确格式的key
// key值算法： 将所有tags的 "key=value" 格式，先按字符排序，再用逗号相连，形成 "k1=v1,k2=v2,k3=v3'的格式
// TODO 使用 bytes 会不会更快? 有空试
func TagEncode(tags map[string]string) (key string) {
	tagList := make([]string, 0)
	for k,v := range tags {
		tagList = append(tagList, k+"="+v)
	}
	sort.Strings(tagList)
	key = strings.Join(tagList, ",")
	return

}

// TagDecode 上面TagEncode的逆过程
func TagDecode(key string) (tags map[string]string) {
	tagList := strings.Split(key, ",")
	tags = make(map[string]string)
	for _, tagValue := range tagList {
		i := strings.Index(tagValue, "=")
		if i == -1  || i == len(tagValue) - 1{
			continue
		}
		k, v := tagValue[:i], tagValue[i+1:]
		// TODO string validation
		tags[k] = v
	}
	return

}
