package db

import (
	"Chat/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	CnString string
}

type sqlCommands struct {
	conn                    *sqlx.DB
	// sqlSelectCurrentUserCmd *sqlx.NamedStmt
	sqlSelectUsersCmd       *sqlx.NamedStmt
	// sqlInsertUserCmd        *sqlx.NamedStmt
	sqlSelectMessagesCmd    *sqlx.Stmt
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
	// if sqlSelectCurrentUserCmd, err := cmds.conn.PrepareNamed("Select \"user\" from public.\"User\" where \"user\"=:name"); err != nil {
	// 	return err
	// } else {
	// 	cmds.sqlSelectCurrentUserCmd = sqlSelectCurrentUserCmd
	// }
	if sqlSelectUsersCmd, err := cmds.conn.Preparex("Select \"user\" from public.\"User\""); err != nil {
		return err
	} else {
		cmds.sqlSelectUsersCmd = sqlSelectUsersCmd
	}
	// if sqlInsertUserCmd, err := cmds.conn.PrepareNamed("Insert into \"User\" values(:name)"); err != nil {
	// 	return err
	// } else {
	// 	cmds.sqlInsertUserCmd = sqlInsertUserCmd
	// }
	if sqlSelectMessagesCmd, err := cmds.conn.Preparex("Select message,\"user\" as user,date,my FROM message ORDER BY date"); err != nil {
		return err
	} else {
		cmds.sqlSelectMessagesCmd = sqlSelectMessagesCmd
	}
	if sqlInsertMessageCmd, err := cmds.conn.PrepareNamed("INSERT INTO message (message,\"user\",date,my) VALUES (:message,:user,:date,:my)"); err != nil {
		return err
	} else {
		cmds.sqlInsertMessageCmd = sqlInsertMessageCmd
	}
	return nil
}

func (cmds *sqlCommands) SelectMessages() ([]*model.Message, error) {
	messages := make([]*model.Message, 0)
	if err := cmds.sqlSelectMessagesCmd.Select(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (cmds *sqlCommands) InsertMessage(message model.Message) error {
	if _, err := cmds.sqlInsertMessageCmd.Exec(message); err != nil {
		return err
	}
	return nil
}
