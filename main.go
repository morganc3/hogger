package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var payloads map[string]string

// Any HTTP method to /redirect?url=<REDIRECT_URL>&status=3XX with params in body or URL
func Redirect(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	statusStr := r.Form.Get("status")
	url := r.Form.Get("url")

	//use 302 if none is provided or invalid code provided
	//redirects to itself if no url is provided

	if len(statusStr) > 0 && statusStr[:1] == "3" {
		status, err := strconv.Atoi(statusStr)
		if err != nil {
			status = 302
		}

		http.Redirect(w, r, url, status)
	} else {
		http.Redirect(w, r, url, 302)
	}
}

// Any HTTP method to /echo with echo=<reflected string here!> in body or URL param
func Echo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	echoStr := r.Form.Get("echo")
	w.Write([]byte(echoStr))
}

func StorePayload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	payloadStr := r.Form.Get("v")
	keyStr := r.Form.Get("k")
	payloads[keyStr] = payloadStr
	//delete payload after 2 minutes to save memory
	go deleteFromMap(keyStr)

}

func GetPayload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	key := r.Form.Get("k")
	if _, ok := payloads[key]; ok {
		w.Write([]byte(payloads[key]))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		//open CORS policy 
		w.Header().Set("Access-Control-Allow-Origin", "*")

		statusStr := r.Form.Get("status")
		headers := strings.Split(r.Form.Get("headers"), ",")

		for _, h := range headers {
			headerAndValue := strings.Split(h, ":")
			if len(headerAndValue) != 2 {
				continue
			}
			w.Header().Set(headerAndValue[0], headerAndValue[1])
		}

		//handle status differently for redirect endpoint
		if r.URL.Path != "/redirect" {
			var status int
			var err error
			if len(statusStr) > 0 {
				status, err = strconv.Atoi(statusStr)
				if err != nil {
					status = 200
				}
				if status != 200 {
					w.WriteHeader(status)
				}
			}
		}
		next.ServeHTTP(w, r)

	})
}

//sleep then delete from map
func deleteFromMap(key string) {
	time.Sleep(time.Minute * 2)
	delete(payloads, key)
}

func main() {
	payloads = make(map[string]string)

	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	router.HandleFunc("/redirect", Redirect)
	router.HandleFunc("/echo", Echo)
	router.HandleFunc("/store", StorePayload)
	router.HandleFunc("/p", GetPayload)
	log.Fatal(http.ListenAndServe(":80", router))
}
