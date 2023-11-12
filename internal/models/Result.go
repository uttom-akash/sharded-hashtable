package models

type ResultStatus int

const (
	Found ResultStatus = iota + 1
	ContinueSearch
	StopSearch
	AddedOrUpdated
	Deleted
)

type Result struct {
	Status ResultStatus
	Value  Value
}

func NewFoundResult(value Value) *Result {
	return &Result{
		Status: Found,
		Value:  value,
	}
}

func NewContinueSearchResult() *Result {
	return &Result{
		Status: ContinueSearch,
	}
}

func NewStopSearchResult() *Result {
	return &Result{
		Status: StopSearch,
	}
}

func NewContinueOrStopSearchResult(continueSearch bool) *Result {

	if continueSearch {
		return NewContinueSearchResult()
	}

	return NewStopSearchResult()
}

func NewAddedOrUpdatedResult(value Value) *Result {
	return &Result{
		Status: AddedOrUpdated,
		Value:  value,
	}
}

func NewDeletedResult() *Result {
	return &Result{
		Status: Deleted,
	}
}
