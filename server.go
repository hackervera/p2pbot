package main

import (
  "http"
  "fmt"
  "websocket"
  "flag"
  "json"

)


//channels
var PingResponse = make(chan string)
var messageChan = make(chan []byte)
var ircChan = make(chan []byte)
var subscriptionChan = make(chan subscription)
var quit = make(chan int)

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


func Subscribe(ws *websocket.Conn){ //called from main() on every websocket opened, also monitors input from webclient
  subscriptionChan <- subscription{ws, true} // add this channel to multiplexer
  buf := make([]byte, 10000)
  for {
    n, err := ws.Read(buf)
    if err != nil {
      fmt.Println(err)
      break
    }
    message := buf[0:n]
    fmt.Println("From webclient:",string(message))
    var packet *Packet
    err = json.Unmarshal(message,&packet)
    if err != nil {
      fmt.Println(err)
    }
    if packet.Type == "username" {
      myUsername = packet.Name
    } else if packet.Type == "tweet" {
      TweetChan <- packet.Tweet
    }
  }
}

func Multiplex(){ // handles websocket subscriptions, and messages sent to webclient
  conns := make(map[*websocket.Conn]int)
  for {
    select {
    case subscription := <-subscriptionChan:
      fmt.Println("got subscription")
      conns[subscription.conn] = 1
    case message := <-messageChan: 
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
  go WriteTweet()
  go WholeThing()
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
