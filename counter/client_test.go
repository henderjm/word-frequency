package counter_test

import (
	"net/http"

	"github.com/jarcoal/httpmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/henderjm/word-frequency/counter"
)

var _ = Describe("Client", func() {

	BeforeEach(func() {
		httpmock.Reset()
		validEP := "https://en.wikipedia.org/w/api.php?action=query&explaintext=true&format=json&pageids=some-id&prop=extracts"
		invalidEP := "https://en.wikipedia.org/w/api.php?action=query&explaintext=true&format=json&pageids=InV72alid&prop=extracts"
		missingEP := "https://en.wikipedia.org/w/api.php?action=query&explaintext=true&format=json&pageids=1234567890&prop=extracts"
		wr := counter.WikiResponse{
			"batch is complete",
			counter.Query{Pages: []byte(`{
													"some-id": { 
													"title": "test title",
													"extract":"test extract"
													}
													}`,
			)},
		}
		errWR := counter.ErrorMediaWiki{
			Err: counter.ErrInfo{
				Code: "badinteger",
				Info: "Invalid value \"87a90\" for integer parameter \"pageids\".",
			},
		}
		missingWR := counter.WikiResponse{
			"batch is complete",
			counter.Query{Pages: []byte(`{
													"some-id": { 
													"pageid": "1234567890",
													"missing":""
													}
													}`,
			)},
		}
		httpmock.RegisterResponder("GET", validEP,
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, wr)
				if err != nil {
					return httpmock.NewStringResponse(500, "Something went wrong"), nil
				}
				return resp, nil
			},
		)
		httpmock.RegisterResponder("GET", invalidEP,
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, errWR)
				if err != nil {
					return httpmock.NewStringResponse(500, "Something went wrong"), nil
				}
				return resp, nil
			},
		)
		httpmock.RegisterResponder("GET", missingEP,
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, missingWR)
				if err != nil {
					return httpmock.NewStringResponse(500, "Something went wrong"), nil
				}
				return resp, nil
			},
		)

	})

	Context("When article exists", func() {
		It("Should be able to unmarshal json response", func() {
			c := counter.NewClient("some-id")
			mw, err := c.FetchMediaWiki()
			Expect(err).ToNot(HaveOccurred())
			Expect(mw.Title).To(Equal("test title"))
			Expect(mw.Extract).To(Equal("test extract"))
		})
	})

	Context("When PageIDs is invalid", func() {
		It("Should return error", func() {
			c := counter.NewClient("InV72alid")
			_, err := c.FetchMediaWiki()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("error code: badinteger\ninfo: Invalid value \"87a90\" for integer parameter \"pageids\"."))
		})
	})

	Context("When extract is missing", func() {
		It("Should return error", func() {
			c := counter.NewClient("1234567890")
			_, err := c.FetchMediaWiki()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("no matching wiki for pageids `1234567890`"))
		})
	})
})
