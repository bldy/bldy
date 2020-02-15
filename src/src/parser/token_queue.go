package parser

import (
	"sync"

	"bldy.build/bldy/src/token"
)

type tq []*token.Token

func (tq *tq) push(t *token.Token) { *tq = append(*tq, t) }

func (tq *tq) pop() *token.Token {
	old := *tq
	n := len(old)
	if n < 1 {
		return nil
	}
	item := old[0]
	*tq = old[1:]
	return item
}

type Queue struct {
	t     chan *token.Token
	ready chan interface{}
	r     *sync.Cond
	w     *sync.Cond

	q tq
}

func (q *Queue) Peek() *token.Token {
	t := *q.q[0]
	return &t
}
func (q *Queue) Next() *token.Token {
	t := q.q.pop()
	q.q.push(<-q.t)
	return t
}

func NewQueue(t chan *token.Token) *Queue {
	q := &Queue{
		ready: make(chan interface{}),
		r:     sync.NewCond(&sync.Mutex{}),
		w:     sync.NewCond(&sync.Mutex{}),
		t:     t,
	}
	q.q.push(<-t)
	q.q.push(<-t)
	return q
}
