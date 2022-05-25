package util

func CreateIdentifier(s ...string) string {
	identifier := ""
	for i := 0; i < len(s); i++ {
		identifier += s[i] + "_"
	}
	return identifier[:len(identifier)-1]
}
