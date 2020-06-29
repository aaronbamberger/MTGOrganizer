// Code generated by "stringer -type=ResponseType"; DO NOT EDIT.

package backend

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ApiTypesResponse-0]
	_ = x[ErrorResponse-1]
	_ = x[LoginChallengeResponse-2]
	_ = x[LoginResponse-3]
	_ = x[ConsentResponse-4]
	_ = x[CardSearchResponse-5]
	_ = x[CardDetailResponse-6]
}

const _ResponseType_name = "ApiTypesResponseErrorResponseLoginChallengeResponseLoginResponseConsentResponseCardSearchResponseCardDetailResponse"

var _ResponseType_index = [...]uint8{0, 16, 29, 51, 64, 79, 97, 115}

func (i ResponseType) String() string {
	if i < 0 || i >= ResponseType(len(_ResponseType_index)-1) {
		return "ResponseType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ResponseType_name[_ResponseType_index[i]:_ResponseType_index[i+1]]
}