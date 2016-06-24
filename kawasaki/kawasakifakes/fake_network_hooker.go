// This file was generated by counterfeiter
package kawasakifakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/guardian/gardener"
	"github.com/cloudfoundry-incubator/guardian/kawasaki"
	"github.com/pivotal-golang/lager"
)

type FakeNetworkHooker struct {
	HooksStub        func(log lager.Logger, containerSpec garden.ContainerSpec) (gardener.Hooks, error)
	hooksMutex       sync.RWMutex
	hooksArgsForCall []struct {
		log           lager.Logger
		containerSpec garden.ContainerSpec
	}
	hooksReturns struct {
		result1 gardener.Hooks
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeNetworkHooker) Hooks(log lager.Logger, containerSpec garden.ContainerSpec) (gardener.Hooks, error) {
	fake.hooksMutex.Lock()
	fake.hooksArgsForCall = append(fake.hooksArgsForCall, struct {
		log           lager.Logger
		containerSpec garden.ContainerSpec
	}{log, containerSpec})
	fake.recordInvocation("Hooks", []interface{}{log, containerSpec})
	fake.hooksMutex.Unlock()
	if fake.HooksStub != nil {
		return fake.HooksStub(log, containerSpec)
	} else {
		return fake.hooksReturns.result1, fake.hooksReturns.result2
	}
}

func (fake *FakeNetworkHooker) HooksCallCount() int {
	fake.hooksMutex.RLock()
	defer fake.hooksMutex.RUnlock()
	return len(fake.hooksArgsForCall)
}

func (fake *FakeNetworkHooker) HooksArgsForCall(i int) (lager.Logger, garden.ContainerSpec) {
	fake.hooksMutex.RLock()
	defer fake.hooksMutex.RUnlock()
	return fake.hooksArgsForCall[i].log, fake.hooksArgsForCall[i].containerSpec
}

func (fake *FakeNetworkHooker) HooksReturns(result1 gardener.Hooks, result2 error) {
	fake.HooksStub = nil
	fake.hooksReturns = struct {
		result1 gardener.Hooks
		result2 error
	}{result1, result2}
}

func (fake *FakeNetworkHooker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.hooksMutex.RLock()
	defer fake.hooksMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeNetworkHooker) recordInvocation(key string, args []interface{}) {
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

var _ kawasaki.NetworkHooker = new(FakeNetworkHooker)