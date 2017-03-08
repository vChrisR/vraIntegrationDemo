package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"vraIntegrationDemo/config"

	"github.com/gorilla/mux"
)

type tLogContent struct {
	TimeStamp string `json:"timeStamp"`
	Message   string `json:"message"`
}

type tSiteContent struct {
	Title string
}

var (
	siteContent      tSiteContent
	logContent       tLogContent
	qChannel         chan chan string
	apiUser, apiPass string
)

func stepReceiver(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != apiUser || pass != apiPass {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	err = json.Unmarshal(body, &logContent)
	if err != nil {
		panic(err)
	}
	logContent.TimeStamp = time.Now().Format(time.RFC3339)

	var contentString []byte
	contentString, err = json.Marshal(logContent)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case responseChan := <-qChannel:
				responseChan <- string(contentString)
			default:
				return
			}
		}
	}()

	fmt.Println(string(contentString))
	w.WriteHeader(201)
}

func getStepLongPoll(w http.ResponseWriter, r *http.Request) {
	timeout, err := strconv.Atoi(r.URL.Query().Get("timeout"))
	if err != nil || timeout > 180000 || timeout < 0 {
		timeout = 60000 // default timeout is 60 seconds
	}

	responseChan := make(chan string)

	select {
	case qChannel <- responseChan:
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		return
	}
	io.WriteString(w, <-responseChan)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index").ParseFiles("index.tmpl")
	if err != nil {
		panic(err)
	}
	t.ExecuteTemplate(w, "index.tmpl", siteContent)
}

func main() {
	router := mux.NewRouter()
	qChannel = make(chan chan string)
	apiUser, apiPass = config.GetAPICreds()

	siteContent.Title = "vRA Integration Demo"

	router.HandleFunc("/api/step/latest/", getStepLongPoll).Methods("GET")
	router.HandleFunc("/api/step/", stepReceiver).Methods("POST")
	router.HandleFunc("/", rootHandler)
	http.Handle("/", router)
	http.ListenAndServe(":"+config.GetPort(), nil)
}
