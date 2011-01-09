package main
import (
  sqlite "gosqlite.googlecode.com/hg/sqlite"
  "fmt"
  "json"
  "crypto/rsa"
  "os"
  "big"
)

var db *sqlite.Conn

var TweetWrite = make(chan *Tweet)
var History []Tweet


type Modder struct {
  Mod string
  Name string
}

type MyKey struct {
  PublicKey
  D string
  P,Q string
}

type PublicKey struct {
  N string
  E int
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
      username := Base64Encode(key.N.Bytes())
      WriteName(string(username))
      WriteKey(key)
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
      tweet := Tweet{Name:author,Message:message,Timestamp:timestamp}
      tweets = append(tweets, tweet)
    }
  }
  return tweets
}


func WriteName(name string){
  myUsername = name
  db.Exec("INSERT INTO username (name) VALUES (?)",name)
}

func WriteFriend(mod string, username string){
  db.Exec("INSERT INTO friends (mod,username) VALUES (?,?)",mod,username)
}

func GetFriends() []byte {
  var friends []Modder
  stmt, _ := db.Prepare("SELECT mod,username FROM friends")
  stmt.Exec()
  for {
    var mod,username string
    if !stmt.Next() { 
      //fmt.Println("Unknown Username")
      //return []byte("")
      break
    }
    stmt.Scan(&mod,&username)
    fmt.Println("got username",username)
    
    friend := &Modder{mod,username}
    friends = append(friends, *friend)
    
  } 
  stmt.Finalize()
  friendsjson,err := json.Marshal(friends)
  
  if err != nil {
    fmt.Println(err)
  }
  //fmt.Println("friends:",friends)
  //fmt.Println("friendjson:",friendsjson)
  return friendsjson
}

func WriteKey(key *rsa.PrivateKey){
  mykey := &MyKey{D:key.D.String(),P:key.P.String(),Q:key.Q.String(), PublicKey:PublicKey{N:key.PublicKey.N.String(),E:key.PublicKey.E}}
  jsonkey,err := json.Marshal(mykey)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  db.Exec("INSERT INTO key (key) VALUES (?)",jsonkey)
}

func GetKey() rsa.PrivateKey {
  
  var placeholder []byte
  var mykey MyKey
  var key rsa.PrivateKey
  stmt,err := db.Prepare("SELECT key FROM key")
  if err != nil{
    fmt.Println("While SELECTing",err)
  }
  err = stmt.Exec()
  if err != nil {
    fmt.Println("While running Exec()",err)
  }
  for {
    if !stmt.Next() { 
      break
    } else {
      
      stmt.Scan(&placeholder)
      //fmt.Println("Getting key:",string(placeholder))
      err = json.Unmarshal(placeholder,&mykey)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }
      d := big.NewInt(0)
      p := big.NewInt(0)
      q := big.NewInt(0)
      n := big.NewInt(0)
      d.SetString(mykey.D,10)
      p.SetString(mykey.P,10)
      q.SetString(mykey.Q,10)
      n.SetString(mykey.PublicKey.N,10)
      
      pubkey := rsa.PublicKey{N:n, E:mykey.PublicKey.E}
      key = rsa.PrivateKey{D:d,P:p,Q:q,PublicKey:pubkey}
      //fmt.Println("KEY:",key)
      err = key.Validate()
      if err != nil {
        fmt.Println("key errors:",err)
      } else {
        //fmt.Println("key looks valid")
      }
    }
  }
  return key
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

func GetRelays() []string{
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



func TweetWriter(){
  for {
    fmt.Println("Waiting for next tweet")
    tweet :=<-TweetWrite
    stmt,perr := db.Prepare("SELECT * FROM tweets WHERE timestamp = ?")
    test := Verify([]byte(tweet.Message), tweet.Sig, tweet.Name)
    if test != true {
      fmt.Println("not verified")
      continue
    }
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
      Tweets <- tweet
      fmt.Println("successfully sent tweet on wire")
    } else {
      fmt.Println("Skipping",tweet)
    }
    stmt.Finalize()
    
  }
}


