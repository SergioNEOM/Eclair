package handles

import (
	"encoding/json"
	"fmt"
	"html" //todo: временно
	"log"
	"net/http"
	//	"net/url"
	//
	"../common"
	//	"../users"
	"../templates"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>Eclair root<br><br><a href='/auth'>Вход</a></body><html>")
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./img/favicon.ico")
}

// func handleImages(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.URL.Path)
// }

func handleAuth(w http.ResponseWriter, r *http.Request) {
	var key, href string
	if r.Method == "POST" {
		// 1.получить логин, пароль из r.Post
		// проверить: если нет, GoAuthorize
		// получить UserId
		// открыть главную страницу (с учетом ролей)
		r.ParseForm()
		key = common.Users.Authorize(r.PostFormValue("Uname"), r.PostFormValue("Upass"))
		if key != "" {
			common.CurrentUser = key
			fmt.Printf("Current user: %s - Role:%d\n", key, common.Users[key].Role)
			switch common.Users[key].Role {
			case common.ROLE_Student:
				href = "/studentview/"
			case common.ROLE_Teacher:
				href = "/teacherview/"
			case common.ROLE_Admin:
				href = "/usersview/" //	"/adminview/"
			default:
				href = "/"
			}
			// проверить параметры входа
			// if users.Authorize(Umane,Upass) {
			// 	SetCurrentUser
			// 	SetSessionId
			// 	SetCookie
			http.SetCookie(w, &http.Cookie{Name: "Eclair1cookie", Value: common.CurrentUser, MaxAge: 60})
			//	Redirect("/mainview")
			http.Redirect(w, r, href, http.StatusTemporaryRedirect /*307*/)
		}
	}

	err := templates.GoAuthorize(&w)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		log.Fatalf("--- Error: %s   -----", err)
	}
}

var handleStudentView = func(w http.ResponseWriter, r *http.Request) {
	var err error
	// co, err := r.Cookie("Eclair1cookie")
	// if err!=nil {
	// 	if common.DebugLevel>0 {log.Println("---Handle MainView: Error on get cookie -  %s   -----",err)}
	// 	http.Redirect(w,r,"/", 307)
	// 	return
	// }
	// есть печенька - надо сравнить (простенькая защита от "умников")
	// if co.Value !=common.CurrentUser {
	// 	co.MaxAge=-1
	// 	co.Value="-"
	// 	http.SetCookie(w, co)
	// 	err = templates.GoAuthorize(&w)
	// 	if err!=nil {
	// 		if common.DebugLevel>0 {log.Println("--- Handle MainView Error: %s   -----",err)}
	// 	}
	// 	return
	// }
	err = templates.GoStudentView(&w)
	if err != nil {
		if common.DebugLevel > 0 {
			log.Printf("--- Handle MainView Error: %s   -----", err)
		}
	}
}

/*
 *
 */
var handleUsersView = func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("action") == "saveusers" {
		_ = common.Users.SaveToFile(common.USERS_FILE, 0644)
		//todo: обработать случай ошибки сохранения
	}
	if r.URL.Query().Get("action") == "edit" {
		s := r.URL.Query().Get("uid")
		fmt.Printf("^Uid:%s-> %s\n", s, common.Users[s].Name)
	}
	err := templates.GoUsersView(&w)
	if err != nil {
		if common.DebugLevel > 0 {
			log.Printf("--- Handle UsersView Error: %s   -----", err)
		}
	}
}

/*
 */

var handleAddUser = func(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method == "POST" {
		r.ParseForm()
		common.Users.AddUser(r.PostFormValue("Ulogin"), r.PostFormValue("Upass"), r.PostFormValue("Uname"), common.ROLE_Student)
		http.Redirect(w, r, "/usersview/", http.StatusTemporaryRedirect /*307*/)
	}

	//	if r.Method=="GET" { ??? templates.ParseFiles("adduser.html")

	//http.Redirect(w,r,"/", 307 )

	err = templates.GoAddUser(&w)
	if err != nil {
		fmt.Fprintf(w, "AddUser Error: %s", err)
		log.Fatalf("--- AddUser Error: %s   -----", err)
	}
}

var handleParaView = func(w http.ResponseWriter, r *http.Request) {
	var err error
	var uid, a string
	var pview common.ParaView
	/*
		1. Показать параграф
		2. по значению action перейти к пред/след параграфу (если не контрольный вопрос или не крайний параграф)

		Если в запросе есть параметр action=start, то проверить тип запроса: для GET - стартовать новый курс (uid!=""),
			 для POST - записать ошибку в журнал, т.к. новый курс можно начать только из списка (studentview)

	*/
	uid = r.URL.Query().Get("uid")
	//--
	a = r.URL.Query().Get("action")
	fmt.Printf("ParaView  handle: METHOD=%s , uid=%s, action=%s\n", r.Method, uid, a)
	//
	if r.Method == http.MethodGet {
		if a == "start" {
			if uid == "" {
				fmt.Println(" --> ParaView Error: empty uid in get request")
				fmt.Fprint(w, "<h1>ParaView Error: empty uid in get request</h1>")
				return
			}
			//todo: начало вывода (или продолжение) курса. Как отловить случайный переход?
			c := common.ListOfCourses.GetCourse(uid)
			if c != nil {
				if common.CurrentCourse.LoadFromFile(c.FName) {
					common.CurrentPara = -1 //todo: а если начать не сначала?
					fmt.Printf("Старт курса id:%s из файла %s\n", uid, c.FName)
					//пришли первый раз - тогда покажем страницу полностью
					// err = templates.GoParaView(&w, common.CurrentCourse.Para[common.CurrentPara])
					pview = common.ParaView{ParaCurNum: -1, PrevBut: false, NextBut: true}
					pview.Header = "Вводная информация"
					pview.Text = fmt.Sprintf("dfd;flgd;lf\ndfsdd\nЫАЫВАЫВАЫВАЫВАЫВА ыв ыВАЫВАЫ\nfsd")
					err = templates.GoParaView(&w, &pview)
					if err != nil {
						fmt.Fprintf(w, "ParaView Error: %s", err)
						log.Fatalf("--- ParaView Error: %s   -----\n", err)
					}
				} else { // не загрузился из файла
					fmt.Fprint(w, "ParaView Error: Course is empty or not loaded")
					fmt.Printf("ParaView: Corse %s  is empty or not loaded\n", uid)
				}
			} else {
				fmt.Printf("ParaView: CoursesList[%s] = nil\n", uid)
			}
		} else {
			fmt.Fprint(w, "ParaView Error: GET method without action=start")
			fmt.Println("ParaView Error: GET method without action=start")
			return
		}
	}
	if r.Method == http.MethodPost {
		if a == "exit" {
			//todo: прервать курс
			fmt.Printf("Прервать курс  id:%s\n", uid)
			http.Redirect(w, r, "/studentview/", http.StatusTemporaryRedirect /*307*/)
			return
		}
		if a == "prev" {
			if common.CurrentPara > 0 {
				common.CurrentPara--
			}
		}
		if a == "next" {
			if common.CurrentPara > len(common.CurrentCourse.Para)-1 {
				fmt.Println("ParaView: course is finished")
				//итоговая форма
				pview = common.ParaView{ParaCurNum: -2, PrevBut: true, NextBut: false}
				pview.Header = "Итоги"
				pview.Text = fmt.Sprintf("dfd;dsdsdflgd;lf\ndfsdd\nЫА--------\nЫВАЫВАЫВАЫВАЫВА ыв ыВАЫВАЫ\nfsd")

			} else {
				common.CurrentPara++
			}
		}
		//флаги доступности кнопок
		if common.CurrentPara >= 0 {
			fmt.Printf("ParaView CurrPara=%d -- len(Para): %d\n", common.CurrentPara, len(common.CurrentCourse.Para))
			pview.Header = common.CurrentCourse.Para[common.CurrentPara].Header
			pview.Text = common.CurrentCourse.Para[common.CurrentPara].Text
			pview.Answer = common.CurrentCourse.Para[common.CurrentPara].Answer
		}
		pview.PrevBut = bool(common.CurrentPara > 0)
		pview.NextBut = bool(common.CurrentPara < (len(common.CurrentCourse.Para) - 1))
		// маршаллим из структуры ParaView
		bytes, err := json.Marshal(pview)
		if err != nil {
			fmt.Println("ParaView: error on marshalling")
			return
		}
		// AJAX-запрос ?  отдать только JSON
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	}

}

/*
SetHandles устанавливает все обработчики запросов к Web-серверу, вызывается из main()
*/
func SetHandles() {

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/favicon.ico", handleFavicon)
	// http.HandleFunc("/img/", HandleImages)
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img")))) // если имя файла не указать, отдаёт список файлов из img !!!

	http.HandleFunc("/auth/", handleAuth)
	http.HandleFunc("/studentview/", handleStudentView)
	http.HandleFunc("/usersview/", handleUsersView)
	http.HandleFunc("/adduser/", handleAddUser)
	http.HandleFunc("/paraview/", handleParaView)

	//todo: !!! убрать потом /src/
	http.HandleFunc("/src/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q! \n %s", html.EscapeString(r.URL.Path), r.Method)
	})

}
