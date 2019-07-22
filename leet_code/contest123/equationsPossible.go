package contest123

import (
	"log"
	"time"
)

func equationsPossible(equations []string) bool {

	var neq []string
	mp := make(map[uint8]*[]int)
	for _, s := range equations {
		if s[1] == '!' {
			neq = append(neq, s)
			continue
		}
		a := s[0]
		b := s[3]
		if valueA, okA := mp[a]; okA {
			if valueB, okB := mp[b]; okB {
				for k, v := range mp {
					if v == valueB {
						mp[k] = valueA
					}
				}
			} else {
				mp[b] = valueA
			}
		} else {
			if valueB, okB := mp[b]; okB {
				mp[a] = valueB
			} else {
				mp[b] = &[]int{}
				//log.Println(s, mp)
				mp[a] = mp[b]
			}
		}
		//log.Println(mp)
	}

	for _, s := range neq {
		a := s[0]
		b := s[3]
		if a == b {
			return false
		}
		valueA, okA := mp[a]
		valueB, okB := mp[b]
		if okA && okB {
			if valueA == valueB {
				return false
			}
		}
	}

	return true
}

func Question2() {
	eq := []string{"e==e", "d!=e", "c==d", "d!=e"}
	log.Println(equationsPossible(eq))

	t := time.NewTicker(time.Second)
	stop := false
	go func() {
		//time.Sleep(time.Second)
		t.Stop()
		stop = true
	}()
	for {
		if stop {
			break
		}
		<-t.C
	}
	log.Println("ok")
}
