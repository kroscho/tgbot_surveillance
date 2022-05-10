package VKApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"telegram_test_bot/BL/Encrypt"
	"telegram_test_bot/Config"
	"telegram_test_bot/Models"
)

func GetFriendsByName(user *Models.UserSession, query string) ([]*Models.User, error) {
	config := Config.GetConfig()
	decToken, err := Encrypt.Decrypt(user.Token, config.Secret)
	if err != nil {
		fmt.Println("error decrypting your encrypted text: ", err)
	}
	q := VK_API_GET_FRIENDS_BY_NAME_URL + "?access_token=" + decToken + "&user_id=" + strconv.Itoa(user.UserID_VK) + "&q=" + query + "&count=10" + "&v=" + VK_API_VERSION
	fmt.Println(q)
	resp, err := http.Get(q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *Models.Response
	json.Unmarshal(body, &response)
	for _, friend := range response.Response.Users {
		fmt.Println(friend.UID, " ", friend.FirstName, " ", friend.LastName)
	}
	return response.Response.Users, nil
}

func GetFriendsFromUser(userSession *Models.UserSession, user Models.User) ([]*Models.User, error) {
	config := Config.GetConfig()
	decToken, err := Encrypt.Decrypt(userSession.Token, config.Secret)
	if err != nil {
		fmt.Println("error decrypting your encrypted text: ", err)
	}
	q := VK_API_GET_FRIENDS_FROM_USER_URL + "?access_token=" + decToken + "&user_id=" + strconv.Itoa(user.UID) + "&fields=nickname" + "&v=" + VK_API_VERSION
	fmt.Println(q)
	resp, err := http.Get(q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *Models.Response
	json.Unmarshal(body, &response)
	return response.Response.Users, nil
}
