package stringqueue

import "sync"

type StringQueue struct {
	queue []string
	mutex *sync.Mutex
}
