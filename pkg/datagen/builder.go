package datagen

import (
	"math/rand"
	"path/filepath"

	windtunnelv1alpha1 "github.com/CarnegieMellon-PlantD/PlantD-operator/api/v1alpha1"
	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"

	"github.com/brianvoe/gofakeit/v6"
)

var (
	path = config.GetString("dataGenerator.path")
)

type ColumnBuilder struct {
	Name          string
	Info          *gofakeit.Info
	InfoMapParams *gofakeit.MapParams
	Formula       Formula
	FormulaArgs   []string
}

type SchemaBuilder struct {
	ColBuilders                    []*ColumnBuilder
	ParentPath                     string
	SchemaName                     string
	NumRecords                     int
	NumberOfFilesPerCompressedFile map[string]int
}

type OutputBuilder struct {
	CompressPerSchema bool
	SchBuilders       []*SchemaBuilder
	Operations        []Operation
	Path              string
	Name              string
}

// PutParams creates a gofakeit.MapParams instance based on the provided column and parameters.
func PutParams(in windtunnelv1alpha1.ColumnSpec, params []gofakeit.Param) *gofakeit.MapParams {
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
		ColBuilders: make([]*ColumnBuilder, numCol),
		SchemaName:  schema.Name,
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

// NewOutputBuilder creates a new OutputBuilder based on the provided output configuration.
func NewOutputBuilder(output *windtunnelv1alpha1.DataSet) (*OutputBuilder, error) {
	var lenOperations int
	if output.Spec.CompressedFileFormat != "" {
		lenOperations = 2
	} else {
		lenOperations = 1
	}
	outputBuilder := &OutputBuilder{
		SchBuilders: make([]*SchemaBuilder, len(output.Spec.Schemas)),

		Operations:        make([]Operation, lenOperations),
		Path:              path,
		Name:              output.Name,
		CompressPerSchema: output.Spec.CompressPerSchema,
	}
	for i, sch := range output.Spec.Schemas {
		outputBuilder.SchBuilders[i] = GetSchemaBuilder(sch.Name)
		if outputBuilder.SchBuilders[i].ColBuilders == nil {
			return nil, SchemaUndefinedError(sch.Name)
		}
		minRec := sch.NumRecords["min"]

		maxRec := sch.NumRecords["max"]

		outputBuilder.SchBuilders[i].NumRecords = gofakeit.Number(minRec, maxRec)
		outputBuilder.SchBuilders[i].ParentPath = filepath.Join(path, sch.Name)
		outputBuilder.SchBuilders[i].NumberOfFilesPerCompressedFile = sch.NumberOfFilesPerCompressedFile
	}

	outputBuilder.Operations[0] = GetOpLookups(output.Spec.FileFormat)
	if outputBuilder.Operations[0] == nil {
		return nil, OperationUndefinedError(output.Spec.FileFormat)
	}

	if output.Spec.CompressedFileFormat != "" {
		outputBuilder.Operations[0] = GetOpLookups(output.Spec.FileFormat + "@cache")
		if outputBuilder.Operations[0] == nil {
			return nil, OperationUndefinedError(output.Spec.FileFormat + "@cache")
		}
		outputBuilder.Operations[1] = GetOpLookups(output.Spec.FileFormat + "->" + output.Spec.CompressedFileFormat)
		if outputBuilder.Operations[1] == nil {
			return nil, OperationUndefinedError(output.Spec.FileFormat + "->" + output.Spec.CompressedFileFormat)
		}
	}

	NewFakeDataCache(outputBuilder)
	return outputBuilder, nil
}

// Build generates fake data based on the provided SchemaBuilder.
func (schBldr *SchemaBuilder) Build(r *rand.Rand) error {
	for _, colBldr := range schBldr.ColBuilders {
		var fakeData interface{}
		var err error
		key := GetKey(schBldr, colBldr)

		if colBldr.Info != nil {
			for i := 0; i < schBldr.NumRecords; i++ {
				fakeData, err = colBldr.Info.Generate(r, colBldr.InfoMapParams, colBldr.Info)
				if err != nil {
					return err
				}
				PutFakeData(key, i, fakeData)
			}
		}

		if colBldr.Formula != nil {
			for i := 0; i < schBldr.NumRecords; i++ {
				fakeData, err = colBldr.Formula(i, colBldr.FormulaArgs...)
				if err != nil {
					return err
				}
				PutFakeData(key, i, fakeData)
			}
		}
	}
	return nil
}
