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

type Detail struct {
	ID        uint
	Title     string
	WrittenBy string
	CreateAt  time.Time
	Content   string
	Images    []string
	Files     []File
}

type File struct {
	Name     string
	Download string
	Preview  string
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
				id, _ := strconv.Atoi(idStr[3])
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

func CrawlPage(url Url, id uint) (*Detail, error) {
	res, err := http.Get(string(url) + "/" + strconv.Itoa(int(id)))
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

	result := &Detail{
		ID: id,
	}
	doc.Find(".viewBox p").Each(func(i int, s *goquery.Selection) {
		txt := strings.ReplaceAll(s.Text(), "<br>", "")
		result.Content += txt + "\n"
	})
	doc.Find(".viewBox img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		result.Images = append(result.Images, "http://school.gy06.net"+src)
	})
	doc.Find(".infoBox li").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			result.WrittenBy = strings.ReplaceAll(s.Text(), "작성자", "")
		case 1:
			t, _ := time.Parse("2006-01-02", strings.ReplaceAll(s.Text(), "작성일", ""))
			result.CreateAt = t
		}
	})
	doc.Find(".fieldBox dd").Each(func(i int, s *goquery.Selection) {
		newFile := File{}
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				newFile.Name = s.Text()
				downloadLink, _ := s.Attr("href")
				newFile.Download = "https://school.gy06.net" + downloadLink
			case 1:
				previewLink, _ := s.Attr("href")
				newFile.Preview = previewLink
			}
		})
		result.Files = append(result.Files, newFile)
	})

	return result, nil
}
