package imports_analyzer

type StringSliceFlag []string

func (i *StringSliceFlag) String() string {
	return "what should this be?"
}

func (i *StringSliceFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}
