package paginate

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

// Params holds pagination query parameters.
type Params struct {
	Page  int `json:"page"  query:"page"`
	Limit int `json:"limit" query:"limit"`
}

// Normalize clamps Page and Limit to valid ranges in-place.
// Call this explicitly before using the params.
func (p *Params) Normalize() {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.Limit < 1 || p.Limit > MaxLimit {
		p.Limit = DefaultLimit
	}
}

// Offset returns the SQL OFFSET value for the current page.
// It reads Page and Limit without mutating the receiver —
// call Normalize() first if inputs are user-supplied.
func (p *Params) Offset() int {
	page := p.Page
	if page < 1 {
		page = DefaultPage
	}
	limit := p.Limit
	if limit < 1 || limit > MaxLimit {
		limit = DefaultLimit
	}
	return (page - 1) * limit
}
