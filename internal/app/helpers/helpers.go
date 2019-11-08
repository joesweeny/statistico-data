package helpers

import (
	sportmonks "github.com/statistico/statistico-sportmonks-go-client"
	"strconv"
)

func NullableUint64(i *int) *uint64 {
	var assist *uint64 = nil

	if i != nil {
		val := *i
		i := uint64(val)
		assist = &i
	}

	return assist
}

func ParseNullableInt(i interface{}) *int {
	if i == nil {
		return nil
	}

	if _, ok := i.(*int); ok {
		val := i.(*int)
		return val
	}

	if _, ok := i.(int); ok {
		val := i.(int)
		return &val
	}

	if x, ok := i.(float64); ok {
		val := int(x)
		return &val
	}
	if _, ok := i.(string); ok {
		val, err := strconv.Atoi(i.(string))

		if err != nil {
			panic(err)
		}

		return &val
	}

	return nil
}

func ParseFlexInt(i *sportmonks.FlexInt) *int {
	if i == nil {
		return nil
	}

	val := int(*i)

	return &val
}
