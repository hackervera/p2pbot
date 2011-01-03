package main

import (
  "strings"
  "fmt"
  "net"
  "json"
)
func PingHandler(){ //triggered on CTCP ping reply
  for {
    res:=<-PingResponse
    host := strings.Split(res," from ",-1)[0]
    fmt.Println("Dialing",host)
    conn,err := net.Dial("tcp","",host+":7878")
    if err != nil {
      fmt.Println(err)
      break
    }
    fmt.Println("Success! Writing host to peer database")
    WritePeers([]string{host})
    message := make([]byte,1000)
    var num int
    num,err = conn.Write([]byte("i can haz peers"))
    num,err = conn.Read(message)
    if err != nil {
      fmt.Println(err)
    }
    if num < 1 {
      fmt.Println("No data received")
    }
    var peers PeerBlob
    fmt.Println("Received data:",string(message))
    err = json.Unmarshal(message[0:num],&peers)
    if err != nil {
      fmt.Println(err)
    }
    fmt.Println(peers)
    WritePeers(peers.Peers)
  }
}

func ListenCall(l *net.TCPListener){
  for {
    conn,_:=l.Accept()
    message := make([]byte, 1000)
    _,err := conn.Read(message)
    if err != nil {
      fmt.Println(err)
    }
    if strings.Contains(string(message),"peer"){
      data := &PeerBlob{GetPeers(),"good"}
      //marshaldata := make([]byte, 1000)
      var marshaldata []byte
      marshaldata,err = json.Marshal(data)
      if err != nil {
        fmt.Println(err)
      }
      conn.Write([]byte(marshaldata))
    } else {
      conn.Write([]byte("no dataz for you!"))
    }
  }
}

