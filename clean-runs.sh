#!/bin/bash
# Clean "failure", "cancelled", or "skipped" runs from all workflows in a repository
org="IDEXX"
repo="poe-mobile-service-graphql"

# Get all workflow IDs
workflow_ids=($(gh api repos/$org/$repo/actions/workflows --paginate | jq '.workflows[].id'))

for workflow_id in "${workflow_ids[@]}"
do
  echo "Listing failed runs for the workflow ID $workflow_id"
  # Get run IDs with status "failure", "cancelled", or "skipped"
  run_ids=( $(gh api repos/$org/$repo/actions/workflows/$workflow_id/runs --paginate | jq '.workflow_runs[] | select(.conclusion == "skipped") | .id') )
  # run_ids=( $(gh api repos/$org/$repo/actions/workflows/$workflow_id/runs --paginate | jq '.workflow_runs[] | select(.conclusion == "failure") | .id') )
  # run_ids=( $(gh api repos/$org/$repo/actions/workflows/$workflow_id/runs --paginate | jq '.workflow_runs[] | select(.conclusion == "cancelled") | .id') )
  for run_id in "${run_ids[@]}"
  do
    echo "Deleting Run ID $run_id"
    gh api repos/$org/$repo/actions/runs/$run_id -X DELETE >/dev/null
  done
done