package main

import (
  //"strings"
  "fmt"
  "net"
  "json"
  "time"
  "os"
)

var Tweets = make(chan *Tweet)
var Relays = make(chan *net.UDPConn)
var Clients = make(chan *net.UDPConn)


func DialRelays(){ 
  relays := GetRelays()
  for _,relay := range relays {
    fmt.Println("Dialing",relay,"...")
    UDPAddr,_ := net.ResolveUDPAddr(relay+":7878")
    Local,_ := net.ResolveUDPAddr("")
    conn,err := net.DialUDP("udp", Local,UDPAddr)
    if err != nil {
      fmt.Println(err)
      continue // skip to next peer on connection error
    }
    jsonbuf,err := json.Marshal(relays)
    if err != nil {
      fmt.Println(err)
    }
    conn.Write(jsonbuf) // send peers to relay
    
    Relays <- conn 
  }
  time.Sleep(10e9)
}

func ListenClients(){ 
  UDPAddr,_ := net.ResolveUDPAddr("0.0.0.0:7878")
  c,err := net.ListenUDP("udp", UDPAddr) 
  var buf [10000]byte
  if err != nil {
    fmt.Println("Error while reading from UDP:",err)
    os.Exit(1)
  }
  
  for {
    n,_,_ := c.ReadFrom(buf[0:])
    if err != nil {
      fmt.Println("Error while reading from UDP:",err)
      os.Exit(1)
    }
    fmt.Println("Incoming: ", buf[0:n])
    Clients <- c
  }
}

func ConnectionMonitor(){ //multiplexer for client connections, tweets, and relay's peers
  conns := make(map[*net.UDPConn]string)
  for {
    select {
    case connection :=<- Relays: 
      conns[connection] = "relay"
    case connection :=<- Clients:
      conns[connection] = "client"
    case tweet :=<-Tweets: // send tweets to connections
      fmt.Println("incoming tweet")
      jsonbuf,err := json.Marshal(tweet)
      if err != nil {
        fmt.Println(err)
      }
      for conn,t := range conns {
        if t == "relay" {
          conn.Write(jsonbuf)
        } else {
          conn.WriteTo(jsonbuf, conn.RemoteAddr())
        }
      }
    }
  }
}
