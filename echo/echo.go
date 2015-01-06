package main

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

const token = "barrage123"

func main() {
	http.HandleFunc("/", weChatHandle)

	addr := ":3080"
	log.Printf("Serer is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func weChatHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("url :", r.URL.String())

	if !checkSignature(r) {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		fmt.Fprintf(w, r.FormValue("echostr"))
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Receive WeChat post message failed:", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	log.Println("body:", string(data))

	// w.WriteHeader(http.StatusCreated)

	textMsg := recvTextMsg(w, data)
	if textMsg != nil {
		replyTextMsg(w, textMsg)
	}
}

func recvTextMsg(w http.ResponseWriter, data []byte) (msg *TextMsg) {
	msg = &TextMsg{}
	if err := xml.Unmarshal(data, msg); err != nil {
		log.Println("Parse WeChat post message failed:", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	log.Printf("TextMsg: %#v", msg)
	return
}

func replyTextMsg(w http.ResponseWriter, m *TextMsg) {
	msgFormat := `<xml>
	<ToUserName><![CDATA[%s]]></ToUserName>
	<FromUserName><![CDATA[%s]]></FromUserName>
	<CreateTime>%d</CreateTime>
	<MsgType><![CDATA[text]]></MsgType>
	<Content><![CDATA[%s]]></Content>
	</xml>`

	fmt.Fprintf(w, msgFormat, m.FromUserName, m.ToUserName, time.Now().Unix(), m.Content)
}

func checkSignature(r *http.Request) bool {
	r.ParseForm()
	var signature = r.FormValue("signature")
	var timestamp = r.FormValue("timestamp")
	var nonce = r.FormValue("nonce")

	strs := sort.StringSlice{token, timestamp, nonce}
	sort.Strings(strs)
	var str string
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil)) == signature
}

// Common message header
type MsgHeader struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
}

type TextMsg struct {
	MsgHeader
	Content string
	MsgId   int64
	Encrypt string
}
