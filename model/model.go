package model

type dbCommands interface {
	SelectMessages(map[string]interface{}) ([]*Message, error)
	InsertMessage(Message) error

	SelectUsers(User) ([]*User, error)
	CheckUser(User) (bool, error)
	AddUser(User) error
}

type Model struct {
	dbCommands
}

type User struct {
	Name string `json:name`
}

type Message struct {
	Message  string `json:"message"`
	Userfrom string `json:"userfrom"`
	Userto   string `json:"userto"`
	Date     int64  `json:"date"`
}

func New(cmds dbCommands) *Model {
	return &Model{dbCommands: cmds}
}

func (m *Model) Messages(dialogUsers map[string]interface{}) ([]*Message, error) {
	return m.SelectMessages(dialogUsers)
}

func (m *Model) CheckCurrentUser(currentUser User) (bool, error) {
	return m.CheckUser(currentUser)
}

func (m *Model) AddCurrentUser(currentUser User) error {
	return m.AddUser(currentUser)
}

func (m *Model) GetUsers(currentUser User) ([]*User, error) {
	return m.SelectUsers(currentUser)
}

func (m *Model) AddMessage(message Message) error {
	return m.InsertMessage(message)
}
