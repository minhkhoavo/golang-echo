package request

type PaginationReq struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

func (p *PaginationReq) GetQueryParams() (offset int, limit int, page int, pageSize int) {
	page = p.Page
	if page < 1 {
		page = 1
	}

	pageSize = p.PageSize
	switch {
	case pageSize < 1:
		pageSize = 10
	case pageSize > 100:
		pageSize = 100
	}

	limit = pageSize
	offset = (page - 1) * pageSize
	return
}
