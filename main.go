package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	Utils "telegram_test_bot/BL"
	"telegram_test_bot/BL/Encrypt"
	"telegram_test_bot/BL/VKApi"
	"telegram_test_bot/Config"
	"telegram_test_bot/Models"
	Postgesql "telegram_test_bot/PostgeSQL"
	"telegram_test_bot/UI/Keyboard"
	"telegram_test_bot/UI/Messages"
	"telegram_test_bot/UI/Pages"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	app := Models.NewServerApplication()

	config := Config.GetConfig()

	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	photoBytes, err := ioutil.ReadFile("UI/images/1.jpg")
	if err != nil {
		panic(err)
	}
	photoFileBytes := tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: photoBytes,
	}

	// В канал updates будут приходить все новые сообщения.
	for update := range updates {
		if update.Message != nil {
			messText := update.Message.Text
			chatID := update.Message.Chat.ID
			messID := update.Message.MessageID
			userID := update.Message.From.ID
			userName := update.Message.From.UserName

			user := app.GetUserSession(userID)
			err := Postgesql.CheckUser(userID, userName, user)
			if err != nil {
				panic(err)
			}
			isToken, err := Postgesql.CheckUserToken(user)
			if err != nil {
				panic(err)
			}

			switch {
			// начало работы бота
			case messText == "/start":
				userName := update.Message.From.UserName
				msg := tgbotapi.NewPhotoUpload(chatID, photoFileBytes)
				if isToken {
					msg.Caption = fmt.Sprintf(`Привет %s, меня зовут mr. Kros, хочу помочь тебе в слежке`, userName)
				} else {
					msg.Caption = fmt.Sprintf(`Привет %s, меня зовут mr. Kros, хочу помочь тебе в слежке. \n%s`, userName, Messages.MAIN_NO_TOKEN_TEXT)
				}
				msg.ReplyMarkup = Keyboard.GetMainKeyboard(isToken)
				fmt.Println(chatID, messID)
				if user.CurMsg, err = bot.Send(msg); err != nil {
					panic(err)
				}
			// ожидаем сообщение о токене
			case user.TypeMessage.TokenVK:
				token, userId := Messages.CutAccessTokenAndUserId(messText)
				fmt.Println(userId, token)
				if token == "" || userId == "" {
					// удаляем сообщение через две секунды
					time.Sleep(1 * time.Second)
					msg := tgbotapi.NewDeleteMessage(chatID, messID)
					if _, err = bot.Send(msg); err != nil {
						panic(err)
					}
					// меняем caption у сообщения, сообщаем о том, что токен возникла ошибка
					msg1 := tgbotapi.NewEditMessageCaption(chatID, user.CurMsg.MessageID, Messages.GET_UNSUCCESS_VK_TOKEN)
					msg1.ReplyMarkup = Keyboard.GetTokenVKKeyboard()
					if _, err = bot.Send(msg1); err != nil {
						panic(err)
					}
				} else {
					// удаляем сообщение через две секунды
					time.Sleep(2 * time.Second)
					msg := tgbotapi.NewDeleteMessage(chatID, messID)
					if _, err = bot.Send(msg); err != nil {
						panic(err)
					}
					encToken, err := Encrypt.Encrypt(token, config.Secret)
					if err != nil {
						fmt.Println("error encrypting your classified text: ", err)
					}
					user.Token = encToken
					user.UserID_VK, err = strconv.Atoi(userId)
					if err != nil {
						panic(err)
					}
					Postgesql.AddUserTokenAndVkID(user.UserID_TG, user.Token, user.UserID_VK)
					// меняем caption у сообщения, сообщаем о том, что токен успешно сохранен
					msg1 := tgbotapi.NewEditMessageCaption(chatID, user.CurMsg.MessageID, Messages.GET_SUCCESS_VK_TOKEN)
					msg1.ReplyMarkup = Keyboard.GetTokenVKKeyboard()
					if _, err = bot.Send(msg1); err != nil {
						panic(err)
					}
					user.TypeMessage.ChangeTypeMessage(false, false, false)
				}

			// ожидаем имя человека, которого хотим найти
			case user.TypeMessage.SearchByName:
				time.Sleep(2 * time.Second)
				msg := tgbotapi.NewDeleteMessage(chatID, messID)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
				friends, err := VKApi.GetFriendsByName(user, messText)
				if err != nil {
					panic(err)

				}
				var msg1 tgbotapi.EditMessageCaptionConfig
				listFriends := make(map[int]Models.User)
				if len(friends) == 0 {
					// меняем caption у сообщения, сообщаем о том, не найдено по запросу ничего
					msg1 = tgbotapi.NewEditMessageCaption(chatID, user.CurMsg.MessageID, Messages.FRIENDS_BY_NAME_EMPTY)
					msg1.ReplyMarkup = Keyboard.GetFriendsByNameKeyboard(listFriends)
				} else {
					for _, friend := range friends {
						fmt.Println(friend.UID, " ", friend.FirstName, " ", friend.LastName)
						listFriends[friend.UID] = *friend
					}
					// меняем caption у сообщения, сообщаем о том, что поиск прошел успешно
					msg1 = tgbotapi.NewEditMessageCaption(chatID, user.CurMsg.MessageID, Messages.FRIENDS_BY_NAME)
					msg1.ReplyMarkup = Keyboard.GetFriendsByNameKeyboard(listFriends)
					user.TypeMessage.ChangeTypeMessage(false, false, false)
				}
				if _, err = bot.Send(msg1); err != nil {
					panic(err)
				}
			// ожидаем имя человека, которого хотим найти
			case user.TypeMessage.SearchByID:
				time.Sleep(2 * time.Second)
				msg := tgbotapi.NewDeleteMessage(chatID, messID)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
				var msg1 tgbotapi.EditMessageCaptionConfig
				listFriends := make(map[int]Models.User)
				idInt, err := strconv.Atoi(messText)
				if err != nil {
					msg1 = tgbotapi.NewEditMessageCaption(chatID, user.CurMsg.MessageID, Messages.ID_NOT_INT_ERROR)
					msg1.ReplyMarkup = Keyboard.GetFriendsByIDKeyboard(listFriends)
				} else {
					friend, err := VKApi.GetUserByID(user, idInt)
					if err != nil {
						panic(err)
					}
					if friend == nil {
						// меняем caption у сообщения, сообщаем о том, что пользователь с таким id не найден
						msg1 = tgbotapi.NewEditMessageCaption(chatID, user.CurMsg.MessageID, Messages.FRIENDS_BY_NAME_EMPTY)
						msg1.ReplyMarkup = Keyboard.GetFriendsByIDKeyboard(listFriends)
					} else {
						fmt.Println(friend.UID, " ", friend.FirstName, " ", friend.LastName)
						listFriends[friend.UID] = *friend
						// меняем caption у сообщения, сообщаем о том, что токен успешно сохранен
						msg1 = tgbotapi.NewEditMessageCaption(chatID, user.CurMsg.MessageID, Messages.FRIENDS_BY_NAME)
						msg1.ReplyMarkup = Keyboard.GetFriendsByIDKeyboard(listFriends)
						user.TypeMessage.ChangeTypeMessage(false, false, false)
					}
				}
				if _, err = bot.Send(msg1); err != nil {
					panic(err)
				}

			// если пользователь отправил сообщение просто так, удаляем его
			default:
				msg := tgbotapi.NewDeleteMessage(chatID, messID)
				if _, err = bot.Send(msg); err != nil {
					fmt.Println(err)
					panic(err)
				}
			}
		} else if update.CallbackQuery != nil {
			user := app.GetUserSession(update.CallbackQuery.From.ID)
			user.TypeMessage.ChangeTypeMessage(false, false, false)

			callQuery := update.CallbackQuery.Data
			chatID := update.CallbackQuery.Message.Chat.ID
			messageID := update.CallbackQuery.Message.MessageID

			switch {

			/* Главная страница */

			case callQuery == Pages.MAIN:
				fmt.Println(Pages.MAIN)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Pages.MAIN)
				msg.ReplyMarkup = Keyboard.GetMainKeyboard(true)
				if _, err = bot.Send(msg); err != nil {
					fmt.Println(err)
					panic(err)
				}
			case callQuery == Pages.SURVEILLANCE:
				fmt.Println(Pages.SURVEILLANCE)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Pages.TRACKED_PEOPLE)
				msg.ReplyMarkup = Keyboard.GetSurvaillanceKeyboard(user)

				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			case callQuery == Pages.SUBSCRIPTION:
				fmt.Println(Pages.SUBSCRIPTION)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Pages.SUBSCRIPTION)
				msg.ReplyMarkup = Keyboard.GetMainKeyboard(true)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			case callQuery == Pages.HELP_MAIN:
				fmt.Println(Pages.HELP_MAIN)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Pages.HELP_MAIN)
				msg.ReplyMarkup = Keyboard.GetMainKeyboard(true)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}

			case callQuery == Pages.TOKENVK:
				fmt.Println(Pages.TOKENVK)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Messages.GET_TOKEN_VK_TEXT)
				msg.ReplyMarkup = Keyboard.GetTokenVKKeyboard()
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
				user.TypeMessage.ChangeTypeMessage(true, false, false)

				/* Страница "Слежка" */

			case callQuery == Pages.HELP_SURVAILlANCE:
				fmt.Println(Pages.HELP_SURVAILlANCE)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Pages.HELP_SURVAILlANCE)
				msg.ReplyMarkup = Keyboard.GetSurvaillanceKeyboard(user)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			case callQuery == Pages.ADD_TO_TRACKED:
				fmt.Println(Pages.ADD_TO_TRACKED)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Pages.ADD_TO_TRACKED)
				msg.ReplyMarkup = Keyboard.GetAddToTrackedKeyboard()
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}

			case callQuery == Pages.ADD_BY_NAME:
				fmt.Println(Pages.ADD_BY_NAME)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Messages.ADD_BY_NAME)
				msg.ReplyMarkup = Keyboard.GetAddByNameKeyboard()
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
				user.TypeMessage.ChangeTypeMessage(false, true, false)

			case strings.Contains(callQuery, Pages.ADD_BY_NAME+"_"):
				fmt.Println(Pages.ADD_BY_NAME)
				runes := []rune(callQuery)
				id := string(runes[12:])
				idInt, err := strconv.Atoi(id)
				fmt.Println("ID: ", idInt)
				if err != nil {
					panic(err)
				}
				userAdd, err := VKApi.GetUserByID(user, idInt)
				if err != nil {
					panic(err)
				}
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Pages.ADD_BY_NAME+" "+strconv.Itoa(userAdd.UID))
				msg.ReplyMarkup = Keyboard.GetAddByNameUserKeyboard(id)
				if _, err = bot.Send(msg); err != nil {
					fmt.Println(err)
					panic(err)
				}

			case callQuery == Pages.ADD_BY_ID:
				fmt.Println(Pages.ADD_BY_ID)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Messages.ADD_BY_ID)
				msg.ReplyMarkup = Keyboard.GetAddByNameKeyboard()
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
				user.TypeMessage.ChangeTypeMessage(false, false, true)

			case strings.Contains(callQuery, Pages.ADD_BY_ID+"_"):
				fmt.Println(Pages.ADD_BY_ID)
				runes := []rune(callQuery)
				id := string(runes[10:])
				idInt, err := strconv.Atoi(id)
				fmt.Println("ID: ", idInt)
				if err != nil {
					panic(err)
				}
				userAdd, err := VKApi.GetUserByID(user, idInt)
				if err != nil {
					panic(err)
				}
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Pages.ADD_BY_NAME+" "+strconv.Itoa(userAdd.UID))
				msg.ReplyMarkup = Keyboard.GetAddByIDUserKeyboard(id)
				if _, err = bot.Send(msg); err != nil {
					fmt.Println(err)
					panic(err)
				}

			case strings.Contains(callQuery, Pages.ADD_IN_TRACKED+"_"):
				fmt.Println(Pages.ADD_IN_TRACKED)
				runes := []rune(callQuery)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, "Идет процесс добавления в отслеживаемые, пожождите.")
				if _, err = bot.Send(msg); err != nil {
					fmt.Println(err)
					panic(err)
				}
				id := string(runes[15:])
				idInt, err := strconv.Atoi(id)
				fmt.Println("ID: ", idInt)
				if err != nil {
					panic(err)
				}
				userAdd, err := VKApi.GetUserByID(user, idInt)
				if err != nil {
					panic(err)
				}
				person, isExist := Utils.AddPersonInTracked(user, userAdd)
				//var msg tgbotapi.EditMessageCaptionConfig
				if isExist {
					msg = tgbotapi.NewEditMessageCaption(chatID, messageID, person.User.FirstName+" "+person.User.LastName+" уже есть в отслеживаемых.")
					msg.ReplyMarkup = Keyboard.GetAddToTrackedKeyboard()
				} else {
					msg = tgbotapi.NewEditMessageCaption(chatID, messageID, person.User.FirstName+" "+person.User.LastName+" Успешно добавлен")
					msg.ReplyMarkup = Keyboard.GetAddedInTrackedKeyboard()
				}
				if _, err = bot.Send(msg); err != nil {
					fmt.Println(err)
					panic(err)
				}

				/* Страница "Отслеживаемый" */

			case strings.Contains(callQuery, "tracked_"):
				fmt.Println(Pages.TRACKED)
				runes := []rune(callQuery)
				id_vk := string(runes[8:])
				idVKInt, err := strconv.Atoi(id_vk)
				if err != nil {
					panic(err)
				}
				trackedPerson, err := VKApi.GetUserByID(user, idVKInt)
				if err != nil {
					panic(err)
				}
				user.CurTrackedPerson = *trackedPerson

				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, Pages.TRACKED+" "+id_vk)

				msg.ReplyMarkup = Keyboard.GetTrackedPersonKeyboard()
				if _, err = bot.Send(msg); err != nil {
					fmt.Println(err)
					panic(err)
				}

				/* Меню о новых/удаленных друзьях */
			case callQuery == Pages.INFO_ABOUT_FRIENDS:
				fmt.Println(Pages.INFO_ABOUT_FRIENDS)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, "Вы можете получить актуальную информацию о друзьях либо посмотреть историю изменений за все время.")
				msg.ReplyMarkup = Keyboard.GetInfoAboutFriendsKeyboard(user.CurTrackedPerson.UID)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}

				/* Запрашвиваем инфу о новых/удаленных друзьях */
			case callQuery == Pages.INFO_ABOUT_FRIENDS_NEW:
				fmt.Println(Pages.INFO_ABOUT_FRIENDS_NEW)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, "Идет процесс проверки друзей, подождите немного. \nВремя зависит от количества друзей.")
				//msg.ReplyMarkup = Keyboard.GetAddedAndDeletedFriendsKeyboard(user.CurTrackedPersonID)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
				newListFriends, err := VKApi.GetFriendsFromUser(user, user.CurTrackedPerson)
				if err != nil {
					panic(err)
				}
				// получить id_vk с таблицы UsersVK
				id_user_vk, err := Postgesql.AddUserInUsersVK(user.CurTrackedPerson)
				if err != nil {
					panic(err)
				}

				// добавить в TrackedPersons
				id_tp, err := Postgesql.AddUserInTrackedPersons(id_user_vk)
				if err != nil {
					panic(err)
				}
				newListFriendsMap := Postgesql.UpdateNewListFriends(id_tp, newListFriends)
				prevListFriends, err := Postgesql.GetPrevListFriends(id_tp)
				if err != nil {
					panic(err)
				}
				addedFriendsIds, deletedFriendsIds := Utils.CheckDeletedAndNewFriends(user, newListFriendsMap, prevListFriends)
				// если есть какие то изменения в друзьях, то обновляем прежний список друзей на новый
				if len(addedFriendsIds) != 0 || len(deletedFriendsIds) != 0 {
					msg = tgbotapi.NewEditMessageCaption(chatID, messageID, "Обновляем данные, осталось немного.")
					if _, err = bot.Send(msg); err != nil {
						panic(err)
					}
					Postgesql.UpdatePrevListFriends(id_tp, newListFriends)
					Postgesql.AddInfoAboutFriendsInHistory(id_tp, addedFriendsIds, deletedFriendsIds)
				}
				text := Utils.GetTextAboutAddedAndDeletedFriends(addedFriendsIds, deletedFriendsIds)
				msg = tgbotapi.NewEditMessageCaption(chatID, messageID, text)
				msg.ReplyMarkup = Keyboard.GetAddedAndDeletedFriendsKeyboard(user.CurTrackedPerson.UID)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}

				/* Запрашвиваем инфу о новых/удаленных друзьях */
			case callQuery == Pages.INFO_ABOUT_FRIENDS_HISTORY:
				fmt.Println(Pages.INFO_ABOUT_FRIENDS_HISTORY)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, "Идет процесс составления истории.")
				//msg.ReplyMarkup = Keyboard.GetAddedAndDeletedFriendsKeyboard(user.CurTrackedPersonID)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
				// получить id_vk с таблицы UsersVK
				id_user_vk, err := Postgesql.AddUserInUsersVK(user.CurTrackedPerson)
				if err != nil {
					panic(err)
				}

				// получить id в TrackedPersons
				id_tp, err := Postgesql.AddUserInTrackedPersons(id_user_vk)
				if err != nil {
					panic(err)
				}

				addedFriends := Postgesql.GetHistoryAboutAddedFriends(id_tp)
				deletedFriends := Postgesql.GetHistoryAboutDeletedFriends(id_tp)
				text := Utils.GetTextHistoryFriends(addedFriends, deletedFriends)
				msg = tgbotapi.NewEditMessageCaption(chatID, messageID, text)
				msg.ReplyMarkup = Keyboard.GetAddedAndDeletedFriendsKeyboard(user.CurTrackedPerson.UID)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}

				/* Убрать из отслеживаемых */
			case callQuery == Pages.DELETE_FROM_SURVEILLANCE:
				fmt.Println(Pages.DELETE_FROM_SURVEILLANCE)
				msg := tgbotapi.NewEditMessageCaption(chatID, messageID, "Идет процесс удаления из отслеживаемых, подождите")
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
				// получить id_vk с таблицы UsersVK
				id_user_vk, err := Postgesql.AddUserInUsersVK(user.CurTrackedPerson)
				if err != nil {
					panic(err)
				}
				// получить id в TrackedPersons
				id_tp, err := Postgesql.AddUserInTrackedPersons(id_user_vk)
				if err != nil {
					panic(err)
				}
				// Взять id пользователя, у которого хотим убрать отслеживаемого
				id_user := Postgesql.GetIDUserFromUsersVK(user.UserID_VK)
				Postgesql.DeleteFromTrackedPerson(id_user, id_tp, id_user_vk)
				Postgesql.AddTrackedPersonsInSession(user)
				msg = tgbotapi.NewEditMessageCaption(chatID, messageID, "Пользователь убран из отслеживаемых успешно.")
				msg.ReplyMarkup = Keyboard.GetDeleteFromTrackedKeyboard()
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			}
		}
	}
}
