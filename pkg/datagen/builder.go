package datagen

import (
	"fmt"
	"path/filepath"

	"github.com/brianvoe/gofakeit/v7"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
)

type ColumnBuilder struct {
	// Name of the ColumnBuilder
	Name string
	// Struct containing reference to gofakeit function
	Info *gofakeit.Info
	// Struct containing parameters for gofakeit function
	InfoMapParams *gofakeit.MapParams
	// Function of formula
	Formula Formula
	// Parameters for formula
	FormulaArgs []string
}

type SchemaBuilder struct {
	// Name of the Schema, also name of the SchemaBuilder
	SchemaName string
	// File path the SchemaBuilder works in
	Path string
	// ColumnBuilders of the SchemaBuilder
	ColBuilders []*ColumnBuilder
	// Number of records per file
	NumRecords int
	// Number of files per compressed file
	NumFilesPerCompressedFile int
	// Total number of records the SchemaBuilder should generate
	TotalNumRecords int
}

type OutputBuilder struct {
	// Name of the OutputBuilder
	Name string
	// File path the OutputBuilder works in
	Path string
	// SchemaBuilders of the OutputBuilder
	SchBuilders []*SchemaBuilder
	// Operations the OutputBuilder should do
	Operations []Operation
	// Whether compressed file should be created per Schema
	CompressPerSchema bool
}

// PutParams creates a gofakeit.MapParams instance based on the provided column and parameters.
func PutParams(in windtunnelv1alpha1.Column, params []gofakeit.Param) *gofakeit.MapParams {
	out := gofakeit.NewMapParams()
	for _, param := range params {
		field := param.Field
		v, ok := in.Params[field]
		if ok {
			out.Add(field, v)
		} else {
			out.Add(field, param.Default)
		}
	}
	return out
}

// NewSchemaBuilder creates a new SchemaBuilder based on the provided schema.
func NewSchemaBuilder(schema *windtunnelv1alpha1.Schema) (*SchemaBuilder, error) {
	numCol := len(schema.Spec.Columns)
	schBldr := SchemaBuilder{
		SchemaName:  schema.Name,
		ColBuilders: make([]*ColumnBuilder, numCol),
	}
	colNames := make([]string, numCol)
	for i, col := range schema.Spec.Columns {
		info := gofakeit.GetFuncLookup(col.Type)
		formula := GetFormulaLookup(col.Formula.Name)
		var infoParams *gofakeit.MapParams
		if info != nil {
			infoParams = PutParams(col, info.Params)
		} else if formula == nil {
			return nil, ColumnError(col.Name)
		}

		schBldr.ColBuilders[i] = &ColumnBuilder{
			Name:          col.Name,
			Info:          info,
			InfoMapParams: infoParams,
			Formula:       formula,
			FormulaArgs:   col.Formula.Args,
		}

		colNames[i] = GetKey(&schBldr, schBldr.ColBuilders[i])
	}

	PutColumnNames(schema.Name, colNames)
	return &schBldr, nil
}

// Build generates fake data based on the provided SchemaBuilder.
func (schBldr *SchemaBuilder) Build(faker *gofakeit.Faker) error {
	for _, colBldr := range schBldr.ColBuilders {
		// Prepare fake data in the cache for this column
		var fakeData interface{}
		var err error
		key := GetKey(schBldr, colBldr)

		if colBldr.Info != nil {
			for i := 0; i < schBldr.TotalNumRecords; i++ {
				fakeData, err = colBldr.Info.Generate(faker, colBldr.InfoMapParams, colBldr.Info)
				if err != nil {
					return err
				}
				PutFakeData(key, i, fakeData)
			}
		}

		if colBldr.Formula != nil {
			for i := 0; i < schBldr.TotalNumRecords; i++ {
				fakeData, err = colBldr.Formula(faker, i, colBldr.FormulaArgs...)
				if err != nil {
					return err
				}
				PutFakeData(key, i, fakeData)
			}
		}
	}
	return nil
}

// NewOutputBuilder creates a new OutputBuilder based on the provided output configuration.
func NewOutputBuilder(dataSet *windtunnelv1alpha1.DataSet, path string) (*OutputBuilder, error) {
	var lenOperations int
	if dataSet.Spec.CompressedFileFormat == "" {
		lenOperations = 1
	} else {
		lenOperations = 2
	}

	outBldr := &OutputBuilder{
		Name:              dataSet.Name,
		Path:              path,
		SchBuilders:       make([]*SchemaBuilder, len(dataSet.Spec.Schemas)),
		Operations:        make([]Operation, lenOperations),
		CompressPerSchema: dataSet.Spec.CompressPerSchema,
	}

	for i, sch := range dataSet.Spec.Schemas {
		outBldr.SchBuilders[i] = GetSchemaBuilder(sch.Name)
		if outBldr.SchBuilders[i].ColBuilders == nil {
			return nil, SchemaUndefinedError(sch.Name)
		}

		outBldr.SchBuilders[i].Path = filepath.Join(path, sch.Name)
	}

	if dataSet.Spec.CompressedFileFormat == "" {
		outBldr.Operations[0] = GetOpLookups(dataSet.Spec.FileFormat)
		if outBldr.Operations[0] == nil {
			return nil, OperationUndefinedError(dataSet.Spec.FileFormat)
		}
	} else {
		op := fmt.Sprintf("%s@cache", dataSet.Spec.FileFormat)
		outBldr.Operations[0] = GetOpLookups(op)
		if outBldr.Operations[0] == nil {
			return nil, OperationUndefinedError(op)
		}

		compressionOp := fmt.Sprintf("%s->%s", dataSet.Spec.FileFormat, dataSet.Spec.CompressedFileFormat)
		outBldr.Operations[1] = GetOpLookups(compressionOp)
		if outBldr.Operations[1] == nil {
			return nil, OperationUndefinedError(compressionOp)
		}
	}

	return outBldr, nil
}

// SetRandomnessAndCache sets the number of records and number of files per compressed file for each SchemaBuilder in
// the OutputBuilder and initializes the fake data cache.
func (outBldr *OutputBuilder) SetRandomnessAndCache(faker *gofakeit.Faker, dataSet *windtunnelv1alpha1.DataSet) {
	for i, sch := range dataSet.Spec.Schemas {
		outBldr.SchBuilders[i].NumRecords = faker.Number(sch.NumRecords.Min, sch.NumRecords.Max)
		outBldr.SchBuilders[i].NumFilesPerCompressedFile = faker.Number(sch.NumRecords.Min, sch.NumRecords.Max)
		if dataSet.Spec.CompressedFileFormat == "" {
			outBldr.SchBuilders[i].TotalNumRecords = outBldr.SchBuilders[i].NumRecords
		} else {
			outBldr.SchBuilders[i].TotalNumRecords = outBldr.SchBuilders[i].NumRecords * outBldr.SchBuilders[i].NumFilesPerCompressedFile
		}
	}

	NewFakeDataCache(outBldr)
}
