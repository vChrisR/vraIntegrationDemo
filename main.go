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
	redis "gopkg.in/redis.v5"
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
	redisPubClient   *redis.Client
)

func stepReceiver(w http.ResponseWriter, r *http.Request) {
	//authentication
	user, pass, ok := r.BasicAuth()
	if !ok || user != apiUser || pass != apiPass {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Get the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	//parse the body
	err = json.Unmarshal(body, &logContent)
	if err != nil {
		panic(err)
	}

	//Add timestamp
	logContent.TimeStamp = time.Now().Format(time.RFC3339)

	//convert back to json string
	var contentString []byte
	contentString, err = json.Marshal(logContent)
	if err != nil {
		panic(err)
	}

	//publish in redis
	err = redisPubClient.Publish("vraintegrationdemo", string(contentString)).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("pub: " + string(contentString))

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
	siteContent.Title = "vRA Integration Demo"
	fmt.Println("Starting ", siteContent.Title)

	redisPubClient = redisNewClient()
	redisSubClient := redisNewClient()
	fmt.Println("Connected to Redis!")

	sub, err := redisSubClient.Subscribe("vraintegrationdemo")
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	fmt.Println("Subscribing to Redis.")
	qChannel = make(chan chan string)
	go func() {
		for {
			message, suberr := sub.ReceiveMessage()
			fmt.Println("sub:" + message.Payload)
			if suberr != nil {
				panic(err)
			}

			select {
			case responseChan := <-qChannel:
				responseChan <- string(message.Payload)
			default:
			}
		}
	}()

	fmt.Println("Starting web server now.")
	apiUser, apiPass = config.GetAPICreds()
	router := mux.NewRouter()
	router.HandleFunc("/api/step/latest/", getStepLongPoll).Methods("GET")
	router.HandleFunc("/api/step/", stepReceiver).Methods("POST")
	router.HandleFunc("/", rootHandler)

	http.Handle("/", router)
	err = http.ListenAndServe(":"+config.GetPort(), nil)
	if err != nil {
		fmt.Println(err)
	}
}
