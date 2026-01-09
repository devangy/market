package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Bot(tgEventC <-chan any) {

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Bot token is empty! Set BOT_TOKEN environment variable.")
	}
	//dddddddddddd
	// context with a timeout for cancellation of request
	ctx := context.Background()

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	botUser, err := bot.GetMe(ctx)
	if err != nil {
		log.Fatal("bot authentication: ", err)
	}
	//deviig
	log.Debug("Bot user: %+v\n", botUser)

	// updates from bot via long polling for testing
	updates, err := bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		log.Fatal("failed to start long polling:", err)
	}

	// bot handler to handle req
	bh, _ := th.NewBotHandler(bot, updates)

	defer func() { _ = bh.Stop() }()

	// Register new handler with match on command `/start`
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		// Send a message with inline keyboard
		_, _ = ctx.Bot().SendMessage(ctx, tu.Messagef(
			tu.ID(message.Chat.ID),
			`Hello %s !  Welcome to DataLog`, message.From.FirstName,
		).WithReplyMarkup(tu.InlineKeyboard(
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("Start").WithCallbackData("all_events"))),
		))
		return nil
	}, th.CommandEqual("start"))

	// Register new handler with match on a call back query with data equal to `go` and non-nil message
	bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {

		// Answer callback query
		_ = bot.AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Done"))

		for event := range tgEventC {
			fmt.Println("tgChanEvents", event)
			_, _ = ctx.Bot().SendMessage(ctx, tu.Messagef(
				tu.ID(query.Message.GetChat().ID),
				"Received: %v", event,
			))
		}
		// }()

		return nil

	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("all_events"))

	bh.Start()

}
