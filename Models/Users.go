package Models

type User struct {
	UID                     int    `json:"id"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Sex                     int    `json:"sex"`
	Nickname                string `json:"nickname"`
	ScreenName              string `json:"screen_name"`
	BDate                   string `json:"bdate"`
	Photo                   string `json:"photo"`
	PhotoMedium             string `json:"photo_medium"`
	PhotoBig                string `json:"photo_big"`
	Photo50                 string `json:"photo_50"`
	Photo100                string `json:"photo_100"`
	HasMobile               int    `json:"has_mobile"`
	Online                  int    `json:"online"`
	CanPost                 int    `json:"can_post"`
	CanSeeAllPosts          int    `json:"can_see_all_posts"`
	CanWritePrivateMessages int    `json:"can_write_private_message"`
	Status                  string `json:"activity"`
	Hidden                  int    `json:"hidden"`
	Deactivated             string `json:"deactivated"`
	Relation                int    `json:"relation"`
}

type Users struct {
	Count int     `json:"count"`
	Users []*User `json:"items"`
}

type ResponseUsers struct {
	Response []*User `json:"response"`
}
