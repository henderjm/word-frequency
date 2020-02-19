package counter

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Tuple struct {
	Key   int
	Value []string
}

type TupleList []Tuple

func (p TupleList) Len() int           { return len(p) }
func (p TupleList) Less(i, j int) bool { return p[i].Key > p[j].Key }
func (p TupleList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func SortMapByKey(m map[int][]string) TupleList {
	var tps TupleList
	for k, v := range m {
		sort.Strings(v)
		tps = append(tps, Tuple{k, v})
	}
	sort.Sort(tps)

	return tps
}

// Invalid pageids
type ErrorMediaWiki struct {
	Err ErrInfo `json:"error"`
}

type ErrInfo struct {
	Code string `json:"code"`
	Info string `json:"info"`
}

func NewMediaWikiError(data []byte) error {
	var emw ErrorMediaWiki
	err := json.Unmarshal(data, &emw)
	if err != nil {
		return err
	}
	return &emw
}

func (e *ErrorMediaWiki) Error() string {
	return fmt.Sprintf("error code: %s\ninfo: %s", e.Err.Code, e.Err.Info)
}

// Missing contents from request
// No matching id to mediawiki
func NewContentsMissingError(pageids string) error {
	return &errorString{fmt.Sprintf("no matching wiki for pageids `%s`", pageids)}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
