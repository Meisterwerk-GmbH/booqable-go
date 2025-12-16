package list

type listMeta struct {
	Page        int `json:"page,omitempty"`
	CurrentPage int `json:"current_page,omitempty"`
	PerPage     int `json:"per_page,omitempty"`
	TotalPages  int `json:"total_pages,omitempty"`
}

type listLinks struct {
	Next string `json:"next,omitempty"`
}

func hasMore(meta *listMeta, links *listLinks, page, per, got int) bool {
	if meta != nil {
		current := meta.Page
		if current == 0 {
			current = meta.CurrentPage
		}
		if meta.TotalPages > 0 && current < meta.TotalPages {
			return true
		}
	}
	if links != nil && links.Next != "" {
		return true
	}
	return got == per && got > 0
}
