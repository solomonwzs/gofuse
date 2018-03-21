package gofuse

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type FuseRequestContext struct {
	kv     map[interface{}]interface{}
	kvLock *sync.RWMutex

	deadline atomic.Value
	raw      []byte

	done       chan struct{}
	doneReason error
	doneLock   *sync.Mutex
}

func newFuseRequestContext() (ctx *FuseRequestContext) {
	ctx = &FuseRequestContext{
		kv:     make(map[interface{}]interface{}),
		kvLock: &sync.RWMutex{},

		raw: nil,

		done:       make(chan struct{}),
		doneReason: nil,
		doneLock:   &sync.Mutex{},
	}
	return
}

func (ctx *FuseRequestContext) IsDone() bool {
	select {
	case <-ctx.done:
		return true
	default:
		return false
	}
}

func (ctx *FuseRequestContext) Deadline() (deadline time.Time, ok bool) {
	deadline, ok = ctx.deadline.Load().(time.Time)
	return
}

func (ctx *FuseRequestContext) setDeadline(t time.Time) error {
	ctx.deadline.Store(t)
	return nil
}

func (ctx *FuseRequestContext) Done() <-chan struct{} {
	return ctx.done
}

func (ctx *FuseRequestContext) Err() error {
	select {
	case <-ctx.done:
		return ctx.doneReason
	default:
		return nil
	}
}

func (ctx *FuseRequestContext) Value(key interface{}) interface{} {
	ctx.kvLock.RLock()
	defer ctx.kvLock.RUnlock()

	return ctx.kv[key]
}

func (ctx *FuseRequestContext) SetValue(key, value interface{}) (
	oldValue interface{}) {
	ctx.kvLock.Lock()
	defer ctx.kvLock.Unlock()

	oldValue, _ = ctx.kv[key]
	ctx.kv[key] = value
	return
}

func (ctx *FuseRequestContext) setDone(reason error) error {
	ctx.doneLock.Lock()
	defer ctx.doneLock.Unlock()

	if ctx.IsDone() {
		return errors.New("gofuse: context was closed")
	}
	ctx.doneReason = reason
	close(ctx.done)

	return nil
}
