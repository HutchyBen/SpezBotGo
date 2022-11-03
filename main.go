package main

import (
	"fmt"
	"spezbot/bot"
	"spezbot/commands"
)

func main() {
	bot, err := bot.NewBot("config.json")
	if err != nil {
		panic(err)
	}
	bot.CH.AddBulk(commands.Commands)
	bot.CH.Register(bot.Client)
	bot.Client.UpdateGameStatus(0, "Clash of Clans")
	fmt.Println("Bot is running with name " + bot.Client.State.User.Username)
	select {}

}
