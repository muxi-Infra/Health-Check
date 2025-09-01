package main

import "time"

func Retry[T any](operation func() (T, error), attempts int, delay time.Duration) (T, error) {
	var (
		err  error
		data T
	)
	currentDelay := delay
	for i := 0; i < attempts; i++ {
		data, err = operation()
		if err == nil {
			return data, nil
		}

		time.Sleep(currentDelay)
	}
	var zero T
	return zero, err
}

func KeepBeat[T any](operation func() (T, error), delay time.Duration, maxDelay time.Duration) {
	var (
		err error
	)
	currentDelay := delay
	var (
		lastOK  bool
		hasLast bool
	)
	for {
		_, err = operation()
		ok := err == nil

		// 状态改变则重置，否则\*2；限制不超过 maxDelay
		if hasLast && ok != lastOK {
			currentDelay = delay
		} else {
			currentDelay *= 2
			if currentDelay > maxDelay {
				currentDelay = maxDelay
			}
		}

		lastOK = ok
		hasLast = true

		time.Sleep(currentDelay)
	}
}
