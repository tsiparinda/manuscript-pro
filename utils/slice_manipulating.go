package utils


func RemoveSliceElementString(s []string, i int) []string {
	ret := make ([]string, 0)
	ret = append(ret, s[:i]...)
	return append(ret, s[i+1:]...)
}
