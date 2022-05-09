package Messages

type TypeMessage struct {
	TokenVK      bool
	SearchByName bool
	SearchByID   bool
}

func NewTypeMessage() *TypeMessage {
	typeMes := TypeMessage{}

	typeMes.TokenVK = false
	typeMes.SearchByName = false
	typeMes.SearchByID = false

	return &typeMes
}

func (typeMes *TypeMessage) ChangeTypeMessage(tok bool, searchByName bool, searchByID bool) {
	typeMes.TokenVK = tok
	typeMes.SearchByName = searchByName
	typeMes.SearchByID = searchByID
}
