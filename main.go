package main

import (
	"crawler/models"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

func search(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	w.Header().Set("Access-Control-Allow-Origin", "*")
	urlParam := r.URL.Query().Get("url")
	url, _ := url.Parse(urlParam)
	depthParam := r.URL.Query().Get("depth")
	threadsParam := r.URL.Query().Get("threads")
	depth := 2
	threads := 2

	if depthParam != "" {
		depth, _ = strconv.Atoi(depthParam)
		depth++
	}

	if threadsParam != "" {
		threads, _ = strconv.Atoi(threadsParam)
	}

	c := colly.NewCollector(
		colly.MaxDepth(depth),
		colly.AllowedDomains(url.Host),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: threads})
	response := &models.Result{ErrorLinks: []models.Link{}, VisitedLinks: []models.Link{}}

	c.OnHTML("a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		link := models.Link{RealURL: r.URL.String(), Depth: r.Depth - 1}
		response.VisitedLinks = append(response.VisitedLinks, link)
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		if r.StatusCode == 404 {
			link := models.Link{RealURL: r.Request.URL.String(), Depth: r.Request.Depth - 1}
			response.ErrorLinks = append(response.ErrorLinks, link)
		}
	})

	startTime := time.Now()
	c.Visit(urlParam)
	c.Wait()
	runtime.ReadMemStats(&m)
	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()
	response.Duration = strconv.FormatFloat(duration, 'f', 6, 64)
	response.MemoryUsage = strconv.FormatFloat(float64(m.HeapAlloc/(1024*1024)), 'f', 2, 64) + "MB"
	response.VisitedLinks = removeDuplicates(response.VisitedLinks)
	response.VisitedLinksCount = len(response.VisitedLinks)
	responseBytes, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
	r.Body.Close()
}

func main() {
	portFlag := flag.Int("p", 9090, "A port")
	flag.Parse()
	port := strconv.Itoa(*portFlag)
	http.DefaultTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   0,
			KeepAlive: 0,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          0,
		IdleConnTimeout:       0,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	server := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Minute,
		Addr:         ":" + port,
	}
	http.HandleFunc("/", search)
	log.Fatal(server.ListenAndServe())
}

func removeDuplicates(elements []models.Link) []models.Link {
	encountered := map[string]bool{}
	result := []models.Link{}

	for v := range elements {
		if encountered[elements[v].RealURL] != true {
			encountered[elements[v].RealURL] = true
			result = append(result, elements[v])
		}
	}
	return result
}
