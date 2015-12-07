package utils

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

type noError struct {
	Pippo    int
	Topolino float64
	Paperino string
}

type wrongField struct {
	Pippo    int
	Topolino float64
	Paperone string
}

type wrongType struct {
	Pippo    string
	Topolino float64
	Paperino string
}

type cantSet struct {
	Pippo    string
	topolino float64
	Paperino string
}

func TestFillStruct(t *testing.T) {
	testMap := map[string]interface{}{
		"Pippo":    5,
		"Topolino": 4.6,
		"Paperino": "ciao",
	}
	noErr := &noError{}
	err1 := &wrongField{}
	err2 := &wrongType{}
	err3 := &cantSet{}

	assert.NoError(t, FillStruct(noErr, testMap))
	assert.Error(t, FillStruct(err1, testMap))
	assert.Error(t, FillStruct(err2, testMap))
	assert.Error(t, FillStruct(err3, testMap))
}
