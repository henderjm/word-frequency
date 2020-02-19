package counter

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type wordFrequency struct {
	NumberOfWords int
	PageIDs       string
	Client        Client
	Tuples        TupleList
	Title         string
}

func NewWordFrequencyCounter(n int, id string) *wordFrequency {
	return &wordFrequency{
		NumberOfWords: n,
		PageIDs:       id,
		Client:        NewClient(id),
	}
}

func (wf *wordFrequency) Run() error {
	if wf.NumberOfWords == 0 {
		return nil
	}

	mw, err := wf.Client.FetchMediaWiki()
	if err != nil {
		return err
	}
	dict, err := createFrequncyDictionary(mw)
	if err != nil {
		return err
	}

	freqs := SortMapByKey(dict)
	var result TupleList
	reported := 0
	pKey := 0
	for _, f := range freqs {
		if f.Key < pKey {
			reported += 1
		}
		if reported == wf.NumberOfWords {
			break
		}
		pKey = f.Key
		result = append(result, f)
	}

	wf.Tuples = result
	wf.Title = mw.Title

	return nil
}

func (wf *wordFrequency) Write(w io.Writer) {
	fmt.Fprintf(w, "URL: %s\n", wf.Client.URL().URL.String())
	fmt.Fprintf(w, "Title:%s\n", wf.Title)
	fmt.Fprintf(w, "Top %d word(s):\n", wf.NumberOfWords)
	for _, t := range wf.Tuples {
		fmt.Fprintf(w, "- %d %s\n", t.Key, strings.Join(t.Value, ", "))
	}
}

func createFrequncyDictionary(mw MediaWiki) (map[int][]string, error) {
	wordToFreq := make(map[string]int)
	reg, err := regexp.Compile(`\w{4,}`)
	if err != nil {
		return nil, err
	}

	matches := reg.FindAllString(mw.Extract, -1)
	matches = append(matches, reg.FindAllString(mw.Title, -1)...)
	for _, m := range matches {
		lm := strings.ToLower(m)
		if _, ok := wordToFreq[lm]; !ok {
			wordToFreq[lm] = 0
		}
		wordToFreq[lm] += 1
	}

	freqToWord := make(map[int][]string)
	for word, frequency := range wordToFreq {
		if _, ok := freqToWord[frequency]; !ok {
			freqToWord[frequency] = []string{}
		}
		freqToWord[frequency] = append(freqToWord[frequency], word)
	}
	return freqToWord, nil
}
