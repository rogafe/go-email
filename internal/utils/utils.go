package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"

	"golang.org/x/oauth2"
)

func IsAfterCurrentTime(t time.Time) bool {
	now := time.Now()
	current := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	return t.After(current)
}

func IsTimeEqualToday(t time.Time) bool {
	// truncate the time to midnight
	t = t.Truncate(24 * time.Hour)

	// get today's date
	today := time.Now().Truncate(24 * time.Hour)

	// compare the two dates
	return t.Equal(today)
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}

// create folder
func CreateFolder(folderName string) {

	err := os.MkdirAll(folderName, 0755)
	if err != nil {
		log.Println(err)
	}
}

func ChanToSlice(ch interface{}) interface{} {
	chv := reflect.ValueOf(ch)
	slv := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(ch).Elem()), 0, 0)
	for {
		v, ok := chv.Recv()
		if !ok {
			return slv.Interface()
		}
		slv = reflect.Append(slv, v)
	}
}

func ContainString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func WriteToken(token *oauth2.Token) {
	tokenJson, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		log.Println(err)
	}

	err = ioutil.WriteFile("./token.json", tokenJson, 0777)
	if err != nil {
		log.Println(err)
	}
}
