package Models

type PostgreSQL struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
	SslMode  string `json:"sslmode"`
}

type Config struct {
	TelegramBotToken       string
	Secret                 string
	MaxCountTrackedPersons int
	PostgreSQL             PostgreSQL
}

// тип слежки
type TypeSurvaillance struct {
	NewFriends    bool // появление новых друзей
	DeleteFriends bool // удаление друзей
	Likes         bool // появление новых лайков
}

// предыдущий, новый списки и дата обновления
type ListUpdates struct {
	PrevList   []int
	NewList    []int
	DateUpdate string
}

// отслеживаемый человек
type Person struct {
	User          User
	NewFriends    []ListUpdates
	DeleteFriends []ListUpdates
	Photos        []Photo
	LikesPersons  []LikesPersonsUpdate // люди, у которых хотим просматривать, лайкает ли наша жертва их фотографии
}

func NewPerson(user User) Person {
	person := Person{}
	person.User = user
	return person
}

// фотография, у коророй прослеживаются лайки
type Photo struct {
	IdImage int
	Likes   []ListUpdates
}

type LikesPersonsUpdate struct {
	IdPerson int
	Photos   []LikesPhotos
}

type LikesPhotos struct {
	IdPhoto int
	IsLike  bool
}
