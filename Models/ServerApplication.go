package Models

type ServerApplication struct {
	UserSession map[int]*UserSession
}

func NewServerApplication() ServerApplication {
	server := ServerApplication{}

	server.UserSession = make(map[int]*UserSession)
	return server
}

func (ServApp *ServerApplication) GetUserSession(userId int) *UserSession {
	_, ok := ServApp.UserSession[userId]
	if !ok {
		user := NewUserSession(userId)
		ServApp.UserSession[userId] = user
	}
	return ServApp.UserSession[userId]
}
