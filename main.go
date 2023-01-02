package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var conf_file string = "bot.conf"
var base_url string = "https://discordapp.com/api"

type Bot struct {
	Name    string `json:"name"`
	Token   string `json:"token"`
	Channel string `json:"channel"`
}

type Author struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type Message struct {
	Id        string `json:"id"`
	Author    Author `json:"author"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

func loadBot() (Bot, error) {
	var bot Bot
	f, err := os.Open(conf_file)
	if err != nil {
		return bot, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&bot); err != nil {
		return bot, err
	}
	return bot, nil
}

func registerBot() Bot {
	var bot Bot
	fmt.Println("Start registaration")
	fmt.Print("Enter Bot name: ")
	fmt.Scan(&bot.Name)
	fmt.Print("Enter Token: ")
	fmt.Scan(&bot.Token)
	fmt.Print("Enter Channel: ")
	fmt.Scan(&bot.Channel)

	bot_json, err := json.MarshalIndent(bot, "", "\t")
	if err != nil {
		fmt.Println(err)
		return bot
	}
	fmt.Println(string(bot_json))

	fmt.Print("y/n: ")
	var confirm string
	fmt.Scan(&confirm)

	if confirm != "y" {
		return bot
	}

	f, err := os.Create(conf_file)
	_, err = f.Write(bot_json)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	return bot
}

func sendMessage(bot Bot) {
	url := fmt.Sprintf("%v/channels/%v/messages", base_url, bot.Channel)

	fmt.Print("Enter message: ")
	var message string
	fmt.Scan(&message)
	jsonStr := fmt.Sprintf(`{"content": "%v"}`, message)
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bot %v", bot.Token))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	fmt.Println(res.Status)
}

func showMessages(bot Bot) {
	url := fmt.Sprintf("%v/channels/%v/messages", base_url, bot.Channel)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bot %v", bot.Token))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	//fmt.Println(string(body))

	var messages []Message
	json.Unmarshal(body, &messages)
	//fmt.Println(messages)

	messages_json, err := json.MarshalIndent(messages, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(messages_json))
}

func main() {
	conf_file = "bot.conf"
	var bot Bot
	fmt.Println("Hello")
	register_flag := flag.Bool("register", false, "register bot")
	send_flag := flag.Bool("send", false, "send a message")
	show_flag := flag.Bool("show", false, "show messages")

	flag.Parse()

	if !*register_flag {
		loaded_bot, err := loadBot()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Please execute with -r option")
			return
		}
		bot = loaded_bot
	}

	if *register_flag {
		registerBot()
	} else if *send_flag {
		sendMessage(bot)
	} else if *show_flag {
		showMessages(bot)
	}
}
