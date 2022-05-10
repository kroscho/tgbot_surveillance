package Utils

import (
	"fmt"
	"telegram_test_bot/BL/VKApi"
	"telegram_test_bot/Models"
	Postgesql "telegram_test_bot/PostgeSQL"
)

// добавить в отслеживаемые если его еще нет в этом списке
func AddPersonInTracked(user *Models.UserSession, userAdd *Models.User) (Models.Person, bool) {
	person := Models.NewPerson(*userAdd)
	isExist := false
	for _, person := range user.Persons {
		if person.User.UID == userAdd.UID {
			isExist = true
			break
		}
	}
	if !isExist {
		user.Persons = append(user.Persons, person)
		Postgesql.AddPersonInTracked(*userAdd, user.UserID_VK, user)
	}
	return person, isExist
}

// выявить новых друзей на основе прежнего и нового списка друзей
func CheckDeletedAndNewFriends(user *Models.UserSession, newListFriends map[int]Models.User, prevListFriends []int, id_tp int) (map[int]Models.User, map[int]Models.User) {
	addedFriendsIds := make(map[int]Models.User)
	deletedFriendsIds := make(map[int]Models.User)
	i := 0
	for key, val := range newListFriends {
		isDeleted := false
		// проверяем удаленных друзей (если человек есть в прежнем списке, но нет в новом)
		if i < len(prevListFriends) {
			_, ok := newListFriends[prevListFriends[i]]
			if !ok {
				friend, err := VKApi.GetUserByID(user, prevListFriends[i])
				if err != nil {
					panic(err)
				}
				if friend != nil {
					fmt.Println(friend.UID)
					deletedFriendsIds[friend.UID] = *friend
					Postgesql.DeleteDeletedUserFromPrevFriends(id_tp, *friend)
					isDeleted = true
				}
			}
		}
		if !isDeleted && !CheckExistInList(key, prevListFriends) {
			addedFriendsIds[key] = val
			Postgesql.AddAddedUserInPrevFriends(id_tp, val)
		}
		i += 1
	}
	// если в прежнем списке еще остались ids, проверяем, есть ли они в новом (проверка удаленных друзей)
	if i < len(prevListFriends) {
		prevListPart := prevListFriends[i:]
		for _, idVK := range prevListPart {
			_, ok := newListFriends[idVK]
			if !ok {
				friend, err := VKApi.GetUserByID(user, prevListFriends[i])
				if err != nil {
					panic(err)
				}
				if friend != nil {
					fmt.Println(friend.UID)
					deletedFriendsIds[friend.UID] = *friend
				}
			}
		}
	}
	return addedFriendsIds, deletedFriendsIds
}

func CheckExistInList(x int, list []int) bool {
	for _, v := range list {
		if v == x {
			return true
		}
	}
	return false
}

// получить текст о добавленных и удаленных друзьях
func GetTextAboutAddedAndDeletedFriends(addedFriendsIds map[int]Models.User, deletedFriendsIds map[int]Models.User) string {
	text := ""
	if len(addedFriendsIds) == 0 && len(deletedFriendsIds) == 0 {
		text += "Пока изменений нет"
	} else {
		if len(addedFriendsIds) != 0 {
			text += "Новые друзья:\n"
		}
		for _, addedFriend := range addedFriendsIds {
			//fmt.Println("Новый друг: ", addedFriend.UID, addedFriend.FirstName, addedFriend.LastName)
			text += "Новый друг: " + addedFriend.FirstName + " " + addedFriend.LastName + "\n"
		}
		if len(deletedFriendsIds) != 0 {
			text += "Удаленные друзья:\n"
		}
		for _, deletedFriend := range deletedFriendsIds {
			//fmt.Println("Удаленный друг: ", deletedFriend.UID, deletedFriend.FirstName, deletedFriend.LastName)
			text += "Удаленный друг: " + deletedFriend.FirstName + " " + deletedFriend.LastName + "\n"
		}
	}
	return text
}

// получить текст истории о добавленных и удаленных друзьях
func GetTextHistoryFriends(addedFriendsIds map[string][]Models.User, deletedFriendsIds map[string][]Models.User) string {
	text := ""
	if len(addedFriendsIds) == 0 && len(deletedFriendsIds) == 0 {
		text += "История пуста"
	} else {
		for date, addedFriends := range addedFriendsIds {
			text += date + "\n Новыe друзья: "
			for _, addedfriend := range addedFriends {
				text += addedfriend.FirstName + "_" + addedfriend.LastName + "  "
			}
		}
		for date, deletedFriends := range deletedFriendsIds {
			text += "\n" + date + "\nУдаленные друзья: "
			for _, addedfriend := range deletedFriends {
				text += addedfriend.FirstName + "_" + addedfriend.LastName + "  "
			}
		}
	}
	return text
}
