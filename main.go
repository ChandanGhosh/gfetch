package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	// here using pflag instead of builtin flag as pflag is a drop-in
	// replacement of flag and has more options of capturing
	// slices from cli arguments
	flag "github.com/spf13/pflag"
)

func usage() {
	fmt.Printf(`fetch - a simple commandline url tester.
Usage: fetch [options] args

Valid options:
	`)
	flag.PrintDefaults()
}

func main() {

	var urls *[]string
	urls = flag.StringSliceP("url", "u", []string{"http://google.com", "http://msn.com"}, "enter urls seperated by comma")
	flag.Usage = usage
	flag.Parse()

	maxLen := findLongestUrl(urls)

	t := time.Now()
	ch := make(chan string)
	if len(*urls) == 0 {
		fmt.Println("Enter some urls to fetch seperated by space")
		os.Exit(1)
	}

	for _, url := range *urls {
		// using goroutine for concurrent executions
		// channels for communication across goroutines
		// in non-blocking way
		go fetch(url, maxLen, ch)
	}

	for i := 0; i < len(*urls); i++ {
		fmt.Println(<-ch)
	}

	fmt.Println(fmt.Sprintf("Total Time taken: %.2fs", time.Since(t).Seconds()))
}

func fetch(url string, maxLen int, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("url: %s\t error: %v", url, err)
		return
	}
	defer resp.Body.Close()

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("url: %s\t error: %v", url, err)
		return
	}
	t := time.Since(start).Seconds()
	// this formatting line is important, here we are dynamically adding formatting
	// place value for url.
	ch <- fmt.Sprintf("timeTaken: %.2fs\turl: %+*s\tdata: %v bytes", t, maxLen, url, nbytes)
}

func findLongestUrl(urls *[]string) int {
	var c int
	for _, url := range *urls {
		if len(url) > c {
			c = len(url)
		}
	}
	return c
}
