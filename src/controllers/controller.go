// Package controllers
// is used to hold all controller functions
// A controller function handles orchestration logic after the routes. Each controller function should simply call
// service functions and leave the complex logic up to them.
package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"harmonia-example.io/src/models"
	exGit "harmonia-example.io/src/services/git"
)

const (
	// statuses for RFC loads
	LOAD_REQUESTED_STATUS = "load_requested"
	NOT_APPLICABLE_STATUS = "not_applicable"
	LOADING_STATUS        = "loading"
	SUCCESSFUL_STATUS     = "successful"
	FAILED_STATUS         = "failed"
)

// CreateRFCIdentifier creates a unique identifier for a new RFC
var CreateRFCIdentifier models.RFCIdentifierCreator = func() *string {
	// Creates identifier based on current time
	epoch := time.Now().Unix()
	identifier := strconv.FormatInt(epoch, 10)
	return &identifier
}

// SubmitRequest orchestrates creating a new RFC branch, making the first commit with the given RFC data and
// opening a pull request. The corresponding branch name is returned.
// Parameters:
//	ctx - standard context
//	git - Git service implementation used to drive interactions
// 	data - RFC to populate
func SubmitRequest(ctx context.Context, git exGit.Git, data *models.RFC) (*string, error) {
	// add hash signatures to incoming data
	rfcSignature, err := data.ToSha()
	if err != nil {
		return nil, err
	}
	data.Signature = *rfcSignature
	for _, action := range data.Actions {
		actionSha, err := action.ToSha()
		if err != nil {
			return nil, err
		}
		action.Signature = *actionSha
	}

	// create new branch identifier
	branch := *CreateRFCIdentifier()

	// <this is a good place to add RFC metadata to logger> //

	if err = git.CreateBranch(ctx, branch, exGit.BASE_BRANCH); err != nil {
		errStr := "Failed to create branch for RFC: %s, please try again"
		fmt.Printf(errStr, branch)
		return nil, err
	}

	// create new RFC file
	if err = git.CreateFile(ctx, branch, branch, data); err != nil {
		errStr := "Failed to write file for RFC: %s to datastore, starting revoke process..."
		fmt.Printf(errStr, branch)
		if revErr := git.DeleteBranch(ctx, branch); revErr == nil {
			infoStr := "Successfully revoked RFC: %s"
			fmt.Printf(infoStr, branch)
		}
		return nil, err
	}

	// open PR
	if err = git.CreatePullRequest(ctx, branch, exGit.BASE_BRANCH); err != nil {
		errStr := "Failed to open Pull Request for RFC: %s, starting revoke process..."
		fmt.Printf(errStr, branch)
		if revErr := git.DeleteBranch(ctx, branch); revErr == nil {
			infoStr := "Successfully revoked RFC: %s"
			fmt.Printf(infoStr, branch)
		}
		return nil, err
	}

	return &branch, nil
}

// UpdateRequest orchestrates the update RFC process, which includes updating an existing RFC, persisting existing
// actions and clearing out existing approvals. The branch name is returned.
// Parameters:
// 	ctx - standard context
// 	git - Git service implementation used to drive interactions
//	data - RFC new data
func UpdateRequest(ctx context.Context, git exGit.Git, data *models.Update) (*string, error) {
	// retrieve pull request
	pr, err := git.GetPullRequest(ctx, data.RFCIdentifier)
	if err != nil {
		return nil, err
	}

	// retrieve existing RFC content
	content, _, err := git.GetRFCContents(ctx, data.RFCIdentifier)
	if err != nil {
		return nil, err
	}

	// format existing RFC into model
	existingRFC := &models.RFC{}
	if err = json.Unmarshal([]byte(*content), existingRFC); err != nil {
		errStr := "unable to unmarshal existing RFC content"
		fmt.Print(errStr)
		return nil, err
	}

	// add action hash signatures
	for _, action := range data.RFC.Actions {
		actionSha, err := action.ToSha()
		if err != nil {
			return nil, err
		}
		action.Signature = *actionSha
	}

	// persist actions from existing RFC to new RFC
	data.RFC.AddPersistentActions(existingRFC)

	// add rfc hash signature
	rfcSignature, err := data.RFC.ToSha()
	if err != nil {
		return nil, err
	}
	data.RFC.Signature = *rfcSignature

	// update existing RFC in repo
	if err = git.UpdateFile(ctx, pr, data.RFC); err != nil {
		return nil, err
	}

	reviews, err := git.GetReviews(ctx, pr)
	if err != nil {
		return nil, err
	}
	if err = git.DismissApprovalReviews(ctx, reviews, pr); err != nil {
		return nil, err
	}

	return &data.RFCIdentifier, nil
}

// ReviewRequest orchestrates submitting a review based on the given data
func ReviewRequest(ctx context.Context, git exGit.Git, gitMachine exGit.Git, data *models.Review) (*string, error) {
	// if the review type is a comment or requesting changes there needs to be some sort of comments associated
	if data.Type == exGit.COMMENT_REVIEW_TYPE || data.Type == exGit.REQUEST_CHANGES_REVIEW_TYPE {
		if data.TopLevelComment == "" && len(data.Comments) == 0 {
			errStr := fmt.Sprintf("Review of type %s must include a top level comment or inline comments", data.Type)
			fmt.Println(errStr)
			return nil, fmt.Errorf(errStr)
		}
	}

	// retrieve PR associated with the given rfcIdentifier
	pr, err := git.GetPullRequest(ctx, data.RFCIdentifier)
	if err != nil {
		return nil, err
	}

	// retrieve current user
	login, err := git.GetUserLogin(ctx)
	if err != nil {
		return nil, err
	}

	// retrieve existing RFC content
	content, _, err := git.GetRFCContents(ctx, data.RFCIdentifier)
	if err != nil {
		return nil, err
	}

	// format existing RFC into model
	rfc := &models.RFC{}
	if err = json.Unmarshal([]byte(*content), rfc); err != nil {
		errStr := "unable to unmarshal existing RFC content"
		fmt.Print(errStr)
		return nil, err
	}

	// add comments to RFC
	if err = rfc.AddComments(data.Comments, *login); err != nil {
		return nil, err
	}

	// we only want to create a review action if this is an approval or request for changes OR there are top level comments
	if data.Type != exGit.COMMENT_REVIEW_TYPE || data.TopLevelComment != "" {
		// our identifier = reviewer, unless this is a comment, then we want commenter
		identifier := models.ReviewerData
		if data.Type == exGit.COMMENT_REVIEW_TYPE {
			identifier = models.CommentData
		}
		action := models.Action{
			ActionType: models.ActionType(strings.ToLower(data.Type)),
			Target: models.Target{
				TargetType:  models.TargetType("rfc"),
				LookupKey:   models.SignatureLookupKey,
				LookupValue: rfc.Signature,
			},
			Data: map[string]interface{}{
				string(identifier): *login,
			},
		}
		// add review comment if necessary
		if data.TopLevelComment != "" {
			action.Data["comment"] = data.TopLevelComment
		}
		// add the review action to the RFC
		if err = rfc.AddAction(action); err != nil {
			return nil, err
		}
	}

	// propagate updated RFC to the repo
	if err = git.UpdateFile(ctx, pr, rfc); err != nil {
		return nil, err
	}

	// create PR review
	if err = git.CreateReview(ctx, pr, data); err != nil {
		return nil, err
	}

	var message string
	// if this was an approval and the user wishes to initiate a load request, then attempt the load and merge process
	if data.Type == exGit.APPROVE_REVIEW_TYPE && data.LoadOnApproval {
		/*
			all admin work to be performed by machine client

			attempt to load and merge request asynchronously
			a new unattached context needs to be created prior to the call because the go routine is not waited on
			and any cancellation will invalidate the child
		*/
		go attemptLoadAndMerge(context.Background(), gitMachine, pr, rfc, data.RFCIdentifier)
		message = fmt.Sprintf(`Successfully approved RFC %s. A load request was submitted. You may query the load status
		through the /status endpoint.`, data.RFCIdentifier)
	} else {
		message = fmt.Sprintf("Successfully reviewed RFC %s with type of '%s'", data.RFCIdentifier, data.Type)
	}

	return &message, nil
}

// MergeRequest orchestrates merging the given RFC and tagging it for tracking, returns a message if successful
func MergeRequest(ctx context.Context, git exGit.Git, data *models.Merge) (*string, error) {
	// init. vars to maintain state beyond "if" statements
	var err error
	var pr exGit.PullRequest

	// get corresponding pr
	if pr, err = git.GetPullRequest(ctx, data.RFCIdentifier); err != nil {
		return nil, err
	}

	// merge request and create tag with the rfc identifier name
	if err = mergeRequest(ctx, git, pr, data.RFCIdentifier); err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Successfully merged and tagged RFC %s", data.RFCIdentifier)
	return &message, nil
}

// LoadRequest orchestrates loading the given RFC data into the backing datastore asynchronously - load status will
// be populated in the RFC file
func LoadRequest(ctx context.Context, git exGit.Git, data *models.Load) error {
	// init. vars to maintain state beyond "if" statements
	var err error
	var pr exGit.PullRequest
	var content *string
	var user *string

	// Get user login for load status update
	if user, err = git.GetUserLogin(ctx); err != nil {
		return err
	}

	// get corresponding pr so content can be fetched
	if pr, err = git.GetPullRequest(ctx, data.RFCIdentifier); err != nil {
		return err
	}

	// retrieve corresponding raw RFC content that will be loaded
	if content, _, err = git.GetRFCContents(ctx, data.RFCIdentifier); err != nil {
		return err
	}

	// format existing content into RFC model so the load status can be manipulated
	rfc := &models.RFC{}
	if err = json.Unmarshal([]byte(*content), rfc); err != nil {
		errStr := "unable to unmarshal existing RFC content in preparation for load, RFC: %s"
		fmt.Printf(errStr, data.RFCIdentifier)
		return err
	}

	// update load status to LOAD_REQUESTED_STATUS so that there is a record of this request
	if err = rfc.UpdateLoadStatus(LOAD_REQUESTED_STATUS, *user); err != nil {
		return err
	}
	if err = git.UpdateFile(ctx, pr, rfc); err != nil {
		return err
	}

	/*
		attempt to load request asynchronously
		a new unattached context needs to be created prior to the call because the go routine is not waited on
		and any cancellation will invalidate the child
	*/
	go loadRequest(context.Background(), git, pr, rfc)

	return err
}

// Status returns the current load status of the given RFC, if any
func Status(ctx context.Context, git exGit.Git, data *models.Status) (*string, error) {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var content *string

	// retrieve corresponding raw RFC content that can be parsed
	if content, _, err = git.GetRFCContents(ctx, data.RFCIdentifier); err != nil {
		return nil, err
	}

	// format existing content into RFC model so the load status can be searched for
	rfc := &models.RFC{}
	if err = json.Unmarshal([]byte(*content), rfc); err != nil {
		errStr := "unable to unmarshal existing RFC content in preparation for status retrieval, RFC: %s"
		fmt.Printf(errStr, data.RFCIdentifier)
		return nil, err
	}

	return rfc.GetLoadStatus(), nil
}

// GetRfcs returns all submitted RFCs based on given data filtering
func GetRfcs(ctx context.Context, git exGit.Git, data *models.GetRfcs) ([]map[string]string, error) {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var prs exGit.PullRequests
	filters := []exGit.FilterOption{git.WithOwner(data.Owner), git.IsMerged(data.Merged)}

	// query for PRs
	if prs, err = git.GetPullRequests(ctx, data.State, data.Count, filters...); err != nil {
		return nil, err
	}

	// retrieve RFC ID and Title map
	return git.GetIdsAndTitles(prs)
}

// GetRfcContents returns the contents of the target RFC
func GetRfcContents(ctx context.Context, git exGit.Git, data *models.GetRfcContents) (*string, error) {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var content *string

	// retrieve corresponding raw RFC content that can be parsed
	if content, _, err = git.GetRFCContents(ctx, data.RFCIdentifier); err != nil {
		return nil, err
	}

	return content, nil
}

// the below methods (not capitalized) exist strictly to be called by other functions within this module, which have
// already performed the boilerplate retrieval of rfc entities like the pull request and rfc content

// attemptLoadAndMerge attempts to load and then merge the given RFC data and corresponding pull request
func attemptLoadAndMerge(ctx context.Context, git exGit.Git, pr exGit.PullRequest, rfc *models.RFC,
	rfcIdentifier string) error {
	// init. vars to maintain state beyond "if" statements
	var err error
	var mergeable *bool
	var user *string

	// Get user login for load status update
	if user, err = git.GetUserLogin(ctx); err != nil {
		return err
	}

	// update load status to LOAD_REQUESTED_STATUS
	if err = rfc.UpdateLoadStatus(LOAD_REQUESTED_STATUS, *user); err != nil {
		return err
	}
	if err = git.UpdateFile(ctx, pr, rfc); err != nil {
		return err
	}

	// determine if the pr can be merged, this is 1:1 with loadability (can't load if we can't merge)
	if mergeable, err = git.GetMergeability(ctx, pr); err != nil {
		return err
	}
	if !*mergeable {
		infoStr := "Attempted to load and merge RFC %s, but it is not mergeable."
		fmt.Printf(infoStr, rfcIdentifier)

		// update load status to NOT_APPLICABLE_STATUS
		if err = rfc.UpdateLoadStatus(NOT_APPLICABLE_STATUS, *user); err != nil {
			return err
		}
		if err = git.UpdateFile(ctx, pr, rfc); err != nil {
			return err
		}

		return nil
	}

	// attempt load
	if err = loadRequest(ctx, git, pr, rfc); err != nil {
		return err
	}

	// mergeability needs to be recalculated here because loadRequest updates the RFC file - CI check
	if mergeable, err = git.GetMergeability(ctx, pr); err != nil {
		return err
	}
	if !*mergeable {
		errStr := "Attempted to merge RFC %s, but it is not mergeable - NOTE: LOADED BUT NOT MERGED."
		fmt.Printf(errStr, rfcIdentifier)
		return fmt.Errorf(errStr, rfcIdentifier)
	}

	// attempt merge
	if err = mergeRequest(ctx, git, pr, rfcIdentifier); err != nil {
		return err
	}

	return nil
}

// loadRequest loads the given rfc content into the backing data store
// The pull request param. seems unnecessary, but it is needed to update the load status periodically
func loadRequest(ctx context.Context, git exGit.Git, pr exGit.PullRequest, rfc *models.RFC) error {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var content []byte
	var user *string

	// Get user login for load status update
	if user, err = git.GetUserLogin(ctx); err != nil {
		return err
	}

	// update load status to LOADING_STATUS
	if err = rfc.UpdateLoadStatus(LOADING_STATUS, *user); err != nil {
		return err
	}
	if err = git.UpdateFile(ctx, pr, rfc); err != nil {
		return err
	}

	// format rfc for loading
	if content, err = json.Marshal(rfc); err != nil {
		errStr := "unable to marshal existing RFC content in preparation for load."
		fmt.Printf(errStr)
		return err
	}

	// call database service with the RFC content to load
	// ...
	fmt.Println(content)
	// ...
	// update file with failed status if there was a load error

	// update load status to SUCCESSFUL_STATUS
	if err = rfc.UpdateLoadStatus(SUCCESSFUL_STATUS, *user); err != nil {
		return err
	}
	if err = git.UpdateFile(ctx, pr, rfc); err != nil {
		return err
	}

	return nil
}

// mergeRequest merges the given pr and creates a tag with the given tag name
func mergeRequest(ctx context.Context, git exGit.Git, pr exGit.PullRequest, tag string) error {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var sha *string

	// merge pr and retrieve resulting sha
	if sha, err = git.MergePullRequest(ctx, pr); err != nil {
		return err
	}

	// create a tag of sha and name it after tag name
	if err = git.CreateTag(ctx, *sha, tag); err != nil {
		return err
	}

	return nil
}
