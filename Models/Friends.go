package Models

type Friends struct {
	Count int     `json:"count"`
	Users []*User `json:"items"`
}

type Response struct {
	Response Friends `json:"response"`
}
