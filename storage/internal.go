package storage

import (
	"errors"
)

var (
	ErrInvalidDataType error = errors.New("Invalid data type")
)

type internal struct {
	statsData     map[string][]byte
	analyticsData map[string][]byte
}

func (p *internal) Name() string {
	return "internal"
}

func (p *internal) Initialize() error {
	p.statsData = make(map[string][]byte)
	p.analyticsData = make(map[string][]byte)
	return nil
}

func (p *internal) StoreData(key string, data []byte, dataType string) error {
	switch dataType {
	case "stats":
		p.statsData[key] = data
	default:
		return ErrInvalidDataType
	}

	return nil
}

func (p *internal) GetData(key string, dataType string) ([]byte, error) {
	switch dataType {
	case "stats":
		return p.statsData[key], nil
	}

	return nil, ErrInvalidDataType
}

func (p *internal) GetAllData(dataType string) (map[string][]byte, error) {
	switch dataType {
	case "stats":
		return p.statsData, nil
	}

	return nil, ErrInvalidDataType
}

func (p *internal) DeleteData(key string, dataType string) error {
	switch dataType {
	case "stats":
		delete(p.statsData, key)
		return nil
	}

	return ErrInvalidDataType
}

func (p *internal) DeleteAllData(dataType string) error {
	switch dataType {
	case "stats":
		p.statsData = make(map[string][]byte)
		return nil
	}

	return ErrInvalidDataType
}
