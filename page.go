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
	MaxPageSize        int    `json:"max_page_size" yaml:"max_page_size"`
	DefaultPageSize    int    `json:"default_page_size" yaml:"default_page_size"`
	PageSizeArgName    string `json:"page_size_arg_name" yaml:"page_size_arg_name"`
	CurrentPageArgName string `json:"current_page_arg_name" yaml:"current_page_arg_name"`
	OrderArgName       string `json:"order_arg_name" yaml:"order_arg_name"`
	NoPageArgName      string `json:"no_page_arg_name" yaml:"no_page_arg_name"`
}

type Pager struct {
	maxPageSize        int
	defaultPageSize    int
	pageSizeArgName    string
	currentPageArgName string
	orderArgName       string
	noPageArgName      string
}

// 根据filter生成对于的查询参数
func (o *Pager) MakeNoPageFilter(filters FilterArgs) *query.Query {
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
func (o *Pager) MakePageFilter(filters FilterArgs) *query.Query {
	currentPage := 1
	pageSize := o.defaultPageSize
	var _querys []string
	var _agrs []interface{}
	var _order []query.OrderBy
	var _noPage bool

	for k, v := range filters {
		if k == o.currentPageArgName {
			if _page, err := strconv.Atoi(v); err == nil {
				currentPage = _page
			}
			continue
		}
		if k == o.pageSizeArgName {
			if _size, err := strconv.Atoi(v); err == nil {
				pageSize = _size
				if pageSize > o.maxPageSize {
					pageSize = o.maxPageSize
				}
			}
			continue
		}

		if k == o.noPageArgName {
			_noPage = true
			continue
		}

		if k == o.orderArgName {
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

// 自动选择是否分页
func (o *Pager) MakeFilter(filters FilterArgs) *query.Query {
	if _, ok := filters[o.noPageArgName]; ok {
		return o.MakeNoPageFilter(filters)
	}
	return o.MakePageFilter(filters)
}

// 分页列表查询器
func (o *Pager) PageQueryResult(db *gorm.DB, filters FilterArgs, results interface{}, preload ...string) (*query.Page, error) {
	query := o.MakePageFilter(filters)
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
func (o *Pager) NoPageQueryResult(db *gorm.DB, filters FilterArgs, results interface{}, preload ...string) error {
	query := o.MakeNoPageFilter(filters)
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
func (o *Pager) QueryResult(db *gorm.DB, filters FilterArgs, results interface{}, preload ...string) (*query.Page, error) {
	if _, ok := filters[o.noPageArgName]; !ok {
		return o.PageQueryResult(db, filters, results, preload...)
	}
	err := o.NoPageQueryResult(db, filters, results, preload...)
	return nil, err
}

func (o *Pager) QueryResultByCommon(db *gorm.DB, filters FilterArgs, results interface{}, preload ...string) (interface{}, error) {
	if _, ok := filters[o.noPageArgName]; !ok {
		return o.PageQueryResult(db, filters, results, preload...)
	}
	err := o.NoPageQueryResult(db, filters, results, preload...)
	return &results, err
}

func NewFilter(options *Options) *Pager {
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

	return &Pager{
		maxPageSize:        options.MaxPageSize,
		defaultPageSize:    options.DefaultPageSize,
		pageSizeArgName:    options.PageSizeArgName,
		currentPageArgName: options.CurrentPageArgName,
		orderArgName:       options.OrderArgName,
		noPageArgName:      options.NoPageArgName,
	}
}
