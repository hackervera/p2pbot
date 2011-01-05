package main

import (
  //"strings"
  "fmt"
  "net"
  "json"
  //"time"
  "os"
)


func BroadcastPeers(){
  peers := GetPeers()
  for _,peer := range peers {
    fmt.Println("Dialing",peer,"...")
    conn,err := net.Dial("udp","",peer+":7878")
    if err != nil {
      fmt.Println(err)
      break
    }
    conn.SetReadTimeout(1e9)
    Connections <- conn
    //fmt.Println(conn)
    peerpacket := &Packet{"peers",peers,nil}
    jsonbuf,jerr := json.Marshal(peerpacket)
    if jerr != nil {
      fmt.Println(jerr)
    }
    conn.Write(jsonbuf)
    var buf [1000]byte
    size,readerr := conn.Read(buf[0:])
    if readerr != nil {
      fmt.Println("read error:",readerr)
    } else {
      
      var packet *Packet
      json.Unmarshal(buf[0:size],&packet)
      if packet.Type == "peers" {
        fmt.Println(packet.Peers)
        WritePeers(packet.Peers) 
      } else if packet.Type == "tweets" {
        fmt.Println(packet.Tweet)
        WriteTweet(packet.Tweet)
      }
      
    }
  }
}


type UDPresponse struct {
  N int
  Buf [10000]byte
  Con net.PacketConn
  Addr net.Addr
  Err os.Error
}

var UDPchan = make(chan *UDPresponse)

func UDPServer(){
  c,cerr := net.ListenPacket("udp", "0.0.0.0:7878")
  if cerr != nil {
    fmt.Println("Error while reading from UDP:",cerr)
    os.Exit(1)
  }

  var buf [10000]byte
  for {
    
    fmt.Println("Blocking: waiting on read")
    n, addr, aerr := c.ReadFrom(buf[0:])
    fmt.Println("Not blocking, read from connection")
    fmt.Println(n, addr, aerr)
    fmt.Println(buf[0:n])
    
    
    res := &UDPresponse{n,buf,c,addr,aerr}
    fmt.Println("Blocking: waiting for channel write")
    UDPchan <- res
    fmt.Println("Not Blocking, Wrote to channel")
  }
}

var Connections = make(chan net.Conn)
var TweetChan = make(chan *Tweet)

func TweetSender(){
  var err os.Error
  conns := make(map[net.Conn]int)
  for {
    select {
    case connection :=<- Connections:
      conns[connection] = 1
    case tweet :=<-TweetChan:
      tweetpacket := &Packet{"tweets",nil,tweet}
      var jsonbuf []byte
      
      jsonbuf,err = json.Marshal(tweetpacket)
      if err != nil {
        fmt.Println(err)
      }
      for conn,_ := range conns {
        conn.Write(jsonbuf)
      }
    }
  }
}

func ProcessUDP(){
  for {
    reply :=<- UDPchan
    if reply.Err != nil {
      fmt.Println("Error while reading from UDP:",reply.Err)
      os.Exit(1)
    }
    //Connections <- reply.Conn
    buf := reply.Buf
    n := reply.N
    //addr := reply.Addr
    fmt.Println("read",reply.N,"bytes")
    fmt.Println("addr:",reply.Addr)
    fmt.Println("Incoming message:",string(buf[0:n]))
    
    var packet *Packet
    err := json.Unmarshal(buf[0:n],&packet)
    if err != nil {
      fmt.Println(err)
    }
    if packet.Type == "peers" {
      WritePeers(packet.Peers)
      
    } else if packet.Type == "tweet" {
      WriteTweet(packet.Tweet)
    }
  }
}

