package yeb_exp

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

func getUser() <-chan string {
	ch := make(chan string, 10)
	go func() {
		ch <- "13823260948"
		return
		var rows []bson.M
		lastMobile := "1379"
		for {
			err := appConfig.mgo.user.Find(bson.M{
				appConfig.MobileColumn: bson.M{"$gt": lastMobile},
			}).Limit(100).Sort(appConfig.MobileColumn).All(&rows)
			if err != nil {
				log.Println(err)
				continue
			}
			lastMobile = rows[len(rows)-1][appConfig.MobileColumn].(string)
			for _, r := range rows {
				ch <- r[appConfig.MobileColumn].(string)
			}
		}
		close(ch)
	}()
	return ch
}
