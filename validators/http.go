package validators

func ValidateResponseStatusCode(code int) bool {
	if code >= 200 && code <= 300 {
		return true
	}

	return false

}
