package main

import (
  "http"
  //"fmt"
  "websocket"
  "flag"
  //"json"

)


//channels
var PingResponse = make(chan string)
var messageChan = make(chan []byte)
var ircChan = make(chan []byte)
var subscriptionChan = make(chan subscription)
var quit = make(chan int)

//global config vars
var hasUsername int
var myUsername string
var websocketPort string
var records []record

//structs
type subscription struct {
    conn      *websocket.Conn
    subscribe bool
}

type Packet struct {
  Type string
  Peers []string
  Tweet *Tweet
  Name string
}

type record struct {
  Author string
  Message string
  Timestamp string
}

type Tweet struct{
  Name string
  Message string
  Timestamp string
}


  
func main(){
  
  SetupDatabase()
  go WriteTweet()
  //go WholeThing()
  go UDPServer()
  go ProcessUDP()
  go ircStuff()
  go TweetSender()
  flag.Parse()
  
  

  http.HandleFunc("/", Index)
  http.HandleFunc("/tweetgrabber", TweetGrabber)
  http.HandleFunc("/insert", InsertTweet)
  http.ListenAndServe(":9999", nil)
}
