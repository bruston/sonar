package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	domain := flag.String("d", "", "domain to find subdomains for")
	noDupes := flag.Bool("nodupes", true, "remove duplicate hosts if true")
	flag.Parse()
	if *domain == "" {
		fmt.Println("must specify a domain with the -d flag")
		return
	}
	const baseURL = "https://sonar.omnisint.io/subdomains/"
	client := &http.Client{Timeout: time.Minute}
	dupes := make(map[string]struct{})
	page := 0
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s?page=%d", baseURL, *domain, page), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing URL: %s\n", err)
			os.Exit(1)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error connecting to omnisint: %s\n", err)
			os.Exit(1)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading response body: %s\n", err)
			os.Exit(1)
		}
		resp.Body.Close()
		if bytes.HasPrefix(b, []byte("null")) {
			break
		}
		var hosts []string
		if err := json.Unmarshal(b, &hosts); err != nil {
			fmt.Fprintf(os.Stderr, "omnisint responded with non-JSON data: %s\n", err)
			os.Exit(1)
		}
		if len(hosts) == 0 {
			break
		}
		for _, v := range hosts {
			if *noDupes {
				if _, ok := dupes[v]; !ok {
					fmt.Println(v)
					dupes[v] = struct{}{}
				}
				continue
			}
			fmt.Println(v)
		}
		page++
	}
}
