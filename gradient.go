package rse

import (
	"log"
	"math"
)

type MatrixTarget = int

const (
	Row MatrixTarget = iota
	Column
)

func gradient(knownMatrix *Matrix, target MatrixTarget, targetIndex int, rowMatrix *Matrix, rowIndex int, columnMatrix *Matrix, columnIndex int) (float64, error) {
	knownValue := knownMatrix.ValueAt(rowIndex, columnIndex)

	rowVector := rowMatrix.GetRow(rowIndex)
	columnVector := columnMatrix.GetColumn(columnIndex)

	predictedValue, err := DotProduct(rowVector, columnVector)
	if err != nil {
		return 0, err
	}
	log.Printf("gradient - known %v, prerdicted %v\n", knownValue, predictedValue)

	switch target {
	case Row:
		element := columnVector[targetIndex]
		return (2 * (knownValue - predictedValue) * element), nil
	case Column:
		element := rowVector[targetIndex]
		return (2 * (knownValue - predictedValue) * element), nil
	default:
		return 0, nil
	}
}

func rowFeatureGratient(isSparce bool, targetIndex int, rowIndex int, knownMAtrix, rowMatrix, columnMatrix *Matrix) (float64, error) {
	var sum float64 = 0
	var count float64 = 0
	for index := 0; index < len(columnMatrix.values[0]); index += 1 {
		if isSparce && math.IsNaN(knownMAtrix.ValueAt(rowIndex, index)) {
			continue
		}

		fGradient, gradientErr := gradient(knownMAtrix, Row, targetIndex, rowMatrix, rowIndex, columnMatrix, index)
		if gradientErr != nil {
			return 0, gradientErr
		}
		count += 1
		sum += fGradient
	}

	return sum / count, nil
}

func columnFeatureGratient(isSparce bool, targetIndex int, columnIndex int, knownMAtrix, rowMatrix, columnMatrix *Matrix) (float64, error) {
	var sum float64 = 0
	var count float64 = 0
	for index := 0; index < len(rowMatrix.values); index += 1 {
		if isSparce && math.IsNaN(knownMAtrix.ValueAt(index, columnIndex)) {
			continue
		}
		fGradient, gradientErr := gradient(knownMAtrix, Column, targetIndex, rowMatrix, index, columnMatrix, columnIndex)
		if gradientErr != nil {
			return 0, gradientErr
		}

		count += 1
		sum += fGradient
	}

	return sum / count, nil
}
