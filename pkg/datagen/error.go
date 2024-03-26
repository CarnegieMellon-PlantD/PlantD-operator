package datagen

import (
	"strconv"
)

type TypeError string
type ColumnError string
type SchemaUndefinedError string
type OperationUndefinedError string
type OutOfIndexError string
type ResourceNotFoundError string
type FormulaArgsError string

type NumParamError int

func (e TypeError) Error() string {
	return "Variable Type Mismatch: " + string(e)
}

func (e ColumnError) Error() string {
	return "Column Bad Defined: " + string(e)
}

func (e SchemaUndefinedError) Error() string {
	return "Schema Undefined: " + string(e)
}

func (e OperationUndefinedError) Error() string {
	return "Operation Undefined: " + string(e)
}

func (e OutOfIndexError) Error() string {
	return "Out of index: " + string(e)
}

func (e ResourceNotFoundError) Error() string {
	return "Resource Not Found: " + string(e)
}

func (e NumParamError) Error() string {
	return "Expect 1 parameter, but got: " + strconv.Itoa(int(e))
}

func (e FormulaArgsError) Error() string {
	return "Formula got wrong arguments: " + string(e)
}
