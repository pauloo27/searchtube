package searchtube

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

func getContent(data []byte, index int) []byte {
	id := fmt.Sprintf("[%d]", index)
	contents, _, _, _ := jsonparser.Get(data, "contents", "twoColumnSearchResultsRenderer", "primaryContents", "sectionListRenderer", "contents", id, "itemSectionRenderer", "contents")
	return contents
}

type SearchResult struct {
	Title, Uploader, URL, RawDuration, ID, Thumbnail string
	Live                                             bool
	duration                                         time.Duration
}

func (s *SearchResult) GetDuration() (duration time.Duration, err error) {
	duration = s.duration
	if s.Live {
		err = errors.New("cannot get duration of a live")
		return
	}
	if duration == 0 {
		str := s.RawDuration + "s"
		if strings.Count(str, ":") == 2 {
			str = strings.Replace(str, ":", "h", 1)
		}
		str = strings.Replace(str, ":", "m", 1)

		duration, err = time.ParseDuration(str)
	}
	return
}

var httpClient = &http.Client{}

func Search(searchTerm string, limit int) (results []*SearchResult, err error) {
	url := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", url.QueryEscape(searchTerm))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Cannot create GET request: %v", err)
	}
	req.Header.Add("Accept-Language", "en")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Cannot get youtube page: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot read body: %v", err)
	}

	body := string(buffer)
	splittedScript := strings.Split(body, `window["ytInitialData"] = `)

	if len(splittedScript) != 2 {
		splittedScript = strings.Split(body, `var ytInitialData = `)
	}

	if len(splittedScript) != 2 {
		if err != nil {
			return nil, fmt.Errorf("Cannot split script: %v", err)
		}
	}

	splittedScript = strings.Split(splittedScript[1], `window["ytInitialPlayerResponse"] = null;`)
	jsonData := []byte(splittedScript[0])

	index := 0
	var contents []byte

	for {
		contents = getContent(jsonData, index)
		_, _, _, err = jsonparser.Get(contents, "[0]", "carouselAdRenderer")

		if err == nil {
			index++
		} else {
			break
		}
	}

	_, err = jsonparser.ArrayEach(contents, func(value []byte, t jsonparser.ValueType, i int, err error) {
		if limit > 0 && len(results) >= limit {
			return
		}

		id, err := jsonparser.GetString(value, "videoRenderer", "videoId")
		if err != nil {
			return
		}

		title, err := jsonparser.GetString(value, "videoRenderer", "title", "runs", "[0]", "text")
		if err != nil {
			return
		}

		uploader, err := jsonparser.GetString(value, "videoRenderer", "ownerText", "runs", "[0]", "text")
		if err != nil {
			return
		}

		live := false
		duration, err := jsonparser.GetString(value, "videoRenderer", "lengthText", "simpleText")

		if err != nil {
			duration = ""
			live = true
		}

		results = append(results, &SearchResult{
			Title:       title,
			Uploader:    uploader,
			RawDuration: duration,
			ID:          id,
			URL:         fmt.Sprintf("https://youtube.com/watch?v=%s", id),
			Live:        live,
			Thumbnail:   fmt.Sprintf("https://i1.ytimg.com/vi/%s/hqdefault.jpg", id),
		})
	})

	return
}
