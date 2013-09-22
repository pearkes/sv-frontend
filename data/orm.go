package data

import (
	"errors"
	"fmt"
	"github.com/eaigner/hood"
	"log"
	"os"
)

var host = os.Getenv("HOST")

// The user table
type User struct {
	Id           hood.Id
	Name         string
	DropboxUid   string `validate:"presence"`
	DropboxToken string `validate:"presence"`
	FolderSum    string // We use this to check for dropbox folder changes
	SettingsRev  int    // We use this to check for dropbox setting changes
}

// Indexes for the user table
func (table *User) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("name_index", "name")
	indexes.AddUnique("dropbox_uid_index", "dropbox_uid")
}

type Orm struct {
	Hd *hood.Hood // The instance of the ORM
}

func NewOrm(connectionString string) *Orm {
	hd, err := hood.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database %s", err))
	}
	o := &Orm{hd}
	return o
}

func (o Orm) Create() {
	tx := o.Hd.Begin()
	err := tx.CreateTableIfNotExists(&User{})
	err = tx.Commit()
	if err != nil {
		panic(fmt.Sprintf("Failed to create user table %s", err))
	}
}

// Creates a new user
func (o Orm) NewUser(token string, uid string) error {
	tx := o.Hd.Begin()

	u := &User{
		Name:         uid + "." + host,
		DropboxToken: token,
		DropboxUid:   uid,
		SettingsRev:  999999, // something it couldn't be
		FolderSum:    "",
	}

	// Save
	_, err := tx.Save(u)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Persist
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Creates a new user
func (o Orm) UpdateName(u User, domain string, rev int) error {
	tx := o.Hd.Begin()

	u.Name = domain
	u.SettingsRev = rev

	// Save
	_, err := tx.Save(&u)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Persist
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (o Orm) UpdateSum(u User, sum string) error {
	tx := o.Hd.Begin()

	u.FolderSum = sum

	// Save
	_, err := tx.Save(&u)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Persist
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (o Orm) UpdateToken(u User, token string) error {
	tx := o.Hd.Begin()

	u.DropboxToken = token
	u.SettingsRev = 999999 // start again

	// Save
	_, err := tx.Save(&u)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Persist
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// A count of all of the users in the database
func (o Orm) UserCount() int64 {
	var count int64
	row := o.Hd.QueryRow("SELECT count(*) FROM \"user\";")
	err := row.Scan(&count)
	if err != nil {
		log.Printf("failed to retrieve user count: %s", err)
	}
	return count
}

// Retrieves the last user in the database
func (o Orm) LastUser() User {
	var us []User
	err := o.Hd.OrderBy("user.id").Desc().Limit(1).Find(&us)
	if err != nil {
		log.Printf("failed to retrieve last user: %s", err)
	}
	return us[0]
}

func (o Orm) FindByName(name string) (User, error) {
	var us []User
	err := o.Hd.Where("name", "=", name).Desc().Limit(1).Find(&us)
	if err != nil {
		log.Printf("failed to find user by name: %s", err)
	}
	if len(us) != 1 {
		// Er, what?
		return User{}, errors.New("find bad results for user query")
	}
	return us[0], nil
}

func (o Orm) FindByUid(uid string) (User, error) {
	var us []User
	err := o.Hd.Where("dropbox_uid", "=", uid).Desc().Limit(1).Find(&us)
	if err != nil {
		log.Printf("failed to find user by uid: %s", err)
		return User{}, err
	}
	if len(us) != 1 {
		// Er, what?
		return User{}, errors.New("find bad results for user query")
	}
	return us[0], nil
}
