package main

import (
  //"strings"
  "fmt"
  "net"
  "json"
  "time"
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
      time.Sleep(1)
    }
  }
}

type UDPresponse struct {
  N int
  Buf [1000]byte
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

  var buf [1000]byte
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
    
    sendpacket := &Packet{"peers",GetPeers(),nil}
    jsonbuf,_ := json.Marshal(sendpacket)
    _,err := reply.Con.WriteTo(jsonbuf,addr)
    if err != nil {
      fmt.Println("writing error:",err)
    }
    var packet *Packet
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
      

