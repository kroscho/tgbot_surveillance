package Postgesql

import (
	"database/sql"
	"fmt"
	"strconv"
	"telegram_test_bot/BL/VKApi"
	"telegram_test_bot/Models"
	"time"

	_ "github.com/lib/pq"
)

// проверить получил ли пользователь токен вк
func CheckUserToken(user *Models.UserSession) (bool, error) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf("select user_token, vk_id from users where tg_id = '%s'", strconv.Itoa(user.UserID_TG))
	var token sql.NullString
	var vk_id sql.NullInt32
	if err := db.QueryRow(query).Scan(&token, &vk_id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	fmt.Println("token: ", vk_id.Int32, token.String)
	user.Token = token.String
	user.UserID_VK = int(vk_id.Int32)
	return token.String != "", nil
}

// получить vkID пользователя с таблицы usersVK
func GetVKID(id_vk int) (int, error) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf("select vk_id from usersvk where id_user_vk = %d", id_vk)
	var vk_id sql.NullInt32
	if err := db.QueryRow(query).Scan(&vk_id); err != nil {
		if err == sql.ErrNoRows {
			return int(vk_id.Int32), nil
		}
		return int(vk_id.Int32), err
	}
	return int(vk_id.Int32), nil
}

// проверить есть ли пльзователь в базе, если нет, добавляем
func CheckUser(userID int, userName string, user *Models.UserSession) error {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf("select id_user from users where tg_id = '%s'", strconv.Itoa(userID))
	var id sql.NullInt32
	if err := db.QueryRow(query).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			_, err := db.Exec("insert into users (tg_id, user_name) values ($1, $2)", strconv.Itoa(userID), userName)
			if err != nil {
				panic(err)
			}
			fmt.Println(userName + " успешно добавлен")
			return nil
		}
		return err
	}
	AddTrackedPersonsInSession(user)
	return nil
}

// заполнить у пользователя токен и id вк
func AddUserTokenAndVkID(tgID int, token string, vkID int) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	result, err := db.Exec("update users set vk_id = $1, user_token = $2 where tg_id = $3", vkID, token, strconv.Itoa(tgID))
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected()) // количество обновленных строк
}

// добавить пользователя в таблицу UsersVK
func AddUserInUsersVK(user Models.User) (int, error) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var id_user_vk sql.NullInt32

	// добавить в таблицу UsersVK
	query := fmt.Sprintf("select id_user_vk from usersvk where vk_id = %d", user.UID)
	if err := db.QueryRow(query).Scan(&id_user_vk); err != nil {
		if err == sql.ErrNoRows {
			err := db.QueryRow("insert into usersvk (vk_id, first_name, last_name) values ($1, $2, $3) returning id_user_vk", user.UID, user.FirstName, user.LastName).Scan(&id_user_vk)
			if err != nil {
				panic(err)
			}
			//fmt.Println(user.FirstName + " " + user.LastName + " успешно добавлен в таблицу usersVK")
			return int(id_user_vk.Int32), nil
		}
		return int(id_user_vk.Int32), err
	}
	return int(id_user_vk.Int32), nil
}

// добавить пользователя в таблицу TrackedPersons
func AddUserInTrackedPersons(id_user_vk int) (int, error) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var id_tp sql.NullInt32

	// добавить в таблицу UsersVK
	query := fmt.Sprintf("select id_tp from trackedpersons where vk_id = %d", id_user_vk)
	if err := db.QueryRow(query).Scan(&id_tp); err != nil {
		if err == sql.ErrNoRows {
			err := db.QueryRow("insert into trackedpersons (vk_id) values ($1) returning id_tp", id_user_vk).Scan(&id_tp)
			if err != nil {
				panic(err)
			}
			fmt.Println("Успешно добавлен в таблицу TrackedPersons")
			return int(id_tp.Int32), nil
		}
		return int(id_tp.Int32), err
	}
	return int(id_tp.Int32), nil
}

// добавить пользователя в таблицу UserPerson
func AddUserInUserPersons(id_user int, id_tp int) (int, error) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var id sql.NullInt32

	// добавить в таблицу UserPerson
	query := fmt.Sprintf("select id from userperson where user_id = %d and tp_id = %d", id_user, id_tp)
	if err := db.QueryRow(query).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			_, err := db.Exec("insert into userperson (user_id, tp_id) values ($1, $2)", id_user, id_tp)
			if err != nil {
				panic(err)
			}
			fmt.Println("Успешно добавлен в таблицу UserPerson")
			return int(id.Int32), nil
		}
		return int(id.Int32), err
	}
	return int(id.Int32), nil
}

// получить пользователей у отслеживаемого
func GetUsersAtTheTracked(id_tp int) (bool, error) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var id sql.NullInt32

	// добавить в таблицу UserPerson
	query := fmt.Sprintf("select id from userperson where tp_id = %d", id_tp)
	if err := db.QueryRow(query).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Никем не отслеживается")
			return false, nil
		}
		return true, err
	}
	return true, nil
}

// получить пользователей у отслеживаемого
func GetHistoryAboutDeletedFriends(id_tp int) map[string][]Models.User {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var date_of_change string
	userDel := Models.User{}
	delFriends := make(map[string][]Models.User)

	// добавить в таблицу UserPerson
	query := "select date_of_change, usvk.vk_id, usvk.first_name, usvk.last_name" +
		" from historyinfoaboutfriendstp as h" +
		" inner join deletedfriendstp as del on h.id = del.history_id" +
		" inner join usersvk as usvk on del.vk_id = usvk.id_user_vk" +
		" where tp_id=" + strconv.Itoa(id_tp)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&date_of_change, &userDel.UID, &userDel.FirstName, &userDel.LastName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		delFriends[date_of_change] = append(delFriends[date_of_change], userDel)
	}
	return delFriends
}

// получить пользователей у отслеживаемого
func GetHistoryAboutAddedFriends(id_tp int) map[string][]Models.User {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var date_of_change string
	userDel := Models.User{}
	addFriends := make(map[string][]Models.User)

	// добавить в таблицу UserPerson
	query := "select date_of_change, usvk.vk_id, usvk.first_name, usvk.last_name" +
		" from historyinfoaboutfriendstp as h" +
		" inner join addedfriendstp as del on h.id = del.history_id" +
		" inner join usersvk as usvk on del.vk_id = usvk.id_user_vk" +
		" where tp_id=" + strconv.Itoa(id_tp)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&date_of_change, &userDel.UID, &userDel.FirstName, &userDel.LastName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		addFriends[date_of_change] = append(addFriends[date_of_change], userDel)
	}
	return addFriends
}

// Удалить данные таблицы прошлого списка друзей пользователя
func DeleteDataPrevFriends(id_tp int) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// удаляем данные с таблицы
	result, err := db.Exec("delete from personprevfriends where tp_id=$1", id_tp)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
}

// Удалить из прежнего списка удаленного друга
func DeleteDeletedUserFromPrevFriends(id_tp int, delUser Models.User) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// получить id удаленного друга с таблицы UsersVK
	id_user_vk, err := AddUserInUsersVK(delUser)
	if err != nil {
		panic(err)
	}
	// удаляем данные с таблицы
	result, err := db.Exec("delete from personprevfriends where tp_id=$1 and vk_id=$2", id_tp, id_user_vk)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
}

// Добавить в прежний список нового друга
func AddAddedUserInPrevFriends(id_tp int, addedUser Models.User) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// получить id удаленного друга с таблицы UsersVK
	id_user_vk, err := AddUserInUsersVK(addedUser)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("insert into personprevfriends (tp_id, vk_id) values ($1, $2)", id_tp, id_user_vk)
	if err != nil {
		panic(err)
	}
}

// Убрать человека из отслеживаемых
func DeleteFromTrackedPerson(user_id int, id_tp int, id_vk int) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// удаляем данные с таблицы UserPerson
	result, err := db.Exec("delete from userperson where user_id=$1 and tp_id=$2", user_id, id_tp)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
	isTracked, err := GetUsersAtTheTracked(id_tp)
	if err != nil {
		panic(err)
	}
	// если больше никем не отслеживается, удаляем все данных с этим персонажем
	if !isTracked {
		DeleteDataPrevFriends(id_tp)
		DeleteDataNewFriends(id_tp)
		// удаляем данные с таблицы TrackedPersons
		result, err = db.Exec("delete from trackedpersons where id_tp=$1 and vk_id=$2", id_tp, id_vk)
		if err != nil {
			panic(err)
		}
		fmt.Println(result.RowsAffected())
	}
}

// Удалить данные таблицы нового списка друзей пользователя
func DeleteDataNewFriends(id_tp int) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// удаляем данные с таблицы
	result, err := db.Exec("delete from personnewfriends where tp_id=$1", id_tp)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
}

// Возобновнить таблицу прошлого списка друзей пользователя
func AddIDPersonInPrevFriends(id_tp int, id_vk int) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("insert into personprevfriends (tp_id, vk_id) values ($1, $2)", id_tp, id_vk)
	if err != nil {
		panic(err)
	}
}

// Возобновнить таблицу нового списка друзей пользователя
func AddIDPersonInNewFriends(id_tp int, id_vk int) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("insert into personnewfriends (tp_id, vk_id) values ($1, $2)", id_tp, id_vk)
	if err != nil {
		panic(err)
	}
}

// Возобновнить таблицу нового списка друзей пользователя
func GetIDUserFromUsersVK(userVKID int) int {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var id_user sql.NullInt32
	// Взять id пользователя, для которого хотим добавить отслеживаемого
	query := fmt.Sprintf("select id_user from users where vk_id = %d", userVKID)
	if err := db.QueryRow(query).Scan(&id_user); err != nil {
		panic(err)
	}
	return int(id_user.Int32)
}

func CheckCountTrackedPersons(userSession *Models.UserSession, maxCount int) bool {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var count sql.NullInt32

	// Взять id пользователя, для которого хотим добавить отслеживаемого
	user_id := GetIDUserFromUsersVK(userSession.UserID_VK)

	// Взять id пользователя, для которого хотим добавить отслеживаемого
	query := fmt.Sprintf("select count(*) from userperson where user_id = %d", user_id)
	if err := db.QueryRow(query).Scan(&count); err != nil {
		panic(err)
	}
	return int(count.Int32) < maxCount
}

// добавить личность в отслеживаемые
func AddPersonInTracked(user Models.User, userVkID int, userSession *Models.UserSession) error {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// добавить в таблицу UsersVK
	id_user_vk, err := AddUserInUsersVK(user)
	if err != nil {
		panic(err)
	}

	// добавить в TrackedPersons
	id_tp, err := AddUserInTrackedPersons(id_user_vk)
	if err != nil {
		panic(err)
	}

	// Взять id пользователя, для которого хотим добавить отслеживаемого
	id_user := GetIDUserFromUsersVK(userSession.UserID_VK)

	// добавить в UserPerson
	_, err = AddUserInUserPersons(id_user, id_tp)
	if err != nil {
		panic(err)
	}

	listFriends, err := VKApi.GetFriendsFromUser(userSession, user)
	if err != nil {
		panic(err)
	}
	UpdateListFriends(id_tp, listFriends)
	return nil
}

// Добавить инфу о добавленных друзьях в историю изменений
func AddInfoAboutAddedFriendsInHistory(id_history int, addedListFriends map[int]Models.User) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, user := range addedListFriends {
		// добавить в таблицу UsersVK
		id_user_vk, err := AddUserInUsersVK(user)
		if err != nil {
			panic(err)
		}
		_, err = db.Exec("insert into addedfriendstp (history_id, vk_id) values ($1, $2)", id_history, id_user_vk)
		if err != nil {
			panic(err)
		}
	}
}

// Добавить инфу об удаленных друзьях в историю изменений
func AddInfoAboutDeletedFriendsInHistory(id_history int, deletedListFriends map[int]Models.User) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, user := range deletedListFriends {
		// добавить в таблицу UsersVK
		id_user_vk, err := AddUserInUsersVK(user)
		if err != nil {
			panic(err)
		}
		_, err = db.Exec("insert into deletedfriendstp (history_id, vk_id) values ($1, $2)", id_history, id_user_vk)
		if err != nil {
			panic(err)
		}
	}
}

// Добавить инфу об изменениях в историю изменений
func AddInfoAboutFriendsInHistory(id_tp int, addedListFriends map[int]Models.User, deletedListFriends map[int]Models.User) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	now := time.Now()
	date_of_change := now.String()[:16]
	fmt.Println(date_of_change)
	var id_history sql.NullInt32
	err = db.QueryRow("insert into historyinfoaboutfriendstp (tp_id, date_of_change) values ($1, $2) returning id", id_tp, date_of_change).Scan(&id_history)
	if err != nil {
		panic(err)
	}
	if len(addedListFriends) != 0 {
		AddInfoAboutAddedFriendsInHistory(int(id_history.Int32), addedListFriends)
	}
	if len(deletedListFriends) != 0 {
		AddInfoAboutDeletedFriendsInHistory(int(id_history.Int32), deletedListFriends)
	}
	fmt.Println("Изменения зафиксированы в истории успешно")
}

// обновить прошлый и новый список друзей отслеживаемого
func UpdateListFriends(id_tp int, listfriends []*Models.User) {
	DeleteDataNewFriends(id_tp)
	DeleteDataPrevFriends(id_tp)
	for _, friend := range listfriends {
		userNew := Models.User{}
		userNew.UID = friend.UID
		userNew.FirstName = friend.FirstName
		userNew.LastName = friend.LastName
		id_new, err := AddUserInUsersVK(userNew)
		if err != nil {
			panic(err)
		}
		AddIDPersonInPrevFriends(id_tp, id_new)
		AddIDPersonInNewFriends(id_tp, id_new)
	}
	fmt.Println("Друзья обновлены")
}

// обновить новый список друзей отслеживаемого
func UpdateNewListFriends(id_tp int, listfriends []*Models.User) map[int]Models.User {
	DeleteDataNewFriends(id_tp)
	newListFriendsMap := make(map[int]Models.User)
	for _, friend := range listfriends {
		userNew := Models.User{}
		userNew.UID = friend.UID
		userNew.FirstName = friend.FirstName
		userNew.LastName = friend.LastName
		newListFriendsMap[userNew.UID] = userNew
		id_new, err := AddUserInUsersVK(userNew)
		if err != nil {
			panic(err)
		}
		AddIDPersonInNewFriends(id_tp, id_new)
	}
	fmt.Println("Новый список друзей обновлен")
	return newListFriendsMap
}

// обновить прежний список друзей отслеживаемого
func UpdatePrevListFriends(id_tp int, listfriends []*Models.User) {
	DeleteDataPrevFriends(id_tp)
	for _, friend := range listfriends {
		userNew := Models.User{}
		userNew.UID = friend.UID
		userNew.FirstName = friend.FirstName
		userNew.LastName = friend.LastName
		id_new, err := AddUserInUsersVK(userNew)
		if err != nil {
			panic(err)
		}
		AddIDPersonInPrevFriends(id_tp, id_new)
	}
	fmt.Println("Прежний список друзей обновлен")
}

// Получить ids отслеживаемых людей у пользователя
func GetIdsTrackedPersons(userID int) ([]int, error) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf("select id_user from users where tg_id = '%s'", strconv.Itoa(userID))
	var id sql.NullInt32
	if err := db.QueryRow(query).Scan(&id); err != nil {
		panic(err)
	}
	// получить id отслеживаемых людей
	query = fmt.Sprintf("select tp_id from userperson where user_id = %d", id.Int32)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	trPersons := []int{}

	for rows.Next() {
		var tp int
		err := rows.Scan(&tp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		trPersons = append(trPersons, tp)
	}
	return trPersons, nil
}

// Добавить/обновить отслеживаемых людей у пользователя в сессию
func AddTrackedPersonsInSession(user *Models.UserSession) error {
	trIdsPersons, err := GetIdsTrackedPersons(user.UserID_TG)
	if err != nil {
		panic(err)
	}
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	user.Persons = nil
	for _, id_tp := range trIdsPersons {
		trPerson := Models.User{}
		query := fmt.Sprintf("select uvk.vk_id, first_name, last_name from trackedpersons as tp inner join usersvk as uvk on tp.vk_id = uvk.id_user_vk where id_tp = %d", id_tp)
		if err := db.QueryRow(query).Scan(&trPerson.UID, &trPerson.FirstName, &trPerson.LastName); err != nil {
			panic(err)
		}
		person := Models.NewPerson(trPerson)
		user.Persons = append(user.Persons, person)
	}
	return nil
}

// Получить прежний список друзей пользователя
func GetPrevListFriends(id_tp int) ([]int, error) {
	db, err := Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// получить список прежних друзей
	query := fmt.Sprintf("select vk_id from personprevfriends where tp_id = %d", id_tp)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	prevListFriends := []int{}

	for rows.Next() {
		var tp int
		err := rows.Scan(&tp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		vk_id, err := GetVKID(tp)
		if err != nil {
			panic(err)
		}
		prevListFriends = append(prevListFriends, vk_id)
	}
	return prevListFriends, nil
}
