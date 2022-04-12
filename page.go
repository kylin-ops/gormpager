package gormpager

import (
	"strconv"
	"strings"

	"gitlab.sheincorp.cn/ops_go_sdk/gormpager/query"
	"gorm.io/gorm"
)

// 过滤参数
type FilterArgs map[string]string


// 过滤器
type Options struct {
	MaxPageSize        int
	DefaultPageSize    int
	PageSizeArgName    string
	CurrentPageArgName string
	OrderArgName       string
	NoPageArgName      string
}

type filter struct {
	Options
}

// 根据filter生成对于的查询参数
func (o *filter)MakeNoPageQuery(filters FilterArgs) *query.Query {
	var _querys []string
	var _agrs []interface{}
	var _order []query.OrderBy

	for k, v := range filters {
		if k == "order" {
			_orders := strings.Split(v, ",")
			for _, field := range _orders {
				if field == "" {
					continue
				}
				if field[:1] == "-" {
					_order = append(_order, query.OrderBy{Field: field[1:], Order: -1})
				} else {
					_order = append(_order, query.OrderBy{Field: field, Order: 1})
				}
			}
			continue
		}

		if k != "" {
			_querys = append(_querys, k+" = ?")
			_agrs = append(_agrs, v)
		}
	}
	return &query.Query{
		Where: strings.Join(_querys, " AND "),
		Args:  _agrs,
		Order: _order,
	}
}

// 根据filter生成对于的查询参数
func (o *filter)MakePageQuery(filters FilterArgs) *query.Query {
	currentPage := 1
	pageSize := o.DefaultPageSize
	var _querys []string
	var _agrs []interface{}
	var _order []query.OrderBy
	var _noPage bool

	for k, v := range filters {
		if k == o.CurrentPageArgName {
			if _page, err := strconv.Atoi(v); err == nil {
				currentPage = _page
			}
			continue
		}
		if k == o.PageSizeArgName {
			if _size, err := strconv.Atoi(v); err == nil {
				pageSize = _size
				if pageSize > o.MaxPageSize {
					pageSize = o.MaxPageSize
				}
			}
			continue
		}

		if k == o.NoPageArgName {
			_noPage = true
			continue
		}

		if k == o.OrderArgName {
			_orders := strings.Split(v, ",")
			for _, field := range _orders {
				if field == "" {
					continue
				}
				if field[:1] == "-" {
					_order = append(_order, query.OrderBy{Field: field[1:], Order: -1})
				} else {
					_order = append(_order, query.OrderBy{Field: field, Order: 1})
				}
			}
			continue
		}

		if k != "" {
			_querys = append(_querys, k+" = ?")
			_agrs = append(_agrs, v)
		}
	}
	return &query.Query{
		Where:  strings.Join(_querys, " AND "),
		Args:   _agrs,
		Order:  _order,
		Page:   currentPage,
		Size:   pageSize,
		NoPage: _noPage,
		Limit:  pageSize,
		Offset: (currentPage - 1) * pageSize,
	}
}

// 分页列表查询器
func (o *filter)PageQueryResult(db *gorm.DB, filters FilterArgs, results interface{}, preload ...string) (*query.Page, error) {
	query := o.MakePageQuery(filters)
	page, err := query.PageQuery(db)
	if err != nil {
		return nil, err
	}
	for _, p := range preload {
		db = db.Preload(p)
	}
	if err := db.Find(results).Error; err != nil {
		return nil, err
	}
	page.Results = results
	return page, nil
}

// 不分页列表查询器
func (o *filter)NoPageQueryResult(db *gorm.DB, filters FilterArgs, results interface{}, preload ...string) error {
	query := o.MakePageQuery(filters)
	db = query.Query(db)

	for _, p := range preload {
		db = db.Preload(p)
	}
	if err := db.Find(results).Error; err != nil {
		return err
	}
	return nil
}

// 自动选择器，通过判断filterArgs里是否有no_page的参数
func (o *filter)QueryResult(db *gorm.DB,filters FilterArgs, results interface{}, preload ...string)(*query.Page, error) {
	if _, ok := filters[o.NoPageArgName]; ok {
		return o.PageQueryResult(db, filters, results, preload...)
	}
	err := o.NoPageQueryResult(db, filters, results, preload...)
	return nil, err
}

func NewFilter(options *Options)*filter{
	if options.MaxPageSize == 0 {
		options.MaxPageSize = 50
	}
	if options.DefaultPageSize == 0 {
		options.DefaultPageSize = 20
	}
	if options.DefaultPageSize > options.MaxPageSize {
		options.DefaultPageSize = options.MaxPageSize
	}
	if options.PageSizeArgName == "" {
		options.PageSizeArgName = "size"
	}
	if options.CurrentPageArgName == "" {
		options.CurrentPageArgName = "page"
	}
	if options.OrderArgName == "" {
		options.OrderArgName = "order"
	}
	if options.NoPageArgName == "" {
		options.NoPageArgName = "no_page"
	}

	return &filter{
		Options: *options,
	}
}