package models


import(
  "github.com/jinzhu/gorm"
// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct{
  gorm.Model
  Name string `gorm:"type:varchar(30)"`
  Login string `gorm:"type:varchar(30);index:logins_idx;unique_index"`
  Pass string `gorm:"type:varchar(30)"`
  Token []byte `gorm:"type:varchar(130)"`
  Role int
}

type Users []User

type DBUsersInterface interface {
  UsersGetList() *Users
  GetUserById(id int) *User
  GetUserByLogin(login string) *User
  //GetUserByToken(t )

}
