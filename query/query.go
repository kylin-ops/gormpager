package query

import (
	"gorm.io/gorm"
)

type OrderBy struct {
	Field string
	Order int8 // 0不排序， >0正序； <0倒序
}

// 定义查询参数
type Query struct {
	Where  interface{}
	Args   []interface{}
	Order  []OrderBy
	Limit  int
	Offset int
	Page   int
	Size   int
	NoPage bool
}

type Page struct {
	// 数据总行数
	TotalRow int `json:"total_row"`
	// 总页数
	TotalPage int `json:"total_page"`
	// 没页多少数据
	PageSize int `json:"page_size"`
	// 当前页码
	CurrentPage int `json:"current_page"`
	// 返回的数据
	Results interface{} `json:"results"`
}


func totalRow(db *gorm.DB) (int, error) {
	var n int64
	err := db.Count(&n).Error
	return int(n), err
}

func totalPage(count, pageSize int) int {
	page := count / pageSize
	mod := count % pageSize
	if mod > 0 {
		page++
	}
	return page
}

// 返回查询拼接语句, gorm需要指定model
func (q *Query) Query(d *gorm.DB) *gorm.DB {
	if q.Where != nil {
		d = d.Where(q.Where, q.Args...)
	}

	if q.Order != nil {
		for _, order := range q.Order {
			if order.Order > 0 {
				d = d.Order(order.Field + " " + "ASC")
			}
			if order.Order < 0 {
				d = d.Order(order.Field + " " + "DESC")
			}
		}
	}
	if q.Offset != 0 {
		d = d.Offset(q.Offset)
	}
	if q.Limit == 0 {
		q.Limit = 50
	}
	d = d.Limit(q.Limit)
	return d
}

// 分页器, gorm需要指定model确定查下的表格，返回分页数据格式有基本数据
func (q *Query) PageQuery(db *gorm.DB) (*Page, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	row, err := totalRow(db)
	q.Query(db)
	if err != nil {
		return nil, err
	}
	totalPage := totalPage(row, q.Size)

	if q.Page > totalPage {
		q.Page = totalPage
	}
	offset := (q.Page - 1) * q.Size
	q.Limit = q.Size
	q.Offset = offset
	_page := &Page{
		TotalRow:    row,
		TotalPage:   totalPage,
		PageSize:    q.Size,
		CurrentPage: q.Page,
	}
	return _page, err
}