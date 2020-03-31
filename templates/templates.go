package templates

import (
	//	"fmt"
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/SergioNEOM/Eclair/common"
)

//todo: вынести блок структур и интерфейс в common_users

type FullUsers struct {
	Id      string
	Name    string
	Account string
	Role    string
}

type UsersContent struct {
	SortBy int
	Sorted []FullUsers
}

func (uc *UsersContent) Len() int           { return len(uc.Sorted) }
func (uc *UsersContent) Less(i, j int) bool { return uc.Sorted[i].Name < uc.Sorted[j].Name }
func (uc *UsersContent) Swap(i, j int)      { uc.Sorted[i], uc.Sorted[j] = uc.Sorted[j], uc.Sorted[i] }
func (uc *UsersContent) Fill() {
	for key := range common.Users {
		uc.Sorted = append(uc.Sorted, FullUsers{Id: key, Name: common.Users[key].Name, Account: common.Users[key].Account, Role: common.RolesNames[common.Users[key].Role]})
	}
}
func (uc *UsersContent) MakeSort() {
	sort.Sort(uc)
}

//todo:  --------------------------------

func GoAuthorize(w *http.ResponseWriter) error {
	t, err := template.ParseFiles("./templates/authform.html" /* , header, footer*/)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(*w, "authform", nil)
	if err != nil {
		return err
	}
	return nil
}

// GoStudentView  отобразить экран для студента
func GoStudentView(w *http.ResponseWriter) error {
	t, err := template.ParseFiles("./templates/studentview.html" /* , header, footer*/)
	if err != nil {
		return err
	}
	fmt.Printf("GoStudentView: %s -- %v\n", common.CurrentUser, common.Users[common.CurrentUser].Plane)
	err = t.ExecuteTemplate(*w, "studentview", common.Users[common.CurrentUser].Plane)
	if err != nil {
		return err
	}
	return nil

}

func GoUsersView(w *http.ResponseWriter) error {
	t, err := template.ParseFiles("./templates/usersview.html" /* , header, footer*/)
	if err != nil {
		return err
	}

	uc := UsersContent{}
	uc.Fill()
	uc.MakeSort()

	err = t.ExecuteTemplate(*w, "usersview", uc.Sorted /*common.Users*/)
	if err != nil {
		return err
	}
	return nil
}

func GoAddUser(w *http.ResponseWriter) error {
	t, err := template.ParseFiles("./templates/adduser.html" /* , header, footer*/)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(*w, "adduser", nil)
	if err != nil {
		return err
	}
	return nil
}

func GoParaView(w *http.ResponseWriter, para *common.ParaView) error {
	t, err := template.ParseFiles("./templates/paraview.html" /* , header, footer*/)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(*w, "paraview", para)
	if err != nil {
		return err
	}
	return nil
}

// GoErrorView  отобразить экран ошибки
func GoErrorView(errMes string, w *http.ResponseWriter) {
	t, _ := template.ParseFiles("./templates/error.html")

	_ = t.ExecuteTemplate(*w, "error", errMes)

	return
}
