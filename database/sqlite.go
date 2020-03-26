package database
// module implements work with SQLite database

import(

  "github.com/SergioNEOM/Eclair/models"

  "github.com/jinzhu/gorm"
 _ "github.com/jinzhu/gorm/dialects/sqlite"
)


func DBConnect() *gorm.DB {
  //todo: файл БД взять из настроек
    db, err := gorm.Open("sqlite3" , "eclair.db")
    if err != nil {
      //todo: Log
      panic("failed to connect database")
    }
    defer db.Close()

    //todo: оставить или убрать потом?
    db.AutoMigrate(&models.Users{})
    return db
}
