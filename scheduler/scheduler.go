package scheduler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"testapp1/bot1/models"
	"time"
)

func StartScheduler(bot *tgbotapi.BotAPI) {
	go func() {
		for {
			CheckReminders(bot)
			time.Sleep(1 * time.Minute)
		}
	}()
}

func CheckReminders(bot *tgbotapi.BotAPI) {
	currentTime := time.Now()
	for _, reminder := range models.GetReminders() {
		reminderTime, err := time.Parse("2006-01-02 15:04:05", reminder.Time)
		if err != nil {
			log.Printf("Ошибка парсинга времени '%s': %v\n", reminder.Time, err)
			continue
		}
		if currentTime.After(reminderTime) {
			msg := tgbotapi.NewMessage(int64(reminder.ID), "Напоминание: "+reminder.Message)
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Ошибка отправки сообщения: %v\n", err)
			}
			models.DeleteReminder(reminder.ID)
		}
	}
}
