#!/usr/bin/env bash
set -euo pipefail

PR_URL="${1:-}"
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "detached")

echo "=== API Review Preflight ==="
echo "Original-Branch: $CURRENT_BRANCH"

if [[ -n "$PR_URL" ]] && [[ "$PR_URL" =~ github\.com.*pull ]]; then
    PR_NUMBER=$(echo "$PR_URL" | grep -oE '[0-9]+$')
    echo "Mode: pr"
    echo "PR: #$PR_NUMBER"

    if ! git diff --quiet || ! git diff --cached --quiet; then
        echo ""
        echo "ERROR: Uncommitted changes detected. Cannot checkout PR."
        git status --porcelain
        exit 1
    fi

    gh pr checkout "$PR_NUMBER"

    echo ""
    echo "=== Changed API Files ==="
    CHANGED_FILES=$(gh pr view "$PR_NUMBER" --json files --jq '.files[].path' \
        | grep '\.go$' \
        | grep -E '/(v1|v1alpha1|v1beta1)/' || true)
else
    echo "Mode: local"

    OPENSHIFT_REMOTE=""
    for remote in $(git remote); do
        remote_url=$(git remote get-url "$remote" 2>/dev/null || echo "")
        if [[ "$remote_url" =~ github\.com[/:]openshift/api(\.git)?$ ]]; then
            OPENSHIFT_REMOTE="$remote"
            break
        fi
    done

    if [[ -z "$OPENSHIFT_REMOTE" ]]; then
        git remote add upstream https://github.com/openshift/api.git 2>/dev/null || true
        OPENSHIFT_REMOTE="upstream"
    fi

    echo "Remote: $OPENSHIFT_REMOTE"
    git fetch "$OPENSHIFT_REMOTE" master 2>/dev/null

    echo ""
    echo "=== Changed API Files ==="
    CHANGED_FILES=$( {
        git diff --name-only "$OPENSHIFT_REMOTE/master...HEAD" || true
        git diff --cached --name-only || true
        git diff --name-only || true
    } | grep '\.go$' | grep -E '/(v1|v1alpha1|v1beta1)/' | sort -u || true)
fi

if [[ -z "$CHANGED_FILES" ]]; then
    echo "NO_API_FILES_CHANGED"
    exit 0
fi

echo "$CHANGED_FILES"

echo ""
echo "=== Lint Results ==="
if [[ "${EVAL_SKIP_LINT:-}" == "1" ]]; then
    echo "SKIPPED (EVAL_SKIP_LINT=1)"
else
    LINT_OUTPUT=""
    if LINT_OUTPUT=$(make lint 2>&1); then
        echo "PASSED"
    else
        echo "FAILED"
        echo "$LINT_OUTPUT"
    fi
fi

echo ""
echo "=== Preflight Complete ==="
