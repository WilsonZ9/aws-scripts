#!/bin/bash
# Clean all runs from a specific workflow named "[Auto] Run tests" in a repository
org="IDEXX"
repo="poe-appointments"
workflow_name="[Auto] Run tests"

# Get the workflow ID for the specified workflow name
workflow_id=$(gh api repos/$org/$repo/actions/workflows --paginate | jq -r --arg name "$workflow_name" '.workflows[] | select(.name == $name) | .id')

if [ -z "$workflow_id" ]; then
  echo "Workflow named '$workflow_name' not found in the repository '$repo'."
  exit 1
fi

echo "Listing all runs for the workflow ID $workflow_id"

# Get all run IDs for the specified workflow
run_ids=( $(gh api repos/$org/$repo/actions/workflows/$workflow_id/runs --paginate | jq -r '.workflow_runs[].id') )

for run_id in "${run_ids[@]}"
do
  echo "Deleting Run ID $run_id"
  gh api repos/$org/$repo/actions/runs/$run_id -X DELETE >/dev/null
done

echo "All runs for the workflow '$workflow_name' have been deleted."