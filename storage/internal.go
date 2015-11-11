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
	ErrInvalidDataType error = errors.New("Invalid data type")
)

type internal struct {
	statsData     map[string][]byte
	analyticsData map[string][]byte
	plansData     map[string][]byte
}

func (p *internal) Name() string {
	return "internal"
}

func (p *internal) Initialize() error {
	p.statsData = make(map[string][]byte)
	p.analyticsData = make(map[string][]byte)
	p.plansData = make(map[string][]byte)
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
	case enum.PLANS:
		mutex_p.Lock()
		p.plansData[key] = data
		mutex_p.Unlock()
	}
	runtime.Gosched()

	return nil
}

func (p *internal) GetData(key string, dataType enum.Datatype) ([]byte, error) {
	var data []byte
	switch dataType {
	case enum.STATS:
		mutex_s.RLock()
		data = p.statsData[key]
		mutex_s.RUnlock()
	case enum.ANALYTICS:
		mutex_a.RLock()
		data = p.analyticsData[key]
		mutex_a.RUnlock()
	case enum.PLANS:
		mutex_p.RLock()
		data = p.plansData[key]
		mutex_p.RUnlock()
	}
	runtime.Gosched()

	return data, nil
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
	case enum.PLANS:
		mutex_p.RLock()
		data = p.plansData
		mutex_p.RUnlock()
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
	case enum.PLANS:
		mutex_p.Lock()
		delete(p.plansData, key)
		mutex_p.Unlock()
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
	case enum.PLANS:
		mutex_p.Lock()
		p.plansData = make(map[string][]byte)
		mutex_p.Unlock()
	}
	runtime.Gosched()

	return nil
}
