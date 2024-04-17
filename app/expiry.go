package main

import (
	"fmt"
	"time"
)

func expire(store Store, key string, time_in_ms int64) {
	expiration_timer := time.NewTimer(time.Duration(time_in_ms) * time.Millisecond)
	<-expiration_timer.C
	fmt.Printf("Expiring %s after %d ms.\n", key, time_in_ms)
	delete(store.data, key)
}
