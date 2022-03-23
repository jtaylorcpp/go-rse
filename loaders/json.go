package loaders

import (
	"encoding/json"
	"rse"
)

type AggregationFunc = int

const (
	SUM AggregationFunc = iota
	NORMALIZE
)

type JSONLoader struct {
	raw    json.RawMessage
	parser []JSONLoaderParser
}

type JSONLoaderParser struct {
	Criteria struct {
		Key    string
		Values []string
	}
	RowIdentifier    string
	ColumnIdentifier string
}

func NewJSONLoader(json json.RawMessage, parser ...JSONLoaderParser) JSONLoader {
	return JSONLoader{json, parser}
}

func (j JSONLoader) GetParameters() (iterations int, features int, learningRate float64) {
	return DefaultParameters()
}

func (j JSONLoader) Load() (*rse.Matrix, error) {
	// first scan all records find row x column labels
	return nil, nil
}
