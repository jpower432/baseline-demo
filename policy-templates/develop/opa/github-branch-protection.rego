package branch_protection

# Assisted by: Gemini 2.5 Flash

default allow := false

allow if {
	pull_request_rule_meets_min_approvals
}

# Helper rule to check if a pull_request rule exists
has_pull_request_rule if {
	input.values # This explicitly guards against 'input.values' being null/undefined/false
	some i
	input.values[i].type == "pull_request"
}

# Check if a pull_request rule exists and meets the minimum approval requirement
pull_request_rule_meets_min_approvals if {
	some i
	rule = input.values[i]
	rule.type == "pull_request"
	required_count := data.main_branch_min_approvals
	rule.parameters.required_approving_review_count >= required_count
}

evaluation_resource_id := input.id
evaluation_resource_name := input.name
evaluation_resource_type := "resource"
policy_id := "github-branch-protection"

violation[msg] if {
	not input.values # If the 'values' array itself is missing or empty
	msg = "No branch protection rules found in the input 'values' array."
}

violation[msg] if {
	input.values # Ensure 'values' array exists for the check
	not has_pull_request_rule
	msg = "No 'pull_request' rule found for the main branch in the input's 'values' array."
}

violation[msg] if {
	required_count := data.main_branch_min_approvals
	not pull_request_rule_meets_min_approvals
	msg = sprintf("Branch protection for 'main' requires pull request reviews but has less than the configured minimum of %v required approving reviews.", [required_count])
}
