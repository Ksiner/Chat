package model

type dbCommands interface {
	SelectMessages() ([]*Message, error)
	InsertMessage(Message) error
}

type Model struct {
	dbCommands
}

// type User struct {
// 	Name string `json:name`
// }

// type FromToUsers struct {
// 	Userfrom User `json:userfrom`
// 	UserTo   User `json:userto`
// }

type Message struct {
	Message     string      `json:"message"`
	User string `json:"user"`
	Date        int64       `json:"date"`
	My      bool   `json:"my"`
}

func New(cmds dbCommands) *Model {
	return &Model{dbCommands: cmds}
}

func (m *Model) Messages() ([]*Message, error) {
	return m.SelectMessages()
}

func (m *Model) AddMessage(message Message) error {
	return m.InsertMessage(message)
}
