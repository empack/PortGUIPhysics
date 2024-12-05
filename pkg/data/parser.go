package data

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type measurement struct {
	Count int
	Time  float64
	Data  []float64
	Error float64
}

func newMeasurment(time float64, error float64, count int, data []float64) measurement {
	return measurement{
		Time:  time,
		Count: count,
		Data:  data,
		Error: error,
	}
}

type FileParser interface {
	tryParse([]byte) ([]measurement, error) //TODO change from Measurement to programState
}

type DatParser struct {
	valueSplitter       [][]byte
	measurementSplitter []byte
}

func DefaultDatParser() *DatParser {
	return &DatParser{
		valueSplitter:       [][]byte{{'\t', ' '}},
		measurementSplitter: []byte{'\n'},
	}
}

func (p *DatParser) tryParse(data []byte) ([]measurement, error) {
	var pre = string(data)
	// https://stackoverflow.com/a/37293398
	reLeadcloseWhtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	reInsideWhtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	pre = reLeadcloseWhtsp.ReplaceAllString(pre, "")
	pre = reInsideWhtsp.ReplaceAllString(pre, " ")
	data = []byte(pre)
	var err = errors.New("DAT_Parsing_Error: No Measurement separator defined")
	var res []measurement = nil

	for _, b := range p.measurementSplitter {
		if err == nil {
			break
		}
		err = nil
		res = nil
		segments := bytes.Split(data, []byte{b})
		if len(segments) < 1 {
			err = errors.New("DAT_Parsing_Error: Couldn't find length indicator at begin of file")
			continue
		}
		mCount, mCountError := tryParseInt(segments[0])
		if mCountError != nil {
			err = errors.New("DAT_Parsing_Error: Couldn't parse length indicator at begin of file")
			continue
		}
		res = make([]measurement, mCount)
		if len(segments) < mCount+1 {
			err = errors.New("DAT_Parsing_Error: Less measurements found than header indicated")
			continue
		}

		for i := range mCount {
			seg := segments[i+1]

			var nums []float64 = nil
			for _, s := range p.valueSplitter {
				if nums != nil {
					break
				}
				var vals = bytes.Split(seg, []byte{s[0]})
				var nextLayer [][]byte
				for _, ms := range s[1:] {
					for _, val := range vals {
						for _, l := range bytes.Split(val, []byte{ms}) {
							nextLayer = append(nextLayer, l)
						}
					}
					vals = nextLayer
				}

				if len(vals) < 3 {
					continue
				}
				nums = make([]float64, len(vals))
				for j, c := range vals {
					num, err := tryParseFloat(c)
					if err != nil {
						err = errors.Join(errors.New(fmt.Sprintf("DAT_Parsing_Error: Failed to parse value at line %d:%d", i+2, bytes.Index(seg, c))), err)
						nums = nil
						break
					}
					nums[j] = num
				}
			}
			if nums == nil {
				err = errors.New(fmt.Sprintf("DAT_Parsing_Error: Failed to parse measurement at line %d", i+2))
				res = nil
				break
			}
			time := nums[0]
			valErr := nums[len(nums)-1]
			count := len(nums) - 2
			measuredData := nums[1 : len(nums)-1]
			res[i] = newMeasurment(time, valErr, count, measuredData)
		}
	}

	return res, err
}

func tryParseInt(data []byte) (int, error) {
	return strconv.Atoi(string(data))
}
func tryParseFloat(data []byte) (float64, error) {
	return strconv.ParseFloat(string(data), 64)
}
