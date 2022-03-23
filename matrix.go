package rse

import (
	"errors"
	"log"
	"math"
	"math/rand"
)

type Labels struct {
	RowLabels    []string
	ColumnLabels []string
}

type Matrix struct {
	labels *Labels
	values [][]float64
}

func NewEmpty(i, j int) *Matrix {
	values := make([][]float64, i)
	for index, _ := range values {
		column := make([]float64, j)
		for columnIndex := range column {
			column[columnIndex] = math.NaN()
		}
		values[index] = column
	}

	return &Matrix{
		values: values,
		labels: nil,
	}
}

type ValueSet struct {
	Value  float64
	CoordI int64
	CoordJ int64
}

func (m *Matrix) WithValues(values ...ValueSet) *Matrix {
	for _, value := range values {
		m.values[value.CoordI][value.CoordJ] = value.Value
	}
	return m
}

type LabelSet struct {
	Value string
	Coord int64
}

func (m *Matrix) WithColumnLabels(labels ...LabelSet) *Matrix {
	if m.labels == nil {
		m.labels = &Labels{
			RowLabels:    make([]string, len(m.values)),
			ColumnLabels: make([]string, len(m.values[0])),
		}
	}

	for _, label := range labels {
		m.labels.ColumnLabels[label.Coord] = label.Value
	}

	return m
}

func (m *Matrix) WithRowLabels(labels ...LabelSet) *Matrix {
	if m.labels == nil {
		m.labels = &Labels{
			RowLabels:    make([]string, len(m.values)),
			ColumnLabels: make([]string, len(m.values[0])),
		}
	}

	for _, label := range labels {
		m.labels.RowLabels[label.Coord] = label.Value
	}

	return m
}

func (m *Matrix) ValueAt(i, j int) float64 {
	return m.values[i][j]
}

func (m *Matrix) Transpose() *Matrix {
	newValues := make([][]float64, len(m.values[0]))
	for index, _ := range newValues {
		newValues[index] = make([]float64, len(m.values))
	}

	for i, _ := range m.values {
		for j, value := range m.values[i] {
			newValues[j][i] = value
		}
	}

	var newLables *Labels = nil
	if m.labels != nil {
		newLables = &Labels{
			RowLabels:    m.labels.ColumnLabels,
			ColumnLabels: m.labels.RowLabels,
		}
	}

	return &Matrix{
		values: newValues,
		labels: newLables,
	}
}

func (m *Matrix) GetRow(index int) []float64 {
	return m.values[index]
}

func (m *Matrix) GetColumn(index int) []float64 {
	column := make([]float64, len(m.values))

	for i, row := range m.values {
		column[i] = row[index]
	}

	return column
}

func DotProduct(xV, yV []float64) (float64, error) {
	if len(xV) != len(yV) {
		return 0, errors.New("vectors are not the same length")
	}

	var dotProduct float64 = 0
	for index, xVal := range xV {
		dotProduct = dotProduct + (xVal * yV[index])
	}

	return dotProduct, nil
}

func Multiply(m1, m2 *Matrix) (*Matrix, error) {
	/*
		n-m * m-y results in n-y size matrix
	*/
	log.Printf("Multiplying Matrices:\n\t%#v\n\t%#v\n", m1.values, m2.values)
	if len(m1.values[0]) != len(m2.values) {
		return nil, errors.New("matrix 1 column and matrix 2 rows are not of same length")
	}

	newMatrixValues := make([][]float64, len(m1.values))
	for index, _ := range newMatrixValues {
		newMatrixValues[index] = make([]float64, len(m2.values[0]))
	}

	var dotMultErr error = nil
	for i, _ := range m1.values {
		for j, _ := range m2.values[0] {
			newMatrixValues[i][j], dotMultErr = DotProduct(m1.GetRow(i), m2.GetColumn(j))
			if dotMultErr != nil {
				return nil, dotMultErr
			}
		}
	}

	labels := &Labels{}
	if m1.labels != nil {
		labels.RowLabels = m1.labels.RowLabels
	}

	if m2.labels != nil {
		labels.ColumnLabels = m2.labels.ColumnLabels
	}

	return &Matrix{
		values: newMatrixValues,
		labels: labels,
	}, nil
}

func MSE(m1, m2 *Matrix) (float64, error) {
	log.Printf("Evaluating MSE of:\n\t%#v\n\t%#v", m1.values, m2.values)
	if len(m1.values) != len(m2.values) ||
		len(m1.values[0]) != len(m2.values[0]) {
		return 0, errors.New("matrices are not same dimensions")
	}
	// error with respect to m1
	var sumError float64 = 0
	for i, _ := range m1.values {
		for j, _ := range m1.values[i] {
			if !math.IsNaN(m1.ValueAt(i, j)) {
				// nan value is possible in sparse matrix and we dont want to eval againt unknown
				elementErr := math.Pow((m1.ValueAt(i, j) - m2.ValueAt(i, j)), 2)
				if !math.IsNaN(elementErr) {
					sumError += elementErr
				} else {
					log.Printf("Recieved NaN MSE at %v,%v\n", i, j)
					log.Println("values: ", m1.ValueAt(i, j), m2.ValueAt(i, j), (m1.ValueAt(i, j) - m2.ValueAt(i, j)), math.Pow((m1.ValueAt(i, j)-m2.ValueAt(i, j)), 2))
				}
			}
		}
	}
	return sumError, nil
}

func (m *Matrix) MaxMin() (float64, float64) {
	max := m.values[0][0]
	min := m.values[0][0]

	for _, row := range m.values {
		for _, column := range row {
			if column > max {
				max = column
			}

			if column < min {
				min = column
			}
		}
	}

	return max, min
}

func (m *Matrix) WithRandomValues(min, max float64) *Matrix {
	for i, _ := range m.values {
		for j, _ := range m.values[i] {
			m.values[i][j] = (rand.Float64() * (max - min)) + min
		}
	}

	return m
}
