package main
import (
  "http"
  "io/ioutil"
)

func Index(w http.ResponseWriter, r *http.Request){
  index,_ := ioutil.ReadFile("index.html")
  w.Write(index)
}

func TweetGrabber(w http.ResponseWriter, r *http.Request){
  w.Write([]byte(GetHistory()))
}

func InsertTweet(w http.ResponseWriter, r *http.Request){
  tweet := &Tweet{Name:"foo",Message:"bar",Timestamp:"tomorrow"}
  TweetWrite <- tweet
  w.Write([]byte("looks good"))
}
