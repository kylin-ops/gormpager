# 作用
gormpage的主要功能是简化数gorm查询是，需要使用复杂的orm域名

# 使用方法
## 1 初始化查询器
```go

// 初始化查询器
p := gormpager.NewFilter(&gormpager.Options{
    MaxPageSize: 50,
    DefaultPageSize: 20,
    PageSizeArgName: "size",   // 指定每页多少条数据的查询参数
    CurrentPageArgName: "page", // 指定当前页的查询参数
    OrderArgName: "order",      // 指定排序查询参数
    NoPageArgName: "no_page",   // 指定不分页的参数
})
```

## 2 结果查询器
根据查询器
```go
// 结果查询器，自动识别分页查询和不分页查询，根据查询参数中是否有NoPageArgName参数决定是否分页
// db：gorm.DB对象， 传入之前使用db.model(&struct)指定查询的表
// filters： gormpager.FilterArgs 的map类型，传入查询的key与value
// likeFields： []string 类型， 指定哪些字段使用like模糊查询
// results：    传入存储结果的指针
// preload：    传入需要预加载的字段
results, err := p.QueryResultByCommon(db, filters, likeFields， &result, preload)


// 结果查询器，不进行分页查询
// db：gorm.DB对象， 传入之前使用db.model(&struct)指定查询的表
// filters： gormpager.FilterArgs 的map类型，传入查询的key与value
// likeFields： []string 类型， 指定哪些字段使用like模糊查询
// results：    传入存储结果的指针
// preload：    传入需要预加载的字段
result = []struct{}
err := p.NoPageQueryResult(db, filters, likeFields， &result, preload)


// 结果查询器，分页查询器
// db：gorm.DB对象， 传入之前使用db.model(&struct)指定查询的表
// filters： gormpager.FilterArgs 的map类型，传入查询的key与value
// likeFields： []string 类型， 指定哪些字段使用like模糊查询
// results：    传入存储结果的指针
// preload：    传入需要预加载的字段
result = []struct{}
r, err := p.NoPageQueryResult(db, filters, likeFields， &result, preload)
```

## 3 过滤条件解析器
```go
// 通用过滤条件解析器，自动根据NoPageArgName判断是否分页
// filters： gormpager.FilterArgs 的map类型，传入查询的key与value
// likeFields： []string 类型， 指定哪些字段使用like模糊查询
// 返回值为 query.Query类型的查询器
query := p.MakeFilter(FilterArgs, likeField)

// 分页条件过滤器，根据NoPageArgName查找过了器中的分页参数
query := p.MakePageFilter(FilterArgs, likeField)

// 不分页条件过滤器
query := p.MakeNoPageFilter(FilterArgs, likeField)
```

## 4 查询器
```go
q := query.Query{
    Where: "id = ? AND name like ?",
    Args: []string{"uid_value", "name_value"},
    Order: []string{"id", "name"},
    Limit: 10,
    Offset: 40,
}
// db 为gorm.DB类型
db = db.model(&struct{})
// 将查询参数自动设置到db中
db = q.Query(db)

```