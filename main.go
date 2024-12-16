package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"testapp1/bot1/commands"
	"testapp1/bot1/models"
	"testapp1/bot1/scheduler"
)

func main() {
	models.SaveRemindersToFile()
	log.Printf("Бот запущен")
	token := os.Getenv("TELEGRAM_BOT_TOKEN") // получения токена
	if token == "" {
		log.Fatalf("Не удалось загрузить токен")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalln("Ошибка инициализации бота: %v", err)
	}

	bot.Debug = true
	log.Printf("Авторизован под именем %s", bot.Self.UserName)

	go scheduler.StartScheduler(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
	if err != nil {
		log.Fatalln("Ошибка получения обновления: %v", err)
	}

	for update := range updates {
		if update.Message != nil {
			go handleUpdate(bot, update)
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.Text != "" {
		switch update.Message.Command() {

		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я - бот-помощник для напоминаний. "+
				"Введите /help, чтобы узнать, что я умею")
			bot.Send(msg)
		case "help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Доступные команды: \n/start- начало работы, "+
				"\n/add - добавить напоминание, \n/help - помощь, "+
				"\n/search -  для поиска закладок, \n/delete - для удаления заклдаки, "+
				"\n/view -  для просмотра Ваших закладок")

			bot.Send(msg)
		case "add":
			commands.HandleAddCommand(bot, update.Message)
		case "search":
			commands.HandleSearchCommand(bot, update.Message)
		case "delete":
			commands.HandleDeleteCommand(bot, update.Message)
		case "view":
			commands.HandleViewCommand(bot, update.Message)
			//remindres := models.GetReminders()
			//if len(remindres) == 0 {
			//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нет текущих напоминаний")
			//	bot.Send(msg)
			//} else {
			//	messageText := "Ваши напоминания: \n"
			//	for _, r := range remindres {
			//		messageText += fmt.Sprintf("ID: %d | Сообщение: %s | Время: %s\n", r.ID, r.Message, r.Time)
			//	}
			//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			//	bot.Send(msg)
			//}
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я пока не понимаю эту команду, попробуйте /help")
			bot.Send(msg)
		}
	}
}
