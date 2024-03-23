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

	"github.com/brianvoe/gofakeit/v6"
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
		filePath := filepath.Join(schBldr.ParentPath, fmt.Sprintf("%s_%s_%d.csv", outputBuilder.Name, schBldr.SchemaName, seqNum))
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
	outFile.Close()
	return nil
}

// Raw2CSVAtCache converts raw data to CSV format and stores it in cache.
func Raw2CSVAtCache(outputBuilder *OutputBuilder, seqNum int) error {
	for _, schBldr := range outputBuilder.SchBuilders {
		err := Raw2CSVAtCacheBySchema(schBldr.NumRecords, schBldr.SchemaName)
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
		filePath := filepath.Join(schBldr.ParentPath, fmt.Sprintf("%s_%s_%d", outputBuilder.Name, schBldr.SchemaName, seqNum))
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
	outFile, err := os.OpenFile(fmt.Sprintf("%s.bin", filePath), os.O_CREATE|os.O_WRONLY, 0644)
	for i := 0; i < numRecords; i++ {

		if err != nil {
			return err
		}
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
		outFile.Close()
	}
	return nil
}

// Raw2BinaryAtCache converts raw data to binary format and stores it in cache.
func Raw2BinaryAtCache(outputBuilder *OutputBuilder, seqNum int) error {
	for _, schBldr := range outputBuilder.SchBuilders {
		err := Raw2BinaryAtCacheBySchema(schBldr.NumRecords, schBldr.SchemaName)
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

			schemaName := schBldr.SchemaName
			numberOfFilesPerCompressedFile := schBldr.NumberOfFilesPerCompressedFile

			zipFilePath := filepath.Join(schBldr.ParentPath, fmt.Sprintf("%s_%s_%d.zip", outputBuilder.Name, schBldr.SchemaName, seqNum))

			numRecords := schBldr.NumRecords

			err := CSVAtCache2ZipAtFileBySchema(schemaName, numberOfFilesPerCompressedFile, zipFilePath, seqNum, numRecords, nil, true)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		var zipFilePath string
		var zipWriter *zip.Writer
		var zipFile *os.File
		var err error
		zipFilePath = filepath.Join(outputBuilder.Path, fmt.Sprintf("%s_%d.zip", outputBuilder.Name, seqNum))
		zipFile, err = os.OpenFile(zipFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		zipWriter = zip.NewWriter(zipFile)
		if err != nil {
			return err
		}

		for _, schBldr := range outputBuilder.SchBuilders {

			schemaName := schBldr.SchemaName
			numberOfFilesPerCompressedFile := schBldr.NumberOfFilesPerCompressedFile
			numRecords := schBldr.NumRecords
			err := CSVAtCache2ZipAtFileBySchema(schemaName, numberOfFilesPerCompressedFile, zipFilePath, seqNum, numRecords, zipWriter, false)
			if err != nil {
				return err
			}
		}

		zipWriter.Close()
		zipFile.Close()

		return nil
	}
}

// CSVAtCache2ZipAtFileBySchema converts CSV data stored in cache to a zip file for a specific schema.
func CSVAtCache2ZipAtFileBySchema(schemaName string, numberOfFilesPerCompressedFile map[string]int, zipFilePath string, seqNum int, numRecords int, zipWriter *zip.Writer, compressPerFile bool) error {
	if compressPerFile {
		zipFile, err := os.OpenFile(zipFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		zipWriter = zip.NewWriter(zipFile)
		if err != nil {
			return err
		}
		defer zipWriter.Close()
		defer zipFile.Close()
	}

	minFilesInZip := numberOfFilesPerCompressedFile["min"]

	maxFilesInZip := numberOfFilesPerCompressedFile["max"]

	numFilesInZip := gofakeit.Number(minFilesInZip, maxFilesInZip)

	for i := 0; i < numFilesInZip; i++ {

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

		var rows []byte
		for j := 0; j < numRecords; j++ {
			row, err := GetFakeData(schemaName, i)
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

		if _, err := fWriter.Write(rows); err != nil {
			return err
		}

	}

	return nil
}

// BinaryAtCache2ZipAtFile converts binary data stored in cache to a zip file.
func BinaryAtCache2ZipAtFile(outputBuilder *OutputBuilder, seqNum int) error {
	if outputBuilder.CompressPerSchema {
		for _, schBldr := range outputBuilder.SchBuilders {
			zipFilePath := filepath.Join(schBldr.ParentPath, fmt.Sprintf("%s_%s_%d.zip", outputBuilder.Name, schBldr.SchemaName, seqNum))

			schemaName := schBldr.SchemaName
			numberOfFilesPerCompressedFile := schBldr.NumberOfFilesPerCompressedFile
			numRecords := schBldr.NumRecords

			err := BinaryAtCache2ZipAtFileBySchema(schemaName, numberOfFilesPerCompressedFile, zipFilePath, seqNum, numRecords, nil, true)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		var zipFilePath string
		var zipWriter *zip.Writer
		var zipFile *os.File
		var err error
		zipFilePath = filepath.Join(outputBuilder.Path, fmt.Sprintf("%s_%d.zip", outputBuilder.Name, seqNum))
		zipFile, err = os.OpenFile(zipFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		zipWriter = zip.NewWriter(zipFile)
		if err != nil {
			return err
		}

		for _, schBldr := range outputBuilder.SchBuilders {

			schemaName := schBldr.SchemaName
			numberOfFilesPerCompressedFile := schBldr.NumberOfFilesPerCompressedFile
			numRecords := schBldr.NumRecords

			err := BinaryAtCache2ZipAtFileBySchema(schemaName, numberOfFilesPerCompressedFile, zipFilePath, seqNum, numRecords, zipWriter, false)
			if err != nil {
				return err
			}
		}

		zipWriter.Close()
		zipFile.Close()

		return nil
	}
}

// BinaryAtCache2ZipAtFileBySchema converts binary data stored in cache to a zip file for a specific schema.
func BinaryAtCache2ZipAtFileBySchema(schemaName string, numberOfFilesPerCompressedFile map[string]int, zipFilePath string, seqNum int, numRecords int, zipWriter *zip.Writer, compressPerFile bool) error {
	var zipFile *os.File
	if compressPerFile {
		zipFile, err := os.OpenFile(zipFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		zipWriter = zip.NewWriter(zipFile)
		if err != nil {
			return err
		}
	}

	colNames := GetColumnNames(schemaName)

	minFilesInZip := numberOfFilesPerCompressedFile["min"]

	maxFilesInZip := numberOfFilesPerCompressedFile["max"]

	numFilesInZip := gofakeit.Number(minFilesInZip, maxFilesInZip)
	rawFileName := fmt.Sprintf("%s_%d_", schemaName, seqNum)
	for i := 0; i < numFilesInZip; i++ {
		var rows []byte
		for j := 0; j < numRecords; j++ {
			var row []byte
			for _, key := range colNames {
				bCol, _ := GetFakeData(key, j)
				if v, ok := bCol.([]byte); ok {
					row = append(row, v...)
				} else {
					return TypeError(key)
				}
			}
			rows = append(rows, row...)
		}

		fileName := fmt.Sprintf("%s%d.bin", rawFileName, i)
		fWriter, err := zipWriter.Create(fileName)
		if err != nil {
			return err
		}
		if _, err := fWriter.Write(rows); err != nil {
			return err
		}
	}

	if compressPerFile {
		zipWriter.Close()
		zipFile.Close()
	}

	return nil
}

// Raw2CSVAtBytesBySchema converts raw data to CSV format for a specific schema and returns it as bytes.
func Raw2CSVAtBytesBySchema(numRecords int, schemaName string) ([]byte, error) {
	colNames := GetColumnNames(schemaName)

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	trimColNames := make([]string, len(colNames))
	prefixPattern := schemaName + "."
	for i, colName := range colNames {
		trimColNames[i] = strings.TrimPrefix(colName, prefixPattern)
	}
	if err := w.Write(trimColNames); err != nil {
		return nil, err
	}

	line := make([]string, len(colNames))

	for i := 0; i < numRecords; i++ {
		for j, key := range colNames {
			if fakeData, err := GetFakeData(key, i); err == nil {
				line[j] = fmt.Sprint(fakeData)
			} else {
				return nil, err
			}
		}
		if err := w.Write(line); err != nil {
			return nil, err
		}
	}

	w.Flush()

	return buf.Bytes(), nil
}

// Raw2BinaryAtBytesBySchema converts raw data to binary format for a specific schema and returns it as bytes.
func Raw2BinaryAtBytesBySchema(numRecords int, schemaName string) ([]byte, error) {
	colNames := GetColumnNames(schemaName)

	var buf bytes.Buffer
	bColLenBuf := make([]byte, 4)

	for i := 0; i < numRecords; i++ {
		for _, key := range colNames {
			if fakeData, err := GetFakeData(key, i); err == nil {
				bCol := []byte(fmt.Sprint(fakeData))
				binary.BigEndian.PutUint32(bColLenBuf, uint32(len(bCol)))
				buf.Write(bColLenBuf)
				buf.Write(bCol)
			} else {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

// CSVAtCache2ZipAtBytes converts CSV data stored in cache to a zip file and returns it as bytes.
func CSVAtCache2ZipAtBytes(schemaNames []string, numberOfFilesPerCompressedFile map[string]map[string]int, seqNum int, numRecordsMap map[string]int) ([]byte, error) {
	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	for _, schemaName := range schemaNames {

		minFilesInZip := numberOfFilesPerCompressedFile[schemaName]["min"]
		maxFilesInZip := numberOfFilesPerCompressedFile[schemaName]["max"]
		numFilesInZip := gofakeit.Number(minFilesInZip, maxFilesInZip)

		for i := 0; i < numFilesInZip; i++ {
			fileBuffer := new(bytes.Buffer)
			fileName := fmt.Sprintf("%s_%d_%d.csv", schemaName, seqNum, i)

			fWriter, err := zipWriter.Create(fileName)
			if err != nil {
				return nil, err
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
				return nil, err
			}
			w.Flush()

			if _, err := fWriter.Write(headerBuff.Bytes()); err != nil {
				return nil, err
			}

			for j := 0; j < numRecordsMap[schemaName]; j++ {
				row, err := GetFakeData(schemaName, i)
				if err != nil {
					return nil, err
				}

				if line, ok := row.([]byte); ok {
					if _, err := fWriter.Write(line); err != nil {
						return nil, err
					}
				} else {
					return nil, TypeError("CSVAtCache2ZipAtBytesBySchema")
				}
			}

			fileWriter, err := zipWriter.Create(fileName)
			if err != nil {
				return nil, err
			}

			_, err = fileWriter.Write(fileBuffer.Bytes())
			if err != nil {
				return nil, err
			}
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return zipBuffer.Bytes(), nil
}

// BinaryAtCache2ZipAtBytes converts binary data stored in cache to a zip file and returns it as bytes.
func BinaryAtCache2ZipAtBytes(schemaNames []string, numberOfFilesPerCompressedFileMap map[string]map[string]int, seqNum int, numRecordsMap map[string]int) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, schemaName := range schemaNames {

		colNames := GetColumnNames(schemaName)

		minFilesInZip := numberOfFilesPerCompressedFileMap[schemaName]["min"]
		maxFilesInZip := numberOfFilesPerCompressedFileMap[schemaName]["max"]
		numFilesInZip := gofakeit.Number(minFilesInZip, maxFilesInZip)

		rawFileName := fmt.Sprintf("%s_%d_", schemaName, seqNum)

		for i := 0; i < numFilesInZip; i++ {
			var rows []byte
			for j := 0; j < numRecordsMap[schemaName]; j++ {
				var row []byte
				for _, key := range colNames {
					bCol, _ := GetFakeData(key, j)
					if v, ok := bCol.([]byte); ok {
						row = append(row, v...)
					} else {
						return nil, TypeError(key)
					}
				}
				rows = append(rows, row...)
			}

			fileName := fmt.Sprintf("%s%d.bin", rawFileName, i)
			fWriter, err := zipWriter.Create(fileName)
			if err != nil {
				return nil, err
			}

			if _, err := fWriter.Write(rows); err != nil {
				return nil, err
			}
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
