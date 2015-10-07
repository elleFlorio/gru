package utils

import (
	"testing"

	"github.com/elleFlorio/gru/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestRAMInBytes(t *testing.T) {
	human := [13]string{
		"32",
		"32b",
		"32B",
		"32k",
		"32K",
		"32kb",
		"32Kb",
		"32Mb",
		"32Gb",
		"32Tb",
		"32Pb",
		"32PB",
		"32P",
	}

	values := [13]int64{
		32,
		32,
		32,
		32 * KiB,
		32 * KiB,
		32 * KiB,
		32 * KiB,
		32 * MiB,
		32 * GiB,
		32 * TiB,
		32 * PiB,
		32 * PiB,
		32 * PiB,
	}

	errors := [9]string{
		"",
		"hello",
		"-32",
		"32.3",
		" 32 ",
		"32.3Kb",
		"32 mb",
		"32m b",
		"32bm",
	}

	for i := 0; i < len(human); i++ {
		ram, _ := RAMInBytes(human[i])
		assert.Equal(t, values[i], ram)
	}

	for i := 0; i < len(errors); i++ {
		_, err := RAMInBytes(errors[i])
		assert.Error(t, err)
	}
}
