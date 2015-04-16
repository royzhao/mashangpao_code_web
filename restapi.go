package main
import (
"bytes"
"fmt"
"log"
"net/http"
"net/url"
)
type Rest struct {
// Header *http.Header
}
func parse_url(r_url string) *url.URL {
	uri, err := url.Parse(r_url)
	if err != nil {
	log.Fatal(err)
	}
	return uri
}
func (r *Rest) Get(r_url string) string {
	fmt.Println("----In restapi Get----")
	uri := parse_url(r_url)
	urlStr := fmt.Sprintf("%v", uri)
	rq, _ := http.NewRequest("GET", urlStr, nil)
	rq.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(rq)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
	}
	return "true"
}
func (r *Rest) Post(r_url string) string {
	fmt.Println("----In restapi Post----")
	uri := parse_url(r_url)
	fmt.Printf("%#v", uri)
	return "true"
}
func (r *Rest) Put(r_url string) string {
	fmt.Println("----In restapi Put----")
	uri := parse_url(r_url)
	fmt.Printf("%#v", uri)
	return "true"
}
func (r *Rest) Delete(r_url string) string {
	fmt.Println("----In restapi Delete----")
	uri := parse_url(r_url)
	fmt.Printf("%#v", uri)
	return "true"
}