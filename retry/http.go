package retry

func RetryHTTPStatus(code int) bool {
	switch code {
	case 408, 425, 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}
