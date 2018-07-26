package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"os"

	"github.com/nlopes/slack"
)

type ContentBody struct {
	Message string `json:"msg"`
}

var (
	slackToken     = os.Getenv("SLACK_TOKENE")
	slackChannelID = os.Getenv("SLACK_CHANNEL_ID")
)

func sendSlackMessage(msg string) error {
	api := slack.New(slackToken)
	api.SetDebug(false)
	params := slack.PostMessageParameters{}
	channelID, timestamp, err := api.PostMessage(slackChannelID, msg, params)
	if err != nil {
		log.Printf("%s\n", err)
		return fmt.Errorf("Failed to send message: %s", err)
	}
	log.Printf("[SLACK] Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var c ContentBody
		if err := json.Unmarshal(body, &c); err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		if len(strings.TrimSpace(c.Message)) == 0 {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		log.Printf("[POST] send message %s to slack...\n", c.Message)
		if err := sendSlackMessage(c.Message); err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, "success")
	}
}

func getMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg := os.Getenv("SIMPLE_MSG")
	if msg == "" {
		msg = ":( SIMPLE_MSG variable not defined"
	}
	fmt.Fprintf(w, "<h1>%s</h1>", msg)
}

func main() {
	addr := flag.String("addr", ":8080", "listen address and port.")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/say", sendMessageHandler)
	mux.HandleFunc("/", getMessageHandler)

	log.Printf("Knative slack app start to listen on %s ...", *addr)
	http.ListenAndServe(*addr, mux)
}