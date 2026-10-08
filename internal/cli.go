package benchkitty

import (
	"flag"
	"fmt"
)

type Config struct {
	URL      string
	Method   string
	Body     string
	Delay    int
	ReqCount int
	Out      string
}

var methods = map[string]bool{
	"GET":    false,
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": false,
}

func Parse() (Config, bool) {
	reqCount := flag.Int("r", 100, "request amount")
	delay := flag.Int("d", 1000, "delay between every request call in ms")
	out := flag.String("o", "", "output path for the benchmark xlsx (defaults to terminal)")
	method := flag.String("t", "GET", "http method used on the url")
	body := flag.String("b", "", "request body (only for POST, PUT, PATCH)")

	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		fmt.Println("Missing url")
		flag.Usage()
		return Config{}, false
	}

	allowsBody, known := methods[*method]
	if !known {
		fmt.Printf("Invalid value %v for command param -t\n", *method)
		return Config{}, false
	}

	if *body != "" && !allowsBody {
		fmt.Printf("Method %v does not accept a body\n", *method)
		return Config{}, false
	}

	return Config{
		URL:      url,
		Method:   *method,
		Body:     *body,
		Delay:    *delay,
		ReqCount: *reqCount,
		Out:      *out,
	}, true
}
