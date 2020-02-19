package counter

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

const wikiApi = "https://en.wikipedia.org/w/api.php?"

type Client struct {
	pageID string
}

type MediaWiki struct {
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

type PagesID map[string]MediaWiki

type WikiResponse struct {
	BatchComplete string `json:"batchcomplete"`
	Query         Query  `json:"query"`
}

type Query struct {
	Pages json.RawMessage
}

func NewClient(pageID string) Client {
	return Client{
		pageID,
	}
}

func (c Client) URL() *http.Request {
	req, err := http.NewRequest("GET", wikiApi, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("action", "query")
	q.Add("prop", "extracts")
	q.Add("explaintext", "true")
	q.Add("format", "json")
	q.Add("pageids", c.pageID)
	req.URL.RawQuery = q.Encode()
	return req
}

func (c Client) FetchMediaWiki() (MediaWiki, error) {
	req := c.URL()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return unmarshalWikiResponse(c.pageID, data)
}

func unmarshalWikiResponse(pageIDs string, data []byte) (MediaWiki, error) {
	var r WikiResponse
	err := json.Unmarshal(data, &r)
	if err != nil {
		return MediaWiki{}, err
	}
	if reflect.DeepEqual(Query{}, r.Query) {
		return MediaWiki{}, NewMediaWikiError(data)
	}

	dec := json.NewDecoder(strings.NewReader(string(r.Query.Pages)))
	var pages PagesID
	for {
		if err := dec.Decode(&pages); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	if pages[pageIDs].Extract == "" && pages[pageIDs].Title == "" {
		return MediaWiki{}, NewContentsMissingError(pageIDs)
	}
	return pages[pageIDs], err
}
