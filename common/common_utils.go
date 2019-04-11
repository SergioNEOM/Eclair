package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

//FileExists - проверка на существование файла
func FileExists(name string) bool {
	fi, err := os.Lstat(name)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

//LoadJSON - загружает данные в структуру из файла
func LoadJSON(filename string, stru interface{}) error {
	if !FileExists(filename) {
		return errors.New("-- File not found: " + filename)
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &stru)
}

//SaveJSON - сохраняет данные структуры в файл
func SaveJSON(stru interface{}, filename string, permiss os.FileMode) error {
	bytes, err := json.MarshalIndent(stru, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, permiss)
}

//Hash - вычисляет CRC32 строки
func Hash(str string) string {
	return fmt.Sprintf("%0X", crc32.ChecksumIEEE([]byte(str)))
}

//GenUID - генератор псевдоуникальных строк для идентификации сущностей
func GenUID(prefix string) string {
	return prefix + Hash(fmt.Sprintf("%d", rand.Int63n(time.Now().Unix()))) + "-" +
		fmt.Sprintf("%d", time.Now().Unix())
}

//StartCourse - выполняет действия по запуску нового курса
/*func StartCourse(uid string) *ParaView {
	var pview *ParaView
	c := ListOfCourses.GetCourse(uid)
	if c != nil {
		if CurrentCourse.LoadFromFile(c.FName) {
			CurrentPara = -1 //todo: а если начать не сначала?
			fmt.Printf("Старт курса id:%s из файла %s\n", uid, c.FName)
			//пришли первый раз - тогда покажем страницу полностью
			// err = templates.GoParaView(&w, common.CurrentCourse.Para[common.CurrentPara])
			pview = &ParaView{ParaCurNum: -1, PrevBut: false, NextBut: true}
			pview.Header = "Вводная информация"
			pview.Text = fmt.Sprintf("dfd;flgd;lf\ndfsdd\nЫАЫВАЫВАЫВАЫВАЫВА ыв ыВАЫВАЫ\nfsd")
			return pview
		}
	}
	// не загрузился из файла
	fmt.Printf("ParaView: Course %s  is empty or not loaded\n", uid)
	return nil
}
*/
//**********

//PrevPara - переход к предыдущему параграфу (если возможно)
func PrevPara() {
	if CurrentPara > 0 {
		CurrentPara--
	}
}

//NextPara - переход к следующему параграфу (если возможно)
func NextPara() {
	if CurrentPara >= len(CurrentCourse.Para)-1 {
		CurrentPara = -2 //finish
	} else {
		CurrentPara++
	}
}

//FillViewJSON - возвращает текущий Параграф, оформленный в JSON
/*func FillViewJSON() string {
	var pview ParaView
	if CurrentPara == -2 {
		//итоговая форма
		fmt.Println("ParaView: course is finished")
		pview = ParaView{ParaCurNum: -2, PrevBut: true, NextBut: false}
		pview.Header = "Итоги"
		pview.Text = fmt.Sprintf("dfd;dsdsdflgd;lf\ndfsdd\nЫА--------\nЫВАЫВАЫВАЫВАЫВА ыв ыВАЫВАЫ\nfsd")
	}
	if CurrentPara >= 0 {
		fmt.Printf("ParaView CurrPara=%d -- len(Para): %d\n", CurrentPara, len(CurrentCourse.Para))
		pview.Header = CurrentCourse.Para[CurrentPara].Header
		pview.Text = CurrentCourse.Para[CurrentPara].Text
		pview.Answer = CurrentCourse.Para[CurrentPara].Answer
	}
	pview.PrevBut = bool(CurrentPara > 0)
	pview.NextBut = bool((CurrentPara != -2) && (CurrentPara <= (len(CurrentCourse.Para) - 1)))
	bytes, err := json.Marshal(pview)
	if err != nil {
		fmt.Println("ParaView: error on marshalling")
		return ""
	}
	return string(bytes)
}
*/
