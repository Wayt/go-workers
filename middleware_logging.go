package workers

import (
	"fmt"
	"time"
)

type MiddlewareLogging struct{}

func (l *MiddlewareLogging) Call(queue string, message *Msg, next func() error) (err error) {
	prefix := fmt.Sprint(queue, " JID-", message.Jid())

	start := time.Now()
	Logger.Println(prefix, "start")
	Logger.Println(prefix, "args: ", message.Args().ToJson())

	if err = next(); err != nil {

		Logger.Println(prefix, "fail:", time.Since(start))
		if !IsFatal(err) {
			if retryAt, err := message.Get("retry_at").String(); err == nil {
				Logger.Println(prefix, "retry at:", retryAt)
			}
		}
	} else {
		Logger.Println(prefix, "done:", time.Since(start))
	}

	return
}
