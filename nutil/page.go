package nutil

// Page -
type Page struct {
	PageNo   int32 `form:"pageNo" json:"pageNo" binding:"gt=0"`
	PageSize int32 `form:"pageSize" json:"pageSize" binding:"gt=0"`
	Total    int64 `json:"total"`
}

// NewPage -
func NewPage(pageNo int32, pageSize int32) *Page {
	return &Page{
		PageNo:   pageNo,
		PageSize: pageSize,
	}
}

// Offset -
func (p *Page) Offset() int {
	return int((p.PageNo - 1) * p.PageSize)
}

// Limit -
func (p *Page) Limit() int {
	return int(p.PageSize)
}
