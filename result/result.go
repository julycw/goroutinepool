package result

type Result struct {
}

func (this *Result) Init() *Result {
	return this
}

func New() *Result {
	return new(Result).Init()
}
