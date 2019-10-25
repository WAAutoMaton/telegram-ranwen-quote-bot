package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

var bot *tgbotapi.BotAPI

func initBotAPI(token, proxy string) (bot *tgbotapi.BotAPI) {
	var err error
	if len(proxy) > 0 {
		client, err := createProxyClient(proxy)
		if err != nil {
			log.Fatalln(err)
		}
		bot, err = tgbotapi.NewBotAPIWithClient(token, client)
		if err != nil {
			log.Fatalf("Some error occur: %s.\n", err)
		}
	} else {
		bot, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return
}

func createProxyClient(proxy string) (client *http.Client, err error) {
	log.Println("verify proxy:", proxy)
	var proxyURL *url.URL
	proxyURL, err = url.Parse(proxy)
	if err == nil {
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
		return
	}
	return
}
type Config struct {
	Token string
	Proxy string
}
var quote []string
func main() {
	readMessages()
	b,err:=ioutil.ReadFile("config.json")
	if err!=nil {
		panic(err)
	}
	config:=&Config{}
	err=json.Unmarshal(b,config)
	if err!=nil {
		panic(err)
	}
	bot = initBotAPI(config.Token,config.Proxy)
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	log.Print("telegram bot started")
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "ranwen":
			msg.Text=quote[rand.Intn(len(quote))]
		}
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
func readMessages() {
	for i:=1; i<=2; i++ {
		sss:=strconv.Itoa(i)
		if i==1 {
			sss=""
		}
		html,err:=ioutil.ReadFile(fmt.Sprintf("messages%s.html",sss))
		if err!=nil {
			panic(err)
		}
		//fmt.Println(string(html))
		r:=regexp.MustCompile(`<div\sclass="from_name">\s*ranwen\s*<span\sclass="details">.+?</span>\s*</div>\s*<div\sclass="text">\s*([\s\S]+?)\s*</div>`)
		//r:=regexp.MustCompile(`d(i)v`)
		t:=r.FindAllStringSubmatch(string(html),-1)
		fmt.Println(len(t))
		for _,j:=range t {
			quote=append(quote,j[1])
		}
	}
}

