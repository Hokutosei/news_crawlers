package utils

// ToUtf8 convert non utf8 to string
func ToUtf8(str string) string {
	nonUtf := strToByteArr(str)
	buf := make([]rune, len(nonUtf))
	for i, b := range nonUtf {
		buf[i] = rune(b)
	}
	return string(buf)
}

func strToByteArr(str string) []byte {
	return []byte(str)
}
