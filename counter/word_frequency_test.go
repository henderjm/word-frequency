package counter_test

import (
	"io"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/jarcoal/httpmock"

	"github.com/henderjm/word-frequency/counter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("WordFrequency", func() {

	It("Should initialise with correct values", func() {
		nfc := counter.NewWordFrequencyCounter(1, "some-id")
		Expect(nfc.NumberOfWords).To(Equal(1))
		Expect(nfc.PageIDs).To(Equal("some-id"))
	})

	Context("When collecting word frequency", func() {

		BeforeEach(func() {
			httpmock.Reset()
			endPoint := "https://en.wikipedia.org/w/api.php?action=query&explaintext=true&format=json&pageids=some-id&prop=extracts"
			wr := counter.WikiResponse{
				"batch is complete",
				counter.Query{Pages: []byte(`{
													"some-id": { 
													"title": "workday is funday",
													"extract":"After after squeezing the chicken cracker, flavor chicken lard, cracker crumps and whipped chicken cream with it in a lard jar."
													}
													}`,
				)},
			}
			httpmock.RegisterResponder("GET", endPoint,
				func(req *http.Request) (*http.Response, error) {
					resp, err := httpmock.NewJsonResponse(200, wr)
					if err != nil {
						return httpmock.NewStringResponse(500, "Something went wrong"), nil
					}
					return resp, nil
				},
			)

		})

		DescribeTable("collect word frequencies",
			func(n int, pageids string, expected counter.TupleList) {
				wf := counter.NewWordFrequencyCounter(n, pageids)
				wf.Run()
				Expect(wf.Tuples).To(ConsistOf(expected))
			},
			Entry("n is 1", 1, "some-id", counter.TupleList{
				{Key: 3, Value: []string{"chicken"}},
			}),
			Entry("multiple same frequencies", 2, "some-id", counter.TupleList{
				{Key: 3, Value: []string{"chicken"}},
				{Key: 2, Value: []string{"after", "cracker", "lard"}},
			}),
			Entry("n > number of available words", math.MaxInt64, "some-id", counter.TupleList{
				{Key: 3, Value: []string{"chicken"}},
				{Key: 2, Value: []string{"after", "cracker", "lard"}},
				{Key: 1, Value: []string{"cream", "crumps", "flavor", "funday", "squeezing", "whipped", "with", "workday"}},
			}),
		)
	})

	Context("When wiki article does not exist", func() {

		BeforeEach(func() {
			httpmock.Reset()
			endPoint := "https://en.wikipedia.org/w/api.php?action=query&explaintext=true&format=json&pageids=some-id&prop=extracts"
			wr := counter.WikiResponse{
				"batch is complete",
				counter.Query{Pages: []byte(`{
													"some-id": { 
													"title": "workday is funday",
													"extract":"After after squeezing the chicken cracker, flavor chicken lard, cracker crumps and whipped chicken cream with it in a lard jar."
													}
													}`,
				)},
			}
			httpmock.RegisterResponder("GET", endPoint,
				func(req *http.Request) (*http.Response, error) {
					resp, err := httpmock.NewJsonResponse(200, wr)
					if err != nil {
						return httpmock.NewStringResponse(500, "Something went wrong"), nil
					}
					return resp, nil
				},
			)

		})
		It("Should return error for invalid pageids semantic", func() {

		})
	})

	Context("When writing", func() {
		It("Should write to std out", func() {
			wf := counter.NewWordFrequencyCounter(1, "some-id")
			wf.Tuples = counter.TupleList{
				{Key: 3, Value: []string{"chicken"}},
				{Key: 2, Value: []string{"after", "cracker", "lard"}},
			}
			wf.Title = "Test Title"
			r, w := io.Pipe()
			go func() {
				wf.Write(w)
				w.Close()
			}()
			stdout, _ := ioutil.ReadAll(r)
			Expect(string(stdout)).To(Equal("URL: https://en.wikipedia.org/w/api.php?action=query&explaintext=true&format=json&pageids=some-id&prop=extracts\n" +
				"Title:Test Title\n" +
				"Top 1 word(s):\n" +
				"- 3 chicken\n" +
				"- 2 after, cracker, lard\n"))
		})
	})

})
