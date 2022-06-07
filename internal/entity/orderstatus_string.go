// Code generated by "stringer -type=OrderStatus -trimprefix OrderStatus"; DO NOT EDIT.

package entity

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OrderStatusUNKNOWN-0]
	_ = x[OrderStatusNEW-1]
	_ = x[OrderStatusREGISTERED-2]
	_ = x[OrderStatusINVALID-3]
	_ = x[OrderStatusPROCESSING-4]
	_ = x[OrderStatusPROCESSED-5]
}

const _OrderStatus_name = "UNKNOWNNEWREGISTEREDINVALIDPROCESSINGPROCESSED"

var _OrderStatus_index = [...]uint8{0, 7, 10, 20, 27, 37, 46}

func (i OrderStatus) String() string {
	if i >= OrderStatus(len(_OrderStatus_index)-1) {
		return "OrderStatus(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OrderStatus_name[_OrderStatus_index[i]:_OrderStatus_index[i+1]]
}
