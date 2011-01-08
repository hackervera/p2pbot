package main

import (
  "http"
  //"fmt"
  //"flag"
  //"json"

)

//channels
var PingResponse = make(chan string)
var messageChan = make(chan []byte)
var ircChan = make(chan []byte)
var quit = make(chan int)

//global config vars
var hasUsername int
var myUsername string
var websocketPort string

//structs

type Tweet struct{
  Name string
  Message string
  Timestamp string
  Sig []byte
}

func main(){
  
  SetupDatabase()
  go TweetWriter()
  go ircStuff()
  go ConnectionMonitor()

  http.HandleFunc("/", Index)
  http.HandleFunc("/tweetgrabber", TweetGrabber)
  http.HandleFunc("/insert", InsertTweet)
  http.HandleFunc("/addressbook", AddressBook)
  http.HandleFunc("/addfriend", AddFriend)
  http.ListenAndServe(":9999", nil)
}
