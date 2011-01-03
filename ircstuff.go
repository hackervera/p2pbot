package main
import (
  "irc"
  "fmt"
  "strings"
  "time"
)
func ircStuff() {
  var nick string
  type payload struct {
    Author string
    Text string
    Timestamp string
  }

  irccon := irc.IRC("testgo", "testgo")
  irccon.Connect("irc.freenode.net:6667")
  
  irccon.AddCallback("JOIN",func(e *irc.IRCEvent){
    fmt.Println(e)
    nick = e.Nick
    irccon.Privmsg("#bootstrap","PING")
  })
  irccon.AddCallback("NOTICE",func(e *irc.IRCEvent){
    if strings.Contains(e.Message, "PING"){ 
      fmt.Println("received response from",e.Host)
      PingResponse <- e.Host
    }
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
