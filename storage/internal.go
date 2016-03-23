package storage

import (
	"errors"
	"runtime"
	"sync"

	"github.com/elleFlorio/gru/enum"
)

var (
	mutex_s                  = sync.RWMutex{}
	mutex_a                  = sync.RWMutex{}
	mutex_p                  = sync.RWMutex{}
	mutex_i                  = sync.RWMutex{}
	ErrInvalidDataType error = errors.New("Invalid data type")
	ErrNoData          error = errors.New("No such data")
)

type internal struct {
	statsData     map[string][]byte
	analyticsData map[string][]byte
	policiesData  map[string][]byte
	sharedData    map[string][]byte
}

func (p *internal) Name() string {
	return "internal"
}

func (p *internal) Initialize() error {
	p.statsData = make(map[string][]byte)
	p.analyticsData = make(map[string][]byte)
	p.policiesData = make(map[string][]byte)
	p.sharedData = make(map[string][]byte)
	return nil
}

func (p *internal) StoreData(key string, data []byte, dataType enum.Datatype) error {
	switch dataType {
	case enum.STATS:
		mutex_s.Lock()
		p.statsData[key] = data
		mutex_s.Unlock()
	case enum.ANALYTICS:
		mutex_a.Lock()
		p.analyticsData[key] = data
		mutex_a.Unlock()
	case enum.POLICIES:
		mutex_p.Lock()
		p.policiesData[key] = data
		mutex_p.Unlock()
	case enum.SHARED:
		mutex_i.Lock()
		p.sharedData[key] = data
		mutex_i.Unlock()
	}
	runtime.Gosched()

	return nil
}

func (p *internal) GetData(key string, dataType enum.Datatype) ([]byte, error) {
	var data []byte
	var ok bool
	switch dataType {
	case enum.STATS:
		mutex_s.RLock()
		data, ok = p.statsData[key]
		mutex_s.RUnlock()
	case enum.ANALYTICS:
		mutex_a.RLock()
		data, ok = p.analyticsData[key]
		mutex_a.RUnlock()
	case enum.POLICIES:
		mutex_p.RLock()
		data, ok = p.policiesData[key]
		mutex_p.RUnlock()
	case enum.SHARED:
		mutex_i.RLock()
		data, ok = p.sharedData[key]
		mutex_i.RUnlock()
	}
	runtime.Gosched()

	if ok {
		return data, nil
	}

	return data, ErrNoData
}

func (p *internal) GetAllData(dataType enum.Datatype) (map[string][]byte, error) {
	var data map[string][]byte
	switch dataType {
	case enum.STATS:
		mutex_s.RLock()
		data = p.statsData
		mutex_s.RUnlock()
	case enum.ANALYTICS:
		mutex_a.RLock()
		data = p.analyticsData
		mutex_a.RUnlock()
	case enum.POLICIES:
		mutex_p.RLock()
		data = p.policiesData
		mutex_p.RUnlock()
	case enum.SHARED:
		mutex_i.RLock()
		data = p.sharedData
		mutex_i.RUnlock()
	}
	runtime.Gosched()

	return data, nil
}

func (p *internal) DeleteData(key string, dataType enum.Datatype) error {
	switch dataType {
	case enum.STATS:
		mutex_s.Lock()
		delete(p.statsData, key)
		mutex_s.Unlock()
	case enum.ANALYTICS:
		mutex_a.Lock()
		delete(p.analyticsData, key)
		mutex_a.Unlock()
	case enum.POLICIES:
		mutex_p.Lock()
		delete(p.policiesData, key)
		mutex_p.Unlock()
	case enum.SHARED:
		mutex_i.Lock()
		delete(p.sharedData, key)
		mutex_i.Unlock()
	}

	return nil
}

func (p *internal) DeleteAllData(dataType enum.Datatype) error {
	switch dataType {
	case enum.STATS:
		mutex_s.Lock()
		p.statsData = make(map[string][]byte)
		mutex_s.Unlock()
	case enum.ANALYTICS:
		mutex_a.Lock()
		p.analyticsData = make(map[string][]byte)
		mutex_a.Unlock()
	case enum.POLICIES:
		mutex_p.Lock()
		p.policiesData = make(map[string][]byte)
		mutex_p.Unlock()
	case enum.SHARED:
		mutex_i.Lock()
		p.sharedData = make(map[string][]byte)
		mutex_i.Unlock()
	}
	runtime.Gosched()

	return nil
}
