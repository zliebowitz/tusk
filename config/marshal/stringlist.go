package marshal

// StringList is a list of strings optionally represented in yaml as a string.
// A single string in yaml will be unmarshalled as the first entry in a list,
// so the internal representation is always a list.
type StringList []string

// UnmarshalYAML unmarshals a string or list of strings always into a list.
func (sl *StringList) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var err error

	var single string
	if err = unmarshal(&single); err == nil {
		*sl = []string{single}
		return nil
	}

	var list []string
	if err = unmarshal(&list); err == nil {
		*sl = list
		return nil
	}

	return err
}
