package helpers

func NullableUint64(i *int) *uint64 {
	var assist *uint64 = nil

	if i != nil {
		val := *i
		i := uint64(val)
		assist = &i
	}

	return assist
}
