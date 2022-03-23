package rse

import (
	"testing"
)

func TestNewMatrixAndTransposeEmpty(t *testing.T) {
	// create matrix with 3 rows, 2 columns
	m := NewEmpty(3, 2)

	if len(m.values) != 3 {
		t.Fatal("matrix should have 3 rows")
	}

	if len(m.values[0]) != 2 {
		t.Fatal("matrix should have 2 columns")
	}

	mT := m.Transpose()
	// should now be 2 rows and 3 columns

	if len(mT.values) != 2 {
		t.Fatal("transposed matrix should ave 2 rows")
	}

	if len(mT.values[0]) != 3 {
		t.Fatal("transposed matrix should have 3 columns")
	}

}

func TestNewMatrixAndTransposeWithValues(t *testing.T) {
	m := NewEmpty(2, 3).WithValues([]ValueSet{
		ValueSet{1, 0, 0},
		ValueSet{4, 1, 0},
		ValueSet{2, 0, 1},
		ValueSet{5, 1, 1},
		ValueSet{3, 0, 2},
		ValueSet{6, 1, 2},
	}...)

	if m.values[0][0] != 1 {
		t.Fatal("value set not used to set matrix properly")
	}
	if m.values[1][1] != 5 {
		t.Fatal("value set not used to set matrix properly")
	}

	transposeValues := [][]float64{
		{1, 4},
		{2, 5},
		{3, 6},
	}

	mT := m.Transpose()

	for i, _ := range transposeValues {
		for j, value := range transposeValues[i] {
			if mT.values[i][j] != value {
				t.Fatal("matrix values not transposed correctly")
			}
		}
	}
}

func TestMatrixValueAt(t *testing.T) {
	m := NewEmpty(2, 3).WithValues([]ValueSet{
		ValueSet{1, 0, 0},
		ValueSet{4, 1, 0},
		ValueSet{2, 0, 1},
		ValueSet{5, 1, 1},
		ValueSet{3, 0, 2},
		ValueSet{6, 1, 2},
	}...)

	if m.ValueAt(0, 0) != 1 {
		t.Fatal("error getting value with ValueAt func")
	}

	if m.ValueAt(1, 2) != 6 {
		t.Fatal("error getting value with ValueAt func")
	}
}

func TestNewMatrixLabels(t *testing.T) {
	m := NewEmpty(1, 1).
		WithColumnLabels(LabelSet{"column", 0}).
		WithRowLabels(LabelSet{"row", 0})

	if m.labels.RowLabels[0] != "row" {
		t.Fatal("row label was not assigned properly")
	}

	if m.labels.ColumnLabels[0] != "column" {
		t.Fatal("column label was not assigned properly")
	}
}

func TestMatrixGettersRowColumn(t *testing.T) {
	m := NewEmpty(2, 3).WithValues([]ValueSet{
		ValueSet{1, 0, 0},
		ValueSet{4, 1, 0},
		ValueSet{2, 0, 1},
		ValueSet{5, 1, 1},
		ValueSet{3, 0, 2},
		ValueSet{6, 1, 2},
	}...)

	row := m.GetRow(0)

	if len(row) != 3 {
		t.Fatal("row should be 3 elements long")
	}

	if row[2] != 3 {
		t.Fatal("row has incorrect value")
	}

	column := m.GetColumn(2)

	if len(column) != 2 {
		t.Logf("%#v\n", column)
		t.Fatal("column should have 2 elements")
	}

	if column[1] != 6 {
		t.Fatal("column has incorrect value")
	}

}

func TestDotProduct(t *testing.T) {
	val, _ := DotProduct([]float64{1, 1}, []float64{1, 1})
	if val != 2 {
		t.Fatal("Did not do dot product correctly")
	}

	_, err := DotProduct([]float64{1}, []float64{1, 1})

	if err == nil {
		t.Fatal("should have errored with different size vectors")
	}
}

func TestMatrixMultiply(t *testing.T) {
	m1 := NewEmpty(2, 3).WithValues([]ValueSet{
		ValueSet{1, 0, 0},
		ValueSet{4, 1, 0},
		ValueSet{2, 0, 1},
		ValueSet{5, 1, 1},
		ValueSet{3, 0, 2},
		ValueSet{6, 1, 2},
	}...)

	m2 := NewEmpty(3, 2).WithValues([]ValueSet{
		ValueSet{7, 0, 0},
		ValueSet{8, 0, 1},
		ValueSet{9, 1, 0},
		ValueSet{10, 1, 1},
		ValueSet{11, 2, 0},
		ValueSet{12, 2, 1},
	}...)

	m3, err := Multiply(m1, m2)
	if err != nil {
		t.Log(err.Error())
		t.Fatal("no error should have occured during mutl")
	}

	if len(m3.values) != 2 {
		t.Fatal("new matrix should be row size 2")
	}

	if len(m3.values[0]) != 2 {
		t.Fatal("new matrix should be column size 2")
	}

	if m3.ValueAt(0, 0) != 58 {
		t.Fatal("incorrect mutliply value")
	}

	if m3.ValueAt(0, 1) != 64 {
		t.Fatal("incorrect mutliply value")
	}

	if m3.ValueAt(1, 0) != 139 {
		t.Fatal("incorrect mutliply value")
	}

	if m3.ValueAt(1, 1) != 154 {
		t.Fatal("incorrect mutliply value")
	}
}

func TestMSE(t *testing.T) {
	m1 := NewEmpty(2, 2).WithValues([]ValueSet{
		ValueSet{2, 0, 0},
	}...)
	m2 := NewEmpty(2, 2).WithValues([]ValueSet{
		ValueSet{1, 0, 0},
	}...)

	val, err := MSE(m1, m2)
	if err != nil {
		t.Log(err.Error())
		t.Fatal("matrices should have been same dimensions")
	}

	if val != 1 {
		t.Logf("%v\n", val)
		t.Fatal("incorrect mse value")
	}
}

func TestNewEmptyRandom(t *testing.T) {
	m := NewEmpty(2, 2).WithRandomValues(0, 10)

	t.Logf("m values: %#v\n", m.values)
}
