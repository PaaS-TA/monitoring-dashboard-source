package model

import (
	"log"
)

type SessionQueue struct {
	sessions map[string]interface{}
}

func NewQueue() *SessionQueue {
	return &SessionQueue{make(map[string]interface{})}
}

func (q *SessionQueue) Push(id string, i interface{}) {
	q.sessions[id] = i
}

func (q *SessionQueue) Pop(id string) {
	delete(q.sessions, id)
}

func (q *SessionQueue) isExist(id string) bool {
	_, ok := q.sessions[id]
	return ok
}

func (q *SessionQueue) PrintQueueList() {
	queueList := make([]string, len(q.sessions))
	idx := 0
	for key := range q.sessions {
		queueList[idx] = key
		idx++
	}
	log.Printf("[system] Queue list : %v\n", queueList)
}