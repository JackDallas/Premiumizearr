package stringqueue

import "sync"

func NewStringQueue() *StringQueue {
	return &StringQueue{queue: make([]string, 0), mutex: &sync.Mutex{}}
}

func (UploadQueue *StringQueue) Len() int {
	UploadQueue.mutex.Lock()
	defer UploadQueue.mutex.Unlock()
	return len(UploadQueue.queue)
}

func (UploadQueue *StringQueue) Add(path string) {
	UploadQueue.mutex.Lock()
	defer UploadQueue.mutex.Unlock()
	UploadQueue.queue = append(UploadQueue.queue, path)
}

func (UploadQueue *StringQueue) PopTopOfQueue() (bool, string) {
	UploadQueue.mutex.Lock()
	defer UploadQueue.mutex.Unlock()
	if len(UploadQueue.queue) > 0 {
		rtn := UploadQueue.queue[0]
		UploadQueue.queue = UploadQueue.queue[1:]
		return true, rtn
	}
	return false, ""
}

func (UploadQueue *StringQueue) GetQueue() []string {
	UploadQueue.mutex.Lock()
	defer UploadQueue.mutex.Unlock()
	return UploadQueue.queue
}
