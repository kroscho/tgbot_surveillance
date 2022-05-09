package Models

import (
	"telegram_test_bot/UI/Messages"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// вся информация об одном юзере в течение сессии
type UserSession struct {
	UserID_TG        int
	UserID_VK        int
	User             User
	Token            string
	Persons          []Person // отслеживаемые люди
	CurMsg           tgbotapi.Message
	TypeMessage      *Messages.TypeMessage
	CurTrackedPerson User //  текущий отслеживаемый человек
	CurSearchByName  bool
}

func NewUserSession(userId int) *UserSession {
	user := UserSession{}
	user.UserID_TG = userId
	user.TypeMessage = Messages.NewTypeMessage()

	return &user
}
