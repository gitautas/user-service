package storage

// Decided to use MySQL as for such a simple schema a relational
// database is very simple and effective. Perhaps I could've used
// an ORM to simplify it further but this is the way that I am familiar with.

import (
	"database/sql"
	"fmt"
	"user-service/src/models"

	_ "github.com/go-sql-driver/mysql" // Imports the MySQL driver.
)

type Database interface {
	Connect() error
	Disconnect() error
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	RemoveUser(userID string) error
	GetUser(userID string) (file *models.User, err error) // Mixing named and unnamed parameters could be construed as inconsistent, but is more readable in my opinion.
	GetUsers(limit int, offset int) (users []*models.User, err error)
}

type Mysql struct {
	UserName     string
	Password     string
	Endpoint     string // IP and port, eg. 127.0.0.1:3306
	DatabaseName string
	Conn         *sql.DB
}

func (m *Mysql) Connect() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", m.UserName, m.Password, m.Endpoint, m.DatabaseName)) // ex: root:test@tcp(127.0.0.1:3306)/faceit?parseTime=true
	if err != nil {
		return err
	}
	m.Conn = db
	return m.Conn.Ping() // Establishes and checks the connection.
}

func (m *Mysql) Disconnect() error {
	if m.Conn != nil {
		return m.Conn.Close()
	}
	return nil
}

func (m *Mysql) CreateUser(user *models.User) (err error) {
	fmt.Println(user)
	_, err = m.Conn.Exec("INSERT INTO `user` (`user_id`, `first_name`, `last_name`, `nickname`, `password`, `email`, `country`) " +
		" VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.FirstName, user.LastName, user.Nickname, user.Password, user.Email, user.Country)
	if err != nil {
		fmt.Println(err) //FIXME
		return err
	}

	return nil
}

func (m *Mysql) UpdateUser(user *models.User) (err error) {
	_, err = m.Conn.Exec("UPDATE `user` " +
		"SET `first_name`=?, `last_name`=?, `nickname`=?, `password`=?, `email`=?, `country`=? " +
		"WHERE `user_id`=?",
		user.FirstName, user.LastName, user.Nickname, user.Password, user.Email, user.Country, user.ID)
	if err != nil {
		fmt.Println(err) //FIXME
		return err
	}
	return nil
}

func (m *Mysql) RemoveUser(userID string) (err error) {
	_, err = m.Conn.Exec("DELETE FROM `user` WHERE `user_id`=?", userID)
	if err != nil {
		fmt.Println(err) //FIXME
		return err
	}
	return nil
}

func (m *Mysql) GetUser(userID string) (user *models.User, err error) {
	rows, err := m.Conn.Query("SELECT " + // I specify which fields to list so that in the case of a schema change this would not break.
		"`first_name`, " +
		"`last_name`, " +
		"`nickname`, " +
		"`password`, " +
		"`email`, " +
		"`country`, " +
		"`created_at`, " +
		"`updated_at` " +
		"FROM `user` WHERE `user_id` = ?", userID)
	if err != nil {
		fmt.Println(err) //FIXME
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		user = &models.User{
			ID:        userID,
		}
		err = rows.Scan(&user.FirstName, &user.LastName, &user.Nickname, &user.Password, &user.Email, &user.Country, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Println(err) //FIXME
			return nil, err
		}
	}

	if rows.Err() != nil {
		fmt.Println(err) //FIXME
		return nil, err
	}

	return user, nil
}

func (m *Mysql) GetUsers(limit int, offset int) (users []*models.User, err error) {
	users = []*models.User{}

	rows, err := m.Conn.Query("SELECT " +
		"`user_id`, " +
		"`first_name`, " +
		"`last_name`, " +
		"`nickname`, " +
		"`password`, " +
		"`email`, " +
		"`country`, " +
		"`created_at`, " +
		"`updated_at` " +
		"FROM `user` LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		fmt.Println(err) //FIXME
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.Password, &user.Email, &user.Country, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Println(err) //FIXME
			return nil, err
		}

		if err != nil {
			fmt.Println(err) //FIXME
			return nil, err
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		fmt.Println(err) //FIXME
		return nil, err
	}

	return users, nil
}
