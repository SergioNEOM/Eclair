package dbase;

import (
	"../common"
	//"errors"
	//--
	"fmt"
	"database/sql"
	//_ "github.com/lib/pq"
)

//
type test1 struct{
    id int
    model string
}


//----- test 1 -----------
func Db_test1() {
//**********
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
	products := []test1{} // сюда будем бросать построчно полученный набор строк типа test1

	for rows.Next(){
        	p := test1{} // одна строка типа test1
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
}
//----------------------------------------------------------------------------------------------
//


func DBConn () {

}


// GetCoursesList(StudentID int)
//
// ...


/* GetParaList - получить список параграфов указанного курса */
func GetParaList(CourseId int) (PL common.ParagraphList, err error) {
  // получить
	return //??
}

/* GetParagraph(ParaID int) (*Paragraph, Error) */
func GetParagraph(ParaID int) (P *common.Paragraph, err error) {
	// получить данные параграфа по его ID
	return //?
}
