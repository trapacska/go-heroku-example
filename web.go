package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"bytes"
	"io/ioutil"
	
)

func main() {
	http.HandleFunc("/", hello)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Println("--------orig req----------")
	b, err := httputil.DumpRequest(req, true)
	fmt.Println(string(b), err)
	fmt.Println("------------------")

	body, _ := ioutil.ReadAll(req.Body)

	newUrl := *req.URL
	newUrl.Host = "developerservices2.apple.com"
	newUrl.Scheme = "https"
	fmt.Println("## to:", newUrl.String())
	fmt.Println("-------proxyreq-----------")
	proxyReq, err := http.NewRequest(req.Method, newUrl.String(), bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
		return
	}
	// We may want to filter some headers, otherwise we could just use a shallow copy
	// proxyReq.Header = req.Header
	proxyReq.Header = make(http.Header)
	for h, val := range req.Header {
		proxyReq.Header[h] = val
	}

	b, err = httputil.DumpRequest(proxyReq, true)
	fmt.Println(string(b), err)

	fmt.Println("--------resp----------")
	resp, err := (&http.Client{}).Do(proxyReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err = httputil.DumpResponse(resp, true)
	fmt.Println(string(b), err)

	fmt.Println("------------------")
	//w.Header = make(http.Header)
	for h, val := range resp.Header {
		w.Header()[h] = val
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	w.Write(respBody)
}
