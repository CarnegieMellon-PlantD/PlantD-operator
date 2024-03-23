package datagen

import (
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

var formulaLookups map[string]Formula

// Formula represents a function that generates data based on a sequence number and arguments.
type Formula func(seqNum int, args ...string) (interface{}, error)

func init() {
	initFormulaLookups()
}

// initFormulaLookups initializes the formula lookups map if it is nil.
func initFormulaLookups() {
	if formulaLookups == nil {
		formulaLookups = make(map[string]Formula)
	}
	PutFormulaLookup("AddInt", AddInt)
	PutFormulaLookup("AddFloat", AddFloat)
	PutFormulaLookup("AddString", AddString)
	PutFormulaLookup("And", And)
	PutFormulaLookup("Or", Or)
	PutFormulaLookup("XOrInt", XOrInt)
	PutFormulaLookup("Copy", Copy)
	PutFormulaLookup("CurrentTimeMs", CurrentTimeMs)
	PutFormulaLookup("ToUnixMilli", ToUnixMilli)
	PutFormulaLookup("AddRandomTimeMs", AddRandomTimeMs)
	PutFormulaLookup("AddRandomNumber", AddRandomNumber)
}

// PutFormulaLookup adds a formula function to the formula lookups map.
func PutFormulaLookup(formulaName string, formula Formula) {
	formulaLookups[formulaName] = formula
}

// GetFormulaLookup retrieves a formula function from the formula lookups map by name.
func GetFormulaLookup(formulaName string) Formula {
	if formula, ok := formulaLookups[formulaName]; ok {
		return formula
	}
	return nil
}

// AddInt calculates the sum of integer values retrieved from the fake data cache.
func AddInt(seqNum int, args ...string) (interface{}, error) {
	sum := 0
	for _, param := range args {
		if fakeData, err := GetFakeData(param, seqNum); err == nil {
			if v, ok := fakeData.(int); ok {
				sum += v
			} else {
				return nil, TypeError(param)
			}
		} else {
			return nil, err
		}
	}
	return sum, nil
}

// AddFloat calculates the sum of float values retrieved from the fake data cache.
func AddFloat(seqNum int, args ...string) (interface{}, error) {
	sum := 0.0
	for _, param := range args {
		if fakeData, err := GetFakeData(param, seqNum); err == nil {
			if v, ok := fakeData.(float64); ok {
				sum += v
			} else {
				return nil, TypeError(param)
			}
		} else {
			return nil, err
		}
	}
	return sum, nil
}

// AddString concatenates string values retrieved from the fake data cache.
func AddString(seqNum int, args ...string) (interface{}, error) {
	sum := ""
	for _, param := range args {
		if fakeData, err := GetFakeData(param, seqNum); err == nil {
			if v, ok := fakeData.(string); ok {
				sum += v
			} else {
				return nil, TypeError(param)
			}
		} else {
			return nil, err
		}
	}
	return sum, nil
}

// And calculates the logical AND operation on boolean values retrieved from the fake data cache.
func And(seqNum int, args ...string) (interface{}, error) {
	res := true
	for _, param := range args {
		if fakeData, err := GetFakeData(param, seqNum); err == nil {
			if v, ok := fakeData.(bool); ok {
				res = res && v
			} else {
				return nil, TypeError(param)
			}
		} else {
			return nil, err
		}
	}
	return res, nil
}

// Or calculates the logical OR operation on boolean values retrieved from the fake data cache.
func Or(seqNum int, args ...string) (interface{}, error) {
	res := false
	for _, param := range args {
		if fakeData, err := GetFakeData(param, seqNum); err == nil {
			if v, ok := fakeData.(bool); ok {
				res = res || v
			} else {
				return nil, TypeError(param)
			}
		} else {
			return nil, err
		}
	}
	return res, nil
}

// XOrInt calculates the XOR operation on integer values retrieved from the fake data cache.
func XOrInt(seqNum int, args ...string) (interface{}, error) {
	res := 0
	for _, param := range args {
		if fakeData, err := GetFakeData(param, seqNum); err == nil {
			if v, ok := fakeData.(int); ok {
				res ^= v
			} else {
				return nil, TypeError(param)
			}
		} else {
			return nil, err
		}
	}
	return res, nil
}

// Copy retrieves a value from the fake data cache.
func Copy(seqNum int, args ...string) (interface{}, error) {
	if len(args) != 1 {
		return nil, NumParamError(len(args))
	}
	if fakeData, err := GetFakeDataFromRandomRecord(args[0]); err == nil {
		return fakeData, nil
	} else {
		return nil, err
	}
}

// CurrentTimeMs returns the current time in milliseconds.
func CurrentTimeMs(seqNum int, args ...string) (interface{}, error) {
	return time.Now().UnixMilli(), nil
}

// ToUnixMilli converts a string representation of a date to Unix milliseconds.
func ToUnixMilli(seqNum int, args ...string) (interface{}, error) {
	if fakeData, err := GetFakeData(args[0], seqNum); err == nil {
		if v, ok := fakeData.(string); ok {
			t, err := time.Parse("2006-01-02", v)
			if err != nil {
				return nil, err
			}
			return t.UnixMilli(), nil
		} else {
			return nil, TypeError(args[0])
		}
	} else {
		return nil, err
	}
}

// AddRandomTimeMs adds a random number of milliseconds to a given value retrieved from the fake data cache.
func AddRandomTimeMs(seqNum int, args ...string) (interface{}, error) {
	min, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, FormulaArgsError("AddRandomTimeMs.min")
	}
	max, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, FormulaArgsError("AddRandomTimeMs.max")
	}
	randomNum := gofakeit.Number(min, max)
	if fakeData, err := GetFakeData(args[0], seqNum); err == nil {
		if v, ok := fakeData.(int64); ok {
			return v + int64(randomNum), nil
		} else {
			return nil, TypeError(args[0])
		}
	} else {
		return nil, err
	}
}

// AddRandomNumber adds a random number to a given value retrieved from the fake data cache.
func AddRandomNumber(seqNum int, args ...string) (interface{}, error) {
	min, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, FormulaArgsError("AddRandomNumber.min")
	}
	max, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, FormulaArgsError("AddRandomNumber.max")
	}
	randomNum := gofakeit.Number(min, max)
	if fakeData, err := GetFakeData(args[0], seqNum); err == nil {
		if v, ok := fakeData.(int); ok {
			return v + randomNum, nil
		} else {
			return nil, TypeError(args[0])
		}
	} else {
		return nil, err
	}
}
