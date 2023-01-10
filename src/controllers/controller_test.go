// This is to hold all tests related to controller.go

package controllers

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"harmonia-example.io/src/models"
	exGit "harmonia-example.io/src/services/git"
	"harmonia-example.io/src/services/set"
)

// gitMockCreator is used to create mocks that implement exGit.Git
// This is done this way so that each test case can have its own mock constructor
type gitMockCreator func() exGit.Git

// mockGit is a base mock that implements exGit.Git
// Each method of exGit.Git is replicated as a lowercase function within the struct so we can override (mock) the
// functionality of the method dynamically for each test case via gitMockCreator
// It is not possible to set the top level uppercase methods dynamically, hence why it is done this way
type mockGit struct {
	// mock.Mock allows us to assert methods were called with certain arguments
	mock.Mock

	createBranch      func(ctx context.Context, branch string, baseBranch string) error
	deleteBranch      func(ctx context.Context, branch string) error
	createFile        func(ctx context.Context, branch string, directory string, data *models.RFC) error
	createPullRequest func(ctx context.Context, branch string, baseBranch string) error
	getRFCContents    func(ctx context.Context, branch string) (*string, *string, error)
	updateFile        func(ctx context.Context, pr exGit.PullRequest, data *models.RFC) error
	getPullRequest    func(ctx context.Context, branch string) (exGit.PullRequest, error)
	getPullRequests   func(ctx context.Context, state string, count int, opts ...exGit.FilterOption) (
		exGit.PullRequests, error)
	getMergeability        func(ctx context.Context, pr exGit.PullRequest) (*bool, error)
	mergePullRequest       func(ctx context.Context, pr exGit.PullRequest) (*string, error)
	getReviews             func(ctx context.Context, pr exGit.PullRequest) (exGit.PullRequestReviews, error)
	createReview           func(ctx context.Context, pr exGit.PullRequest, data *models.Review) error
	dismissApprovalReviews func(ctx context.Context, reviews exGit.PullRequestReviews, pr exGit.PullRequest) error
	getUserLogin           func(ctx context.Context) (*string, error)
	getUserTeams           func(ctx context.Context) (set.Set[string], error)
	createTag              func(ctx context.Context, sha string, name string) error

	getIdsAndTitles func(prs exGit.PullRequests) (exGit.IdsAndTitles, error)

	withOwner func(owner *string) exGit.FilterOption
	isMerged  func(merged *bool) exGit.FilterOption
}

// Each method below simply calls the struct lowercase version that is manipulated per test
// In these methods is where mock.Mock calls should be made because the submethods don't have access to the struct

// CreateBranch calls mg.createBranch
func (mg *mockGit) CreateBranch(ctx context.Context, branch string, baseBranch string) error {
	// ignore ctx for mocking purposes
	// we are ignoring ctx because it is altered by the underlying method and we would have to build one to match
	mg.On("CreateBranch", branch, baseBranch).Return()
	mg.Called(branch, baseBranch)

	return mg.createBranch(ctx, branch, baseBranch)
}

// DelateBranch calls mg.deleteBranch
func (mg *mockGit) DeleteBranch(ctx context.Context, branch string) error {
	// ignore ctx for mocking purposes
	// we are ignoring ctx because it is altered by the underlying method and we would have to build one to match
	mg.On("DeleteBranch", branch).Return()
	mg.Called(branch)

	return mg.deleteBranch(ctx, branch)
}

// CreateFile calls mg.createFile
func (mg *mockGit) CreateFile(ctx context.Context, branch string, directory string, data *models.RFC) error {
	// ignore ctx for mocking purposes
	// we are ignoring ctx because it is altered by the underlying method and we would have to build one to match
	mg.On("CreateFile", branch, directory, data).Return()
	mg.Called(branch, directory, data)

	return mg.createFile(ctx, branch, directory, data)
}

// CreatePullRequest calls mg.createPullRequest
func (mg *mockGit) CreatePullRequest(ctx context.Context, branch string, baseBranch string) error {
	// ignore ctx for mocking purposes
	// we are ignoring ctx because it is altered by the underlying method and we would have to build one to match
	mg.On("CreatePullRequest", branch, baseBranch).Return()
	mg.Called(branch, baseBranch)

	return mg.createPullRequest(ctx, branch, baseBranch)
}

// GetRFCContents calls mg.getRFCContents
func (mg *mockGit) GetRFCContents(ctx context.Context, branch string) (*string, *string, error) {
	// ignore ctx for mocking purposes
	// we are ignoring ctx because it is altered by the underlying method and we would have to build one to match
	mg.On("GetRFCContents", branch).Return()
	mg.Called(branch)

	return mg.getRFCContents(ctx, branch)
}

// UpdateFile calls mg.updateFile
func (mg *mockGit) UpdateFile(ctx context.Context, pr exGit.PullRequest, data *models.RFC) error {
	// ignore ctx for mocking purposes
	// we are ignoring ctx because it is altered by the underlying method and we would have to build one to match
	mg.On("UpdateFile", pr, data).Return()
	mg.Called(pr, data)
	fmt.Println(pr)
	fmt.Println(*data)
	return mg.updateFile(ctx, pr, data)
}

// GetPullRequest calls mg.getPullRequest
func (mg *mockGit) GetPullRequest(ctx context.Context, branch string) (exGit.PullRequest, error) {
	// ignore ctx for mocking purposes
	// we are ignoring ctx because it is altered by the underlying method and we would have to build one to match
	mg.On("GetPullRequest", branch).Return()
	mg.Called(branch)

	return mg.getPullRequest(ctx, branch)
}

// GetPullRequests calls mg.getPullRequests
func (mg *mockGit) GetPullRequests(ctx context.Context, state string, count int, opts ...exGit.FilterOption) (
	exGit.PullRequests, error) {
	return mg.getPullRequests(ctx, state, count, opts...)
}

// GetMergeability calls mg.getMergeability
func (mg *mockGit) GetMergeability(ctx context.Context, pr exGit.PullRequest) (*bool, error) {
	return mg.getMergeability(ctx, pr)
}

// MergePullRequest calls mg.mergePullRequest
func (mg *mockGit) MergePullRequest(ctx context.Context, pr exGit.PullRequest) (*string, error) {
	return mg.mergePullRequest(ctx, pr)
}

// GetReviews calls mg.getReviews
func (mg *mockGit) GetReviews(ctx context.Context, pr exGit.PullRequest) (exGit.PullRequestReviews, error) {
	return mg.getReviews(ctx, pr)
}

// CreateReview calls mg.createReview
func (mg *mockGit) CreateReview(ctx context.Context, pr exGit.PullRequest, data *models.Review) error {
	return mg.createReview(ctx, pr, data)
}

// DismissApprovalReviews calls mg.dismissApprovalReviews
func (mg *mockGit) DismissApprovalReviews(ctx context.Context, reviews exGit.PullRequestReviews,
	pr exGit.PullRequest) error {
	return mg.dismissApprovalReviews(ctx, reviews, pr)
}

// GetUserLogin calls mg.getUserLogin
func (mg *mockGit) GetUserLogin(ctx context.Context) (*string, error) {
	return mg.getUserLogin(ctx)
}

// GetUserTeams calls mg.getUserTeams
func (mg *mockGit) GetUserTeams(ctx context.Context) (set.Set[string], error) {
	return mg.getUserTeams(ctx)
}

// CreateTag calls mg.createTag
func (mg *mockGit) CreateTag(ctx context.Context, sha string, name string) error {
	return mg.createTag(ctx, sha, name)
}

// GetIdsAndTitles calls mg.getIdsAndTitles
func (mg *mockGit) GetIdsAndTitles(prs exGit.PullRequests) (exGit.IdsAndTitles, error) {
	return mg.getIdsAndTitles(prs)
}

// WithOwner calls mg.withOwner
func (mg *mockGit) WithOwner(owner *string) exGit.FilterOption {
	return mg.withOwner(owner)
}

// IsMerged calls mg.isMerged
func (mg *mockGit) IsMerged(merged *bool) exGit.FilterOption {
	return mg.isMerged(merged)
}

// call is a type used to assist in asserting certain methods/functions were called with the given arguments
type call struct {
	// function name
	name      string
	arguments []interface{}
}

// getStringPointer is a helper function that returns a pointer to the given string
func getStringPointer(target string) *string {
	return &target
}

// setup returns common variables used across many tests
// returns an identifier and a RFCIdentifierCreator
func setup() (string, models.RFCIdentifierCreator) {
	identifier := "test-identifier"
	createRFCIdentifier := func() *string { return &identifier }

	return identifier, createRFCIdentifier
}

// commonAsserter fails the test if any common assertions fail
// This currently is assuming expected and actual to be *strings, and will have to be shifted accordingly in the future
// if necessary to be more open
func commonAsserter(t *testing.T, expected *string, actual *string, expectedErr *string, actualErr error) {
	// check response - avoid dereferencing if nil
	if expected != nil && actual == nil {
		t.Errorf("expected != actual. expected: %s\n actual: %p", *expected, actual)
	} else if expected == nil && actual != nil {
		t.Errorf("expected != actual. expected: %p\n actual: %s", expected, *actual)
	} else if expected != nil && actual != nil && *expected != *actual {
		t.Errorf("expected != actual. expected: %s\n actual: %s", *expected, *actual)
	}

	// check error response - avoid dereferencing if nil
	if expectedErr != nil && actualErr == nil {
		t.Errorf("expected error != actual error. expected: %s\n actual: %p", *expectedErr, actualErr)
	} else if expectedErr == nil && actualErr != nil {
		t.Errorf("expected error != actual error. expected: %p\n actual: %s", expectedErr, actualErr.Error())
	} else if expectedErr != nil && actualErr != nil && *expectedErr != actualErr.Error() {
		t.Errorf("expected error != actual error. expected: %s\n actual: %s", *expectedErr, actualErr.Error())
	}
}

// TestSubmitRequest tests the SubmitRequest function
func TestSubmitRequest(t *testing.T) {
	// initialize
	identifier, createRFCIdentifier := setup()
	CreateRFCIdentifier = createRFCIdentifier

	// initialize test cases
	testCases := []struct {
		mockCreator   gitMockCreator
		data          *models.RFC
		expected      *string
		expectedErr   *string
		expectedCalls []call
	}{
		// failed to create branch
		{
			mockCreator: func() exGit.Git {
				cb := func(ctx context.Context, branch string, baseBranch string) error {
					return fmt.Errorf("create branch error")
				}
				return &mockGit{createBranch: cb}
			},
			data:        &models.RFC{},
			expected:    nil,
			expectedErr: getStringPointer("create branch error"),
			expectedCalls: []call{
				{
					name:      "CreateBranch",
					arguments: []interface{}{identifier, exGit.BASE_BRANCH},
				},
			},
		},
		// failed to create file
		{
			mockCreator: func() exGit.Git {
				cb := func(ctx context.Context, branch string, baseBranch string) error {
					return nil
				}
				cf := func(ctx context.Context, branch string, directory string, data *models.RFC) error {
					return fmt.Errorf("create file error")
				}
				db := func(ctx context.Context, branch string) error {
					return nil
				}
				return &mockGit{createBranch: cb, createFile: cf, deleteBranch: db}
			},
			data: &models.RFC{
				Actions: models.Actions{
					&models.Action{
						ActionType: models.AddAction,
						Target: models.Target{
							TargetType:       models.ItemTarget,
							TargetDescriptor: "entity",
						},
						Data: map[string]interface{}{
							"id": "123",
						},
					},
				},
			},
			expected:    nil,
			expectedErr: getStringPointer("create file error"),
			expectedCalls: []call{
				{
					name: "CreateFile",
					arguments: []interface{}{
						identifier,
						identifier,
						&models.RFC{
							Actions: models.Actions{
								&models.Action{
									ActionType: models.AddAction,
									Target: models.Target{
										TargetType:       models.ItemTarget,
										TargetDescriptor: "entity",
									},
									Data: map[string]interface{}{
										"id": "123",
									},
									Signature: "49991c32fc001d99b9c5908005509686aff6ba7d16a14cd3ecaebc5d6d916cf0",
								},
							},
							Signature: "7fe5c325b99df102515c1f8d5e1cdde084dc9beabec4a346f07dcd90d4ddb4b1",
						},
					},
				},
			},
		},
		// failed create file and delete branch
		{
			mockCreator: func() exGit.Git {
				cb := func(ctx context.Context, branch string, baseBranch string) error {
					return nil
				}
				cf := func(ctx context.Context, branch string, directory string, data *models.RFC) error {
					return fmt.Errorf("create file error")
				}
				db := func(ctx context.Context, branch string) error {
					return fmt.Errorf("delete branch error")
				}
				return &mockGit{createBranch: cb, createFile: cf, deleteBranch: db}
			},
			// already asserted call in test case above
			data:        &models.RFC{},
			expected:    nil,
			expectedErr: getStringPointer("create file error"),
			expectedCalls: []call{
				{
					name:      "DeleteBranch",
					arguments: []interface{}{identifier},
				},
			},
		},
		// failed to create pull request, successfully deleted branch
		{
			mockCreator: func() exGit.Git {
				cb := func(ctx context.Context, branch string, baseBranch string) error {
					return nil
				}
				cf := func(ctx context.Context, branch string, directory string, data *models.RFC) error {
					return nil
				}
				db := func(ctx context.Context, branch string) error {
					return nil
				}
				cpr := func(ctx context.Context, branch string, baseBranch string) error {
					return fmt.Errorf("create pull request error")
				}
				return &mockGit{createBranch: cb, createFile: cf, deleteBranch: db, createPullRequest: cpr}
			},
			data:        &models.RFC{},
			expected:    nil,
			expectedErr: getStringPointer("create pull request error"),
			expectedCalls: []call{
				{
					name:      "CreatePullRequest",
					arguments: []interface{}{identifier, exGit.BASE_BRANCH},
				},
			},
		},
		// failed to create pull request and delete branch
		{
			mockCreator: func() exGit.Git {
				cb := func(ctx context.Context, branch string, baseBranch string) error {
					return nil
				}
				cf := func(ctx context.Context, branch string, directory string, data *models.RFC) error {
					return nil
				}
				db := func(ctx context.Context, branch string) error {
					return fmt.Errorf("delete branch error")
				}
				cpr := func(ctx context.Context, branch string, baseBranch string) error {
					return fmt.Errorf("create pull request error")
				}
				return &mockGit{createBranch: cb, deleteBranch: db, createFile: cf, createPullRequest: cpr}
			},
			data:        &models.RFC{},
			expected:    nil,
			expectedErr: getStringPointer("create pull request error"),
			// calls were already asserted in test cases above
			expectedCalls: []call{},
		},
		// success
		{
			mockCreator: func() exGit.Git {
				cb := func(ctx context.Context, branch string, baseBranch string) error {
					return nil
				}
				db := func(ctx context.Context, branch string) error {
					return nil
				}
				cf := func(ctx context.Context, branch string, directory string, data *models.RFC) error {
					return nil
				}
				cpr := func(ctx context.Context, branch string, baseBranch string) error {
					return nil
				}
				return &mockGit{createBranch: cb, deleteBranch: db, createFile: cf, createPullRequest: cpr}
			},
			data:          &models.RFC{},
			expected:      &identifier,
			expectedErr:   nil,
			expectedCalls: []call{},
		},
	}

	// assert
	for _, testCase := range testCases {
		gitInstance := testCase.mockCreator()

		actual, actualErr := SubmitRequest(context.Background(), gitInstance, testCase.data)

		commonAsserter(t, testCase.expected, actual, testCase.expectedErr, actualErr)
		if len(testCase.expectedCalls) > 0 {
			mgInstance, ok := gitInstance.(*mockGit)
			if !ok {
				t.Errorf("git instance not of type mockGit, which is necessary for mock assertions!")
			} else {
				for _, c := range testCase.expectedCalls {
					mgInstance.AssertCalled(t, c.name, c.arguments...)
				}
			}
		}
	}
}

// TestUpdateRequest tests the UpdateRequest function
func TestUpdateRequest(t *testing.T) {
	// initialize
	identifier, createRFCIdentifier := setup()
	CreateRFCIdentifier = createRFCIdentifier

	// initialize test cases
	testCases := []struct {
		mockCreator   gitMockCreator
		data          *models.Update
		expected      *string
		expectedErr   *string
		expectedCalls []call
	}{
		// failed to get pull request
		{
			mockCreator: func() exGit.Git {
				gpr := func(ctx context.Context, branch string) (exGit.PullRequest, error) {
					return nil, fmt.Errorf("get pull request error")
				}
				return &mockGit{getPullRequest: gpr}
			},
			data:        &models.Update{RFC: &models.RFC{}, RFCIdentifier: identifier},
			expected:    nil,
			expectedErr: getStringPointer("get pull request error"),
			expectedCalls: []call{
				{
					name:      "GetPullRequest",
					arguments: []interface{}{identifier},
				},
			},
		},
		// failed to get RFC contents
		{
			mockCreator: func() exGit.Git {
				gpr := func(ctx context.Context, branch string) (exGit.PullRequest, error) { return nil, nil }
				grfc := func(ctx context.Context, branch string) (*string, *string, error) {
					return nil, nil, fmt.Errorf("get rfc contents error")
				}
				return &mockGit{getPullRequest: gpr, getRFCContents: grfc}
			},
			data:        &models.Update{RFC: &models.RFC{}, RFCIdentifier: identifier},
			expected:    nil,
			expectedErr: getStringPointer("get rfc contents error"),
			expectedCalls: []call{
				{
					name:      "GetRFCContents",
					arguments: []interface{}{identifier},
				},
			},
		},
		// marshal error due to bad data
		{
			mockCreator: func() exGit.Git {
				gpr := func(ctx context.Context, branch string) (exGit.PullRequest, error) { return nil, nil }
				grfc := func(ctx context.Context, branch string) (*string, *string, error) {
					return getStringPointer("junk-data"), getStringPointer("junk-sha"), nil
				}
				return &mockGit{getPullRequest: gpr, getRFCContents: grfc}
			},
			data:          &models.Update{RFC: &models.RFC{}, RFCIdentifier: identifier},
			expected:      nil,
			expectedErr:   getStringPointer("invalid character 'j' looking for beginning of value"),
			expectedCalls: []call{},
		},
		// failed to update file
		{
			mockCreator: func() exGit.Git {
				gpr := func(ctx context.Context, branch string) (exGit.PullRequest, error) { return nil, nil }
				grfc := func(ctx context.Context, branch string) (*string, *string, error) {
					existingRfc := `{
						"actions": [
							{"actionType": "comment", "data": {"test": true}},
							{"actionType": "add", "data": {"test": true}}
						]
					}`
					return &existingRfc, getStringPointer("junk-sha"), nil
				}
				uf := func(ctx context.Context, pr exGit.PullRequest, data *models.RFC) error {
					return fmt.Errorf("error updating file")
				}
				return &mockGit{getPullRequest: gpr, getRFCContents: grfc, updateFile: uf}
			},
			data:        &models.Update{RFC: &models.RFC{}, RFCIdentifier: identifier},
			expected:    nil,
			expectedErr: getStringPointer("error updating file"),
			expectedCalls: []call{
				{
					name: "UpdateFile",
					arguments: []interface{}{
						nil,
						&models.RFC{
							Actions: []*models.Action{
								{
									ActionType: models.CommentAction,
									Data: map[string]interface{}{
										"test": true,
									},
									Signature: "",
								},
							},
							Signature: "a02e316df3bc6f8b3da979fd5cdb5c070962fc03c8fbd46345a7eac682a26f0a",
						},
					},
				},
			},
		},
		// success
		{
			mockCreator: func() exGit.Git {
				gpr := func(ctx context.Context, branch string) (exGit.PullRequest, error) { return nil, nil }
				grfc := func(ctx context.Context, branch string) (*string, *string, error) {
					existingRfc := `{}`
					return &existingRfc, getStringPointer("junk-sha"), nil
				}
				uf := func(ctx context.Context, pr exGit.PullRequest, data *models.RFC) error { return nil }
				gr := func(ctx context.Context, pr exGit.PullRequest) (exGit.PullRequestReviews, error) {
					return nil, nil
				}
				dar := func(ctx context.Context, reviews exGit.PullRequestReviews, pr exGit.PullRequest) error {
					return nil
				}
				return &mockGit{
					getPullRequest:         gpr,
					getRFCContents:         grfc,
					updateFile:             uf,
					getReviews:             gr,
					dismissApprovalReviews: dar,
				}
			},
			data:          &models.Update{RFC: &models.RFC{}, RFCIdentifier: identifier},
			expected:      &identifier,
			expectedErr:   nil,
			expectedCalls: []call{},
		},
	}

	// assert
	for _, testCase := range testCases {
		gitInstance := testCase.mockCreator()

		actual, actualErr := UpdateRequest(context.Background(), gitInstance, testCase.data)

		commonAsserter(t, testCase.expected, actual, testCase.expectedErr, actualErr)
		if len(testCase.expectedCalls) > 0 {
			mgInstance, ok := gitInstance.(*mockGit)
			if !ok {
				t.Errorf("git instance not of type mockGit, which is necessary for mock assertions!")
			} else {
				for _, c := range testCase.expectedCalls {
					mgInstance.AssertCalled(t, c.name, c.arguments...)
				}
			}
		}
	}
}
