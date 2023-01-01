package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type Bot struct {
	Name    string `json:"name"`
	Token   string `json:"token"`
	Channel string `json:"channel"`
}

func registerBot() {
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
		return
	}
	fmt.Println(string(bot_json))

	fmt.Print("y/n: ")
	var confirm string
	fmt.Scan(&confirm)

	if confirm != "y" {
		return
	}

	f, err := os.Create("bot.conf")
	_, err = f.Write(bot_json)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

}

func sendMessage() {
	f, err := os.Open("bot.conf")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	var bot Bot
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&bot); err != nil {
		fmt.Println(err)
	}

	url := fmt.Sprintf("https://discordapp.com/api/channels/%v/messages", bot.Channel)

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

func main() {
	fmt.Println("Hello")

	register_flag := flag.Bool("r", false, "register bot")
	send_flag := flag.Bool("s", false, "send a message")

	flag.Parse()

	if *register_flag {
		registerBot()
	} else if *send_flag {
		sendMessage()
	}
}
