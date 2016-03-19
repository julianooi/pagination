package pagination

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

var (
	// pageSize determines the size of each page, default is 20
	pageSize = 20

	// pagePadding defines the number of pages to show before and after the current page
	pagePadding = 4

	// pageParam defines the parameter to look for in the given request
	pageParam = "_page"

	// showFirst determines if the first page should be shown regardless of the current page
	showFirst = true

	// showLast determines if the last page should be shown regardless of the current page
	showLast = true
)

// Config is the configuration required to create a page
type Config struct {
	PageSize    int
	PagePadding int
	PageParam   string
	ShowFirst   bool
	ShowLast    bool
}

// Page struct represents the current page with all of the values required to generate a page link
type Page struct {
	Page  int
	Total int
	URI   string

	PageLinks []PageLink
}

// PageLink represents the details of the page
type PageLink struct {
	Page int

	Current bool
	First   bool
	Last    bool
}

// NewPage creates a Page struct used for pagination
func NewPage(r *http.Request, total int, conf Config) Page {
	p := getCurrentPage(r, conf.PageParam)
	return Page{
		Page:  p,
		Total: total,
		URI:   uriWithoutPage(r, conf.PageParam),

		PageLinks: generatePageURIs(p, total, conf),
	}
}

// getCurrentPage returns an integer representing the current page, valid values are >= 1
func getCurrentPage(r *http.Request, pageParam string) int {
	strPage := r.URL.Query().Get(pageParam)

	if page, err := strconv.Atoi(strPage); err == nil {
		return page
	}

	return 1
}

func uriWithoutPage(r *http.Request, pageParam string) string {
	u := r.URL
	path := u.Path
	q := u.Query()

	// create a copy of the query values
	nmap := new(url.Values)

	for k, v := range q {
		if k == pageParam {
			continue
		}
		(*nmap)[k] = v
	}

	return fmt.Sprintf("%s?%s", path, nmap.Encode())
}

func generatePageURIs(current int, total int, conf Config) []PageLink {
	totalPages := int(math.Ceil(float64(total / conf.PageSize)))

	min := 1
	extraMax := 0
	min = current - conf.PagePadding
	if min < 1 {
		extraMax = int(math.Abs(float64(min))) + 1
		min = 1
	}

	max := current + conf.PagePadding + extraMax

	pageLinks := []PageLink{}

	for p := 1; p <= totalPages; p++ {
		if p == current {
			// it's the current page, mark as current
			pageLinks = appendPageLink(pageLinks, p, true, false, false)
			continue
		}

		if p >= min && p <= max ||
			// if it is only 1 away from first page
			(p == 1 && current-conf.PagePadding-1 == 1) ||
			// if it is only 1 away from last page
			(p == totalPages && p+conf.PagePadding+1 == totalPages) {

			// append page link
			pageLinks = appendPageLink(pageLinks, p, false, false, false)
		} else {
			if p == 1 && conf.ShowFirst {
				// first page
				pageLinks = appendPageLink(pageLinks, p, false, true, false)
			} else if p == totalPages && conf.ShowLast {
				pageLinks = appendPageLink(pageLinks, p, false, false, true)
			}
		}
	}

	return pageLinks
}

func appendPageLink(pageLinks []PageLink, page int, current bool, first bool, last bool) []PageLink {
	return append(pageLinks, PageLink{
		Page:    page,
		Current: current,
		First:   first,
		Last:    last,
	})
}
