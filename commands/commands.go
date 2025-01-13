package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
	"testapp1/bot1/models"
	"time"
)

func HandleAddCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста, укажите напоминание в формате: <текст>, <время>")
		bot.Send(msg)
		return
	}
	parts := strings.SplitN(args, " , ", 2)
	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Напоминание и время должны быть указаны: Напоминание 18:00")
		bot.Send(msg)
		return
	}

	messageText := strings.TrimSpace(parts[0])
	timeStr := strings.TrimSpace(parts[0])[1]

	reminderTime, err := time.Parse("02-01-2006 15:04", strconv.Itoa(int(timeStr)))
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Ошибка с форматом времени: %v", err))
		bot.Send(msg)
		return
	}
	models.AddReminder(messageText, reminderTime)

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Напоминание '%s' на %v добавлено!", messageText, reminderTime))
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
