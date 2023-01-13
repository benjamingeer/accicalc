package dataset

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

func WriteCsv(objs []any, path *string) error {
	if len(objs) == 0 {
		return nil
	}

	var file *os.File
	var err error

	if path != nil {
		file, err = os.Create(*path)

		if err != nil {
			return err
		}

		defer file.Close()
	} else {
		file = os.Stdout
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(toCsvHeader(objs[0])); err != nil {
		return err
	}

	for _, obj := range objs {
		if err := writer.Write(toCsvRow(obj)); err != nil {
			return err
		}
	}

	return nil
}

func toCsvHeader(obj any) []string {
	var header []string
	value := reflect.ValueOf(obj)
	objType := value.Type()

	for index := 0; index < value.NumField(); index++ {
		header = append(header, camelCaseToHeading(objType.Field(index).Name))
	}

	return header
}

func toCsvRow(obj any) []string {
	var row []string
	value := reflect.ValueOf(obj)

	for index := 0; index < value.NumField(); index++ {
		row = append(row, fmt.Sprint(value.Field(index).Interface()))
	}

	return row
}

func camelCaseToHeading(str string) string {
	if len(str) == 1 {
		return strings.ToUpper(str)
	}

	re := regexp.MustCompile(`([A-Z])`)
	str = re.ReplaceAllString(str, ` $1`)
	str = strings.Trim(str, " ")
	lowerCaseStr := strings.ToLower(str)
	runes := []rune(lowerCaseStr)
	return string(append([]rune{unicode.ToUpper(runes[0])}, runes[1:]...))
}
