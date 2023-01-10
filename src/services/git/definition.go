// This is strictly to hold the Git interface definition and common constants used in Git interactions
// All Git specific API implementations (GitHub, BitBucket...) should be in this package but outside of this file
package git

import (
	"context"

	"harmonia-example.io/src/models"
	"harmonia-example.io/src/services/set"
)

// Common constants that will be used across all Git implementations and interactions
const (
	OWNER                       string = "<repository-owner>"
	BASE_BRANCH                 string = "main"
	RFC_FILE_NAME               string = "RFC.json"
	BASE_RFC_DIRECTORY_NAME     string = "RFC"
	APPROVED_STATE              string = "APPROVED"
	OPEN_STATE                  string = "open"
	APPROVE_REVIEW_TYPE         string = "APPROVE"
	REQUEST_CHANGES_REVIEW_TYPE string = "REQUEST_CHANGES"
	COMMENT_REVIEW_TYPE         string = "COMMENT"
	MERGEABILITY_CLEAN_STATE    string = "clean"
	MERGEABILITY_PENDING_STATE  string = "pending"
	MERGEABILITY_UNKNOWN_STATE  string = "unknown"
	MERGEABILITY_RETRY_COUNT    int    = 3
	MERGEABILITY_WAIT_TIME      int    = 10
	ALL_PR_FILTER               string = "all"
)

// PullRequest is a generic Git type used to generalize implementations
type PullRequest interface{}

// PullRequests represents a mapping of RFC ID to PR title for display and UX
type PullRequests []interface{}

// PullRequestReview is a generic Git type used to generalize implementation
type PullRequestReview interface{}

// PullRequestReviews is a generic Git type used to generalize implementation
type PullRequestReviews interface{}

// IdsAndTitles is an aliased type meant to represent an ordered list of pairs of strings
// the key is the ID of an RFC and the value is the title.
type IdsAndTitles []map[string]string

type FilterOption func(PullRequest) bool

// Git defines all methods necessary for Harmonia Git interactions
// All git types (GitHub, BitBucket...) should implement this interface
type Git interface {
	// CreateBranch creates a new branch with the given name from the given base branch
	CreateBranch(ctx context.Context, branch string, baseBranch string) error
	// DeleteBranch deletes the branch with the given name
	DeleteBranch(ctx context.Context, branch string) error
	// CreateFile creates an RFC file on the given branch in the given directory using the given data
	CreateFile(ctx context.Context, branch string, directory string, data *models.RFC) error
	// CreatePullRequest opens a new pull request of the given branch towards the given base branch
	CreatePullRequest(ctx context.Context, branch string, baseBranch string) error
	// GetRFCContents returns the current contents of the RFC for the given pull request
	// The sha of the file is also returned
	GetRFCContents(ctx context.Context, branch string) (*string, *string, error)
	// UpdateFile creates a commit to the RFC file of the given PR using the given data
	UpdateFile(ctx context.Context, pr PullRequest, data *models.RFC) error
	// GetPullRequest returns the most recent open pull request for the given branch
	GetPullRequest(ctx context.Context, branch string) (PullRequest, error)
	// GetPullRequests returns all pull requests with the given state and filters
	GetPullRequests(ctx context.Context, state string, count int, opts ...FilterOption) (PullRequests, error)
	// GetMergeability determines if the given pull request is mergeable (approvals, conflicts, ci...)
	GetMergeability(ctx context.Context, pr PullRequest) (*bool, error)
	// MergePullRequest merges the given pull request and returns the sha
	MergePullRequest(ctx context.Context, pr PullRequest) (*string, error)
	// GetReviews returns all pull request reviews related to the given pull request
	// TODO: interface temporary
	GetReviews(ctx context.Context, pr PullRequest) (PullRequestReviews, error)
	// CreateReview generates a pull request review on the given pull request using the given data
	CreateReview(ctx context.Context, pr PullRequest, data *models.Review) error
	// DismissApprovalReviews dismisses only the "approval" reviews in the given reviews from the given pull request
	DismissApprovalReviews(ctx context.Context, reviews PullRequestReviews, pr PullRequest) error
	// GetUserLogin returns the Git username defined by the client
	GetUserLogin(ctx context.Context) (*string, error)
	// GetUserTeams returns a set of teams for the current authenticated user in the form "<org-name>/<team-name>"
	GetUserTeams(ctx context.Context) (set.Set[string], error)
	// CreateTag tags the given sha with the given name
	CreateTag(ctx context.Context, sha string, name string) error

	// GetIdsAndTitles is meant to retrieve the RFC ID and Title returned from GetPullRequests
	GetIdsAndTitles(prs PullRequests) (IdsAndTitles, error)

	// The following are functions that are meant to support filtering queries like e.g. GetPullRequests
	WithOwner(owner *string) FilterOption
	IsMerged(merged *bool) FilterOption
}
