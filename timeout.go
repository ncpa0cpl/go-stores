package stores

import "time"

func timeout(fn func(), ms int) {
	go (func() {
		time.Sleep(time.Duration(ms) * (time.Millisecond))
		fn()
	})()
}
