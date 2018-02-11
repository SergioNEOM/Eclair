package common

import (
	"fmt"
	"os"
	"errors"
	"io/ioutil"
	"encoding/json"
	"hash/crc32"
	"time"
	"math/rand"
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
		return errors.New("-- File not found: "+filename)
	}
        bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
 	return  json.Unmarshal(bytes, &stru)
}

//************
func SaveJSON(stru interface{}, filename string, permiss os.FileMode) error {
	bytes, err := json.MarshalIndent(stru,""," ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, permiss)
}

//************
func Hash(str string) string {
	return fmt.Sprintf("%0X",crc32.ChecksumIEEE([]byte(str)) )
}

//************
func GenUID(prefix string) string {
	return prefix + Hash(fmt.Sprintf("%d",rand.Int63n(time.Now().Unix())))+"-"+
			fmt.Sprintf("%d",time.Now().Unix()) 
}

//************
