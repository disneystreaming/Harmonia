// package models
// holds data definitions
package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// RFCIdentifierCreator is a function type that returns a custom RFC identifier string, for example, a branch name
type RFCIdentifierCreator func() *string

// RFC contains a set of actions that, in total, represent a proposal for change
type RFC struct {
	Actions    Actions `json:"actions" binding:"required"`
	Signature  string  `json:"signature,omitempty" swaggerignore:"true"`
	Identifier string  `json:"identifier,omitempty" swaggerignore:"true"`
} // @name RFC

// Actions is a slice of *Action types used to hold all RFC actions
type Actions []*Action

// ActionType represents a specific action
type ActionType string

// Comment represents comment intent
var CommentAction ActionType = "comment"
var LoadAction ActionType = "load"
var AddAction ActionType = "add"

// DataKey represents an attribute key within the Action Data object.
type DataKey string

var CommentData DataKey = "comment"
var CommenterData DataKey = "commenter"
var NoteData DataKey = "note"
var LoadStatus DataKey = "status"
var LoadRequester DataKey = "requester"
var ReviewerData DataKey = "reviewer"

// Action is a struct that represents a single schema action
type Action struct {
	ActionType ActionType             `json:"actionType" example:"add" binding:"required"`
	Target     Target                 `json:"target" swaggertype:"object,string" example:"targetType:item,targetDescriptor:EntityType" binding:"required"`
	Signature  string                 `json:"signature,omitempty" swaggerignore:"true"`
	Data       map[string]interface{} `json:"data,omitempty" swaggertype:"object,string" example:"id:MyData"`
} // @name Action

// TargetType represents the type of entity being targeted (item, action, rfc...)
type TargetType string //@name TargetType
var ActionTarget TargetType = "action"
var RfcTarget TargetType = "rfc"
var ItemTarget TargetType = "item"

// Target is a struct that represents data used to locate a given item within the system
type Target struct {
	TargetType       TargetType `copier:"-" json:"targetType" enums:"item,action,rfc" example:"item" binding:"required"`
	TargetDescriptor string     `copier:"-" json:"targetDescriptor" example:"Event" binding:"required"`
	LookupKey        string     `copier:"LookupKey" json:"lookupKey,omitempty" example:"name"`
	LookupValue      string     `copier:"LookupValue" json:"lookupValue,omitempty" example:"MyNewEvent"`
} // @name Target

// SignatureLookupKey is used to target the signature attributes
var SignatureLookupKey string = `signature`

// ToSha enables an `RFC` to return a SHA256 hash of itself
func (rfc *RFC) ToSha() (*string, error) {
	// init. vars to maintain state beyond "if" statements
	var err error
	var jsonBytes []byte

	// build JSON string
	if jsonBytes, err = json.Marshal(rfc); err != nil {
		errStr := "json marshal rfc error"
		fmt.Println(errStr)
		return nil, err
	}

	// hash
	h := sha256.New()
	if _, err = h.Write(jsonBytes); err != nil {
		errStr := "rfc hash generation error"
		fmt.Println(errStr)
		return nil, err
	}

	hashStr := fmt.Sprintf("%x", h.Sum(nil))
	return &hashStr, nil
}

// AddPersistentActions adds the actions that are deemed persistent from the given "old" RFC to "this" RFC
func (rfc *RFC) AddPersistentActions(oldRFC *RFC) {
	// copy persistent actions over
	for _, action := range oldRFC.Actions {
		// currently statically using "comment", but can augment later
		if action.ActionType == CommentAction {
			rfc.Actions = append(rfc.Actions, action)
		}
	}
}

// AddAction adds the given action to the actions defined by this RFC
func (rfc *RFC) AddAction(action Action) error {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var actionSha *string

	// calculate sha
	if actionSha, err = action.ToSha(); err != nil {
		return err
	}

	// add action
	action.Signature = *actionSha
	rfc.Actions = append(rfc.Actions, &action)

	return nil
}

// "comments" is a map of key/value pairs that are detailed below:
// key = RFC or action signature that is being targeted for the comment
// value = the corresponding array of comment strings to add
// AddComments adds the given comments to this RFC, attributing them to the given commenter
func (rfc *RFC) AddComments(comments map[string][]string, commenter string) error {
	// NOTE: it may more straightforward to add the action signatures to a map at the beginning and then loop
	// through the comments

	// this holds built out comment actions to add to the RFC
	// key = target signature
	// value = comment actions
	processed := map[string][]Action{}

	// iterate over RFC actions and create a comment action if one exists for that target
	for _, action := range rfc.Actions {
		if cmts, ok := comments[action.Signature]; ok {
			for _, cmt := range cmts {
				comment := Action{
					ActionType: CommentAction,
					Target: Target{
						TargetType:  ActionTarget,
						LookupKey:   SignatureLookupKey,
						LookupValue: action.Signature,
					},
					Data: map[string]interface{}{
						string(CommentData):   cmt,
						string(CommenterData): commenter,
					},
				}

				processed[action.Signature] = append(processed[action.Signature], comment)
			}
		}
	}

	// handle overall RFC or dangling comments
	for target, cmts := range comments {
		// only create if we haven't processed already
		if _, ok := processed[target]; !ok {
			for _, cmt := range cmts {
				comment := Action{
					ActionType: CommentAction,
					Target: Target{
						TargetType:  RfcTarget,
						LookupKey:   SignatureLookupKey,
						LookupValue: rfc.Signature,
					},
					Data: map[string]interface{}{
						string(CommentData):   cmt,
						string(CommenterData): commenter,
					},
				}

				// dangling note
				if target != rfc.Signature {
					comment.Data[string(NoteData)] = fmt.Sprintf("Target with signature %s was not found in this RFC",
						target)
				}

				processed[target] = append(processed[target], comment)
			}
		}
	}

	// add processed comments to RFC
	for _, comments := range processed {
		for _, comment := range comments {
			if err := rfc.AddAction(comment); err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateLoadStatus updates the RFC load status action to the given status string and attributes it to the given
// requester
func (rfc *RFC) UpdateLoadStatus(status string, requester string) error {
	// init. vars to maintain state beyond "if" statements
	var err error
	var sha *string

	// find if load action already exists and update if so
	for _, action := range rfc.Actions {
		if action.ActionType == LoadAction {
			action.Data[string(LoadStatus)] = status
			action.Data[string(LoadRequester)] = requester
			if sha, err = action.ToSha(); err != nil {
				return err
			} else {
				action.Signature = *sha
			}
			return err
		}
	}

	// add new load action
	loadAction := Action{ActionType: LoadAction, Data: map[string]interface{}{string(LoadStatus): status,
		string(LoadRequester): requester}}
	err = rfc.AddAction(loadAction)

	return err
}

// GetLoadStatus gets the current RFC load status, if any, nil is returned otherwise
func (rfc *RFC) GetLoadStatus() *string {
	// find if load status exists, if so return it
	for _, action := range rfc.Actions {
		if action.ActionType == LoadAction {
			status := fmt.Sprint(action.Data[string(LoadStatus)])
			return &status
		}
	}

	return nil
}

// ToSha enables an `Action` to return a SHA256 hash of itself
func (action *Action) ToSha() (*string, error) {
	// init. vars to maintain state beyond "if" statements
	var err error
	var jsonBytes []byte

	// build JSON string
	if jsonBytes, err = json.Marshal(action); err != nil {
		errStr := "json marshal action error"
		fmt.Println(errStr)
		return nil, err
	}

	// hash
	h := sha256.New()
	if _, err = h.Write(jsonBytes); err != nil {
		errStr := "action hash generation error"
		fmt.Println(errStr)
		return nil, err
	}

	hashStr := fmt.Sprintf("%x", h.Sum(nil))
	return &hashStr, nil
}

//Utility function to pretty print arrays of Actions
func (actions Actions) String() string {
	s := "["
	for i, action := range actions {
		if i > 0 {
			s += ", "
		}
		s += action.String()
	}
	return s + "]"
}

//Utility function to pretty print a single Action
//Purposefully leaving out the signature
func (action Action) String() string {
	s := "{"
	if action.ActionType != "" {
		s += fmt.Sprintf("ActionType: %v ", action.ActionType)
	}
	if action.Target != (Target{}) {
		s += fmt.Sprintf("Target: %v ", action.Target)
	}
	if len(action.Data) > 0 {
		s += fmt.Sprintf("Data: %v", action.Data)
	}

	return s + "}"
}
