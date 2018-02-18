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

//************
func FileExists(name string) bool {
	fi, err := os.Lstat(name)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

//************
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

//************
func SaveJSON(stru interface{}, filename string, permiss os.FileMode) error {
	bytes, err := json.MarshalIndent(stru, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, permiss)
}

//************
func Hash(str string) string {
	return fmt.Sprintf("%0X", crc32.ChecksumIEEE([]byte(str)))
}

//************
func GenUID(prefix string) string {
	return prefix + Hash(fmt.Sprintf("%d", rand.Int63n(time.Now().Unix()))) + "-" +
		fmt.Sprintf("%d", time.Now().Unix())
}

//************
func StartCourse(uid string) *ParaView {
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
