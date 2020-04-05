package handles

import (
	"fmt"
	"html" //todo: временно
	"log"
	"net/http"

	//	"net/url"
	//
	"github.com/SergioNEOM/Eclair/common"
	//	"../users"
	"github.com/SergioNEOM/Eclair/templates"
)

const COOKIE_NAME = "Eclair-+-Token"
const refreshJWTPath = "/refresh/"

//
//

// SetHandles устанавливает все обработчики запросов к Web-серверу, вызывается из main()
//
func SetHandles() {

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/favicon.ico", handleFavicon)
	// http.HandleFunc("/img/", HandleImages)
	//todo: потом сделать параметр в конфиге: в случае проксирующего Веб-сервера тот пусть статику сам отдает
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img")))) // если имя файла не указать, отдаёт список файлов из img !!!
	//
	http.HandleFunc("/auth/", handleAuth)
	// перевыпуск токенов
	http.HandleFunc(refreshJWTPath, handleMakeNewJWT)

	// студенты
	http.HandleFunc("/studentview/", handleStudentView)

	http.HandleFunc("/paraview/", handleParaView)

	// для администраторов:
	http.HandleFunc("/admin/", handleUsersView) // users list
	http.HandleFunc("/admin/adduser/", handleAddUser)
	//http.HandleFunc("/admin/edituser/", handleEditUser)
	//http.HandleFunc("/admin/deluser/", handleDelUser)

	//todo: !!! убрать потом /src/
	http.HandleFunc("/src/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q! \n %s", html.EscapeString(r.URL.Path), r.Method)
	})

}

//
func getCookieValue(r *http.Request) string {
	// проверяем, был ли выдан токен (и записан в куки)
	c, err := r.Cookie(COOKIE_NAME)
	//	log.Printf("%v %v\n", COOKIE_NAME, c)
	if err != nil {
		return ""
	}
	return c.Value
}

func handleMakeNewJWT(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>Eclair token refreshing<br><br><a href='/auth'>Вход</a></body><html>")
}

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
			fmt.Println("token checked?")
			switch common.Users[key].Role {
			case common.ROLE_Student:
				href = "/studentview/"
			case common.ROLE_Teacher:
				href = "/teacherview/"
			case common.ROLE_Admin:
				href = "/admin/" //	"/adminview/"
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
	// ветка else т.к. выше был редирект
	// например, r.Method == GET
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
 * admin - view users list
 */
var handleUsersView = func(w http.ResponseWriter, r *http.Request) {
	CheckJWToken(w, r)
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
		// получли данные формы, надо дальше смотреть...

		//todo: обработать ошибку ?
		CheckJWToken(w, r)

		common.Users.AddUser(r.PostFormValue("Ulogin"), r.PostFormValue("Upass"), r.PostFormValue("Uname"), common.ROLE_Student)
		http.Redirect(w, r, "/admin/", http.StatusTemporaryRedirect /*307*/)
		// редирект на корень своей роли "/admin"
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
	// var pview common.ParaView
	/*
		1. Показать параграф
		2. по значению action перейти к пред/след параграфу (если не контрольный вопрос или не крайний параграф)

		Если в запросе есть параметр action=start, то проверить тип запроса: для GET - стартовать новый курс (uid!=""),
			 для POST - записать ошибку в журнал, т.к. новый курс можно начать только из списка (studentview)

	*/
	uid = r.URL.Query().Get("uid")
	a = r.URL.Query().Get("action")
	//
	if r.Method == http.MethodGet {
		if a == "start" {
			if uid == "" {
				fmt.Println(" --> ParaView Error: empty uid in get request")
				fmt.Fprint(w, "<h1>ParaView Error: empty uid in get request</h1>")
				return
			}
			//todo: начало вывода курса. Как отловить случайный переход?
			//err = templates.GoParaView(&w, common.StartCourse(uid))
			if err != nil {
				fmt.Fprintf(w, "ParaView Error: %s", err)
				log.Fatalf("--- ParaView Error: %s   -----\n", err)
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
			common.PrevPara()
		}
		if a == "next" {
			common.NextPara()
		}

		// AJAX-запрос ?  отдать только JSON
		w.Header().Set("Content-Type", "application/json")
		//	fmt.Fprint(w, common.FillViewJSON())
		// fmt.Fprint(w, string(bytes))
	}

}
