package database
// module implements work with SQLite database

import(

  "github.com/SergioNEOM/Eclair/models"

  "github.com/jinzhu/gorm"

)

var GDB *gorm.DB

type DBInterface interface {
    DBConnect() *gorm.DB
}
