package model

// PaginationParams 分页请求参数
type PaginationParams struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// GetOffset 计算偏移量
func (p *PaginationParams) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Page > 1000000 {
		p.Page = 1000000
	}
	return (p.Page - 1) * p.GetLimit()
}

// GetLimit 获取限制数量
func (p *PaginationParams) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// PaginatedResponse 分页响应结构
type PaginatedResponse struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Pagination 分页信息
type Pagination struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// NewPaginatedResponse 创建分页响应
func NewPaginatedResponse(data any, page, pageSize int, total int64) PaginatedResponse {
	if page <= 0 {
		page = 1
	}
	if page > 1000000 {
		page = 1000000
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}
	return PaginatedResponse{
		Data: data,
		Pagination: Pagination{
			Page:      page,
			PageSize:  pageSize,
			Total:     total,
			TotalPage: totalPage,
		},
	}
}
