package listener

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var errContextIsDone = errors.New("context is done")
var errRetriesNumberIsExceeded = errors.New("retry number is exceeded")

type page func(context.Context, *uint64) error
type watch func(context.Context) error

func paginateWithRetries(ctx context.Context, log logrus.FieldLogger, retryPeriod time.Duration, maxRetries int64, startAt *uint64, group *sync.WaitGroup, f page) {
	defer group.Done()
	log = log.WithFields(logrus.Fields{
		"obj":         "page_with_retries",
		"max_retries": maxRetries,
	})
	if isDone(ctx) {
		log.Error(errContextIsDone)
		return
	}

	var (
		err      error
		retryNum = int64(1)
	)
	for retryNum <= maxRetries {
		if isDone(ctx) {
			log.Error(errContextIsDone)
			return
		}

		if err = f(ctx, startAt); err != nil {
			if errors.Is(err, errContextIsDone) {
				log.Error(errContextIsDone)
				return
			}

			log = log.WithFields(logrus.Fields{
				"start_at":  startAt,
				"retry_num": retryNum,
			})
			log.Error(err)
			retryNum++
		}

		time.Sleep(retryPeriod)
	}
	log.Panic(fmt.Errorf("%s last error: %s", errRetriesNumberIsExceeded.Error(), err.Error()))
}

func watchWithRetries(ctx context.Context, log logrus.FieldLogger, retryPeriod time.Duration, maxRetries int64, group *sync.WaitGroup, f watch) {
	defer group.Done()
	log = log.WithFields(logrus.Fields{
		"obj":         "watch_with_retries",
		"max_retries": maxRetries,
	})
	if isDone(ctx) {
		log.Error(errContextIsDone)
		return
	}

	var (
		err      error
		retryNum = int64(1)
	)
	for retryNum <= maxRetries {
		if isDone(ctx) {
			log.Error(errContextIsDone)
			return
		}

		if err = f(ctx); err != nil {
			if errors.Is(err, errContextIsDone) {
				log.Error(errContextIsDone)
				return
			}

			log = log.WithFields(logrus.Fields{
				"retry_num": retryNum,
			})
			log.Error(err)
			retryNum++
		}

		time.Sleep(retryPeriod)
	}
	log.Panic(fmt.Errorf("%s last error: %s", errRetriesNumberIsExceeded.Error(), err.Error()))
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
