// this holds response objects that are populated upon HTTP response
package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// holds health message
type Healthy struct {
	Message string `json:"message" example:"healthy"`
} // @name Healthy

// holds errors
type Error struct {
	Error string `json:"error" example:"whoops!"`
} // @name Error

// holds RFC unique identifier
type RFCIdentifier struct {
	RFCIdentifier string `json:"rfcIdentifier" example:"woo-hoo123"`
} //@name RFCIdentifier

// holds a success message
type Success struct {
	Success string `json:"success" example:"Success!"`
} //@name Success

// holds a load request response message
type LoadRequest struct {
	Message string `json:"message" example:"submitted load request for 12345, check status via the /status endpoint!"`
} //@name LoadRequest

// holds a status response message
type StatusResponse struct {
	Status string `json:"status" example:"loading"`
} //@name Status

type RFCs struct {
	RFCs  []map[string]string `json:"rfcs" swaggertype:"object,string" example:"1234:Example RFC title"`
	Count *int                `json:"count,omitempty" example:"10"`
}

type RFCContents struct {
	Body string `json:"body" binding:"required"`
}

// Implement Marshaler interface to make the output more compact while retaining meaning of an ordered set of key
// value pairs
func (r *RFCs) MarshalJSON() ([]byte, error) {
	rfcs := r.RFCs
	var marshaled []byte

	marshaled = append(marshaled, []byte(`{"rfcs": {`)...) // key and open brace
	for i, rfc := range rfcs {
		entryJson, err := json.Marshal(rfc)
		if err != nil {
			return nil, err
		}
		marshaled = append(marshaled, bytes.Trim(entryJson, "{}")...) // trim off leading { and trailing }
		if i < len(rfcs)-1 {
			marshaled = append(marshaled, []byte(`,`)...) // comma between entries
		}
	}
	marshaled = append(marshaled, []byte(`}`)...)
	if r.Count != nil {
		c := strconv.Itoa(*r.Count)
		marshaled = append(marshaled, []byte(fmt.Sprintf(`, "count": %v`, c))...) // add count if it exists
	}
	marshaled = append(marshaled, []byte(`}`)...) // close braces
	return marshaled, nil
}
