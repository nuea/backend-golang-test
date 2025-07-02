package util

func MapToSlice[I, O any](mapper func(I) (O, error), input []I) ([]O, error) {
	output := make([]O, 0)
	for _, i := range input {
		o, err := mapper(i)
		if err != nil {
			return nil, err
		}
		output = append(output, o)
	}
	return output, nil
}
