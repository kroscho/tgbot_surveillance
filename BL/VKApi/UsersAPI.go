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

func GetUserByID(user *Models.UserSession, userID int) (*Models.User, error) {
	config := Config.GetConfig()
	decToken, err := Encrypt.Decrypt(user.Token, config.Secret)
	if err != nil {
		fmt.Println("error decrypting your encrypted text: ", err)
	}
	q := VK_API_GET_USER_BY_ID_URL + "?access_token=" + decToken + "&user_ids=" + strconv.Itoa(userID) + "&v=" + VK_API_VERSION
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

	var response *Models.ResponseUsers
	json.Unmarshal(body, &response)
	if len(response.Response) != 0 {
		fmt.Println(response.Response[0].FirstName + " " + response.Response[0].LastName)
		return response.Response[0], nil
	} else {
		fmt.Println("Пользователь с таким id не найден")
		return nil, err
	}
}
