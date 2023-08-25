// Code generated by "enumerator -type LpaType -linecomment -trimprefix -empty"; DO NOT EDIT.
package page

import (
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LpaTypeHealthWelfare-1]
	_ = x[LpaTypePropertyFinance-2]
}

const _LpaType_name = "hwpfa"

var _LpaType_index = [...]uint8{0, 2, 5}

func (i LpaType) String() string {
	i -= 1
	if i >= LpaType(len(_LpaType_index)-1) {
		return "LpaType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _LpaType_name[_LpaType_index[i]:_LpaType_index[i+1]]
}

func (i LpaType) IsHealthWelfare() bool {
	return i == LpaTypeHealthWelfare
}

func (i LpaType) IsPropertyFinance() bool {
	return i == LpaTypePropertyFinance
}

func ParseLpaType(s string) (LpaType, error) {
	switch s {
	case "hw":
		return LpaTypeHealthWelfare, nil
	case "pfa":
		return LpaTypePropertyFinance, nil
	default:
		return LpaType(0), fmt.Errorf("invalid LpaType '%s'", s)
	}
}

type LpaTypeOptions struct {
	HealthWelfare   LpaType
	PropertyFinance LpaType
}

var LpaTypeValues = LpaTypeOptions{
	HealthWelfare:   LpaTypeHealthWelfare,
	PropertyFinance: LpaTypePropertyFinance,
}

func (i LpaType) Empty() bool {
	return i == LpaType(0)
}