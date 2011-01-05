package main
import (
  "irc"
  "fmt"
  "strings"
  "time"
)
func ircStuff() {
  var nick string
  var peers []string
  type payload struct {
    Author string
    Text string
    Timestamp string
  }

  irccon := irc.IRC("testgo", "testgo")
  irccon.Connect("irc.freenode.net:6667")
  
  irccon.AddCallback("JOIN",func(e *irc.IRCEvent){
    fmt.Println(e)
    if nick == e.Nick {
      irccon.Privmsg("#bootstrap",nick + " has arrived!")
      //irccon.Privmsg("#bootstrap","PING")
      go irccon.SendRaw("who #bootstrap")
      
      time.Sleep(5e9)
      fmt.Println("About to write peers:",peers)
      WritePeers(peers)
      go BroadcastPeers()
      //go BroadcastTweets()
      
    }
  })
  irccon.AddCallback("NOTICE",func(e *irc.IRCEvent){
    if strings.Contains(e.Message, "PING"){ 
      fmt.Println("received response from",e.Host)
      PingResponse <- e.Host
    }
  })
  irccon.AddCallback("311",func(e *irc.IRCEvent){
    //irccon.Privmsg("#bootstrap",e.Arguments[3])
  })
  irccon.AddCallback("352",func(e *irc.IRCEvent){ // who response from server
    //irccon.Privmsg("#bootstrap",e.Arguments[3])
    peers = append(peers,e.Arguments[3])
    fmt.Println("GOT WHO RESPONSE:",e)
    //WritePeer(e.Arguments[3])
  })
  irccon.AddCallback("372",func(e *irc.IRCEvent){
    //fmt.Println(e)
  })
  
  irccon.AddCallback("353",func(e *irc.IRCEvent){ //channel name list
    names := strings.Split(e.Message," ",-1)
    fmt.Println(e)
    //fmt.Println(e.Message)
    for i := range names {
      if names[i] != nick {
        fmt.Println(names[i])
      }
    }
  })
  
  irccon.AddCallback("001", func(e *irc.IRCEvent) { 
    irccon.Join("#bootstrap") 
    fmt.Println(e)
    nick = e.Arguments[0]
  })
    
  go func(){
    for {
      v:=<- ircChan
      timestamp := time.LocalTime().String()
      author := myUsername
      payload := `{"author":"`+author+`","text":"`+string(v)+`","timestamp":"`+timestamp+`"}`
      msg := `{"type":"newtweet","tweetbody":`+payload+`}`
      irccon.Privmsg("#bootstrap",msg) 
    }
  }()
}
