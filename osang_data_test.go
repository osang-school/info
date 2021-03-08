package osangdata

import "testing"

func TestCrawl(t *testing.T) {
	result, err := CrawlList(UrlNotice, 1)
	if err != nil {
		t.Error(err)
	}
	t.Logf("list: %+v\n", result[0])
}
