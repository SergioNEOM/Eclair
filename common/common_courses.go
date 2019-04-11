package common

import (
	"fmt"
	"os"
	"sort"
)

const (
	ANSWERTypeNone = iota
	ANSWERTypeInt
	ANSWERTypeString
	ANSWERTypeRadio
	ANSWERTypeCheckSingle
	ANSWERTypeCheckMulti
)

/*
 * Курс - упорядоченный(!) набор параграфов ( []Paragraph ) или лекций
 * Параграф - отображаемый элемент, стостоящий из Текста и Ответа на контрольный вопрос (опционально),
 *      если он предусмотрен в Тексте
 *
 * При отображении обычного параграфа ниже показываются кнопки со стрелками для листания вперед/назад (у крайних элементов - по одной)
 *
 * В случае с ответом на контрольный вопрос - стрелка "вперед" не отображается до получения правильного ответа
 * Для ввода ответа запланировать поле ввода и кнопку
 *
*/
/* тело параграфа - собственно блок Текста*/

type ParaBody struct {
	/*Header string // - в след. версии */
	ParaText   string
	// images -  в след.версии
}

type AnswerString struct {
	AnsText		string
	AnsGrade	int				// весовой коэф. (если полож - один из прав.ответ, отриц - неправ.)
}

/*
Paragraph (экспортируемая структура) - визуальная единица курса
*/
type Paragraph struct {
	ParaID	int
	ParaCurNum int	 /* № параграфа в курсе */
	ParaBody
	ParaAns	[]AnswerString
}

type ParagraphList struct {
	CourseID int
	ParaList []int
}

 /* Для передачи данных в шаблоны потребуется структура экранной формы: */
type ParaView struct {
	Paragraph      /* унаследовать поля	*/
	PrevBut    bool
	NextBut    bool
}

/*
Каждый курс физически располагается в отдельном файле (Предварительно - в одном каталоге, задаваемом в параметрах программы).
Имя файла - uid курса, присваеваемый при создании.
При запуске программа сканирует каталог курсов и заполняет
*/

type Block struct {
	Header string
	Text   string
}

/*
type Paragraph struct {
	Block
	AnsType int //ANSWERType const
	Answer  string
}
*/
/*
Course (экспортируемая структура) - представление курса в памяти, подгружаемое/выгружаемое из/в файл(а)
*/
type Course struct {
	Title   string
	Author  string
	Comment string
	Para    []*Paragraph
}

/*
SetHeader - заполняет заголовок курса (операция in-memory, требует последующего сохранения в файл)
*/
func (c *Course) SetHeader(title, author, comment string) {
	c.Title, c.Author, c.Comment = title, author, comment
	Flag_CourseChanged = true
	//todo: log
}

/*
AddPara - добавляет один параграф к курсу (операция in-memory, требует последующего сохранения в файл)
*/
/*func (c *Course) AddPara(header, text, answer string, anstype int) {
	c.Para = append(c.Para, &Paragraph{Block{header, text}, anstype, answer})
	Flag_CourseChanged = true
	//todo: log
}
*/

/*
LoadFromFile - загружает курс из файла в память (в подготовленную переменную)
*/
func (c *Course) LoadFromFile(filename string) bool {
	err := LoadJSON(filename, c)
	if err == nil {
		return true
	}
	//todo:  log :
	fmt.Fprintf(os.Stderr, "*** Error on load courses from file: %s\n", err)
	return false
}

/*
SaveToFile - выгружает курс в файл из памяти (permiss - права на создаваемый файл, потом можно вынести в конфиг)
*/
func (c *Course) SaveToFile(filename string, permiss os.FileMode) bool {
	err := SaveJSON(c, filename, permiss)
	if err == nil {
		return true // success
	}
	fmt.Fprintf(os.Stderr, "*** Error on save users to file: %s\n", err)
	//todo: выводить в журнал
	return false
}

/*
CourseLocation (экспортируемая структура) - формат справочника курсов
*/
type CourseLocation struct {
	Id    string
	Title string
	FName string
}

/*
CoursesList (экспортируемая структура) - справочник курсов
*/
type CoursesList []*CourseLocation

// Len, Swap и Less - для реализации интерфейса сортировки
func (c CoursesList) Len() int {
	return c.Len()
}
func (c CoursesList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c CoursesList) Less(i, j int) bool {
	return c[i].Title < c[j].Title
}

/*
ListCourses - формирует список курсов отсортированный по Title
*/
func (c CoursesList) ListCourses() []string {
	s := []string{}
	sort.Sort(&c)
	for i := range c {
		s = append(s, c[i].Title)
	}
	return s
}

/*
GetCourse - получить ссылку на курс по его uid
*/
func (c CoursesList) GetCourse(uid string) *CourseLocation {
	for i := range c {
		if c[i].Id == uid {
			return c[i]
		}
	}
	return nil
}

/*
AddCourse - добавляет в список курсов ссылку на файл
*/
func (c *CoursesList) AddCourse(title, filename string) string {
	uid := GenUID("CL")
	//todo: надо ли проверять уникальность id в списке???
	*c = append(*c, &CourseLocation{uid, title, filename})
	Flag_CoursesListChanged = true
	return uid
}

/*
*
*
*
*
 */

/*
 * TrainingPlan - учебный план (курс, назначенный конкретному студенту, тесты и результаты прохождения)
 *
 * каждый план привязан к студенту по UID студента
 *
 * Состоит из списка курсов (0 или более) и/или тестов (0 или более)
 *
 * Список курсов - список следующего вида:
 *
 *	- UID курса (ссылка на курс)
 *  - текущая позиция(№ тек.параграфа) в нём.
 *
 *            ВАЖНО: не путать со списком курсов приложения
 *
 * Список тестов - список следующего вида:
 *
 *  - UID теста (ссылка на файл с набором вопросов и другими параметрами, например,
 * 					допустимое число попыток ответа, последовательный или случайный показ вопросов, к-во вопросов в выборке и т.п.)
 *  - Протоколы попыток (упорядоченный список - см. ниже AttemptLog)
 */

type AttemptLog struct {
	//todo: зафиксровать список вопросов, отобранных для попытки ? Например, в Rates сразу занести записи с нулевыми оценками?
	DateTimeBegin int64
	DateTimeEnd   int64
	Rates         map[string]int //pairs [taskid]:rate
	Result        int
}

type CourseRec struct {
	Title       string        //todo: нужно ли поле, если есть ссылка на курс (в шаблон брать оттуда)
	CourseLink  string        // ссылка на курс (UID)
	MaxAttempts int           //todo: будет ли индивидуально у каждого, или задаётся в каждом курсе для всех одинаковая?
	Summary     []*AttemptLog //протокол попыток. Rates и Result для курса не имеют значения (ну, или в Result может лежать число открытых параграфов)
	Closed      bool          //todo: необходимость поля пока не ясна
	Moderator   string        //todo: необходимость поля пока не ясна
}

type TaskRec struct {
	Title    string        //todo: пока для совместимости, потом - для studentview.template заполнять из названия курса перед работой с шаблоном
	TestLink string        //ссылка на тест (UID)
	Summary  []*AttemptLog //протокол попыток

}

type TrainingPlan struct {
	Courses []CourseRec
	Tests   []TaskRec
}

func NewTrainingPlan() *TrainingPlan { //? ? ?
	tp := new(TrainingPlan)
	tp.Courses = make([]CourseRec, 0)
	tp.Tests = make([]TaskRec, 0)
	fmt.Println("newtrainigplan=", tp)
	return tp
}

func (tp *TrainingPlan) AddCourse(c CourseRec) {
	tp.Courses = append(tp.Courses, c)
}
