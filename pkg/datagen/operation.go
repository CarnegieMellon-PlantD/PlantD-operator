package datagen

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// opLookups maps operation names to their corresponding functions.
var opLookups map[string]Operation

// Operation represents a function that performs a data generation operation.
type Operation func(outputBuilder *OutputBuilder, seqNum int) error

func init() {
	initOpLookups()
}

// initOpLookups initializes the opLookups map with supported operations and their corresponding functions.
func initOpLookups() {
	if opLookups == nil {
		opLookups = make(map[string]Operation)
	}

	// Register operation functions
	PutOpLookups("csv", Raw2CSVAtFile)
	PutOpLookups("csv@cache", Raw2CSVAtCache)
	PutOpLookups("binary", Raw2BinaryAtFile)
	PutOpLookups("binary@cache", Raw2BinaryAtCache)
	PutOpLookups("csv->zip", CSVAtCache2ZipAtFile)
	PutOpLookups("binary->zip", BinaryAtCache2ZipAtFile)
}

// PutOpLookups registers an operation function with a name in the opLookups map.
func PutOpLookups(name string, op Operation) {
	opLookups[name] = op
}

// GetOpLookups retrieves an operation function by its name from the opLookups map.
func GetOpLookups(name string) Operation {
	if op, ok := opLookups[name]; ok {
		return op
	}
	return nil
}

// Raw2CSVAtFile converts raw data to CSV format and writes it to a file.
func Raw2CSVAtFile(outputBuilder *OutputBuilder, seqNum int) error {
	for _, schBldr := range outputBuilder.SchBuilders {
		// Get random number of records for the Schema
		filePath := filepath.Join(schBldr.Path, fmt.Sprintf("%s_%s_%d.csv", outputBuilder.Name, schBldr.SchemaName, seqNum))
		err := Raw2CSVAtFileBySchema(schBldr.NumRecords, schBldr.SchemaName, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// Raw2CSVAtFileBySchema converts raw data to CSV format for a specific schema and writes it to a file.
func Raw2CSVAtFileBySchema(numRecords int, schemaName, filePath string) error {
	colNames := GetColumnNames(schemaName)
	outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	line := make([]string, len(colNames))
	w := csv.NewWriter(outFile)

	trimColNames := make([]string, len(colNames))
	prefixPattern := schemaName + "."
	for i, colName := range colNames {
		trimColNames[i] = strings.TrimPrefix(colName, prefixPattern)
	}
	err = w.Write(trimColNames)
	if err != nil {
		return err
	}

	for i := 0; i < numRecords; i++ {
		for j, key := range colNames {
			if fakeData, err := GetFakeData(key, i); err == nil {
				line[j] = fmt.Sprint(fakeData)
			} else {
				return err
			}
		}
		if err = w.Write(line); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

// Raw2CSVAtCache converts raw data to CSV format and stores it in cache.
func Raw2CSVAtCache(outputBuilder *OutputBuilder, seqNum int) error {
	for _, schBldr := range outputBuilder.SchBuilders {
		err := Raw2CSVAtCacheBySchema(schBldr.TotalNumRecords, schBldr.SchemaName)
		if err != nil {
			return err
		}
	}
	return nil
}

// Raw2CSVAtCacheBySchema converts raw data to CSV format for a specific schema and stores it in cache.
func Raw2CSVAtCacheBySchema(numRecords int, schemaName string) error {
	colNames := GetColumnNames(schemaName)
	line := make([]string, len(colNames))

	for i := 0; i < numRecords; i++ {
		var buff bytes.Buffer
		w := csv.NewWriter(&buff)
		for j, key := range colNames {
			if fakeData, err := GetFakeData(key, i); err == nil {
				line[j] = fmt.Sprint(fakeData)
			} else {
				return err
			}
		}
		if err := w.Write(line); err != nil {
			return err
		}
		w.Flush()
		PutFakeData(schemaName, i, buff.Bytes())
	}
	return nil
}

// Raw2BinaryAtFile converts raw data to binary format and writes it to a file.
func Raw2BinaryAtFile(outputBuilder *OutputBuilder, seqNum int) error {
	for _, schBldr := range outputBuilder.SchBuilders {
		filePath := filepath.Join(schBldr.Path, fmt.Sprintf("%s_%s_%d.bin", outputBuilder.Name, schBldr.SchemaName, seqNum))
		err := Raw2BinaryAtFileBySchema(schBldr.NumRecords, schBldr.SchemaName, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// Raw2BinaryAtFileBySchema converts raw data to binary format for a specific schema and writes it to a file.
func Raw2BinaryAtFileBySchema(numRecords int, schemaName, filePath string) error {
	colNames := GetColumnNames(schemaName)
	bColLenBuf := make([]byte, 4)
	outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for i := 0; i < numRecords; i++ {
		for _, key := range colNames {
			if fakeData, err := GetFakeData(key, i); err == nil {
				bCol := []byte(fmt.Sprint(fakeData))
				binary.BigEndian.PutUint32(bColLenBuf, uint32(len(bCol)))
				outFile.Write(bColLenBuf)
				outFile.Write(bCol)
			} else {
				return err
			}

		}
	}
	return nil
}

// Raw2BinaryAtCache converts raw data to binary format and stores it in cache.
func Raw2BinaryAtCache(outputBuilder *OutputBuilder, seqNum int) error {
	for _, schBldr := range outputBuilder.SchBuilders {
		err := Raw2BinaryAtCacheBySchema(schBldr.TotalNumRecords, schBldr.SchemaName)
		if err != nil {
			return err
		}
	}
	return nil
}

// Raw2BinaryAtCacheBySchema converts raw data to binary format for a specific schema and stores it in cache.
func Raw2BinaryAtCacheBySchema(numRecords int, schemaName string) error {
	colNames := GetColumnNames(schemaName)
	bColLenBuf := make([]byte, 4)

	for i := 0; i < numRecords; i++ {
		for _, key := range colNames {
			if fakeData, err := GetFakeData(key, i); err == nil {
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				err := enc.Encode(fakeData)
				if err != nil {
					return err
				}
				bCol := buf.Bytes()
				binary.BigEndian.PutUint32(bColLenBuf, uint32(len(bCol)))
				PutFakeData(key, i, append(bColLenBuf, bCol...))
			} else {
				return err
			}
		}
	}
	return nil
}

// CSVAtCache2ZipAtFile converts CSV data stored in cache to a zip file.
func CSVAtCache2ZipAtFile(outputBuilder *OutputBuilder, seqNum int) error {
	if outputBuilder.CompressPerSchema {
		for _, schBldr := range outputBuilder.SchBuilders {
			zipFilePath := filepath.Join(outputBuilder.Path, fmt.Sprintf("%s_%s_%d.zip", outputBuilder.Name, schBldr.SchemaName, seqNum))
			if err := CSVAtCache2ZipAtFileBySchema(
				schBldr.SchemaName,
				seqNum,
				schBldr.NumFilesPerCompressedFile,
				schBldr.NumRecords,
				true,
				zipFilePath,
				nil,
			); err != nil {
				return err
			}
		}
		return nil
	} else {
		zipFilePath := filepath.Join(outputBuilder.Path, fmt.Sprintf("%s_%d.zip", outputBuilder.Name, seqNum))
		zipFile, err := os.OpenFile(zipFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer zipFile.Close()
		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		for _, schBldr := range outputBuilder.SchBuilders {
			if err := CSVAtCache2ZipAtFileBySchema(
				schBldr.SchemaName,
				seqNum,
				schBldr.NumFilesPerCompressedFile,
				schBldr.NumRecords,
				false,
				"",
				zipWriter,
			); err != nil {
				return err
			}
		}

		return nil
	}
}

// CSVAtCache2ZipAtFileBySchema converts CSV data stored in cache to a zip file for a specific schema.
func CSVAtCache2ZipAtFileBySchema(schemaName string, seqNum int, numFilesPerCompressedFile int, numRecords int,
	compressPerFile bool, zipFilePath string, zipWriter *zip.Writer) error {
	if compressPerFile {
		zipFile, err := os.OpenFile(zipFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer zipFile.Close()
		zipWriter = zip.NewWriter(zipFile)
		defer zipWriter.Close()
	}

	for i := 0; i < numFilesPerCompressedFile; i++ {
		fileName := fmt.Sprintf("%s_%d_%d.csv", schemaName, seqNum, i)
		fWriter, err := zipWriter.Create(fileName)
		if err != nil {
			return err
		}

		colNames := GetColumnNames(schemaName)
		trimColNames := make([]string, len(colNames))
		prefixPattern := schemaName + "."
		for i, colName := range colNames {
			trimColNames[i] = strings.TrimPrefix(colName, prefixPattern)
		}

		var headerBuff bytes.Buffer
		w := csv.NewWriter(&headerBuff)
		if err := w.Write(trimColNames); err != nil {
			return err
		}
		w.Flush()

		if _, err := fWriter.Write(headerBuff.Bytes()); err != nil {
			return err
		}

		for j := 0; j < numRecords; j++ {
			row, err := GetFakeData(schemaName, i*numRecords+j)
			if err != nil {
				return err
			}

			if line, ok := row.([]byte); ok {
				if _, err := fWriter.Write(line); err != nil {
					return err
				}
			} else {
				return TypeError("CSVAtCache2ZipAtFile")
			}
		}
	}

	return nil
}

// BinaryAtCache2ZipAtFile converts binary data stored in cache to a zip file.
func BinaryAtCache2ZipAtFile(outputBuilder *OutputBuilder, seqNum int) error {
	if outputBuilder.CompressPerSchema {
		for _, schBldr := range outputBuilder.SchBuilders {
			zipFilePath := filepath.Join(outputBuilder.Path, fmt.Sprintf("%s_%s_%d.zip", outputBuilder.Name, schBldr.SchemaName, seqNum))
			if err := BinaryAtCache2ZipAtFileBySchema(
				schBldr.SchemaName,
				seqNum,
				schBldr.NumFilesPerCompressedFile,
				schBldr.NumRecords,
				true,
				zipFilePath,
				nil,
			); err != nil {
				return err
			}
		}
		return nil
	} else {
		zipFilePath := filepath.Join(outputBuilder.Path, fmt.Sprintf("%s_%d.zip", outputBuilder.Name, seqNum))
		zipFile, err := os.OpenFile(zipFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer zipFile.Close()
		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		for _, schBldr := range outputBuilder.SchBuilders {
			if err := BinaryAtCache2ZipAtFileBySchema(
				schBldr.SchemaName,
				seqNum,
				schBldr.NumFilesPerCompressedFile,
				schBldr.NumRecords,
				false,
				"",
				zipWriter,
			); err != nil {
				return err
			}
		}

		return nil
	}
}

// BinaryAtCache2ZipAtFileBySchema converts binary data stored in cache to a zip file for a specific schema.
func BinaryAtCache2ZipAtFileBySchema(schemaName string, seqNum int, numFilesPerCompressedFile int, numRecords int,
	compressPerFile bool, zipFilePath string, zipWriter *zip.Writer) error {
	if compressPerFile {
		zipFile, err := os.OpenFile(zipFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer zipFile.Close()
		zipWriter = zip.NewWriter(zipFile)
		defer zipWriter.Close()
	}

	colNames := GetColumnNames(schemaName)

	for i := 0; i < numFilesPerCompressedFile; i++ {
		var rows []byte
		for j := 0; j < numRecords; j++ {
			var row []byte
			for _, key := range colNames {
				bCol, _ := GetFakeData(key, i*numRecords+j)
				if v, ok := bCol.([]byte); ok {
					row = append(row, v...)
				} else {
					return TypeError(key)
				}
			}
			rows = append(rows, row...)
		}

		fileName := fmt.Sprintf("%s_%d_%d.bin", schemaName, seqNum, i)
		fWriter, err := zipWriter.Create(fileName)
		if err != nil {
			return err
		}
		if _, err := fWriter.Write(rows); err != nil {
			return err
		}
	}

	return nil
}
