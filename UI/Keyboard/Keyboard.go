package Keyboard

import (
	"strconv"
	"telegram_test_bot/Models"
	"telegram_test_bot/UI/Pages"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// клавиатура для "Главная"
func GetMainKeyboard(isToken bool) *tgbotapi.InlineKeyboardMarkup {
	var keyboardMain tgbotapi.InlineKeyboardMarkup
	if isToken {
		keyboardMain = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Обновить токен ВК", Pages.TOKENVK),
				tgbotapi.NewInlineKeyboardButtonData("Слежка", Pages.SURVEILLANCE),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Подписка", Pages.SUBSCRIPTION),
				tgbotapi.NewInlineKeyboardButtonData("Помощь", Pages.HELP_MAIN),
			),
		)
	} else {
		keyboardMain = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Токен ВК", Pages.TOKENVK),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Подписка", Pages.SUBSCRIPTION),
				tgbotapi.NewInlineKeyboardButtonData("Помощь", Pages.HELP_MAIN),
			),
		)
	}
	return &keyboardMain
}

// клавиатура для "Слежка"
func GetSurvaillanceKeyboard(user *Models.UserSession) *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance tgbotapi.InlineKeyboardMarkup

	listButtons1 := []tgbotapi.InlineKeyboardButton{}
	listButtons2 := []tgbotapi.InlineKeyboardButton{}
	listButtons3 := []tgbotapi.InlineKeyboardButton{}
	i := 1
	for _, user := range user.Persons {
		nameFriend := user.User.FirstName + " " + user.User.LastName
		data := Pages.TRACKED + "_" + strconv.Itoa(user.User.UID)
		if i <= 2 {
			listButtons1 = append(listButtons1, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
		} else if i > 2 && i <= 4 {
			listButtons2 = append(listButtons2, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
		} else if i > 4 && i <= 6 {
			listButtons3 = append(listButtons3, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
		}
		i += 1
	}

	if len(user.Persons) != 0 {
		keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
			listButtons1[:],
			listButtons2[:],
			listButtons3[:],
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Добавить в отслеживаемые", Pages.ADD_TO_TRACKED),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Помощь", Pages.HELP_SURVAILlANCE),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.MAIN),
				tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
			),
		)
	} else {
		keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Добавить в отслеживаемые", Pages.ADD_TO_TRACKED),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Помощь", Pages.HELP_SURVAILlANCE),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.MAIN),
				tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
			),
		)
	}
	return &keyboardSurvaillance
}

// клавиатура для "Друзья по имени"
func GetFriendsByNameKeyboard(listFriends map[int]Models.User) *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance tgbotapi.InlineKeyboardMarkup

	if len(listFriends) != 0 {
		listButtons1 := []tgbotapi.InlineKeyboardButton{}
		listButtons2 := []tgbotapi.InlineKeyboardButton{}
		listButtons3 := []tgbotapi.InlineKeyboardButton{}
		listButtons4 := []tgbotapi.InlineKeyboardButton{}
		listButtons5 := []tgbotapi.InlineKeyboardButton{}
		i := 1
		for key, user := range listFriends {
			nameFriend := user.FirstName + " " + user.LastName
			data := Pages.ADD_BY_NAME + "_" + strconv.Itoa(key)
			if i <= 2 {
				listButtons1 = append(listButtons1, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
			} else if i > 2 && i <= 4 {
				listButtons2 = append(listButtons2, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
			} else if i > 4 && i <= 6 {
				listButtons3 = append(listButtons3, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
			} else if i > 6 && i <= 8 {
				listButtons4 = append(listButtons4, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
			} else if i > 8 && i <= 10 {
				listButtons5 = append(listButtons5, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
			}
			i += 1
		}

		keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
			listButtons1[:],
			listButtons2[:],
			listButtons3[:],
			listButtons4[:],
			listButtons5[:],
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.ADD_BY_NAME),
				tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
			),
		)
	} else {
		keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.ADD_BY_NAME),
				tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
			),
		)
	}
	return &keyboardSurvaillance
}

// клавиатура для "Друзья по id"
func GetFriendsByIDKeyboard(listFriends map[int]Models.User) *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance tgbotapi.InlineKeyboardMarkup

	if len(listFriends) != 0 {
		listButtons := []tgbotapi.InlineKeyboardButton{}

		for key, user := range listFriends {
			nameFriend := user.FirstName + " " + user.LastName
			data := Pages.ADD_BY_ID + "_" + strconv.Itoa(key)
			listButtons = append(listButtons, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
		}

		keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(listButtons...),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.ADD_BY_ID),
				tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
			),
		)
	} else {
		keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.ADD_BY_ID),
				tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
			),
		)
	}
	return &keyboardSurvaillance
}

// клавиатура для "Добавить в отслеживаемые"
func GetAddToTrackedKeyboard() *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("По имени (в твоих друзьях)", Pages.ADD_BY_NAME),
			tgbotapi.NewInlineKeyboardButtonData("По ID", Pages.ADD_BY_ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Помощь", Pages.HELP_SURVAILlANCE),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.SURVEILLANCE),
			tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
		),
	)
	return &keyboardSurvaillance
}

// клавиатура для "Добавить в отслеживаемые"
func GetAddedAndDeletedFriendsKeyboard(tpID int) *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.INFO_ABOUT_FRIENDS),
			tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
		),
	)
	return &keyboardSurvaillance
}

// клавиатура для "Убрать из отслеживаемых"
func GetDeleteFromTrackedKeyboard() *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.SURVEILLANCE),
			tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
		),
	)
	return &keyboardSurvaillance
}

// клавиатура для "Добавить по имени"
func GetAddByNameKeyboard() *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.ADD_TO_TRACKED),
			tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
		),
	)
	return &keyboardSurvaillance
}

// клавиатура для "Добавить по имени определенного человека"
func GetAddByNameUserKeyboard(id string) *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить", Pages.ADD_IN_TRACKED+"_"+id),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.ADD_BY_NAME),
			tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
		),
	)
	return &keyboardSurvaillance
}

// клавиатура для "Добавить по id определенного человека"
func GetAddByIDUserKeyboard(id string) *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить", Pages.ADD_IN_TRACKED+"_"+id),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.ADD_BY_ID),
			tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
		),
	)
	return &keyboardSurvaillance
}

// клавиатура для "Человек добавлен в отслеживаемые"
func GetAddedInTrackedKeyboard() *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.ADD_BY_NAME),
			tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
		),
	)
	return &keyboardSurvaillance
}

// клавиатура для "Человек добавлен в отслеживаемые"
func GetInfoAboutFriendsKeyboard(id int) *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Новый запрос", Pages.INFO_ABOUT_FRIENDS_NEW),
			tgbotapi.NewInlineKeyboardButtonData("История изменений", Pages.INFO_ABOUT_FRIENDS_HISTORY),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.TRACKED+"_"+strconv.Itoa(id)),
			tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
		),
	)
	return &keyboardSurvaillance
}

// клавиатура для "Отслеживаемый"
func GetTrackedPersonKeyboard() *tgbotapi.InlineKeyboardMarkup {
	var keyboardSurvaillance = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Друзья", Pages.INFO_ABOUT_FRIENDS),
			tgbotapi.NewInlineKeyboardButtonData("Лайки", Pages.LIKES),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Убрать из слежки", Pages.DELETE_FROM_SURVEILLANCE),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.SURVEILLANCE),
			tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
		),
	)
	return &keyboardSurvaillance
}

// клавиатура для "Авторизация вк"
func GetTokenVKKeyboard() *tgbotapi.InlineKeyboardMarkup {
	var keyboardMain = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Получить токен ВК", "https://oauth.vk.com/authorize?client_id=8157407&redirect_uri=https://oauth.vk.com/blank.html&scope=friends,photos,offline&response_type=token&v=5.131"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", Pages.MAIN),
			tgbotapi.NewInlineKeyboardButtonData("На главную", Pages.MAIN),
		),
	)
	return &keyboardMain
}
