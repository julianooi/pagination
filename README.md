# Pagination

Used to create links used for webapps and stuff.

###Usage

Create new page from *http.Request
```
// a pagingation configuration
config := pagination.Config{
    PageSize:    20,
    PagePadding: 3,
    PageParam:   "_page",
    ShowFirst:   true,
    ShowLast:    true,
}
// total number of items
totalItems := 200

page := pagination.NewPage(request, totalItems, config)
```
