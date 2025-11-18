package query

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type issueParamsErr struct {
	history    bool
	watching   bool
	resolution bool
	issueType  bool
	labels     bool
	status     bool
}

type issueFlagParser struct {
	err           issueParamsErr
	noHistory     bool
	noWatching    bool
	orderDesc     bool
	emptyType     bool
	labels        []string
	status        []string
	withCreated   bool
	withUpdated   bool
	created       string
	updated       string
	createdAfter  string
	createdBefore string
	updatedAfter  string
	updatedBefore string
	jql           string
	orderBy       string
}

func (tfp *issueFlagParser) GetBool(name string) (bool, error) {
	if tfp.err.history && name == "history" {
		return false, fmt.Errorf("oops! couldn't fetch history flag")
	}
	if tfp.err.watching && name == "watching" {
		return false, fmt.Errorf("oops! couldn't fetch watching flag")
	}
	if tfp.noHistory && name == "history" {
		return false, nil
	}
	if tfp.noWatching && name == "watching" {
		return false, nil
	}
	if tfp.orderDesc && name == "reverse" {
		return false, nil
	}
	return true, nil
}

//nolint:gocyclo
func (tfp *issueFlagParser) GetString(name string) (string, error) {
	if tfp.err.resolution && name == "resolution" {
		return "", fmt.Errorf("oops! couldn't fetch resolution flag")
	}
	if tfp.err.issueType && name == "type" {
		return "", fmt.Errorf("oops! couldn't fetch type flag")
	}
	if tfp.created != "" && name == "created" {
		return tfp.created, nil
	}
	if tfp.updated != "" && name == "updated" {
		return tfp.updated, nil
	}
	if tfp.emptyType && name == "type" {
		return "", nil
	}
	if name == "jql" {
		return tfp.jql, nil
	}
	if tfp.orderBy == "" && name == "order-by" {
		return "created", nil
	}
	if strings.HasPrefix(name, "created") {
		if tfp.withCreated {
			switch name {
			case "created-after":
				return tfp.createdAfter, nil
			case "created-before":
				return tfp.createdBefore, nil
			}
		}
		return "", nil
	}
	if strings.HasPrefix(name, "updated") {
		if tfp.withUpdated {
			switch name {
			case "updated-after":
				return tfp.updatedAfter, nil
			case "updated-before":
				return tfp.updatedBefore, nil
			}
		}
		return "", nil
	}
	if name == "paginate" {
		return "", nil
	}
	return "test", nil
}

func (tfp *issueFlagParser) GetStringArray(name string) ([]string, error) {
	if tfp.err.labels && name == "label" {
		return []string{}, fmt.Errorf("oops! couldn't fetch label flag")
	}
	if tfp.err.status && name == "status" {
		return []string{}, fmt.Errorf("oops! couldn't fetch status flag")
	}
	if name == "status" {
		return tfp.status, nil
	}
	return tfp.labels, nil
}

func (*issueFlagParser) GetStringToString(string) (map[string]string, error) { return nil, nil }
func (*issueFlagParser) GetUint(string) (uint, error)                        { return 100, nil }
func (*issueFlagParser) Set(string, string) error                            { return nil }

func TestIssueGet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		initialize func() *Issue
		expected   string
	}{
		{
			name: "query with default parameters",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY lastViewed ASC`,
		},
		{
			name: "query without issue history parameter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{noHistory: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY created ASC`,
		},
		{
			name: "query only with fields filter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{noHistory: true, noWatching: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY created ASC`,
		},
		{
			name: "query with error when fetching history flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{err: issueParamsErr{
					history: true,
				}})
				assert.Error(t, err)
				return i
			},
			expected: "",
		},
		{
			name: "query with error when fetching watching flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{err: issueParamsErr{
					watching: true,
				}})
				assert.Error(t, err)
				return i
			},
			expected: "",
		},
		{
			name: "query with error when fetching resolution flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{err: issueParamsErr{
					resolution: true,
				}})
				assert.Error(t, err)
				return i
			},
			expected: "",
		},
		{
			name: "query with error when fetching labels flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{err: issueParamsErr{
					labels: true,
				}})
				assert.Error(t, err)
				return i
			},
			expected: "",
		},
		{
			name: "query with error when fetching status flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{err: issueParamsErr{
					status: true,
				}})
				assert.Error(t, err)
				return i
			},
			expected: "",
		},
		{
			name: "query with error when fetching type flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{err: issueParamsErr{
					issueType: true,
				}})
				assert.Error(t, err)
				return i
			},
			expected: "",
		},
		{
			name: "query without issue type flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{emptyType: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" ORDER BY lastViewed ASC`,
		},
		{
			name: "query with reverse set to true",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{orderDesc: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY lastViewed DESC`,
		},
		{
			name: "query with labels",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{labels: []string{"first", "second", "third"}})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" AND labels IN ("first", "second", "third") ORDER BY lastViewed ASC`,
		},
		{
			name: "query with status",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{status: []string{"first", "second", "~third"}})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" AND status IN ("first", "second") AND status NOT IN ("third") ORDER BY lastViewed ASC`,
		},
		{
			name: "query with created and updated today filter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{created: "today", updated: "today"})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" AND createdDate>=startOfDay() AND updatedDate>=startOfDay() ORDER BY lastViewed ASC`,
		},
		{
			name: "query with created and updated week filter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{created: "week", updated: "week"})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" AND createdDate>=startOfWeek() AND updatedDate>=startOfWeek() ORDER BY lastViewed ASC`,
		},
		{
			name: "query with created and updated month filter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{created: "month", updated: "month"})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" AND createdDate>=startOfMonth() AND updatedDate>=startOfMonth() ORDER BY lastViewed ASC`,
		},
		{
			name: "query with created and updated year filter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{created: "year", updated: "year"})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" AND createdDate>=startOfYear() AND updatedDate>=startOfYear() ORDER BY lastViewed ASC`,
		},
		{
			name: "query with created and updated filter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{created: "2020-12-31", updated: "2020-12-31"})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" AND component="test" ` +
				`AND parent="test" AND createdDate>="2020-12-31" AND createdDate<"2021-01-01" AND updatedDate>="2020-12-31" AND updatedDate<"2021-01-01" ` +
				`ORDER BY lastViewed ASC`,
		},
		{
			name: "created and updated filter with incorrect date format",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{created: "2020-15-31", updated: "2020-12-31 10:30:30"})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" AND createdDate>="2020-15-31" AND updatedDate>="2020-12-31 10:30:30" ORDER BY lastViewed ASC`,
		},
		{
			name: "query with created-after and created-before filter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{createdAfter: "2020-12-01", createdBefore: "2020-12-31", withCreated: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" AND createdDate>"2020-12-01" AND createdDate<"2020-12-31" ORDER BY lastViewed ASC`,
		},
		{
			name: "query with updated-after and updated-before filter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{updatedAfter: "2020-12-01", updatedBefore: "2020-12-31", withUpdated: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" AND updatedDate>"2020-12-01" AND updatedDate<"2020-12-31" ORDER BY lastViewed ASC`,
		},
		{
			name: "created and updated flags gets precedence",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{
					created:       "2020-11-01",
					updated:       "-10d",
					createdAfter:  "2020-12-01",
					updatedBefore: "2020-12-31",
					withCreated:   true,
					withUpdated:   true,
				})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`AND component="test" AND parent="test" AND createdDate>="2020-11-01" AND createdDate<"2020-11-02" AND updatedDate>="-10d" ` +
				`ORDER BY lastViewed ASC`,
		},
		{
			name: "created order gets priority over updated flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{
					created:       "2020-11-01",
					updated:       "-10d",
					createdAfter:  "2020-12-01",
					updatedBefore: "2020-12-31",
					withCreated:   true,
					withUpdated:   true,
					noHistory:     true,
				})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN watchedIssues() AND type="test" AND resolution="test" ` +
				`AND priority="test" AND reporter="test" AND assignee="test" AND component="test" ` +
				`AND parent="test" AND createdDate>="2020-11-01" AND createdDate<"2020-11-02" AND updatedDate>="-10d" ` +
				`ORDER BY created ASC`,
		},
		{
			name: "it orders by updated if only updated flags are present",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{
					updatedAfter:  "2020-11-31",
					updatedBefore: "2020-12-31",
					withCreated:   false,
					withUpdated:   true,
					noHistory:     true,
				})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND issue IN watchedIssues() AND type="test" AND resolution="test" ` +
				`AND priority="test" AND reporter="test" AND assignee="test" AND component="test" ` +
				`AND parent="test" AND updatedDate>"2020-11-31" AND updatedDate<"2020-12-31" ` +
				`ORDER BY updated ASC`,
		},
		{
			name: "query with jql parameter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{jql: "summary ~ cli OR x = y"})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND summary ~ cli OR x = y AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY lastViewed ASC`,
		},
		{
			name: "query with jql parameter containing ORDER BY should bypass default ordering",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{jql: "summary ~ cli ORDER BY priority DESC"})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND summary ~ cli AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY priority DESC`,
		},
		{
			name: "query with jql parameter containing ORDER BY (case insensitive) should bypass default ordering",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{jql: "status = Done order by updated ASC", noHistory: true, noWatching: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND status = Done AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY updated ASC`,
		},
		{
			name: "query with jql parameter containing ORDER BY with multiple spaces should bypass default ordering",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{jql: "assignee = currentUser() ORDER  BY  rank ASC", noHistory: true, noWatching: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND assignee = currentUser() AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY rank ASC`,
		},
		{
			name: "query with jql containing ORDER BY in different case variations",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{jql: "status != Closed OrDeR bY created DESC", noHistory: true, noWatching: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND status != Closed AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY created DESC`,
		},
		{
			name: "query with jql containing multiple ORDER BY conditions",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{jql: "status = Done ORDER BY priority DESC, updated DESC", noHistory: true, noWatching: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND status = Done AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY priority DESC, updated DESC`,
		},
		{
			name: "query with invalid jql (ORDER BY in middle) - will produce malformed query that Jira will reject",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &issueFlagParser{jql: "status = Done ORDER BY priority DESC AND assignee = currentUser()", noHistory: true, noWatching: true})
				assert.NoError(t, err)
				return i
			},
			expected: `project="TEST" AND status = Done AND ` +
				`type="test" AND resolution="test" AND priority="test" AND reporter="test" ` +
				`AND assignee="test" AND component="test" AND parent="test" ORDER BY priority DESC AND assignee = currentUser()`,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			q := tc.initialize()
			if q != nil {
				assert.Equal(t, tc.expected, q.Get())
			}
		})
	}
}

func TestHasOrderBy(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		jql      string
		expected bool
	}{
		{
			name:     "empty string should return false",
			jql:      "",
			expected: false,
		},
		{
			name:     "jql without ORDER BY should return false",
			jql:      "status = Done AND assignee = currentUser()",
			expected: false,
		},
		{
			name:     "jql with ORDER BY uppercase should return true",
			jql:      "status = Done ORDER BY created DESC",
			expected: true,
		},
		{
			name:     "jql with order by lowercase should return true",
			jql:      "status = Done order by created asc",
			expected: true,
		},
		{
			name:     "jql with Order By mixed case should return true",
			jql:      "status = Done Order By priority",
			expected: true,
		},
		{
			name:     "jql with ORDER BY and multiple spaces should return true",
			jql:      "status = Done ORDER  BY  rank DESC",
			expected: true,
		},
		{
			name:     "jql with ORDER BY and tabs should return true",
			jql:      "status = Done ORDER	BY created",
			expected: true,
		},
		{
			name:     "partial match 'ORDER' only should return false",
			jql:      "summary ~ 'ORDER' AND status = Done",
			expected: false,
		},
		{
			name:     "partial match 'BY' only should return false",
			jql:      "summary ~ 'sorted BY' AND status = Done",
			expected: false,
		},
		{
			name:     "ORDER BY at the beginning should return true, even though it's invalid JQL",
			jql:      "ORDER BY created DESC",
			expected: true,
		},
		{
			name:     "complex jql with ORDER BY should return true",
			jql:      "project = TEST AND (status = 'In Progress' OR status = Done) AND assignee = currentUser() ORDER BY updated DESC, priority ASC",
			expected: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := hasOrderBy(tc.jql)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExtractOrderBy(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		jql             string
		expectedJQL     string
		expectedOrderBy string
	}{
		{
			name:            "empty string should return empty",
			jql:             "",
			expectedJQL:     "",
			expectedOrderBy: "",
		},
		{
			name:            "jql without ORDER BY should return original jql",
			jql:             "status = Done AND assignee = currentUser()",
			expectedJQL:     "status = Done AND assignee = currentUser()",
			expectedOrderBy: "",
		},
		{
			name:            "jql with ORDER BY uppercase",
			jql:             "status = Done ORDER BY created DESC",
			expectedJQL:     "status = Done",
			expectedOrderBy: "created DESC",
		},
		{
			name:            "jql with order by lowercase",
			jql:             "status = Done order by priority asc",
			expectedJQL:     "status = Done",
			expectedOrderBy: "priority asc",
		},
		{
			name:            "jql with Order By mixed case",
			jql:             "assignee = currentUser() Order By rank",
			expectedJQL:     "assignee = currentUser()",
			expectedOrderBy: "rank",
		},
		{
			name:            "jql with ORDER BY and multiple spaces",
			jql:             "status = Done ORDER  BY  rank DESC",
			expectedJQL:     "status = Done",
			expectedOrderBy: "rank DESC",
		},
		{
			name:            "complex jql with ORDER BY",
			jql:             "project = TEST AND (status = 'In Progress' OR status = Done) ORDER BY updated DESC, priority ASC",
			expectedJQL:     "project = TEST AND (status = 'In Progress' OR status = Done)",
			expectedOrderBy: "updated DESC, priority ASC",
		},
		{
			name:            "jql with mixed case ORDER BY variations",
			jql:             "status != Closed OrDeR bY created DESC",
			expectedJQL:     "status != Closed",
			expectedOrderBy: "created DESC",
		},
		{
			name:            "simple ORDER BY only (invalid JQL, but gracefully handled)",
			jql:             "ORDER BY created",
			expectedJQL:     "",
			expectedOrderBy: "created",
		},
		{
			name:            "multiple ORDER BY conditions",
			jql:             "status = Done ORDER BY priority DESC, updated DESC",
			expectedJQL:     "status = Done",
			expectedOrderBy: "priority DESC, updated DESC",
		},
		{
			name:            "multiple ORDER BY conditions with multiple fields",
			jql:             "project = TEST AND assignee = currentUser() ORDER BY rank ASC, priority DESC, created ASC",
			expectedJQL:     "project = TEST AND assignee = currentUser()",
			expectedOrderBy: "rank ASC, priority DESC, created ASC",
		},
		{
			name:            "ORDER BY in middle of query (invalid JQL) - extracts everything after first ORDER BY",
			jql:             "status = Done ORDER BY priority DESC AND assignee = currentUser()",
			expectedJQL:     "status = Done",
			expectedOrderBy: "priority DESC AND assignee = currentUser()",
		},
		{
			name:            "multiple ORDER BY clauses (invalid JQL) - extracts everything after first ORDER BY",
			jql:             "status = Done ORDER BY priority DESC ORDER BY updated ASC",
			expectedJQL:     "status = Done",
			expectedOrderBy: "priority DESC ORDER BY updated ASC",
		},
		{
			name:            "ORDER BY in middle with more text (invalid JQL) - extracts everything after first ORDER BY",
			jql:             "status = Done ORDER BY priority AND assignee = currentUser() ORDER BY rank ASC",
			expectedJQL:     "status = Done",
			expectedOrderBy: "priority AND assignee = currentUser() ORDER BY rank ASC",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jql, orderBy := extractOrderBy(tc.jql)
			assert.Equal(t, tc.expectedJQL, jql, "JQL part mismatch")
			assert.Equal(t, tc.expectedOrderBy, orderBy, "ORDER BY part mismatch")
		})
	}
}
