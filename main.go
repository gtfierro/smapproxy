package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/goji/httpauth"
	"github.com/parnurzeal/gorequest"
	"html/template"
	"log"
	"net/http"
	"time"
)

var port = flag.String("p", "1212", "Port to serve")
var user = flag.String("user", "", "HTTP basic auth user")
var pass = flag.String("pass", "", "HTTP basic auth pass")
var archiver = flag.String("archiver", "http://localhost:8079/api/query", "Query URL to proxy")

func index(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	template, err := template.ParseFiles("index.template")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	// get a token
	h := md5.New()
	seed := make([]byte, 16)
	binary.PutVarint(seed, time.Now().UnixNano())
	h.Write(seed)
	token := fmt.Sprintf("%x", h.Sum(nil))
	template.Execute(w, token)
}

func query(w http.ResponseWriter, req *http.Request) {
	var (
		err       error
		errstring = ""
		input     = make(map[string]string)
	)
	defer req.Body.Close()
	query := req.FormValue("query")
	input["query"] = query
	fmt.Println(query)

	request := gorequest.New()
	resp, _, errs := request.Post(*archiver).Type("text").Send(query).End()
	defer resp.Body.Close()
	if errs != nil {
		err = errs[0]
		fmt.Println(err)
		input["error"] = err.Error()
		render(w, input)
		return
	}

	var v interface{}
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	err = dec.Decode(&v)
	if err != nil {
		fmt.Println(err)
		input["error"] = err.Error()
		render(w, input)
		return
	}

	querybytes, err := json.MarshalIndent(v, "", "   ")
	if err != nil {
		fmt.Println(err)
		input["error"] = err.Error()
		render(w, input)
		return
	}

	h := md5.New()
	seed := make([]byte, 16)
	binary.PutVarint(seed, time.Now().UnixNano())
	h.Write(seed)
	token := fmt.Sprintf("%x", h.Sum(nil))

	if err != nil {
		errstring = err.Error()
	}
	input["token"] = token
	input["query"] = query
	input["result"] = string(querybytes)
	input["error"] = errstring
	render(w, input)
	return
}

func render(w http.ResponseWriter, input map[string]string) {
	template, err := template.ParseFiles("query.template")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	template.Execute(w, map[string]interface{}{"map": input})
}

func main() {
	flag.Parse()

	http.Handle("/", httpauth.SimpleBasicAuth(*user, *pass)(http.HandlerFunc(index)))
	http.Handle("/query", http.HandlerFunc(query))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Printf("Serving on %s...\n", ":"+*port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
