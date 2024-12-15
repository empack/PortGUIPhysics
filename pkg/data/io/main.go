package io

import (
	"fmt"
	"physicsGUI/pkg/function"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type Measurement struct {
	Qz    float64
	Data  float64
	Error float64
}

var (
	// floatPattern is a regex pattern that matches a float number
	// it accepts scientific notation and normal floats
	// f.e. 1.0, 1.0e-10, 1.0e+10
	floatPattern = `(\d+\.\d+(?:[eE][+-]\d+)?)`

	// re is a regex pattern that matches a line with 3 floats
	// separated by whitespace
	// f.e. 1.0 2.0 3.0, 1.0e-10 2.0e-10 3.0e-10
	re = regexp.MustCompile(floatPattern + `\s+` + floatPattern + `\s+` + floatPattern)
)

func Parse(data []byte) ([]*Measurement, error) {
	lines := strings.Split(string(data), "\n")

	measurements := make([]*Measurement, 0)
	expectedLength := -1

	for i, v := range lines {
		if i == 0 {
			lengthString, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("parse error: expected int in first line: %v", err)
			}

			expectedLength = lengthString
			continue
		}

		// skip empty lines (possible at the end)
		if v == "" {
			continue
		}

		matches := re.FindStringSubmatch(v)
		if len(matches) != 4 {
			spew.Dump(matches)
			return nil, fmt.Errorf("parse error: expected '<FLOAT> <FLOAT> <FLOAT>' got '%s'", v)
		}

		// parse qz value
		qz, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return nil, fmt.Errorf("parse error: expected float in first column: %v", err)
		}

		// parse data/signal value
		data, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return nil, fmt.Errorf("parse error: expected float between first and last column: %v", err)
		}

		// parse error value
		ev, err := strconv.ParseFloat(matches[3], 64)
		if err != nil {
			return nil, fmt.Errorf("parse error: expected float in last column: %v", err)
		}

		// ? maybe change whole measurement to point already?
		measurements = append(measurements, &Measurement{
			Qz:    qz,
			Data:  data,
			Error: ev,
		})
	}

	if expectedLength != len(measurements) {
		return nil, fmt.Errorf("parse error: expected length (%d) does not match actual length (%d)", expectedLength, len(measurements))
	}

	return measurements, nil
}

// converts measurement to point
func (m *Measurement) ToPoint() *function.Point {
	return &function.Point{
		X:     m.Qz,
		Y:     m.Data,
		Error: m.Error,
	}
}
