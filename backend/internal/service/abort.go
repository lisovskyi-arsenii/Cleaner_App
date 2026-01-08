package service

import (
	"context"
	"sync"
)

// AbortManager
type AbortManager struct {
	mutex    	sync.RWMutex
	cancelFunc 	context.CancelFunc
	isRunning	bool
}

var globalAbortManager = &AbortManager{}

func (am *AbortManager) SetOperation(cancel context.CancelFunc) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if am.cancelFunc != nil && am.isRunning {
		am.cancelFunc()
	}

	am.cancelFunc = cancel
	am.isRunning = true
}

func (am *AbortManager) Abort() bool {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if am.cancelFunc != nil && am.isRunning {
		am.cancelFunc()
		am.isRunning = false
		return true
	}

	return false
}

func (am *AbortManager) Clear() {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.cancelFunc = nil
	am.isRunning = false
}

func (am *AbortManager) IsRunning() bool {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	return am.isRunning
}

func GetAbortManager() *AbortManager {
	return globalAbortManager
}
