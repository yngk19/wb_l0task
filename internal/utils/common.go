package utils

import "time"

func DoWithTries(fn func() error, maxAttempts int, delay time.Duration) (err error) {
	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay * time.Second)
			maxAttempts--
			continue
		}
		return nil
	}
	return nil
}
