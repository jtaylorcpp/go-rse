package loaders

type Parameters struct {
	Iterations   int
	Features     int
	LearningRate float64
}

func DefaultParameters() (iterations int, features int, learningRate float64) {
	return 1000, 2, 0.01
}
