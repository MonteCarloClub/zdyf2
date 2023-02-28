package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	gPort *int
)

func init() {
	gPort = flag.Int("port", 9001, "ca server port.")
	flag.Parse()
}

func SingleGenerate(w http.ResponseWriter, r *http.Request) {
	postData, _ := ioutil.ReadAll(r.Body)
	sign := Generate(string(postData))
	bData, _ := json.Marshal(sign)
	_, _ = fmt.Fprintf(w, string(bData))
}

func SingleVerify(w http.ResponseWriter, r *http.Request) {
	postData, _ := ioutil.ReadAll(r.Body)

	var sign *ABSSignature
	if err := json.Unmarshal(postData, &sign); err != nil {
		http.Error(w, "parse response error.", 500)
		return
	}
	if verify := Verify(sign); verify {
		http.Error(w, "OK.", 200)
	} else {
		http.Error(w, "Failed to Verify.", 500)
	}
}

func main() {
	http.HandleFunc("/SingleGenerate", SingleGenerate)
	http.HandleFunc("/SingleVerify", SingleVerify)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *gPort), nil); err != nil {
		log.Fatalln(err)
	}
}
