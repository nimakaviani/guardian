// This file was generated by counterfeiter
package rundmcfakes

import (
	"sync"

	"code.cloudfoundry.org/guardian/rundmc"
)

type FakeStateStore struct {
	StoreStoppedStub        func(handle string)
	storeStoppedMutex       sync.RWMutex
	storeStoppedArgsForCall []struct {
		handle string
	}
	IsStoppedStub        func(handle string) bool
	isStoppedMutex       sync.RWMutex
	isStoppedArgsForCall []struct {
		handle string
	}
	isStoppedReturns struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStateStore) StoreStopped(handle string) {
	fake.storeStoppedMutex.Lock()
	fake.storeStoppedArgsForCall = append(fake.storeStoppedArgsForCall, struct {
		handle string
	}{handle})
	fake.recordInvocation("StoreStopped", []interface{}{handle})
	fake.storeStoppedMutex.Unlock()
	if fake.StoreStoppedStub != nil {
		fake.StoreStoppedStub(handle)
	}
}

func (fake *FakeStateStore) StoreStoppedCallCount() int {
	fake.storeStoppedMutex.RLock()
	defer fake.storeStoppedMutex.RUnlock()
	return len(fake.storeStoppedArgsForCall)
}

func (fake *FakeStateStore) StoreStoppedArgsForCall(i int) string {
	fake.storeStoppedMutex.RLock()
	defer fake.storeStoppedMutex.RUnlock()
	return fake.storeStoppedArgsForCall[i].handle
}

func (fake *FakeStateStore) IsStopped(handle string) bool {
	fake.isStoppedMutex.Lock()
	fake.isStoppedArgsForCall = append(fake.isStoppedArgsForCall, struct {
		handle string
	}{handle})
	fake.recordInvocation("IsStopped", []interface{}{handle})
	fake.isStoppedMutex.Unlock()
	if fake.IsStoppedStub != nil {
		return fake.IsStoppedStub(handle)
	} else {
		return fake.isStoppedReturns.result1
	}
}

func (fake *FakeStateStore) IsStoppedCallCount() int {
	fake.isStoppedMutex.RLock()
	defer fake.isStoppedMutex.RUnlock()
	return len(fake.isStoppedArgsForCall)
}

func (fake *FakeStateStore) IsStoppedArgsForCall(i int) string {
	fake.isStoppedMutex.RLock()
	defer fake.isStoppedMutex.RUnlock()
	return fake.isStoppedArgsForCall[i].handle
}

func (fake *FakeStateStore) IsStoppedReturns(result1 bool) {
	fake.IsStoppedStub = nil
	fake.isStoppedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeStateStore) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.storeStoppedMutex.RLock()
	defer fake.storeStoppedMutex.RUnlock()
	fake.isStoppedMutex.RLock()
	defer fake.isStoppedMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeStateStore) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ rundmc.StateStore = new(FakeStateStore)
