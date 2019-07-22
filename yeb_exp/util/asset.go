package util

import (
	"fmt"
	_ "funny/yeb_exp/util/statik"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//go:generate statik -src=assets -dest=.

var hfs http.FileSystem

func init() {
	hfs, _ = fs.New()
}
func GetAsset(path string) (bin []byte, err error) {
	f, err := hfs.Open(path)
	if err != nil {
		log.Println(err, hfs)
		return
	}
	return ioutil.ReadAll(f)
}

func GetAssetStr(path string) string {
	bin, _ := GetAsset(path)
	return string(bin)
}

func StrReplace(str string, mp map[string]interface{}) string {
	for k, v := range mp {
		str = strings.Replace(str, k, fmt.Sprint(v), -1)
	}
	return str
}

func LogErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
