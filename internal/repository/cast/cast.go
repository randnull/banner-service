package cast

import "github.com/lib/pq"


func ConvertIntToInterface(ptr *int) interface{} {
	if ptr != nil {
		return *ptr
	}
	return nil
}


func ConvertArrayToInterface(ptr *[]int) interface{} {
	if ptr != nil {
		return pq.Array(*ptr)
	}
	return nil
}


func ConvertStringToInterface(ptr *string) interface{} {
	if ptr != nil {
		return *ptr
	}
	return nil
}


func ConvertBoolToInterface(ptr *bool) interface{} {
	if ptr != nil {
		return *ptr
	}
	return nil
}


func CastPqArrayToInt(arr pq.Int64Array) []int {
	int_arr := make([]int, len(arr))

	for i, value := range arr {
		int_arr[i] = int(value)
	}

	return int_arr
}
