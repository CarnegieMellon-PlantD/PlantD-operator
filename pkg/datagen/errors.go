package datagen

import (
	"strconv"
)

type TypeError string
type ColumnError string
type SchemaUndefinedError string
type OutputBuilderUndefinedError string
type OperationUndefinedError string
type OutOfIndexError string
type ResourceNotFoundError string
type VolumeUndefinedError string
type ComponentUndefinedError string
type DuplicateIDError string
type FormulaArgsError string
type NotImplementedError string
type MountPathNotExistError string

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

func (e OutputBuilderUndefinedError) Error() string {
	return "Cannot find Data Generator Configuration: " + string(e)
}

func (e VolumeUndefinedError) Error() string {
	return "Cannot find Volume: " + string(e)
}

func (e DuplicateIDError) Error() string {
	return "Experiment ID already exists: " + string(e)
}

func (e FormulaArgsError) Error() string {
	return "Formula got wrong arguments: " + string(e)
}

func (e NotImplementedError) Error() string {
	return "Feature not implemented: " + string(e)
}

func (e MountPathNotExistError) Error() string {
	return "Mount path does not exist: " + string(e)
}

func (e ComponentUndefinedError) Error() string {
	return "Component undefined: " + string(e)
}
