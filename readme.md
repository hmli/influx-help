# Influx Help

## TODO List

* 增加 Insert (One, Many)
* 继续测试并适配 Influx 版本的 SQL 语句
* 完善 test cases
* 代码中的其它 "TODO"


Influx cilent 返回的 `models.Response` 结构有点复杂，不太好用，比如官方 demo 中的代码：

 ```
count := res[0].Series[0].Values[0][1]
for i, row := range res[0].Series[0].Values {}
```
所以我写了几个函数提供一些辅助功能, 尽可能方便地从 `response` 中获取数据。 详情见代码注释。

因为我自己的需求中一次只会发一个 `select` 语句， 所以本项目中所有的函数也都仅限于这种情况。

正在参考 `xorm` 和 `gorm` 的语法， 逐步实现成一个简单的 ORM. 目前 拼接 SQL 的功能正在开发中。