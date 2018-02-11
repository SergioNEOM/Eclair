package common

import (
	"fmt"
	"os"
)

/*
 *  Пользователи
 */

type UserRec struct {
	Account  string //hash
	Password string //hash
	Name     string
	Role     int
	Plane    *TrainingPlan
}

type UsersList map[string]*UserRec

//JSON-file
//   map[UserId]:UserRec
//   UserId -unique identifier

// ??? надо ли ???
func (u *UserRec) SetAccount(acc string) {
	u.Account = acc
}
func (u *UserRec) SetPassword(pass string) {
	u.Password = pass
}
func (u *UserRec) SetName(name string) {
	u.Name = name
}
func (u *UserRec) SetRole(role int) {
	u.Role = role
}

// ???

func (us UsersList) AddUser(acc, pass, name string, role int) string {
	//todo: возвращать error ?
	uid := GenUID(USERS_PREFIX)
	x := us[uid]
	if x != nil {
		return ""
	}
	us[uid] = &UserRec{acc, Hash(pass), name, role, NewTrainingPlan()}
	Flag_UsersChanged = true
	return uid
}

func (us UsersList) SetUserName(id, uname string) bool {
	x := us[id]
	if x != nil {
		return false
	}
	us[id].Name = uname
	Flag_UsersChanged = true
	return true
}

func (us UsersList) SetUserRole(id string, role int) bool {
	x := us[id]
	if x != nil {
		return false
	}
	us[id].Role = role
	Flag_UsersChanged = true
	return true
}

func (us UsersList) ListUserNames() []string {
	//todo: возвращать error ?
	s := []string{}
	for i := range us {
		s = append(s, us[i].Name)
	}
	return s
}

//todo: Find user by account --> ???  or (account+password)
func (us UsersList) Authorize(acc, pass string) string { //return user id or empty string
	fmt.Printf("Authorize:  acc=%s, pass=%s\n", acc, pass)
	for key := range us {
		if acc == us[key].Account && Hash(pass) == us[key].Password {
			return key
		}
	}
	return ""
}

func (us UsersList) LoadFromFile(filename string) bool {
	err := LoadJSON(filename, &us)
	if err == nil {
		Flag_UsersChanged = false
		return true
	}
	//todo:  log :
	fmt.Fprintf(os.Stderr, "*** Error on load users from file: %s\n", err)
	return false

}

func (us UsersList) SaveToFile(filename string, permiss os.FileMode) bool {
	err := SaveJSON(us, filename, permiss)
	if err == nil {
		Flag_UsersChanged = false
		return true // success
	}
	fmt.Fprintf(os.Stderr, "*** Error on save users to file: %s\n", err)
	//todo: выводить в журнал
	return false
}

func (us UsersList) LoadFromDb() bool {
	//todo:
	return false
}

func (us UsersList) SaveToDb() bool {
	//todo:
	return false
}
