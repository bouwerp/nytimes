package nytimes

type QueryMustBeProvided struct {
}

func (e QueryMustBeProvided) Error() string {
	return "a query must be provided"
}

type IncorrectTagType struct {
	Type TagType
}

func (e IncorrectTagType) Error() string {
	if e.Type == "" {
		return "incorrect tag type"
	} else {
		return "incorrect tag type - must be '" + string(e.Type) + "'"
	}
}
