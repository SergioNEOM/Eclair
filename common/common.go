package common

//import (
//
//
//)
/*
 *
 */

const (
	ROLE_Student = iota
	ROLE_Teacher
	ROLE_Coucher
	ROLE_Admin
)

var RolesNames = map[int]string{ROLE_Student: "Студент", ROLE_Teacher: "Преподаватель", ROLE_Coucher: "Модератор", ROLE_Admin: "Администратор"}

const APP_NAME = "Eclair 1.0"
const APP_PREFIX = "E1"
const USERS_PREFIX = "EU"
const COURSES_PREFIX = "EC"
const USERS_FILE = "users.json"
const COURSES_DIR = "./courses" //todo: потом брать из главного конфига

var DebugLevel int
var CurrentUser string //UserId
var AppRootDir string
var SessionId string
var CurrentPara int
var Flag_UsersChanged bool = false
var Flag_CourseChanged bool = false
var Flag_CoursesListChanged bool = false

// init() -?
var CurrentCourse = Course{}

var Users = make(UsersList, 1)

//var Planes = make(map[string]*TrainingPlan, 1)

var ListOfCourses CoursesList

/*
 * Работа с Cookies
 */
