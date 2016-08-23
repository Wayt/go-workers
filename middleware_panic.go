package workers

import (
	"fmt"
)

type MiddlewarePanic struct{}

func (p *MiddlewarePanic) Call(queue string, message *Msg, next func() error) (err error) {

	defer func() {
		if e := recover(); e != nil {
			// Should it be a Fatal error ?
			err = fmt.Errorf("%v", e)
		}
	}()

	err = next()
	return
}
