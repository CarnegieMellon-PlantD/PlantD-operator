package datagen

type SchemaBuilderCache map[string]*SchemaBuilder
type ColumnNamesCache map[string][]string
type FakeDataCache map[string][]interface{} // key is schema name + column name

var schemaBuilderCache SchemaBuilderCache
var columnNamesCache ColumnNamesCache
var fakeDataCache FakeDataCache

// init initializes the caches during package initialization.
func init() {
	initSchemaCache()
	initSeqCache()
}

// initSchemaCache initializes the schema builder cache if it is nil.
func initSchemaCache() {
	if schemaBuilderCache == nil {
		schemaBuilderCache = make(SchemaBuilderCache)
	}
}

// initSeqCache initializes the column names cache if it is nil.
func initSeqCache() {
	if columnNamesCache == nil {
		columnNamesCache = make(ColumnNamesCache)
	}
}

// GetKey generates a unique key for a schema builder and column builder combination.
func GetKey(schBldr *SchemaBuilder, colBldr *ColumnBuilder) string {
	return schBldr.SchemaName + "." + colBldr.Name
}

// NewFakeDataCache creates a new fake data cache based on the provided output builder.
func NewFakeDataCache(outputBuilder *OutputBuilder) {
	mapLen := 0
	for _, schBldr := range outputBuilder.SchBuilders {
		mapLen += len(schBldr.ColBuilders) + 1
	}
	fakeDataCache = make(FakeDataCache, mapLen)
	for _, schBldr := range outputBuilder.SchBuilders {
		for _, colBldr := range schBldr.ColBuilders {
			fakeDataCache[GetKey(schBldr, colBldr)] = make([]interface{}, schBldr.NumRecords)
		}
		fakeDataCache[schBldr.SchemaName] = make([]interface{}, schBldr.NumRecords)
	}
}

// PutSchemaBuilder adds a schema builder to the schema builder cache.
func PutSchemaBuilder(name string, schBldr *SchemaBuilder) {
	schemaBuilderCache[name] = schBldr
}

// GetSchemaBuilder retrieves a schema builder from the schema builder cache by name.
func GetSchemaBuilder(name string) *SchemaBuilder {
	if schBldr, ok := schemaBuilderCache[name]; ok {
		return schBldr
	}
	return nil
}

// PutColumnNames adds column names to the column names cache for a specific schema.
func PutColumnNames(schemaName string, columnNames []string) {
	columnNamesCache[schemaName] = columnNames
}

// GetColumnNames retrieves column names from the column names cache for a specific schema.
func GetColumnNames(schemaName string) []string {
	if columnNames, ok := columnNamesCache[schemaName]; ok {
		return columnNames
	}
	return nil
}

// PutFakeData stores fake data in the fake data cache for a specific key and record ID.
func PutFakeData(key string, recordID int, v interface{}) {
	fakeDataCache[key][recordID] = v
}

// GetFakeData retrieves fake data from the fake data cache for a specific key and record ID.
func GetFakeData(key string, recordID int) (interface{}, error) {
	if colDataList, ok := fakeDataCache[key]; ok {
		if recordID >= len(colDataList) {
			return nil, OutOfIndexError(key)
		}
		return colDataList[recordID], nil
	}
	return nil, ResourceNotFoundError(key)
}
