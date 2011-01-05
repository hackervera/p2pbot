package main

import (
  //"strings"
  "fmt"
  "net"
  "json"
  "time"
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
      fmt.Println(packet.Peers)
      WritePeers(packet.Peers) 
      
    }
  }
}

func BroadcastTweets(){
  peers := GetPeers()
  tweets := GetTweets()
  for _,peer := range peers {
    fmt.Println("Dialing",peer,"...")
    conn,err := net.Dial("udp","",peer+":7878")
    
    if err != nil {
      fmt.Println(err)
      break
    }
    conn.SetReadTimeout(1e9)
    //fmt.Println(conn)
    peerpacket := &Packet{"tweets",nil,tweets}
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
      for _,v := range packet.Tweets {
        WriteTweet(v)
        messageChan <- []byte("window.location.reload();")
      }
      time.Sleep(1)
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

func ProcessUDP(){
  for {
    reply :=<- UDPchan
    if reply.Err != nil {
      fmt.Println("Error while reading from UDP:",reply.Err)
      os.Exit(1)
    }
    buf := reply.Buf
    n := reply.N
    addr := reply.Addr
    fmt.Println("read",reply.N,"bytes")
    fmt.Println("addr:",reply.Addr)
    fmt.Println("Incoming message:",string(buf[0:n]))
    
    var packet *Packet
    err := json.Unmarshal(buf[0:n],&packet)
    if err != nil {
      fmt.Println(err)
    }
    var sendpacket *Packet
    if packet.Type == "peers" {
      WritePeers(packet.Peers)
      sendpacket = &Packet{"peers",GetPeers(),nil}
    } else {
      sendpacket = &Packet{"peers",GetPeers(),nil}
      for _,v := range packet.Tweets {
        WriteTweet(v)
        messageChan <- []byte("window.location.reload();")
      }
    }
    jsonbuf,_ := json.Marshal(sendpacket)
    _,err = reply.Con.WriteTo(jsonbuf,addr)
    if err != nil {
      fmt.Println("writing error:",err)
    }
   
  }
}

