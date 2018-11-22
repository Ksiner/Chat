package db

import (
	"Chat/model"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	CnString string
}

type sqlCommands struct {
	conn                    *sqlx.DB
	sqlSelectCurrentUserCmd *sqlx.NamedStmt
	sqlSelectUsersCmd       *sqlx.NamedStmt
	sqlInsertUserCmd        *sqlx.NamedStmt
	sqlSelectMessagesCmd    *sqlx.NamedStmt
	sqlInsertMessageCmd     *sqlx.NamedStmt
}

func InitDbCmds(cfg Config) (*sqlCommands, error) {
	if dbConn, err := sqlx.Connect("postgres", cfg.CnString); err != nil {
		return nil, err
	} else {
		if err := dbConn.Ping(); err != nil {
			return nil, err
		}
		cmds := &sqlCommands{conn: dbConn}
		if err := cmds.PrepareStatements(); err != nil {
			return nil, err
		}
		return cmds, nil
	}
}

func (cmds *sqlCommands) PrepareStatements() error {
	if sqlSelectCurrentUserCmd, err := cmds.conn.PrepareNamed("Select true from public.\"User\" where \"user\"=:name limit 1"); err != nil {
		return err
	} else {
		cmds.sqlSelectCurrentUserCmd = sqlSelectCurrentUserCmd
	}
	if sqlSelectUsersCmd, err := cmds.conn.PrepareNamed("Select \"user\" as name from public.\"User\" where \"user\"!=:name"); err != nil {
		return err
	} else {
		cmds.sqlSelectUsersCmd = sqlSelectUsersCmd
	}
	if sqlInsertUserCmd, err := cmds.conn.PrepareNamed("Insert into \"User\" values(:name)"); err != nil {
		return err
	} else {
		cmds.sqlInsertUserCmd = sqlInsertUserCmd
	}
	if sqlSelectMessagesCmd, err := cmds.conn.PrepareNamed("Select message,\"userFrom\" as userfrom,\"userTo\" as userto,date FROM message where (\"userFrom\"=:currentuser and \"userTo\"=:targetuser) or (\"userFrom\"=:targetuser and \"userTo\"=:currentuser) ORDER BY date"); err != nil {
		return err
	} else {
		cmds.sqlSelectMessagesCmd = sqlSelectMessagesCmd
	}
	if sqlInsertMessageCmd, err := cmds.conn.PrepareNamed("INSERT INTO message (message,\"userFrom\",\"userTo\",date) VALUES (:message,:userfrom,:userto,:date)"); err != nil {
		return err
	} else {
		cmds.sqlInsertMessageCmd = sqlInsertMessageCmd
	}
	return nil
}

func (cmds *sqlCommands) SelectMessages(dialogUsers map[string]interface{}) ([]*model.Message, error) {
	messages := make([]*model.Message, 0)
	if err := cmds.sqlSelectMessagesCmd.Select(&messages, dialogUsers); err != nil {
		return nil, err
	}
	return messages, nil
}

func (cmds *sqlCommands) CheckUser(user model.User) (bool, error) {
	check := make([]bool, 0)
	if err := cmds.sqlSelectCurrentUserCmd.Select(&check, user); err != nil {
		return false, err
	}
	if cap(check) == 0 {
		return false, nil
	}
	return true, nil
}

func (cmds *sqlCommands) AddUser(user model.User) error {
	if _, err := cmds.sqlInsertUserCmd.Exec(user); err != nil {
		return err
	}
	return nil
}

func (cmds *sqlCommands) SelectUsers(currentUser model.User) ([]*model.User, error) {
	users := make([]*model.User, 0)
	if err := cmds.sqlSelectUsersCmd.Select(&users, currentUser); err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return users, nil
}

func (cmds *sqlCommands) InsertMessage(message model.Message) error {
	if _, err := cmds.sqlInsertMessageCmd.Exec(message); err != nil {
		return err
	}
	return nil
}
