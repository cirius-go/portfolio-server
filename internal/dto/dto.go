package dto

type (
	// ListingReq is the request data of listing.
	ListingReq struct {
		Page    int    `json:"p" query:"p"`
		PerPage int    `json:"pp" query:"pp"`
		Sort    string `json:"s" query:"s"`
	}

	// ListingRes is the response data of listing.
	ListingRes[I any] struct {
		Recs  []*I  `json:"recs"`
		Total int64 `json:"total"`
	}
)

type (
	// ErrorRes represents the default error response.
	ErrorRes struct {
		Type     string `json:"type"`
		Message  string `json:"message"`
		Internal string `json:"internal"`
	}
)
