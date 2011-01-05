package main
import (
  "http"
  "strings"
)
func hello(res http.ResponseWriter, req *http.Request) {
  tweets := GetTweets()
  var tweet []string
  for _,v := range tweets {
    tweet = append(tweet, "<p>"+v.Name + "said: <b>" + v.Message + "</b> @"+v.Timestamp+"</p>")
  }
  tweetstring := strings.Join(tweet, " ")
  var div, blocker string
  if hasUsername != 1{
    div = `
    <div id='user-form'>
  Enter your name<br>
  <input type='text' id='username'><br>
  <a href='#' onCLick ="var json = {}; json.type = 'username'; json.name= $('#username').val(); websocket.send(JSON.stringify(json)); $('#user-form').hide(); $('#enter-data').show()">Update Name</a>
  </div>
  `
    blocker = `$('#enter-data').hide()`
  }
    html := `
<!DOCTYPE html>

<meta charset="utf-8" />

<title>WebSocket Test</title>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.4.4/jquery.min.js"></script>
<script language="javascript" type="text/javascript">

  var wsUri = "ws://`+hostName+`:`+portNumber+`/socket";
  var output;

  function init()
  {
    output = document.getElementById("output");
    testWebSocket();
  }

  function testWebSocket()
  {
    websocket = new WebSocket(wsUri);
    websocket.onopen = function(evt) { onOpen(evt) };
    websocket.onclose = function(evt) { onClose(evt) };
    websocket.onmessage = function(evt) { onMessage(evt) };
    websocket.onerror = function(evt) { onError(evt) };
  }

  function onOpen(evt)
  {
    writeToScreen("CONNECTED");
    //doSend("WebSocket rocks");
  }

  function onClose(evt)
  {
    writeToScreen("DISCONNECTED");
  }

  function onMessage(evt)
  {
    eval(evt.data);
    //websocket.close();
  }

  function onError(evt)
  {
    writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
  }

  function doSend(message)
  {
    writeToScreen("SENT: " + message); 
    websocket.send(message);
  }

  function writeToScreen(message)
  {
    var pre = document.createElement("p");
    pre.style.wordWrap = "break-word";
    pre.innerHTML = message;
    output.insertBefore(pre,output.childNodes[0]);
  }

  window.addEventListener("load", init, false);

</script>

<h2>WebSocket Test</h2>
` + div + `
<div id='enter-data'>
<input type='text' id='status-update'>
<a href='#' onClick="var json = {}; json.type = 'update'; json.msg = document.getElementById('status-update').value; websocket.send(JSON.stringify(json)); $('#output').prepend('<p>'+document.getElementById('status-update').value+'</p>'); $('#status-update').val(''); "><br>Send Message</a> 
</div>
<div id="output"></div>
<div id='review'>`+tweetstring+`</div>
<script>` + blocker + `</script>
</html> 
`

res.Write([]byte(html))
}
