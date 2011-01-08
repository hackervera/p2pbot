package main

import (
  //"strings"
  "fmt"
  "net"
  "json"
  "time"
  "os"
)

var Connections = make(chan *UDPresponse)
var TweetChan = make(chan *Tweet)
var Federation = make(chan net.Conn)
var UDPchan = make(chan *UDPresponse)
var NoSelf = make(chan int) // Don't send data to ourselves and create feedback loop

type UDPresponse struct {
  Buf []byte
  Con net.PacketConn
  Addr net.Addr
}



func BroadcastPeers(){ //Try to connect to relays and broadcast peers
  peers := GetPeers()
  for _,peer := range peers {
    fmt.Println("Dialing",peer,"...")
    conn,err := net.Dial("udp","",peer+":7878")
    if err != nil {
      fmt.Println(err)
      continue // skip to next peer
    }
    //conn.SetReadTimeout(1e9)
    
    //fmt.Println(conn)
    peerpacket := &Packet{Type:"peers",Peers:peers}
    var jsonbuf []byte
    jsonbuf,err = json.Marshal(peerpacket)
    if err != nil {
      fmt.Println(err)
    }
    conn.Write(jsonbuf) // send peers to relay
    
    Federation <- conn // send relay's net.Conn *interface* to federation channel. Read by TweetSender()
  }
  time.Sleep(10e9)
  NoSelf <- 1
}



func UDPServer(){ // This function makes the bot act as a relay. It means its a public interface.
  <-NoSelf
  c,err := net.ListenPacket("udp", "0.0.0.0:7878") // c is a net.PacketConn *interface* for the client connecting to this relay
  if err != nil {
    fmt.Println("Error while reading from UDP:",err)
    os.Exit(1)
  }
 
  var buf [10000]byte
  for {
    
    var n int
    var addr net.Addr // *interface* Network() network name, String() string address
    n, addr, err = c.ReadFrom(buf[0:])
    if err != nil {
      fmt.Println("Error while reading from UDP:",err)
      os.Exit(1)
    }

    res := &UDPresponse{buf[0:n],c,addr}
    UDPchan <- res // send client's connection information to UDPchan in ProcessUDP()
  }
}


func TweetSender(){ //multiplexer for client connections, tweets, and relay's peers
  var err os.Error
  conns := make(map[*UDPresponse]int)
  peers := make(map[net.Conn]int)
  for {
    select {
    case connection :=<- Connections: //adds client connection to map
      conns[connection] = 1
    case tweet :=<-TweetChan: // send tweets to clients
      fmt.Println("incoming tweet")
      tweetpacket := &Packet{Type:"tweet",Tweet:tweet}
      var jsonbuf []byte
      jsonbuf,err = json.Marshal(tweetpacket)
      if err != nil {
        fmt.Println(err)
      }
      for response,_ := range conns {
        fmt.Println("Writing",string(jsonbuf),"to",response.Con)
        response.Con.WriteTo(jsonbuf,response.Addr)
      }
      for peer,_ := range peers {
        peer.Write(jsonbuf)
      }
    case peer :=<- Federation: // receive relay's net.Conn interface from BroadcastPeers()
      peers[peer] = 1 // add relay net.Conn interface to peers map
      go Read(peer) // waits for data to be written to client from relay
    }
  }
}

func ProcessUDP(){
   
  for {
    reply :=<- UDPchan // client's connection information
    Connections <- reply //send client's connection information to Connections in TweetSender()
    buf := reply.Buf
    fmt.Println("Client(",reply.Addr,")","just sent:",string(buf))
    
    var packet *Packet
    err := json.Unmarshal(buf,&packet) // unmarshal client's sent json to Packet
    if err != nil {
      fmt.Println(err)
      continue
    }
    if packet.Type == "peers" {
      WritePeers(packet.Peers)
      
    } else if packet.Type == "tweet" {
      //WriteTweet(packet.Tweet)
      TweetWrite <- packet.Tweet
      
    }
  }
}

func WholeThing(){ // periodically send entire database
  for {
    time.Sleep(10e9)
    tweets := GetTweets()
    for _,tweet := range tweets {
      TweetChan <- &tweet
    }
  }
}

func Read(conn net.Conn){ // receives relay's net.Conn interface, Write() goes to relay, Read() reads from relay
  var buf [1000]byte
  var err os.Error
  var size int
  go func(){
    for {
      time.Sleep(5e9)
      //fmt.Println("Sending ping to peers")
      //conn.Write([]byte("ping"))
    }
  }()
  for {
    size,err = conn.Read(buf[0:]) //block until relay sends us data
    if err != nil {
      fmt.Println("read error:",err)
    } else {
      fmt.Println("Relay(",conn.RemoteAddr(),")","just sent:",string(buf[0:size]))
      var packet *Packet
      err = json.Unmarshal(buf[0:size],&packet) // unmarshal json string from relay
      if err != nil {
        fmt.Println(err)
        continue //skip until next read from relay
      }
      if packet.Type == "peers" { //if relay sends us a "peers" type json, write to database
        fmt.Println(conn.RemoteAddr(),"(RELAY) sent us some peers",packet.Peers)
        WritePeers(packet.Peers) 
      } else if packet.Type == "tweet" { // if relay sends us a tweet, write to db and broadcast to local webclient
        fmt.Println(conn.RemoteAddr(),"(RELAY) sent us some tweets",packet.Tweet)
        //WriteTweet(packet.Tweet)
        TweetWrite <- packet.Tweet
      }
      
    }
  }
}
