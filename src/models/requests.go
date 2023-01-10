// this holds request objects that are populated upon HTTP request
package models

// incoming request structure for loads
type Load struct {
	RFCIdentifier string `json:"rfcIdentifier" binding:"required"`
} // @name Load

// incoming request structure for merges
type Merge struct {
	RFCIdentifier string `json:"rfcIdentifier" binding:"required"`
} // @name Merge

// incoming request structure for reveiws
type Review struct {
	RFCIdentifier   string `json:"rfcIdentifier" binding:"required" example:"123456"`
	Type            string `json:"type" binding:"required" example:"COMMENT"`
	TopLevelComment string `json:"topLevelComment,omitempty" example:"This is my review comment!"`
	// this was not made into its own struct so that we can efficiently look up targets using the power of maps
	Comments       map[string][]string `json:"comments,omitempty" swaggertype:"object,array,string"`
	LoadOnApproval bool                `json:"loadOnApproval,omitempty" swaggerignore:"true"`
} // @name Review

// incoming request structure for load status requests
type Status struct {
	RFCIdentifier string `json:"rfcIdentifier" binding:"required" example:"123456"`
} // @name Status

// incoming request structure for updates
type Update struct {
	RFC           *RFC   `json:"rfc" binding:"required"`
	RFCIdentifier string `json:"rfcIdentifier" binding:"required"`
} // @name Update

// incoming request structure for getRfcs requests
type GetRfcs struct {
	Count int    `json:"count" example:"100" binding:"required"` //Number of requests wanted. If count is -1, return all requests. Required
	State string `json:"state" example:"open"`                   //State of the request, one of "open", "closed", or "all". Default: "all"

	// The following are options used to filter the returned PRs, the default value for all is to not filter
	Owner  *string `json:"owner" example:"tstark"` //Username of the owner of the requests.
	Merged *bool   `json:"merged" example:"false"` //Merged status of the RFC. A closed RFC that has Merged:false indicates that the change was rejected.
} // @name GetRfcs

// incoming request structure for getRfcContents requests
type GetRfcContents struct {
	RFCIdentifier string `json:"rfcIdentifier" binding:"required" example:"123456"`
} // @name GetRfcContents
