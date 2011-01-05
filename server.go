package main

import (
  "http"
  "fmt"
  //"strings"
  "websocket"
  "time"
  "json"
  //"io"
  //"os"
  //"bytes"
  "flag"
  //"net"
  //sqlite "gosqlite.googlecode.com/hg/sqlite"

)


//channels
var PingResponse = make(chan string)
var messageChan = make(chan []byte)
var ircChan = make(chan []byte)
var subscriptionChan = make(chan subscription)

//global config vars
var portNumber string = "9999"
var hostName string = "localhost"
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


func genMessage(text string) string { //create javascript to send to websocket client
  message := `var pre = document.createElement("p");
  pre.style.wordWrap = "break-word";
  pre.innerHTML = "`+text+`";
  output.insertBefore(pre,output.childNodes[0]);`
  return message
}


func Subscribe(ws *websocket.Conn){ //message to bot from websocket clientside
  subscriptionChan <- subscription{ws, true} // add this channel to multiplexer
  fmt.Println("just sent subscription message to channel")
  for {
    buf := make([]byte, 1000)
    n, err := ws.Read(buf)
    if err != nil {
      fmt.Println(err)
      break
    }
    msg := buf[0:n]
    var incoming struct {
      Type string
      Msg string
      Name string
    }
    json.Unmarshal(msg,&incoming)
    if incoming.Type == "update" {
      timestamp := time.LocalTime().String()
      //fmt.Println(timestamp)
      tweet := &Tweet{incoming.Name, incoming.Msg, timestamp}
      TweetChan <- tweet
      WriteTweet(tweet)
      //BroadcastTweets()
      //SendTweet(*tweet)
    }
    if incoming.Type == "username" {
      myUsername = incoming.Name
      WriteName(incoming.Name)
    }
  }
}

func Multiplex(){ // handles websocket connections
  conns := make(map[*websocket.Conn]int)
  for {
    select {
    case subscription := <-subscriptionChan:
      fmt.Println("got subscription")
      conns[subscription.conn] = 1
    case message := <-messageChan: // to websocket client from bot
      
      fmt.Println("got message:", string(message))
      
      for conn, _ := range conns {
        if _, err := conn.Write(message); err != nil {
          conn.Close()
        }
      }
    }
  }
}



  
func main(){
  
  SetupDatabase()
  go UDPServer()
  go ProcessUDP()
  go Multiplex()
  go ircStuff()
  go TweetSender()
  
  //go GetDataFromPeers()
  
  flag.StringVar(&portNumber,"port", "9999", "port number for web client")
  flag.StringVar(&hostName,"hostname", "localhost", "hostname for web client")
  flag.Parse()
  
  

  http.Handle("/socket", websocket.Handler(func(ws *websocket.Conn){ Subscribe(ws) }))
  http.HandleFunc("/",hello)
  http.ListenAndServe(":"+portNumber, nil)
}
