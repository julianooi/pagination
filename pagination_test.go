package pagination

import (
	"net/http"
	"strconv"
	"testing"
)

func initConf() Config {
	return Config{
		PageSize:    20,
		PagePadding: 3,
		PageParam:   "_page",
		ShowFirst:   true,
		ShowLast:    true,
	}
}

func createRequest(t *testing.T, page int) *http.Request {
	r, err := http.NewRequest("GET", "http://pagination.test/home?_page="+strconv.Itoa(page), nil)
	if err != nil {
		t.Errorf("Unable to create request, error: %s", err.Error())
	}
	return r
}

func TestPaginationLength(t *testing.T) {
	conf := initConf()

	r := createRequest(t, 5)
	p := NewPage(r, 220, conf)

	// test for page with 200 total items
	// length must be equal to page padding * 2 (for left and right side) + 3 (the current, first and last page)
	if len(p.PageLinks) != (conf.PagePadding*2 + 3) {
		t.Errorf("PageLinks length is %d, expected %d", len(p.PageLinks), conf.PagePadding*2+3)
	}

	r = createRequest(t, 3)
	p = NewPage(r, 200, conf)
	// length must be equal to page padding * 2 (for left and right side) + 2 (the current and last page)
	if len(p.PageLinks) != (conf.PagePadding*2 + 2) {
		t.Errorf("PageLinks length is %d, expected %d", len(p.PageLinks), conf.PagePadding*2+2)
	}
}

// TestPaginationMiddle tests from the middle page
func TestPaginationMiddle(t *testing.T) {
	conf := initConf()

	r := createRequest(t, 6)
	// expected pages = 11
	p := NewPage(r, 220, conf)

	for idx, pg := range p.PageLinks {
		if idx == 0 {
			// check for first page
			checkIsPage(t, pg, 1)
			if !pg.First {
				t.Errorf("Expected to be first page, but got %+v", pg)
			}
			continue
		} else if idx == 8 {
			// check for last page
			checkIsPage(t, pg, 11)
			if !pg.Last {
				t.Errorf("Expected to be last page, but got %+v", pg)
			}
			break
		}

		checkIsPage(t, pg, idx+2)
	}
}

func TestPagination(t *testing.T) {
	conf := initConf()

	r := createRequest(t, 1)

	p := NewPage(r, 200, conf)

	for idx, pg := range p.PageLinks {
		if idx == 7 {
			// check if it's the last page
			checkIsPage(t, pg, 10)
			if !pg.Last {
				t.Errorf("Expected to be last page, but got %+v", pg)
			}
			break
		}

		checkIsPage(t, pg, idx+1)
	}
}

func checkIsPage(t *testing.T, p PageLink, expected int) {
	if p.Page != expected {
		t.Errorf("Expected page to be %d, got %+v instead", expected, p)
	}
}
