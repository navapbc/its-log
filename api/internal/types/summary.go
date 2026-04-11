package types

type SummaryPostBody struct {
	Date      string  `json:"date" binding:"required"`
	Operation string  `json:"operation" binding:"required"`
	Tags      string  `json:"tags"`
	Value     string  `json:"value"`
	Count     float64 `json:"count" binding:"required"`
}
