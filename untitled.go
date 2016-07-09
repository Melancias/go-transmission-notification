package main

import (

	"net/http"
	"net/url"
	"bytes"
	"github.com/deckarep/gosx-notifier"
	"github.com/Jeffail/gabs"
	//"time"
	//"fmt"
	"time"
)

const urlRPC string = ""
const username string = ""
const password string = ""

func notifys(texto, titulo, subtitulo string) {
	note := gosxnotifier.NewNotification(texto)
	note.Title = titulo
	note.Subtitle = subtitulo
	note.Sound = gosxnotifier.Basso
	note.Group = "com.unique.yourapp.identifier"
	note.Sender = "com.apple.Safari"
	note.Link = "http://www.yahoo.com"
	note.AppIcon = "gopher.png"
	note.ContentImage = "gopher.png"
	note.Push()
}

func testRequest(whitelist []string,token string) []string{
	client := &http.Client{}
	data:=url.Values{}
	data.Set("method","torrent-get")
	jsonParsed, _ := gabs.ParseJSON([]byte(`{
                "arguments": {
                    "fields": [ "id", "name", "percentDone" ]
                },
                "method": "torrent-get",
                "tag": 39693
             }`))
	req, _ := http.NewRequest("POST", urlRPC, bytes.NewBuffer(jsonParsed.Bytes()))
	req.SetBasicAuth(username,password)
	req.Header.Set("X-Transmission-Session-Id",token)
	resp,_:=client.Do(req)
	jsonValues,_:=gabs.ParseJSONBuffer(resp.Body)
	//print(jsonValues.Path("arguments.torrents").String())
	//print(jsonValues.S("arguments","torrents").String())
	children, _ := jsonValues.S("arguments","torrents").Children()
	for _, child := range children {
		torrent:=child.Path("name").String()
		id:=child.Path("id").String()
		percent:=child.Path("percentDone").Data().(float64)
		if(percent==1.0 && !stringInSlice(id,whitelist)){
			whitelist=append(whitelist,id)
			println("Sending notification about: ", torrent)
			notifys(torrent,"Download Acabou","YAY")
		}

	}

	return whitelist
}
func main() {
	data:=url.Values{}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlRPC, bytes.NewBufferString(data.Encode()))
	req.SetBasicAuth(username,password)
	resp,_:=client.Do(req)
	auth:= resp.Header.Get("X-Transmission-Session-Id")
	whitelist:=setupWhitelist(auth)
	println("starting the daemon at 10 second intervals")
	for {

		time.Sleep(10 * time.Second)
		whitelist=testRequest(whitelist,auth)
	}



	//If necessary, check error
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func setupWhitelist(token string) []string{
	client := &http.Client{}
	data:=url.Values{}
	data.Set("method","torrent-get")
	jsonParsed, _ := gabs.ParseJSON([]byte(`{
                "arguments": {
                    "fields": [ "id", "name", "percentDone" ]
                },
                "method": "torrent-get",
                "tag": 39693
             }`))
	req, _ := http.NewRequest("POST", urlRPC, bytes.NewBuffer(jsonParsed.Bytes()))
	req.SetBasicAuth(username,password)
	req.Header.Set("X-Transmission-Session-Id",token)
	resp,_:=client.Do(req)
	jsonValues,_:=gabs.ParseJSONBuffer(resp.Body)
	var whitelist []string
	children,_:=jsonValues.S("arguments","torrents","id").Children()
	for _, child := range children {
		println("Adding ",child.String())
		whitelist=append(whitelist,child.String())
	}
	return whitelist
}
