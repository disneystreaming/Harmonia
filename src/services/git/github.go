// This is the GitHub implementation of the Git interface found in definition.go
package git

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
	"harmonia-example.io/src/models"
	"harmonia-example.io/src/services/config"
	"harmonia-example.io/src/services/set"
)

const (
	trackingRepositoryEnvVar = "TRACKING_REPOSITORY"
)

// GitHub type implements the Git interface for GitHub
type GitHub struct {
	AccessToken        *string
	client             *github.Client
	trackingRepository *string
}

// NewGitHub returns a GitHub Git implementation
func NewGitHub(ctx context.Context, accessToken string) (*GitHub, error) {
	// create instance with new client
	g := &GitHub{AccessToken: &accessToken}
	if err := g.setClient(ctx); err != nil {
		return nil, err
	}

	// set tracking repository - env var if local, else AWS param
	repo, err := config.GetTrackingRepo()
	if err != nil {
		return nil, err
	}
	g.trackingRepository = repo

	return g, nil
}

// setClient sets a Go-GitHub client on the caller that can be used to interact with GitHub
func (g *GitHub) setClient(ctx context.Context) error {
	// establish token config for git
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *g.AccessToken})
	tc := oauth2.NewClient(ctx, ts)

	// establish client
	g.client = github.NewClient(tc)

	return nil
}

// CreateBranch creates a new branch with the given name from the given base branch
func (g *GitHub) CreateBranch(ctx context.Context, branch string, baseBranch string) error {
	// init. vars to maintain scope beyond "if" statements
	var base *github.Branch
	var err error

	// get a reference to the base branch
	if base, _, err = g.client.Repositories.GetBranch(ctx, OWNER, *g.trackingRepository, baseBranch, true); err != nil {
		errStr := "error retrieving base branch"
		fmt.Println(errStr)
		return err
	}

	// create branch with the given name
	targetRef := fmt.Sprintf("refs/heads/%s", branch)
	if _, _, err = g.client.Git.CreateRef(
		ctx,
		OWNER,
		*g.trackingRepository,
		&github.Reference{Ref: &targetRef, Object: &github.GitObject{SHA: base.Commit.SHA}},
	); err != nil {
		errStr := "error creating new branch: %s"
		fmt.Println(errStr)
		return err
	}

	return nil
}

// DeleteBranch deletes the branch with the given name
func (g *GitHub) DeleteBranch(ctx context.Context, branch string) error {
	// init. vars to maintain scope beyond "if" statements
	var err error

	// delete branch
	targetRef := fmt.Sprintf("heads/%s", branch)
	if _, err = g.client.Git.DeleteRef(
		ctx,
		OWNER,
		*g.trackingRepository,
		targetRef,
	); err != nil {
		errStr := "Unable to automatically delete branch: %s, please delete manually"
		fmt.Println(errStr)
		return err
	}

	return nil
}

// CreateFile creates an RFC file on the given branch in the given directory using the given data
func (g *GitHub) CreateFile(ctx context.Context, branch string, directory string, data *models.RFC) error {
	// base message
	commitMessage := "init."

	// init. vars to maintain scope beyond "if" statements
	var err error
	var jsonBytes []byte

	// transform data to bytes, which API accepts
	if jsonBytes, err = json.Marshal(data); err != nil {
		errStr := "json data marshal error"
		fmt.Println(errStr)
		return err
	}

	// file creation
	path := fmt.Sprintf("%s/%s/%s", BASE_RFC_DIRECTORY_NAME, directory, RFC_FILE_NAME)
	if _, _, err = g.client.Repositories.CreateFile(
		ctx,
		OWNER,
		*g.trackingRepository,
		path,
		&github.RepositoryContentFileOptions{
			Message: &commitMessage,
			Content: jsonBytes,
			Branch:  &branch,
		},
	); err != nil {
		errStr := "GitHub file creation error"
		fmt.Println(errStr)
		return err
	}

	return nil
}

// CreatePullRequest opens a new pull request of the given branch towards the given base branch
func (g *GitHub) CreatePullRequest(ctx context.Context, branch string, baseBranch string) error {
	// init. vars to maintain scope beyond "if" statements
	var err error

	// PR title
	title := fmt.Sprintf("RFC: %s", branch)
	body := fmt.Sprintf("Automated creation of RFC %s PR", branch)

	// open PR
	if _, _, err = g.client.PullRequests.Create(
		ctx,
		OWNER,
		*g.trackingRepository,
		&github.NewPullRequest{
			Title: &title,
			Head:  &branch,
			Base:  &baseBranch,
			Body:  &body,
		},
	); err != nil {
		errStr := "GitHub PR creation error for branch: %s"
		fmt.Printf(errStr, branch)
		return err
	}

	return nil
}

// GetRFCContents returns the current contents of the RFC on the given branch in the given directory
// The sha of the file is also returned
func (g *GitHub) GetRFCContents(ctx context.Context, branch string) (*string, *string, error) {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var repositoryContent *github.RepositoryContent
	var content string

	// retrieve file contents
	path := fmt.Sprintf("%s/%s/%s", BASE_RFC_DIRECTORY_NAME, branch, RFC_FILE_NAME)
	if repositoryContent, _, _, err = g.client.Repositories.GetContents(
		ctx,
		OWNER,
		*g.trackingRepository,
		path,
		&github.RepositoryContentGetOptions{
			Ref: branch,
		},
	); err != nil {
		errStr := "unable to retrieve repository content"
		fmt.Println(errStr)
		return nil, nil, err
	}

	// extract content for file and retrieve sha
	if content, err = repositoryContent.GetContent(); err != nil {
		errStr := "unable to extract file content from repository content"
		fmt.Println(errStr)
		return nil, nil, err
	}
	sha := repositoryContent.GetSHA()

	return &content, &sha, nil
}

// GetFileSha returns the current RFC file sha for the given pull request
func (g *GitHub) getFileSha(ctx context.Context, pr PullRequest) (*string, error) {
	// ensure given pr is of github type
	githubPr, ok := pr.(*github.PullRequest)
	if !ok {
		errStr := "given pull request is not of type github.PullRequest"
		fmt.Println(errStr)
		return nil, fmt.Errorf(errStr)
	}

	// init. vars to maintain scope beyond "if" statements
	var err error
	var repositoryContent *github.RepositoryContent

	// retrieve file contents so sha can be extracted
	path := fmt.Sprintf("%s/%s/%s", BASE_RFC_DIRECTORY_NAME, *githubPr.Head.Ref, RFC_FILE_NAME)
	if repositoryContent, _, _, err = g.client.Repositories.GetContents(
		ctx,
		OWNER,
		*g.trackingRepository,
		path,
		&github.RepositoryContentGetOptions{
			Ref: *githubPr.Head.Ref,
		},
	); err != nil {
		errStr := "unable to retrieve repository content for sha extraction"
		fmt.Println(errStr)
		return nil, err
	}

	return repositoryContent.SHA, err
}

// UpdateFile creates a commit to the RFC file of the given PR using the given data
func (g *GitHub) UpdateFile(ctx context.Context, pr PullRequest, data *models.RFC) error {
	commitMessage := "update."

	// init. vars to maintain scope beyond "if" statements
	var err error
	var sha *string
	var jsonBytes []byte

	// ensure given pr is of github type
	githubPr, ok := pr.(*github.PullRequest)
	if !ok {
		errStr := "given pull request is not of type github.PullRequest"
		fmt.Println(errStr)
		return fmt.Errorf(errStr)
	}

	// retrieve file sha - necessary for update request
	if sha, err = g.getFileSha(ctx, pr); err != nil {
		return err
	}

	// transform data to bytes, which API accepts
	if jsonBytes, err = json.Marshal(data); err != nil {
		errStr := "json data marshal error"
		fmt.Println(errStr)
		return err
	}

	// update the file in the repo
	path := fmt.Sprintf("%s/%s/%s", BASE_RFC_DIRECTORY_NAME, *githubPr.Head.Ref, RFC_FILE_NAME)
	if _, _, err = g.client.Repositories.UpdateFile(
		ctx,
		OWNER,
		*g.trackingRepository,
		path,
		&github.RepositoryContentFileOptions{
			Message: &commitMessage,
			Content: jsonBytes,
			Branch:  githubPr.Head.Ref,
			SHA:     sha,
		},
	); err != nil {
		errStr := "GitHub update file error"
		fmt.Println(errStr)
		return err
	}

	return nil
}

// GetPullRequest returns the corresponding pull request for the given branch
func (g *GitHub) GetPullRequest(ctx context.Context, branch string) (PullRequest, error) {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var prs []*github.PullRequest

	// retrieve PRs
	if prs, _, err = g.client.PullRequests.List(
		ctx,
		OWNER,
		*g.trackingRepository,
		&github.PullRequestListOptions{
			State: ALL_PR_FILTER,
			Head:  fmt.Sprintf("%s:%s", OWNER, branch),
		},
	); err != nil {
		errStr := "unable to fetch PRs"
		fmt.Println(errStr)
		return nil, err
	}

	// assert we only got 1 PR back
	if len(prs) != 1 {
		errStr := "exactly one PR was NOT returned"
		fmt.Println(errStr)
		return nil, fmt.Errorf(errStr)
	}

	return prs[0], nil
}

// GetPullRequests returns all pull requests with the given state. Paginated output
func (g *GitHub) GetPullRequests(ctx context.Context, state string, count int, opts ...FilterOption) (PullRequests, error) {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var results []*github.PullRequest
	var response *github.Response
	var prs PullRequests

	retrieved := 0
	pageNumber := 1
	perPage := 100
	// Min isn't defined for integers for some reason
	min := func(a int, b int) int {
		if a < b {
			return a
		}
		return b
	}
	if count != -1 {
		perPage = min(count, 100)
	}

	// Default behavior for PR state
	if state == "" {
		state = ALL_PR_FILTER
	}

	// retrieve PRs
	for retrieved < count || count == -1 { // loop until results are exhausted if count is -1
		if results, response, err = g.client.PullRequests.List(
			ctx,
			OWNER,
			*g.trackingRepository,
			&github.PullRequestListOptions{
				State: state,
				ListOptions: github.ListOptions{
					Page:    pageNumber,
					PerPage: perPage,
				},
			},
		); err != nil {
			errStr := "unable to fetch PRs"
			fmt.Println(errStr)
			return nil, err
		}

		// serialize
		var isValid bool
		for _, result := range results {
			// filter
			isValid = true
			for _, opt := range opts {
				isValid = isValid && opt(result)
			}
			if isValid && (len(prs) < count || count == -1) {
				prs = append(prs, result)
				retrieved++
			}
		}

		// go to next page
		pageNumber = response.NextPage

		// 0 value indicates there is no next page and the results are exhausted
		if pageNumber == 0 {
			break
		}
	}

	return prs, nil
}

// GetMergeability determines if the given pull request is mergeable (approvals, conflicts, ci...)
func (g *GitHub) GetMergeability(ctx context.Context, pr PullRequest) (*bool, error) {
	// ensure given pr is of github type
	githubPr, ok := pr.(*github.PullRequest)
	if !ok {
		errStr := "given pull request is not of type github.PullRequest"
		fmt.Println(errStr)
		return nil, fmt.Errorf(errStr)
	}

	// init. vars to maintain state beyond "if" statements
	var err error
	var status *github.CombinedStatus

	// poll for commit status and allow time for it to stabilize, within reason
	for retryCount := 0; retryCount < MERGEABILITY_RETRY_COUNT; retryCount++ {
		// get combined status - this represents overall status, taking all checks into account
		if status, _, err = g.client.Repositories.GetCombinedStatus(
			ctx,
			OWNER,
			*g.trackingRepository,
			*githubPr.Head.Ref,
			&github.ListOptions{},
		); err != nil {
			errStr := "unable to retrieve ref combined status"
			fmt.Println(errStr)
			return nil, err
		}

		// check and see if the state is still pending, if so, wait a set amount of time and a re-poll
		if status.State != nil && *status.State == MERGEABILITY_PENDING_STATE {
			time.Sleep(time.Duration(MERGEABILITY_WAIT_TIME) * time.Second)
			continue
		}

		break
	}

	// retrieve pr
	// this is unfortunate, but the pr has to be refetched to be able to pull the recalculated mergeable state off of
	// it. According to the docs, mergeable state is calculated in the background by GitHub so polling is necessary here
	// as well.
	// https://docs.github.com/en/rest/reference/pulls#get-a-pull-request
	for retryCount := 0; retryCount < MERGEABILITY_RETRY_COUNT; retryCount++ {
		// not using the "getPullRequest" function here because it uses the list functionality, which doesn't calculate
		// the mergeable state
		if githubPr, _, err = g.client.PullRequests.Get(
			ctx,
			OWNER,
			*g.trackingRepository,
			*githubPr.Number,
		); err != nil {
			errStr := "unable to retrieve pr for mergeability check"
			fmt.Println(errStr)
			return nil, err
		}

		// if still calculating, wait and re-poll
		if githubPr.MergeableState == nil || *githubPr.MergeableState == MERGEABILITY_UNKNOWN_STATE {
			time.Sleep(time.Duration(MERGEABILITY_WAIT_TIME) * time.Second)
			continue
		}

		break
	}

	// mergeability was never able to be determined
	if githubPr.MergeableState == nil || *githubPr.MergeableState == MERGEABILITY_UNKNOWN_STATE {
		errStr := "unable to determine mergeability of rfc"
		fmt.Println(errStr)
		return nil, fmt.Errorf(errStr)
	}

	mergeable := *githubPr.MergeableState == MERGEABILITY_CLEAN_STATE
	return &mergeable, nil
}

// MergePullRequest merges the given pull request and returns the sha
func (g *GitHub) MergePullRequest(ctx context.Context, pr PullRequest) (*string, error) {
	// ensure given pr is of github type
	githubPr, ok := pr.(*github.PullRequest)
	if !ok {
		errStr := "given pull request is not of type github.PullRequest"
		fmt.Println(errStr)
		return nil, fmt.Errorf(errStr)
	}

	// pull request commit message
	message := ""

	// init. vars to maintain scope beyond "if" statements
	var err error
	var res *github.PullRequestMergeResult

	// merge
	if res, _, err = g.client.PullRequests.Merge(
		ctx,
		OWNER,
		*g.trackingRepository,
		*githubPr.Number,
		message,
		&github.PullRequestOptions{
			DontDefaultIfBlank: false,
		},
	); err != nil {
		errStr := "unable to merge pull request"
		fmt.Println(errStr)
		return nil, err
	}

	return res.SHA, nil
}

// GetReviews returns all pull request reviews related to the given pull request
func (g *GitHub) GetReviews(ctx context.Context, pr PullRequest) (PullRequestReviews, error) {
	// ensure given pr is of github type
	githubPr, ok := pr.(*github.PullRequest)
	if !ok {
		errStr := "given pull request is not of type github.PullRequest"
		fmt.Println(errStr)
		return nil, fmt.Errorf(errStr)
	}

	// init. vars to maintain scope beyond "if" statements
	var err error
	var reviews []*github.PullRequestReview

	// retrieve reviews
	if reviews, _, err = g.client.PullRequests.ListReviews(
		ctx,
		OWNER,
		*g.trackingRepository,
		*githubPr.Number,
		&github.ListOptions{
			PerPage: 100,
		},
	); err != nil {
		errStr := "GitHub list reviews error"
		fmt.Println(errStr)
		return nil, err
	}

	return reviews, nil
}

// CreateReview generates a pull request review on the given pull request using the given data
func (g *GitHub) CreateReview(ctx context.Context, pr PullRequest, data *models.Review) error {
	// ensure given pr is of github type
	githubPr, ok := pr.(*github.PullRequest)
	if !ok {
		errStr := "given pull request is not of type github.PullRequest"
		fmt.Println(errStr)
		return fmt.Errorf(errStr)
	}

	// the file to target for review comments
	path := fmt.Sprintf("%s/%s/%s", BASE_RFC_DIRECTORY_NAME, data.RFCIdentifier, RFC_FILE_NAME)
	// all comments relate to the only line in the RFC
	position := 1

	// generate comment structure to attach to the review
	comments := []*github.DraftReviewComment{}
	for _, cmts := range data.Comments {
		for _, cmt := range cmts {
			commentBody := cmt
			comment := &github.DraftReviewComment{
				Path:     &path,
				Body:     &commentBody,
				Position: &position,
			}

			comments = append(comments, comment)
		}
	}

	// pre-generate param so body can be added if necessary
	param := &github.PullRequestReviewRequest{
		Event:    &data.Type,
		Comments: comments,
	}

	// add body if appropriate
	if data.TopLevelComment != "" {
		param.Body = &data.TopLevelComment
	}

	// generate review
	if _, _, err := g.client.PullRequests.CreateReview(
		ctx,
		OWNER,
		*g.trackingRepository,
		*githubPr.Number,
		param,
	); err != nil {
		errStr := "unable to create review"
		fmt.Println(errStr)
		return err
	}

	return nil
}

// DismissApprovalReviews dismisses only the "approval" reviews in the given reviews from the given pull request
func (g *GitHub) DismissApprovalReviews(ctx context.Context, reviews PullRequestReviews, pr PullRequest) error {
	// ensure given reviews are of github type
	githubPrReviews, ok := reviews.([]*github.PullRequestReview)
	if !ok {
		errStr := "given pull request reviews is not of type []github.PullRequestReview"
		fmt.Println(errStr)
		return fmt.Errorf(errStr)
	}
	// ensure given pr is of github type
	githubPr, ok := pr.(*github.PullRequest)
	if !ok {
		errStr := "given pull request is not of type github.PullRequest"
		fmt.Println(errStr)
		return fmt.Errorf(errStr)
	}

	// dismissed message
	message := "dismissed."

	// only operate on approvals
	for _, review := range githubPrReviews {
		// only dismiss approvals
		if *review.State == APPROVED_STATE {
			// dismiss review
			if _, _, err := g.client.PullRequests.DismissReview(
				ctx,
				OWNER,
				*g.trackingRepository,
				*githubPr.Number,
				*review.ID,
				&github.PullRequestReviewDismissalRequest{
					Message: &message,
				},
			); err != nil {
				errStr := "GitHub dismiss review error"
				fmt.Println(errStr)
				return err
			}
		}
	}

	return nil
}

// GetUserLogin returns the Git username defined by the client
func (g *GitHub) GetUserLogin(ctx context.Context) (*string, error) {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var user *github.User

	// retrieve user
	if user, _, err = g.client.Users.Get(ctx, ""); err != nil {
		errStr := "unable to fetch user"
		fmt.Println(errStr)
		return nil, err
	}

	return user.Login, nil
}

// GetUserTeams returns a set of teams for the current authenticated user
func (g *GitHub) GetUserTeams(ctx context.Context) (set.Set[string], error) {
	// init. vars to maintain scope beyond "if" statements
	var err error
	var ghTeams []*github.Team
	var response *github.Response
	teams := set.NewSet[string]()
	page := 1
	perPage := 100

	// get user teams, paginated for users with many teams
	for page != 0 {
		if ghTeams, response, err = g.client.Teams.ListUserTeams(
			ctx,
			&github.ListOptions{
				PerPage: perPage,
				Page:    page,
			},
		); err != nil {
			errStr := "unable to retrieve user teams"
			fmt.Println(errStr)
			return nil, err
		}

		// add to teams set
		for _, team := range ghTeams {
			teams.Add(*team.Name)
		}

		// check what the next page is, terminate if none left
		page = response.NextPage
	}

	return teams, nil
}

// CreateTag tags the given sha with the given name
func (g *GitHub) CreateTag(ctx context.Context, sha string, tag string) error {
	// tag resource
	targetRef := fmt.Sprintf("refs/tags/%s", tag)
	if _, _, err := g.client.Git.CreateRef(
		ctx,
		OWNER,
		*g.trackingRepository,
		&github.Reference{
			Ref:    &targetRef,
			Object: &github.GitObject{SHA: &sha},
		},
	); err != nil {
		errStr := "unable to create tag"
		fmt.Println(errStr)
		return err
	}

	return nil
}

// GetIdsAndTitles is a helper method used to retrieve UI data from an array of Pull Requests
func (g *GitHub) GetIdsAndTitles(prs PullRequests) (IdsAndTitles, error) {
	idsAndTitles := make([]map[string]string, len(prs))
	for i, pr := range prs {
		githubPr, ok := pr.(*github.PullRequest)
		if !ok {
			return nil, fmt.Errorf("cannot convert given pull request to github.PullRequest")
		}
		idsAndTitles[i] = map[string]string{*githubPr.Head.Ref: *githubPr.Title}
	}

	return idsAndTitles, nil
}

// Returns a FilterOption that:
// 	returns true if a given PR is owned by the given user. If no user is given, returns true.
func (g *GitHub) WithOwner(owner *string) FilterOption {
	return func(pr PullRequest) bool {
		githubPr, ok := pr.(*github.PullRequest)
		if !ok {
			return false
		}

		if owner != nil {
			if githubPr.User == nil || githubPr.User.Login == nil {
				return false
			}

			return *owner == *githubPr.User.Login
		}

		return true
	}
}

// Returns a FilterOption that:
//	returns true if a given PR has a merged state equal to the provided state. If no state is given, returns true.
func (g *GitHub) IsMerged(merged *bool) FilterOption {
	return func(pr PullRequest) bool {
		githubPr, ok := pr.(*github.PullRequest)
		if !ok {
			return false
		}

		if merged != nil {
			if githubPr.Merged == nil {
				return !*merged
			}

			return *merged == *githubPr.Merged
		}

		return true
	}
}
