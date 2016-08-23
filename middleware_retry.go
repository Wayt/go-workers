package workers

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	DEFAULT_MAX_RETRY = 25
	LAYOUT            = "2006-01-02 15:04:05 MST"
)

type MiddlewareRetry struct{}

func (r *MiddlewareRetry) Call(queue string, message *Msg, next func() error) (err error) {
	defer func() {

		if err == nil {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}

		if err != nil {

			message.Set("queue", queue)
			message.Set("error_message", err.Error())

			if !IsFatal(err) && retry(message) {

				conn := Config.Pool.Get()
				defer conn.Close()

				retryCount := incrementRetry(message)

				waitDuration := durationToSecondsWithNanoPrecision(
					time.Duration(
						secondsToDelay(message, retryCount),
					) * time.Second,
				)

				message.Set("retry_at", time.Now().UTC().Add(time.Duration(waitDuration)).Format(LAYOUT))

				_, err2 := conn.Do(
					"zadd",
					Config.Namespace+RETRY_KEY,
					nowToSecondsWithNanoPrecision()+waitDuration,
					message.ToJson(),
				)

				// If we can't add the job to the retry queue,
				// then we shouldn't acknowledge the job, otherwise
				// it'll disappear into the void.
				if err2 != nil {
					err = err2
				}
			} else {
				message.Set("failed_at", time.Now().UTC().Format(LAYOUT))
				err = Fatal(err)
			}
		}
	}()

	err = next()

	return
}

func retry(message *Msg) bool {
	retry := false
	max := DEFAULT_MAX_RETRY

	if param, err := message.Get("retry").Bool(); err == nil {
		retry = param
	} else if param, err := message.Get("retry").Int(); err == nil {
		max = param
		retry = true
	}

	count, _ := message.Get("retry_count").Int()

	return retry && count < max
}

func incrementRetry(message *Msg) (retryCount int) {
	retryCount = 0

	if count, err := message.Get("retry_count").Int(); err != nil {
		message.Set("failed_at", time.Now().UTC().Format(LAYOUT))
	} else {
		message.Set("retried_at", time.Now().UTC().Format(LAYOUT))
		retryCount = count + 1
	}

	message.Set("retry_count", retryCount)

	return
}

func secondsToDelay(message *Msg, count int) int {

	if param, err := message.Get("retry_interval").Int(); err == nil && param > 0 {
		return param
	} else {
		power := math.Pow(float64(count), 4)
		return int(power) + 15 + (rand.Intn(30) * (count + 1))
	}
}
