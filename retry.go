package main

import "time"

func Retry(operation func() error, attempts int, delay time.Duration) error {
	var err error
	currentDelay := delay
	for i := 0; i < attempts; i++ {
		err = operation()
		if err == nil {
			return nil
		}

		time.Sleep(currentDelay)
	}
	return err
}
