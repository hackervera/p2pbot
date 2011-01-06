package main
import (
  sqlite "gosqlite.googlecode.com/hg/sqlite"
  "fmt"
)

var db *sqlite.Conn

var TweetWrite = make(chan *Tweet)

func WriteName(name string){
  myUsername = name
  db.Exec("INSERT INTO username (name) VALUES (?)",name)
}

func WritePeers(peers []string){
  
  for i := range peers {
    ip := peers[i]
    stmt,perr := db.Prepare("SELECT * FROM peers WHERE ip = ?")
    if perr != nil{
      fmt.Println("While SELECTing",perr)
    }
    eerr := stmt.Exec(ip)
    if eerr != nil {
      fmt.Println("While running Exec()",eerr)
    }
    if !stmt.Next() { 
      fmt.Println("Inserting:",ip,"into database")
      db.Exec("INSERT INTO peers (ip) VALUES (?)",ip)
      continue
    } else {
      fmt.Println("Skipping",ip)
    }
    stmt.Finalize()
  }
}

func GetPeers() []string{
  stmt,perr := db.Prepare("SELECT * FROM peers")
  if perr != nil{
    fmt.Println("While SELECTing",perr)
  }
  eerr := stmt.Exec()
  if eerr != nil {
    fmt.Println("While running Exec()",eerr)
  }
  var ips []string
  for {
    if !stmt.Next() { 
      break
    } else {
      var ip string
      stmt.Scan(&ip)
      ips = append(ips, ip)
    }
  }
  return ips
}

func WriteTweet(){
  for {
    tweet :=<-TweetWrite
    stmt,perr := db.Prepare("SELECT * FROM tweets WHERE timestamp = ?")
    if perr != nil{
      fmt.Println("While SELECTing",perr)
    }
    eerr := stmt.Exec(tweet.Timestamp)
    if eerr != nil {
      fmt.Println("While running Exec()",eerr)
    }
    if !stmt.Next() { 
      fmt.Println("Inserting:",tweet,"into database")
      db.Exec("INSERT INTO tweets (author,message,timestamp) VALUES (?,?,?)",tweet.Name,tweet.Message,tweet.Timestamp)
      messageChan <- []byte(tweet.Name + " said: " + tweet.Message + " [" + tweet.Timestamp + "]") //write client's message to webclient
    } else {
      fmt.Println("Skipping",tweet)
    }
    stmt.Finalize()
  }
}

func SetupDatabase(){
  db, _ = sqlite.Open("foo.db") 
  db.Exec("CREATE TABLE tweets (author, message, timestamp)")
  db.Exec("CREATE TABLE username (name)")
  db.Exec("CREATE TABLE peers (ip)")
  stmt, _ := db.Prepare("SELECT name FROM username")
  stmt.Exec()
  for {
    if !stmt.Next() { 
      fmt.Println("no usernames found") 
      break
    }
    stmt.Scan(&myUsername)
    fmt.Println("found username:",myUsername)
    break
  } 
  stmt.Finalize()
}

func GetTweets() []Tweet{
  stmt,perr := db.Prepare("SELECT * FROM tweets")
  if perr != nil{
    fmt.Println("While SELECTing",perr)
  }
  eerr := stmt.Exec()
  if eerr != nil {
    fmt.Println("While running Exec()",eerr)
  }
  var tweets []Tweet
  for {
    if !stmt.Next() { 
      break
    } else {
      var author,message,timestamp string
      stmt.Scan(&author,&message,&timestamp)
      tweet := &Tweet{author,message,timestamp}
      tweets = append(tweets, *tweet)
    }
  }
  return tweets
}

