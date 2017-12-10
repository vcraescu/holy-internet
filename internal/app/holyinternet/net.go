package holyinternet

import (
	"github.com/vcraescu/holy-internet/internal/pkg/util"
	"log"
	"math/rand"
	"time"
	"github.com/spf13/viper"
)

func pingHosts(hosts []string) (bool) {
	failures := 0
	for _, host := range hosts {
		ok, err := util.Ping(host)
		if ok {
			continue
		}

		if err != nil {
			log.Println(err)
		}

		failures++
	}

	return failures != len(hosts)
}

func pickHosts(hosts []string) []string {
	hc := viper.GetInt("pray.count")
	if len(hosts) < hc {
		hc = len(hosts)
	}

	indexes := make([]int, 0)
	rand.Seed(time.Now().UnixNano())

	for len(indexes) < hc {
		index := rand.Intn(len(hosts))

		already := false
		for _, i := range indexes {
			if i == index {
				already = true
				break
			}
		}

		if already {
			continue
		}

		indexes = append(indexes, index)
	}

	pickedHosts := make([]string, 0)
	for _, index := range indexes {
		pickedHosts = append(pickedHosts, hosts[index])
	}

	return pickedHosts
}

func IsInternetOK(hosts []string) (bool) {
	checkHosts := pickHosts(hosts)
	return pingHosts(checkHosts)
}
