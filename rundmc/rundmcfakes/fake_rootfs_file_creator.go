// This file was generated by counterfeiter
package rundmcfakes

import (
	"sync"

	"code.cloudfoundry.org/guardian/rundmc"
)

type FakeRootfsFileCreator struct {
	CreateFilesStub        func(rootFSPath string, pathsToCreate ...string) error
	createFilesMutex       sync.RWMutex
	createFilesArgsForCall []struct {
		rootFSPath    string
		pathsToCreate []string
	}
	createFilesReturns struct {
		result1 error
	}
	createFilesReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRootfsFileCreator) CreateFiles(rootFSPath string, pathsToCreate ...string) error {
	fake.createFilesMutex.Lock()
	ret, specificReturn := fake.createFilesReturnsOnCall[len(fake.createFilesArgsForCall)]
	fake.createFilesArgsForCall = append(fake.createFilesArgsForCall, struct {
		rootFSPath    string
		pathsToCreate []string
	}{rootFSPath, pathsToCreate})
	fake.recordInvocation("CreateFiles", []interface{}{rootFSPath, pathsToCreate})
	fake.createFilesMutex.Unlock()
	if fake.CreateFilesStub != nil {
		return fake.CreateFilesStub(rootFSPath, pathsToCreate...)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.createFilesReturns.result1
}

func (fake *FakeRootfsFileCreator) CreateFilesCallCount() int {
	fake.createFilesMutex.RLock()
	defer fake.createFilesMutex.RUnlock()
	return len(fake.createFilesArgsForCall)
}

func (fake *FakeRootfsFileCreator) CreateFilesArgsForCall(i int) (string, []string) {
	fake.createFilesMutex.RLock()
	defer fake.createFilesMutex.RUnlock()
	return fake.createFilesArgsForCall[i].rootFSPath, fake.createFilesArgsForCall[i].pathsToCreate
}

func (fake *FakeRootfsFileCreator) CreateFilesReturns(result1 error) {
	fake.CreateFilesStub = nil
	fake.createFilesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRootfsFileCreator) CreateFilesReturnsOnCall(i int, result1 error) {
	fake.CreateFilesStub = nil
	if fake.createFilesReturnsOnCall == nil {
		fake.createFilesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createFilesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRootfsFileCreator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createFilesMutex.RLock()
	defer fake.createFilesMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeRootfsFileCreator) recordInvocation(key string, args []interface{}) {
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

var _ rundmc.RootfsFileCreator = new(FakeRootfsFileCreator)