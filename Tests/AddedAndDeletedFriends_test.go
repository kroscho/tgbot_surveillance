package Tests

import (
	"fmt"
	Utils "telegram_test_bot/BL"
	"telegram_test_bot/Models"
	"testing"
)

func TestGetAddedAndDeletedFriends(t *testing.T) {
	newListFriends := make(map[int]Models.User)
	prevListFriends := []int{}
	user1 := Models.User{UID: 1116939, FirstName: "Андрей", LastName: "Дураков"}
	user2 := Models.User{UID: 2750942, FirstName: "Сергей", LastName: "Байгулов"}
	user3 := Models.User{UID: 3561824, FirstName: "Алексей", LastName: "Тютин"}
	user4 := Models.User{UID: 3708321, FirstName: "Михаил", LastName: "Крутиков"}
	user5 := Models.User{UID: 4311683, FirstName: "Иван", LastName: "Семериков"}
	user6 := Models.User{UID: 4723587, FirstName: "Михаил", LastName: "Пономарёв"}
	user7 := Models.User{UID: 10102303, FirstName: "Алексей", LastName: "Федосеев"}
	user8 := Models.User{UID: 5928203, FirstName: "Сергей", LastName: "Воронцов"}
	newListFriends[user1.UID] = user1
	newListFriends[user2.UID] = user2
	newListFriends[user3.UID] = user3
	newListFriends[user4.UID] = user4
	newListFriends[user5.UID] = user4
	newListFriends[user6.UID] = user6
	newListFriends[user7.UID] = user7
	newListFriends[user8.UID] = user8
	prevListFriends = append(prevListFriends, user1.UID)
	prevListFriends = append(prevListFriends, user2.UID)
	prevListFriends = append(prevListFriends, user3.UID)
	prevListFriends = append(prevListFriends, user4.UID)
	prevListFriends = append(prevListFriends, user5.UID)
	prevListFriends = append(prevListFriends, user6.UID)

	app := Models.NewServerApplication()
	user := app.GetUserSession(401948312)

	addedFriendsIds, deletedFriendsIds := Utils.CheckDeletedAndNewFriends(user, newListFriends, prevListFriends)

	for _, addedFriend := range addedFriendsIds {
		fmt.Println("Новый друг: ", addedFriend.UID, addedFriend.FirstName, addedFriend.LastName)
	}
	for _, deletedFriend := range deletedFriendsIds {
		fmt.Println("Удаленный друг: ", deletedFriend.UID, deletedFriend.FirstName, deletedFriend.LastName)
	}
}
