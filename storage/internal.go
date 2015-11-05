package storage

import (
	"errors"

	"github.com/elleFlorio/gru/enum"
)

var (
	ErrInvalidDataType error = errors.New("Invalid data type")
)

type internal struct {
	statsData     map[string][]byte
	analyticsData map[string][]byte
	plansData     map[string][]byte
	metricsData   map[string][]byte
}

func (p *internal) Name() string {
	return "internal"
}

func (p *internal) Initialize() error {
	p.statsData = make(map[string][]byte)
	p.analyticsData = make(map[string][]byte)
	p.plansData = make(map[string][]byte)
	p.metricsData = make(map[string][]byte)
	return nil
}

func (p *internal) StoreData(key string, data []byte, dataType enum.Datatype) error {
	switch dataType {
	case enum.STATS:
		p.statsData[key] = data
	case enum.ANALYTICS:
		p.analyticsData[key] = data
	case enum.PLANS:
		p.plansData[key] = data
	case enum.METRICS:
		p.metricsData[key] = data
	}

	return nil
}

func (p *internal) GetData(key string, dataType enum.Datatype) ([]byte, error) {
	var data []byte
	switch dataType {
	case enum.STATS:
		data = p.statsData[key]
	case enum.ANALYTICS:
		data = p.analyticsData[key]
	case enum.PLANS:
		data = p.plansData[key]
	case enum.METRICS:
		data = p.metricsData[key]
	}

	return data, nil
}

func (p *internal) GetAllData(dataType enum.Datatype) (map[string][]byte, error) {
	var data map[string][]byte
	switch dataType {
	case enum.STATS:
		data = p.statsData
	case enum.ANALYTICS:
		data = p.analyticsData
	case enum.PLANS:
		data = p.plansData
	case enum.METRICS:
		data = p.metricsData
	}

	return data, nil
}

func (p *internal) DeleteData(key string, dataType enum.Datatype) error {
	switch dataType {
	case enum.STATS:
		delete(p.statsData, key)
	case enum.ANALYTICS:
		delete(p.analyticsData, key)
	case enum.PLANS:
		delete(p.plansData, key)
	case enum.METRICS:
		delete(p.metricsData, key)
	}

	return nil
}

func (p *internal) DeleteAllData(dataType enum.Datatype) error {
	switch dataType {
	case enum.STATS:
		p.statsData = make(map[string][]byte)
	case enum.ANALYTICS:
		p.analyticsData = make(map[string][]byte)
	case enum.PLANS:
		p.plansData = make(map[string][]byte)
	case enum.METRICS:
		p.metricsData = make(map[string][]byte)
	}

	return nil
}
