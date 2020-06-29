// Code generated by "stringer -type=RequestType"; DO NOT EDIT.

package backend

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ApiTypesRequest-0]
	_ = x[LoginChallengeCheck-1]
	_ = x[ConsentChallengeCheck-2]
	_ = x[LoginRequest-3]
	_ = x[CardSearchRequest-4]
	_ = x[CardDetailRequest-5]
}

const _RequestType_name = "ApiTypesRequestLoginChallengeCheckConsentChallengeCheckLoginRequestCardSearchRequestCardDetailRequest"

var _RequestType_index = [...]uint8{0, 15, 34, 55, 67, 84, 101}

func (i RequestType) String() string {
	if i < 0 || i >= RequestType(len(_RequestType_index)-1) {
		return "RequestType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RequestType_name[_RequestType_index[i]:_RequestType_index[i+1]]
}