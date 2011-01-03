package main

import (
  "strings"
  "fmt"
  "net"
  "json"
  "time"
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
    fmt.Println("closing connection")
    conn.Close()
  }
}

func ListenCall(l *net.TCPListener){
  for {
    conn,_:=l.Accept()
    fmt.Println("Incoming connection!")
    message := make([]byte, 1000)
    _,err := conn.Read(message)
    if err != nil {
      fmt.Println(err)
    }
    var marshaldata []byte
    if strings.Contains(string(message),"peer"){
      data := &PeerBlob{GetPeers(),"good"}
      marshaldata,err = json.Marshal(data)
      if err != nil {
        fmt.Println(err)
      }
      conn.Write([]byte(marshaldata))
    } else if strings.Contains(string(message),"data"){
      timestamp := time.LocalTime().String()
      fmt.Println(timestamp)
      data := GetTweets()
      marshaldata,err = json.Marshal(data)
      if err != nil {
        fmt.Println(err)
      }
      conn.Write([]byte(marshaldata))
    } else if strings.Contains(string(message),"new tweet"){
      fmt.Println("trying to read tweet")
      var marshaldata = make([]byte,1000)
      conn.Write([]byte("flush"))
      num,_ := conn.Read(marshaldata)
      fmt.Println(num)
      var tweet Tweet
      fmt.Println("incoming tweet", string(marshaldata))
      merr := json.Unmarshal(marshaldata,&tweet)
      if merr != nil {
        fmt.Println(merr)
      }
      WriteTweet(tweet)
      messageChan <- []byte(genMessage(tweet.Name + tweet.Message))
    } else {
      conn.Write([]byte("no dataz for you!"))
    }
    fmt.Println("closing connection")
    conn.Close()
  }
}

func GetDataFromPeers(){
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
      fmt.Println("Success! Reading data")
      _,err = conn.Write([]byte("i can haz data"))
      message := make([]byte,1000)
      num,err := conn.Read(message)
      if err != nil {
        fmt.Println(err)
      }
      var tweets []Tweet
      err = json.Unmarshal(message[0:num],&tweets)
      if err != nil {
        fmt.Println(err)
      }
      fmt.Println("Writing tweets to database",tweets)
      for _,tweet := range tweets {
        fmt.Println("Writing",tweet)
        WriteTweet(tweet)
      }
      fmt.Println("closing connection")
      conn.Close()
      
    }(ip)
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
      

