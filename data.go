package main
import (
  sqlite "gosqlite.googlecode.com/hg/sqlite"
  "fmt"
)

var db *sqlite.Conn

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
      //fmt.Println("no usernames found") 
      fmt.Println("Inserting:",ip,"into database")
      db.Exec("INSERT INTO peers (ip) VALUES (?)",ip)
      break
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
  } 
  stmt.Finalize()
}
