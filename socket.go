package main

import (
  //"strings"
  "fmt"
  "net"
  "json"
  //"time"
  "os"
)

func BroadcastPeers(peers []string){
  for _,peer := range peers {
    fmt.Println("Dialing",peer,"...")
    conn,err := net.Dial("udp","",peer+":7878")
    if err != nil {
      fmt.Println(err)
      break
    }
    //fmt.Println(conn)
    peerpacket := &Packet{"peers",peers,nil}
    jsonbuf,jerr := json.Marshal(peerpacket)
    if jerr != nil {
      fmt.Println(jerr)
    }
    conn.Write(jsonbuf)
  }
}


func UDPServer(){
  c,cerr := net.ListenPacket("udp", "0.0.0.0:7878")
  if cerr != nil {
    fmt.Println("Error while reading from UDP:",cerr)
    os.Exit(1)
  }

  var buf [1000]byte
  for {
    n, addr, aerr := c.ReadFrom(buf[0:])
    if aerr != nil {
      fmt.Println("Error while reading from UDP:",aerr)
      os.Exit(1)
    }
    fmt.Println("read",n,"bytes")
    fmt.Println("addr:",addr)
    fmt.Println("Incoming message:",string(buf[0:n]))
    _,werr := c.WriteTo([]byte("Got your message"),addr)
    if werr != nil {
      fmt.Println("write error:",werr)
    }
    var packet Packet
    json.Unmarshal(buf[0:n],&packet)
    if packet.Type == "peers" {
      fmt.Println(packet.Peers)
      fmt.Println("Adding peers")
      WritePeers(packet.Peers)
    }
    
  }
}


func SendTweet(tweet Tweet){
  peers := GetPeers()
  for i:= range peers {
    ip:= peers[i]
    go func(ip string){
      fmt.Println("Dialing",ip)
      conn,err := net.Dial("tcp","",ip+":7878")
      if err != nil {
        fmt.Println(err)
        return
      }
      fmt.Println("Success! Sending tweet")
      conn.Write([]byte("I haz new tweet"))
      message,merr := json.Marshal(tweet)
      if merr != nil {
        fmt.Println(merr)
      }
      conn.Write(message)
      fmt.Println(string(message))
      fmt.Println("closing connection")
      conn.Close()
    }(ip)
  }
}
      

