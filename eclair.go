package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"./common"

	"./handles"
	//	"./users"
	//----------
	"database/sql"
	_ "github.com/lib/pq"
)

type test1 struct{
    id int
    model string
}


func main() {
//----------------
	db, err := sql.Open("postgres", "user=eclair password=eclair dbname=eclair_db host=localhost sslmode=disable")
//	db, err := sql.Open("postgres", "user=eclair_boss password=author dbname=eclair_db host=localhost sslmode=disable")
	if err != nil {
		println("-*-*-*-*-*-   Error open database *-*-*-*-");
		return
	}
	defer db.Close()
/*	_, err := db.Exec("insert into eclair_schema.test1 (f1,f2) values ($1, $2)",11,"11-11-11-11-11-")
	if err != nil{
	        panic(err)
    	}
	println("result=",result)
*/
	rows, err := db.Query("select * from test1")
	if err != nil {
        	panic(err)
	}
    	defer rows.Close()
	products := []test1{}
     
	for rows.Next(){
        	p := test1{}
	        err := rows.Scan(&p.id, &p.model)
	        if err != nil{
        	    fmt.Println(err)
	            continue
	        }
        	products = append(products, p)
    	}
	for _, p := range products{
        	fmt.Println(p.id, p.model)
   	}
//---
	if err = db.Ping(); err != nil {
		println("-*-*-*-*-*-   Error ping database *-*-*-*-");
		return
	}
//----------------
	// 1. открыть конфиг
	// 2. открыть Users (в памяти хранить или по необходимости открывать файл ?)
	common.MainConf.SetDefaultConf()
	//	common.MainConf.SaveConf("defconf.cfg",0644)

	common.DebugLevel = 1

	common.CurrentUser = "UE555BF06-1516216311" //"123"
	common.AppRootDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	//	common.Users = make(map[string]*common.UserRec)
	if !common.Users.LoadFromFile("users.cfg") {
		common.Users.AddUser("123", "456", "бла-бла-бла", common.ROLE_Student)
		common.Users.AddUser("456", "900", "dfgdfgdfбла-бла-бла", common.ROLE_Student)
		common.Users.AddUser("789", "012", "-------а", common.ROLE_Student)
		common.Users.AddUser("admin", "admin", "admin-admin", common.ROLE_Admin)
	}
	fmt.Printf("%v\n", common.Users)

	//common.Courses[common.CurrentUser].LoadFromFile("courses.json")
	common.CurrentCourse.SetHeader("заголовок", "автор", "коммент")
	common.CurrentCourse.AddPara("Заголовок параграфа 1", "длинный-длинный текст", "", 0)
	common.CurrentCourse.AddPara("Заголовок параграфа 2", "очень и очень длинный текст", "1", 0)
	//common.CurrentCourse.SaveToFile("Course-1.json", 0644)

	fmt.Println("Users[cu]->", common.Users[common.CurrentUser])
	if common.Users[common.CurrentUser] != nil {
		u := common.ListOfCourses.AddCourse("111", "./courses/1.json")
		if common.Users[common.CurrentUser].Plane == nil {
			common.Users[common.CurrentUser].Plane = common.NewTrainingPlan()
		}
		common.Users[common.CurrentUser].Plane.AddCourse(common.CourseRec{Title: "Курс 123456", CourseLink: u})

		u = common.ListOfCourses.AddCourse("222", "./courses/1.json")
		//fmt.Println(common.ListOfCourses)
		//common.Users[common.CurrentUser].Plane.AddCourse(common.CourseRec{Title: "Курс 987654---###", CourseLink: u})
		common.Users[common.CurrentUser].Plane.AddCourse(common.CourseRec{Title: "Курс 987654---###", CourseLink: u})
		fmt.Println("users-->", common.Users[common.CurrentUser])
	}

	handles.SetHandles()

	fmt.Printf("Web server was started on %s...\n\t\t To stop press Ctrl-C\n", common.MainConf.Port)

	log.Fatal(http.ListenAndServe(common.MainConf.Port, nil))

}
