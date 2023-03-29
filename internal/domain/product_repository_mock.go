// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package domain

import (
	"sync"
)

// Ensure, that ProductRepositoryMock does implement ProductRepository.
// If this is not the case, regenerate this file with moq.
var _ ProductRepository = &ProductRepositoryMock{}

// ProductRepositoryMock is a mock implementation of ProductRepository.
//
//	func TestSomethingThatUsesProductRepository(t *testing.T) {
//
//		// make and configure a mocked ProductRepository
//		mockedProductRepository := &ProductRepositoryMock{
//			FindByIdFunc: func(id string) Product {
//				panic("mock out the FindById method")
//			},
//			GetAllFunc: func() []Product {
//				panic("mock out the GetAll method")
//			},
//			StoreFunc: func(product Product) error {
//				panic("mock out the Store method")
//			},
//		}
//
//		// use mockedProductRepository in code that requires ProductRepository
//		// and then make assertions.
//
//	}
type ProductRepositoryMock struct {
	// FindByIdFunc mocks the FindById method.
	FindByIdFunc func(id string) Product

	// GetAllFunc mocks the GetAll method.
	GetAllFunc func() []Product

	// StoreFunc mocks the Store method.
	StoreFunc func(product Product) error

	// calls tracks calls to the methods.
	calls struct {
		// FindById holds details about calls to the FindById method.
		FindById []struct {
			// ID is the id argument value.
			ID string
		}
		// GetAll holds details about calls to the GetAll method.
		GetAll []struct {
		}
		// Store holds details about calls to the Store method.
		Store []struct {
			// Product is the product argument value.
			Product Product
		}
	}
	lockFindById sync.RWMutex
	lockGetAll   sync.RWMutex
	lockStore    sync.RWMutex
}

// FindById calls FindByIdFunc.
func (mock *ProductRepositoryMock) FindById(id string) Product {
	if mock.FindByIdFunc == nil {
		panic("ProductRepositoryMock.FindByIdFunc: method is nil but ProductRepository.FindById was just called")
	}
	callInfo := struct {
		ID string
	}{
		ID: id,
	}
	mock.lockFindById.Lock()
	mock.calls.FindById = append(mock.calls.FindById, callInfo)
	mock.lockFindById.Unlock()
	return mock.FindByIdFunc(id)
}

// FindByIdCalls gets all the calls that were made to FindById.
// Check the length with:
//
//	len(mockedProductRepository.FindByIdCalls())
func (mock *ProductRepositoryMock) FindByIdCalls() []struct {
	ID string
} {
	var calls []struct {
		ID string
	}
	mock.lockFindById.RLock()
	calls = mock.calls.FindById
	mock.lockFindById.RUnlock()
	return calls
}

// GetAll calls GetAllFunc.
func (mock *ProductRepositoryMock) GetAll() []Product {
	if mock.GetAllFunc == nil {
		panic("ProductRepositoryMock.GetAllFunc: method is nil but ProductRepository.GetAll was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetAll.Lock()
	mock.calls.GetAll = append(mock.calls.GetAll, callInfo)
	mock.lockGetAll.Unlock()
	return mock.GetAllFunc()
}

// GetAllCalls gets all the calls that were made to GetAll.
// Check the length with:
//
//	len(mockedProductRepository.GetAllCalls())
func (mock *ProductRepositoryMock) GetAllCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetAll.RLock()
	calls = mock.calls.GetAll
	mock.lockGetAll.RUnlock()
	return calls
}

// Store calls StoreFunc.
func (mock *ProductRepositoryMock) Store(product Product) error {
	if mock.StoreFunc == nil {
		panic("ProductRepositoryMock.StoreFunc: method is nil but ProductRepository.Store was just called")
	}
	callInfo := struct {
		Product Product
	}{
		Product: product,
	}
	mock.lockStore.Lock()
	mock.calls.Store = append(mock.calls.Store, callInfo)
	mock.lockStore.Unlock()
	return mock.StoreFunc(product)
}

// StoreCalls gets all the calls that were made to Store.
// Check the length with:
//
//	len(mockedProductRepository.StoreCalls())
func (mock *ProductRepositoryMock) StoreCalls() []struct {
	Product Product
} {
	var calls []struct {
		Product Product
	}
	mock.lockStore.RLock()
	calls = mock.calls.Store
	mock.lockStore.RUnlock()
	return calls
}
