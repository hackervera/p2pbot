package main
import (
  sqlite "gosqlite.googlecode.com/hg/sqlite"
  "fmt"
  "json"
  "crypto/rsa"
)

var db *sqlite.Conn

var TweetWrite = make(chan *Tweet)
var History []Tweet

type PQDN struct {
  P string
  Q string
  D string
  N string
}


func WriteName(name string){
  myUsername = name
  db.Exec("INSERT INTO username (name) VALUES (?)",name)
}

func WriteFriend(mod string, username string){
  db.Exec("INSERT INTO friends (mod,username) VALUES (?,?)",mod,username)
}

func GetFriends() []byte {
  friends := make(map[string]string)
  stmt, _ := db.Prepare("SELECT mod,username FROM friends")
  stmt.Exec()
  for {
    var mod,username string
    if !stmt.Next() { 
      fmt.Println("Unknown Username")
      return []byte("")
      break
    }
    stmt.Scan(&mod,&username)
    friends[mod] = username
    
  } 
  stmt.Finalize()
  friendsjson,err := json.Marshal(friends)
  if err != nil {
    fmt.Println(err)
  }
  return friendsjson
}

func WriteKey(key string){
  db.Exec("INSERT INTO key (key) VALUES (?)",key)
}

func GetUsername() string{
  stmt, _ := db.Prepare("SELECT name FROM username")
  stmt.Exec()
  var username string
  for {
    if !stmt.Next() { 
      fmt.Println("Unknown Username")
      break
    }
    stmt.Scan(&username)
    stmt.Finalize()
    return username
  } 
  return ""
}

func GetHistory() []byte {
  tweets := GetTweets()
  TweetsJSON,err := json.Marshal(tweets)
  if err != nil {
    fmt.Println(err)
  }
  return TweetsJSON
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
  db.Exec("CREATE TABLE key (key)")
  db.Exec("CREATE TABLE peers (ip)")
  db.Exec("CREATE TABLE friends (mod, username)")
  stmt, _ := db.Prepare("SELECT name FROM username")
  stmt.Exec()
  for {
    if !stmt.Next() { 
      fmt.Println("no key found in database, generating... this may take a second") 
      key := GenKey()
      JsonKey,_ := json.Marshal(&PQDN{key.P.String(),key.Q.String(),key.D.String(),key.PublicKey.N.String()})
      fmt.Println(JsonKey)
      username := Base64Encode(key.N.Bytes())
      WriteKey(string(JsonKey))
      WriteName(string(username))
      break
    }
    var unmarshalled rsa.PrivateKey
    var marshalled []byte
    stmt.Scan(&marshalled)
    json.Unmarshal(marshalled, unmarshalled)
    myUsername = "foo"
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
      tweet := Tweet{author,message,timestamp}
      tweets = append(tweets, tweet)
    }
  }
  return tweets
}

