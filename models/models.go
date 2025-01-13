package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Reminder struct {
	ID      int
	Message string
	Time    time.Time
}

var Remindres []Reminder

func AddReminder(message string, reminderTime time.Time) { // функция добавления напоминаний
	id := len(Remindres) + 1
	reminder := Reminder{
		ID:      id,
		Message: message,
		Time:    reminderTime,
	}
	Remindres = append(Remindres, reminder)
	SaveRemindersToFile()
	fmt.Printf("Напоминание '%s' на %s добавлено!\n", reminder.Message, reminder.Time)
}

func SaveRemindersToFile() { // функция сохранения напоминаний в JSON
	data, err := json.MarshalIndent(Remindres, "", "  ")
	if err != nil {
		log.Printf("Ошибка при преобразовании в JSON:", err)
		return
	}

	err = os.WriteFile("reminders.json", data, 0644)
	if err != nil {
		log.Printf("Ощибка при сохранении напоминаний в файл:", err)
	}
}

func LoadRemindersFromFile() { // функция загрузки напоминаний из файла JSON
	data, err := os.ReadFile("reminders.json")
	if err != nil {
		log.Printf("Ошибка при чтении файла с напоминаниями", err)
		return
	}
	err = json.Unmarshal(data, &Remindres)
	if err != nil {
		log.Printf("Ошибка при загрузке напоминаний из JSON:", err)
	}
}

func GetReminders() []Reminder { // функция возвращения напоминаний

	return Remindres
} // функция получения напоминаний

func DeleteReminder(id int) bool { // функция удаления напоминаний
	for i, reminder := range Remindres {
		if reminder.ID == id {
			Remindres = append(Remindres[:i], Remindres[i+1:]...)
			SaveRemindersToFile()
			return true
		}
	}
	return false
}

func UpdateReminder(id int, newMessage string, newTime string) bool { // функция обновления напоминаний
	parsedTime, err := time.Parse("2006-01-02 15:04", newTime)
	if err != nil {
		fmt.Println("Ошибка преобразовния времени: ", err)
		return false
	}
	for i, reminder := range Remindres {
		if reminder.ID == id {
			Remindres[i].Message = newMessage
			Remindres[i].Time = parsedTime
			SaveRemindersToFile()
			return true
		}
	}
	return false
}

func FindReminderByID(id int) (*Reminder, bool) { // функция поиска напоминаний
	for _, reminder := range Remindres {
		if reminder.ID == id {
			return &reminder, true
		}
	}
	return nil, false
}
