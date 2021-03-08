package osangdata

import "testing"

func TestCrawl(t *testing.T) {
	result, err := CrawlList(UrlNotice, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("list: %+v\n", result[0])
	detail, err := CrawlPage(UrlNotice, result[0].ID)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("detail: %+v\n", detail)
}
