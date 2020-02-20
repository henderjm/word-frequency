# word-frequency
## Usage
```
Usage:
  word-frequency [OPTIONS]

Application Options:
  -n, --number-of-words= Report a number of most common words on a wiki page; it must be positive
  -i, --page-ids=        page id of a wiki page

Help Options:
  -h, --help             Show this help message
```

## To Run from source
* Be set up to use go mod
* Golang 13.4.*
* Clone the repository to your local machine
* Change directory to repository
* run "go run main.go -n 5 -i 21721040"

## Use provided stable binary
* run "./word-frequency-osx -n 5 -i 21721040"

## Running tests
* Download [ginkgo](https://onsi.github.io/ginkgo/) binary `go get github.com/onsi/ginkgo/ginkgo` 
* `ginkgo -r ./`
