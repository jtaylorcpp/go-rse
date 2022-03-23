package rse

import (
	"fmt"
	"log"
	"math"
)

type MatrixFactorization struct {
	GivenMatrix  *Matrix
	RowMatrix    *Matrix
	ColumnMatrix *Matrix
	Features     int
	LearningRate float64
	Iterations   int
	isSparce     bool
}

func NewMatrixFactorization(initialMatrix *Matrix, features int, learningRate float64, iterations int) *MatrixFactorization {
	max, min := initialMatrix.MaxMin()

	// assuming given matrix is MxN and desired number of latent features F
	// creates 2 matrices
	// rowMatrix of NxF
	// columnMatrix of FxM
	// create row matrix
	rowMatrix := NewEmpty(len(initialMatrix.values), features).WithRandomValues(min, max)
	// create column matrix
	columnMatrix := NewEmpty(features, len(initialMatrix.values[0])).WithRandomValues(min, max)

	if initialMatrix.labels != nil {
		featureLables := make([]string, features)
		for index, _ := range featureLables {
			featureLables[index] = fmt.Sprintf("Feature-%v", index)
		}

		rowMatrix.labels = &Labels{
			RowLabels:    initialMatrix.labels.RowLabels,
			ColumnLabels: featureLables,
		}

		columnMatrix.labels = &Labels{
			RowLabels:    featureLables,
			ColumnLabels: initialMatrix.labels.ColumnLabels,
		}
	}

	isSparce := false
	for i := range initialMatrix.values {
		for j := range initialMatrix.values[i] {
			if math.IsNaN(initialMatrix.ValueAt(i, j)) {
				isSparce = true
				goto foundSparceness
			}
		}
	}
foundSparceness:
	return &MatrixFactorization{
		GivenMatrix:  initialMatrix,
		RowMatrix:    rowMatrix,
		ColumnMatrix: columnMatrix,
		Features:     features,
		LearningRate: learningRate,
		Iterations:   iterations,
		isSparce:     isSparce,
	}
}

func NewMatrixFactorizationFromLoader(loader MatrixLoader) (*MatrixFactorization, error) {
	initialMatrix, loadErr := loader.Load()
	if loadErr != nil {
		return nil, loadErr
	}

	iterations, features, learningRate := loader.GetParameters()

	return NewMatrixFactorization(initialMatrix, features, learningRate, iterations), nil
}

func (mf *MatrixFactorization) updateRowFeatures() (*MatrixFactorization, error) {
	for rowIndex := range mf.RowMatrix.values {
		for featureIndex := range mf.RowMatrix.values[rowIndex] {
			grad, gradErr := rowFeatureGratient(
				mf.isSparce,
				featureIndex,
				rowIndex,
				mf.GivenMatrix,
				mf.RowMatrix,
				mf.ColumnMatrix,
			)
			if gradErr != nil {
				return mf, gradErr
			}

			mf.RowMatrix.values[rowIndex][featureIndex] += (mf.LearningRate * grad)
		}
	}

	return mf, nil
}

func (mf *MatrixFactorization) updateColumnFeatures() (*MatrixFactorization, error) {
	for featureIndex := range mf.ColumnMatrix.values {
		for columnIndex := range mf.ColumnMatrix.values[featureIndex] {
			grad, gradErr := columnFeatureGratient(
				mf.isSparce,
				featureIndex,
				columnIndex,
				mf.GivenMatrix,
				mf.RowMatrix,
				mf.ColumnMatrix,
			)

			if gradErr != nil {
				return mf, gradErr
			}

			mf.ColumnMatrix.values[featureIndex][columnIndex] += (mf.LearningRate * grad)
		}
	}

	return mf, nil
}

func (mf *MatrixFactorization) Train() error {
	var trainErr error = nil
	for iter := 0; iter < mf.Iterations; iter += 1 {
		mf, trainErr = mf.updateRowFeatures()
		if trainErr != nil {
			return trainErr
		}
		log.Printf("Iteration  %v - updating row %#v\n", iter, mf.RowMatrix.values)

		mf, trainErr = mf.updateColumnFeatures()
		if trainErr != nil {
			return trainErr
		}
		log.Printf("Iteration  %v - updating column %#v\n", iter, mf.ColumnMatrix.values)

		if iter%1 == 0 {
			newMatrix, multErr := Multiply(mf.RowMatrix, mf.ColumnMatrix)
			if multErr != nil {
				return multErr
			}

			mse, mseErr := MSE(mf.GivenMatrix, newMatrix)
			if mseErr != nil {
				log.Printf("Iteration %v - MSE Error: %s\n", iter, mseErr.Error())
			}
			log.Printf("Iteration %v - MSE: %v\n", iter, mse)
		}

	}

	return nil
}

func (mf *MatrixFactorization) MSE() (float64, error) {
	newMatrix, multErr := Multiply(mf.RowMatrix, mf.ColumnMatrix)
	if multErr != nil {
		return 0, multErr
	}

	mse, mseErr := MSE(mf.GivenMatrix, newMatrix)
	if mseErr != nil {
		log.Printf("MSE Error: %s\n", mseErr.Error())
		return 0, mseErr
	}

	return mse, nil
}

type MatrixLoader interface {
	GetParameters() (iterations int, features int, learningRate float64)
	Load() (*Matrix, error)
}
