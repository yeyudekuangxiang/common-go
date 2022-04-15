package system

import "time"

type MergeRequestApprovalState struct {
	ApprovalRulesOverwritten bool `json:"approval_rules_overwritten"`
	Rules                    []struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		RuleType          string `json:"rule_type"`
		EligibleApprovers []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Username  string `json:"username"`
			State     string `json:"state"`
			AvatarURL string `json:"avatar_url"`
			WebURL    string `json:"web_url"`
		} `json:"eligible_approvers"`
		ApprovalsRequired int `json:"approvals_required"`
		Users             []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Username  string `json:"username"`
			State     string `json:"state"`
			AvatarURL string `json:"avatar_url"`
			WebURL    string `json:"web_url"`
		} `json:"users"`
		Groups               []interface{} `json:"groups"`
		ContainsHiddenGroups bool          `json:"contains_hidden_groups"`
		ApprovedBy           []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Username  string `json:"username"`
			State     string `json:"state"`
			AvatarURL string `json:"avatar_url"`
			WebURL    string `json:"web_url"`
		} `json:"approved_by"`
		SourceRule interface{} `json:"source_rule"`
		Approved   bool        `json:"approved"`
		Overridden bool        `json:"overridden"`
	} `json:"rules"`
}
type MergeRequest struct {
	ID           int       `json:"id"`
	Iid          int       `json:"iid"`
	ProjectID    int       `json:"project_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TargetBranch string    `json:"target_branch"`
	SourceBranch string    `json:"source_branch"`
	Upvotes      int       `json:"upvotes"`
	Downvotes    int       `json:"downvotes"`
	Author       struct {
		ID        int         `json:"id"`
		Name      string      `json:"name"`
		Username  string      `json:"username"`
		State     string      `json:"state"`
		AvatarURL interface{} `json:"avatar_url"`
		WebURL    string      `json:"web_url"`
	} `json:"author"`
	User struct {
		CanMerge bool `json:"can_merge"`
	} `json:"user"`
	Assignee struct {
		ID        int         `json:"id"`
		Name      string      `json:"name"`
		Username  string      `json:"username"`
		State     string      `json:"state"`
		AvatarURL interface{} `json:"avatar_url"`
		WebURL    string      `json:"web_url"`
	} `json:"assignee"`
	Assignees []struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		ID        int    `json:"id"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"assignees"`
	Reviewers []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"reviewers"`
	SourceProjectID int      `json:"source_project_id"`
	TargetProjectID int      `json:"target_project_id"`
	Labels          []string `json:"labels"`
	Draft           bool     `json:"draft"`
	WorkInProgress  bool     `json:"work_in_progress"`
	Milestone       struct {
		ID          int       `json:"id"`
		Iid         int       `json:"iid"`
		ProjectID   int       `json:"project_id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		State       string    `json:"state"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		DueDate     string    `json:"due_date"`
		StartDate   string    `json:"start_date"`
		WebURL      string    `json:"web_url"`
	} `json:"milestone"`
	MergeWhenPipelineSucceeds bool        `json:"merge_when_pipeline_succeeds"`
	MergeStatus               string      `json:"merge_status"`
	MergeError                interface{} `json:"merge_error"`
	Sha                       string      `json:"sha"`
	MergeCommitSha            interface{} `json:"merge_commit_sha"`
	SquashCommitSha           interface{} `json:"squash_commit_sha"`
	UserNotesCount            int         `json:"user_notes_count"`
	DiscussionLocked          interface{} `json:"discussion_locked"`
	ShouldRemoveSourceBranch  bool        `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch   bool        `json:"force_remove_source_branch"`
	AllowCollaboration        bool        `json:"allow_collaboration"`
	AllowMaintainerToPush     bool        `json:"allow_maintainer_to_push"`
	WebURL                    string      `json:"web_url"`
	References                struct {
		Short    string `json:"short"`
		Relative string `json:"relative"`
		Full     string `json:"full"`
	} `json:"references"`
	TimeStats struct {
		TimeEstimate        int         `json:"time_estimate"`
		TotalTimeSpent      int         `json:"total_time_spent"`
		HumanTimeEstimate   interface{} `json:"human_time_estimate"`
		HumanTotalTimeSpent interface{} `json:"human_total_time_spent"`
	} `json:"time_stats"`
	Squash       bool   `json:"squash"`
	Subscribed   bool   `json:"subscribed"`
	ChangesCount string `json:"changes_count"`
	MergedBy     struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"merged_by"`
	MergeUser struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"merge_user"`
	MergedAt                    time.Time   `json:"merged_at"`
	ClosedBy                    interface{} `json:"closed_by"`
	ClosedAt                    interface{} `json:"closed_at"`
	LatestBuildStartedAt        time.Time   `json:"latest_build_started_at"`
	LatestBuildFinishedAt       time.Time   `json:"latest_build_finished_at"`
	FirstDeployedToProductionAt interface{} `json:"first_deployed_to_production_at"`
	Pipeline                    struct {
		ID     int    `json:"id"`
		Sha    string `json:"sha"`
		Ref    string `json:"ref"`
		Status string `json:"status"`
		WebURL string `json:"web_url"`
	} `json:"pipeline"`
	DiffRefs struct {
		BaseSha  string `json:"base_sha"`
		HeadSha  string `json:"head_sha"`
		StartSha string `json:"start_sha"`
	} `json:"diff_refs"`
	DivergedCommitsCount int  `json:"diverged_commits_count"`
	RebaseInProgress     bool `json:"rebase_in_progress"`
	FirstContribution    bool `json:"first_contribution"`
	TaskCompletionStatus struct {
		Count          int `json:"count"`
		CompletedCount int `json:"completed_count"`
	} `json:"task_completion_status"`
	HasConflicts                bool `json:"has_conflicts"`
	BlockingDiscussionsResolved bool `json:"blocking_discussions_resolved"`
}
