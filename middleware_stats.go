package workers

import (
	"time"
)

type MiddlewareStats struct{}

func (l *MiddlewareStats) Call(queue string, message *Msg, next func() error) (err error) {

	if err = next(); err != nil {
		incrementStats("failed")
	} else {
		incrementStats("processed")
	}

	return
}

func incrementStats(metric string) {
	conn := Config.Pool.Get()
	defer conn.Close()

	today := time.Now().UTC().Format("2006-01-02")

	conn.Send("multi")
	conn.Send("incr", Config.Namespace+"stat:"+metric)
	conn.Send("incr", Config.Namespace+"stat:"+metric+":"+today)

	if _, err := conn.Do("exec"); err != nil {
		Logger.Println("couldn't save stats:", err)
	}
}
