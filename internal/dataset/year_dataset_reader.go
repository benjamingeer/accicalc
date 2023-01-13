package dataset

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

type YearDatasetReader interface {
	ReadCharacteristics(year uint, dataPath string) (accidents []*Accident, err error)
	ReadPlaces(year uint, dataPath string) (places []*Lieu, err error)
	ReadVehicles(year uint, dataPath string) (vehicles []*Véhicule, err error)
	ReadUsers(year uint, dataPath string) (users []*Usager, err error)
}

var (
	YearDatasetReaders map[uint]YearDatasetReader
	Years              []uint
	FirstYear          uint
	LastYear           uint
)

func init() {
	YearDatasetReaders = make(map[uint]YearDatasetReader)
	yearDatasetReader1 := YearDatasetReader1{}
	yearDatasetReader2 := YearDatasetReader2{}

	for year := uint(2005); year <= 2018; year++ {
		YearDatasetReaders[year] = &yearDatasetReader1
	}

	for year := uint(2019); year <= 2021; year++ {
		YearDatasetReaders[year] = &yearDatasetReader2
	}

	Years = maps.Keys(YearDatasetReaders)

	sort.Slice(Years, func(i, j int) bool {
		return Years[i] < Years[j]
	})

	FirstYear = Years[0]
	LastYear = Years[len(Years)-1]
}

func readCsvFile[T interface{}](path string, delimiter rune, convertRow func(row map[string]string) (*T, error)) ([]*T, error) {
	var items []*T
	var header []string

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	readHeader := true

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if readHeader {
			lowerCaseRow := make([]string, len(row))

			for index, columnName := range row {
				lowerCaseRow[index] = strings.ToLower(columnName)
			}

			header = append(header, lowerCaseRow...)
			readHeader = false
			continue
		}

		rowMap := make(map[string]string)

		for columnIndex, columnValue := range row {
			rowMap[header[columnIndex]] = strings.TrimSpace(columnValue)
		}

		item, err := convertRow(rowMap)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func readColumn(row map[string]string, columnName string, path string) (string, error) {
	if maybeValue, ok := row[strings.ToLower(columnName)]; ok {
		return maybeValue, nil
	} else {
		return "", fmt.Errorf("column '%v' missing in %v", columnName, path)
	}
}
