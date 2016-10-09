package librevault

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var luhnMod58Lookup = make([]int, 255)

func init() {
	for i := range luhnMod58Lookup {
		luhnMod58Lookup[i] = 255
	}
	for i, chr := range alphabet {
		luhnMod58Lookup[chr] = i
	}
}

func luhnModValue(data []byte) byte {
	factor := 2
	sum := 0
	var add, cp int
	for i := len(data) - 1; i > 0; i-- {
		cp = luhnMod58Lookup[data[i]]
		add = factor * cp
		if factor == 2 {
			factor = 1
		} else {
			factor = 2
		}
		add = (add / 58) + (add % 58)
		sum += add
	}

	rem := sum % 58
	check := (58 - rem) % 58
	return alphabet[check]
}

func appendChecksum(s string) string {
	return s + string(luhnModValue([]byte(s)))
}

func validateChecksum(s string) bool {
	front := s[:len(s)-1]
	return luhnModValue([]byte(front)) == s[len(s)-1]
}
