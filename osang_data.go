package osangdata

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Url string

const (
	baseUrl           Url = "http://school.gyo6.net/osangms"
	UrlNotice             = baseUrl + "/0301/board/70879"
	UrlPrints             = baseUrl + "/0302/board/70880"
	UrlRule               = baseUrl + "/020301/board/70875"
	UrlEvaluationPlan     = baseUrl + "/141482/board/70873"
	UrlAdministration     = baseUrl + "/0303/board/70881"
)

type List struct {
	ID        uint
	Number    uint
	Title     string
	WrittenBy string
	CreateAt  time.Time
}

func CrawlList(url Url, page uint) ([]*List, error) {
	res, err := http.Get(string(url) + "?page=" + strconv.Itoa(int(page)))
	if err != nil {
		return nil, fmt.Errorf("Error Loading")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error Loading")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result []*List
	doc.Find("table[class!=\"md\"] tbody tr").Each(func(i int, s *goquery.Selection) {
		if i < 1 || i > 10 {
			return
		}
		newItem := &List{}
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				val, _ := strconv.Atoi(s.Text())
				newItem.Number = uint(val)
			case 1:
				title, _ := s.Find("a").First().Attr("title")
				onclick, _ := s.Find("a").First().Attr("onclick")
				idStr := strings.Split(onclick, "'")
				newItem.Title = title
				id, _ := strconv.Atoi(idStr[4])
				newItem.ID = uint(id)
			case 2:
				newItem.WrittenBy = s.Text()
			case 4:
				t, _ := time.Parse("2006-01-02", s.Text())
				newItem.CreateAt = t
			}
		})
		result = append(result, newItem)
	})
	return result, nil
}
