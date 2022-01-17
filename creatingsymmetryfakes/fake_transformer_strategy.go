// Code generated by counterfeiter. DO NOT EDIT.
package creatingsymmetryfakes

import (
	"io"
	"sync"

	creating_symmetry "github.com/Chadius/creating-symmetry"
)

type FakeTransformerStrategy struct {
	ApplyFormulaToTransformImageStub        func(io.Reader, io.Reader, io.Reader, io.Writer) error
	applyFormulaToTransformImageMutex       sync.RWMutex
	applyFormulaToTransformImageArgsForCall []struct {
		arg1 io.Reader
		arg2 io.Reader
		arg3 io.Reader
		arg4 io.Writer
	}
	applyFormulaToTransformImageReturns struct {
		result1 error
	}
	applyFormulaToTransformImageReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTransformerStrategy) ApplyFormulaToTransformImage(arg1 io.Reader, arg2 io.Reader, arg3 io.Reader, arg4 io.Writer) error {
	fake.applyFormulaToTransformImageMutex.Lock()
	ret, specificReturn := fake.applyFormulaToTransformImageReturnsOnCall[len(fake.applyFormulaToTransformImageArgsForCall)]
	fake.applyFormulaToTransformImageArgsForCall = append(fake.applyFormulaToTransformImageArgsForCall, struct {
		arg1 io.Reader
		arg2 io.Reader
		arg3 io.Reader
		arg4 io.Writer
	}{arg1, arg2, arg3, arg4})
	stub := fake.ApplyFormulaToTransformImageStub
	fakeReturns := fake.applyFormulaToTransformImageReturns
	fake.recordInvocation("ApplyFormulaToTransformImage", []interface{}{arg1, arg2, arg3, arg4})
	fake.applyFormulaToTransformImageMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeTransformerStrategy) ApplyFormulaToTransformImageCallCount() int {
	fake.applyFormulaToTransformImageMutex.RLock()
	defer fake.applyFormulaToTransformImageMutex.RUnlock()
	return len(fake.applyFormulaToTransformImageArgsForCall)
}

func (fake *FakeTransformerStrategy) ApplyFormulaToTransformImageCalls(stub func(io.Reader, io.Reader, io.Reader, io.Writer) error) {
	fake.applyFormulaToTransformImageMutex.Lock()
	defer fake.applyFormulaToTransformImageMutex.Unlock()
	fake.ApplyFormulaToTransformImageStub = stub
}

func (fake *FakeTransformerStrategy) ApplyFormulaToTransformImageArgsForCall(i int) (io.Reader, io.Reader, io.Reader, io.Writer) {
	fake.applyFormulaToTransformImageMutex.RLock()
	defer fake.applyFormulaToTransformImageMutex.RUnlock()
	argsForCall := fake.applyFormulaToTransformImageArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeTransformerStrategy) ApplyFormulaToTransformImageReturns(result1 error) {
	fake.applyFormulaToTransformImageMutex.Lock()
	defer fake.applyFormulaToTransformImageMutex.Unlock()
	fake.ApplyFormulaToTransformImageStub = nil
	fake.applyFormulaToTransformImageReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTransformerStrategy) ApplyFormulaToTransformImageReturnsOnCall(i int, result1 error) {
	fake.applyFormulaToTransformImageMutex.Lock()
	defer fake.applyFormulaToTransformImageMutex.Unlock()
	fake.ApplyFormulaToTransformImageStub = nil
	if fake.applyFormulaToTransformImageReturnsOnCall == nil {
		fake.applyFormulaToTransformImageReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.applyFormulaToTransformImageReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTransformerStrategy) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.applyFormulaToTransformImageMutex.RLock()
	defer fake.applyFormulaToTransformImageMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTransformerStrategy) recordInvocation(key string, args []interface{}) {
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

var _ creating_symmetry.TransformerStrategy = new(FakeTransformerStrategy)