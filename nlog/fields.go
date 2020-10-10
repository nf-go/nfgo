package nlog

// Fields -
type Fields map[string]interface{}

// NewFields -
func NewFields(keyVals ...string) Fields {
	fields := Fields{}
	if len(keyVals) == 0 {
		return fields
	}
	keys := []string{}
	vals := []string{}
	for i, arg := range keyVals {
		if i%2 == 0 {
			keys = append(keys, arg)
		} else {
			vals = append(vals, arg)
		}
	}
	for i, val := range vals {
		if val != "" {
			fields[keys[i]] = val
		}
	}
	return fields
}
