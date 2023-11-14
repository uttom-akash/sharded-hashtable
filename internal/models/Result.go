package models

import "strconv"

type ResultStatus int

const (
	Found ResultStatus = iota + 1
	ContinueSearch
	StopSearch
	AddedOrUpdated
	Deleted
)

type Result struct {
	Status   ResultStatus
	Value    Value
	BucketId *int
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

func NewFoundResultWithBucket(value Value, bucketId int) *Result {
	return &Result{
		Status:   Found,
		Value:    value,
		BucketId: &bucketId,
	}
}

func NewContinueSearchResultWithBucket(bucketId int) *Result {
	return &Result{
		Status:   ContinueSearch,
		BucketId: &bucketId,
	}
}

func NewStopSearchResultWithBucket(bucketId int) *Result {
	return &Result{
		Status:   StopSearch,
		BucketId: &bucketId,
	}
}

func (result *Result) ToString() string {
	return resultStatusNames[ResultStatus(result.Status)] + " : " + strconv.Itoa(int(result.Value.Value))
}

var resultStatusNames = map[ResultStatus]string{
	Found:          "Found",
	ContinueSearch: "ContinueSearch",
	StopSearch:     "StopSearch",
	AddedOrUpdated: "AddedOrUpdated",
	Deleted:        "Deleted",
}
