# Bookmarks

## List

```go
	svc := instapaper.BookmarkServiceOp{
				Client: apiClient,
			}
			bookmarkList, err := svc.List(instapaper.BookmarkListRequestParams{
				Limit: 5,
			})
```