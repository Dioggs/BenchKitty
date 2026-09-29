package benchkitty

import (
	"flag"
	"fmt"
	"slices"
)

type Config struct {
	URL      string
	Method   string
	Delay    int
	ReqCount int
	Out      string
}

var methods = []string{
	"GET",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
}

func isValidHTTPMethod(s string) bool {
	return slices.Contains(methods, s)
}

func Parse() (Config, bool) {
	reqCount := flag.Int("r", 100, "request amount")
	delay := flag.Int("d", 1000, "delay between every request call in ms")
	out := flag.String("o", "", "output path for the benchmark xlsx (defaults to terminal)")
	method := flag.String("t", "GET", "http method used on the url")

	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		fmt.Println("Missing url")
		flag.Usage()
		return Config{}, false
	}

	if !isValidHTTPMethod(*method) {
		fmt.Printf("Invalid value %v for command param -t\n", *method)
		return Config{}, false
	}

	return Config{
		URL:      url,
		Method:   *method,
		Delay:    *delay,
		ReqCount: *reqCount,
		Out:      *out,
	}, true
}
