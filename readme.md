# Influx Help

一个正在开发中的 Influx ORM, 目前能提供部分不算太复杂的查询功能和单条 insert 功能。

## TODO List

* 增加 InsertMany
* 继续测试并适配 InfluxSQL
* 完善 test cases
* 代码中的其它 "TODO"
* 英文 readme


## res.go

Influx cilent 返回的 `models.Response` 结构有点复杂，不太好用，比如官方 demo 中的代码：

 ```
count := res[0].Series[0].Values[0][1]
for i, row := range res[0].Series[0].Values {}
```
所以 res.go 中提供了几个函数提供一些辅助功能, 尽可能方便地从 `response` 中获取数据。 详情见代码注释。

因为我自己的需求中一次只会发一个 `select` 语句， 所以本项目中所有的函数也都仅限于这种情况。


## 其它

参考了 `xorm` 中的一些实现， 而且拼接SQL的过程直接使用了 `xorm` 中的模块。