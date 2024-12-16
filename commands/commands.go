package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
	"testapp1/bot1/models"
)

func HandleAddCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста, укажите напоминание в формате: /add <текст>, <время>")
		bot.Send(msg)
		return
	}
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Напоминание и время должны быть указаны: /add Напоминание 18:00")
		bot.Send(msg)
		return
	}

	messageText := parts[0]
	time := parts[1]

	models.AddReminder(messageText, time)

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Напоминание '%s' на %s добавлено!", messageText, time))
	bot.Send(msg)
}

func HandleViewCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	reminders := models.GetReminders()

	if len(reminders) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Нет текущих напоминаний")
		bot.Send(msg)
		return
	}

	messageText := "Ваши напоминания:\n"
	for _, r := range reminders {
		messageText += fmt.Sprintf("ID: %d | Сообщение: %s | Время: %s\n", r.ID, r.Message, r.Time)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
	bot.Send(msg)
}

func HandleDeleteCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := strings.TrimSpace(message.CommandArguments())
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста, укажите ID напоминания для удаления. Например: /delete 1 ")
		bot.Send(msg)
		return
	}
	id, err := strconv.Atoi(args)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "ID должен быть числом. Например: /delete 1")
		bot.Send(msg)
		return
	}

	success := models.DeleteReminder(id)
	if success {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Напоминание с ID %d успешно удалено", id))
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Напоминание с ID %d не найдено", id))
		bot.Send(msg)
	}
}

func HandleSearchCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	keyword := strings.TrimSpace(message.CommandArguments())
	if keyword == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста , укажите ключевое слово для поиска. Например: /search важное")
		bot.Send(msg)
		return
	}

	var results []models.Reminder
	for _, reminder := range models.Remindres {
		if strings.Contains(strings.ToLower(reminder.Message), strings.ToLower(keyword)) {
			results = append(results, reminder)
		}
	}
	if len(results) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Напоминаний с ключевым словом \"%s\" не найдено.", keyword))
		bot.Send(msg)
	} else {
		messageText := fmt.Sprintf("Найдено %d напоминаний с ключевым словом \"%s\":\n", len(results), keyword)
		for _, r := range results {
			messageText += fmt.Sprintf("ID: %d | Сообщение: %s | Время: %s\n", r.ID, r.Message, r.Time)
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
		bot.Send(msg)
	}
}
