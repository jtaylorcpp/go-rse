package rse

import "testing"

func TestTrain(t *testing.T) {
	m := NewEmpty(3, 3).WithValues([]ValueSet{
		ValueSet{5, 0, 0},
		ValueSet{3, 0, 1},
		ValueSet{1, 0, 2},
		ValueSet{1, 1, 0},
		ValueSet{3, 1, 1},
		ValueSet{5, 1, 2},
		ValueSet{3, 2, 0},
		ValueSet{5, 2, 1},
		ValueSet{1, 2, 2},
	}...)

	mf := NewMatrixFactorization(m, 4, 0.01, 500)

	mf.Train()

	if mseVal, _ := mf.MSE(); mseVal > 1 {
		t.Fatal("Too much error")
	}
}

func TestTrainSparse(t *testing.T) {
	m := NewEmpty(3, 3).WithValues([]ValueSet{
		ValueSet{5, 0, 0},
		ValueSet{1, 1, 0},
		ValueSet{3, 1, 1},
		ValueSet{5, 1, 2},
		ValueSet{5, 2, 1},
		ValueSet{1, 2, 2},
	}...)

	mf := NewMatrixFactorization(m, 4, 0.01, 50)

	mf.Train()

	if mseVal, _ := mf.MSE(); mseVal > 1 {
		t.Fatal("Too much error")
	}
}
