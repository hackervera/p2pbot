package main
import (
  "http"
  "io/ioutil"
  "regexp"
)

func Index(w http.ResponseWriter, r *http.Request){
  index,_ := ioutil.ReadFile("index.html")
  re := regexp.MustCompile("NAME")
  name := GetUsername()
  index = re.ReplaceAll(index, []byte(name))
  w.Write(index)
}

func TweetGrabber(w http.ResponseWriter, r *http.Request){
  w.Write([]byte(GetHistory()))
}

func InsertTweet(w http.ResponseWriter, r *http.Request){
  query,_ := http.ParseQuery(r.URL.RawQuery)
  tweet := &Tweet{Name:query["name"][0],Message:query["message"][0],Timestamp:query["date"][0]}
  TweetWrite <- tweet
  w.Write([]byte(`"looks good"`))
}

func AddressBook(w http.ResponseWriter, r *http.Request){
  friends := GetFriends()
  w.Write(friends)
}

func AddFriend(w http.ResponseWriter, r *http.Request){
  query,_ := http.ParseQuery(r.URL.RawQuery)
  mod := query["mod"][0]
  username := query["username"][0]
  WriteFriend(mod,username)
}
