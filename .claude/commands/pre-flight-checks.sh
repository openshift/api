#!/bin/bash

# Save current branch
CURRENT_BRANCH=$(git branch --show-current)
echo "📍 Current branch: $CURRENT_BRANCH"

# Check if an argument was provided
if [ $# -gt 0 ]; then
    REVIEW_MODE="pr"

    # We expect the argument to be either a PR number, or the URL of a PR,
    # which will end in the PR number
    PR_NUMBER=$(echo "$1" | grep -oE '[0-9]+$')
    echo "🔍 PR review mode: Reviewing PR #$PR_NUMBER"

    # For PR review, check for uncommitted changes
    if ! git diff --quiet || ! git diff --cached --quiet; then
        echo "❌ ERROR: Uncommitted changes detected. Cannot proceed with PR review."
        echo "Please commit or stash your changes before running the API review."
        git status --porcelain
        exit 1
    fi
    echo "✅ No uncommitted changes detected. Safe to proceed with PR review."

    # PR Review: Checkout the PR and get changed files
    echo "🔄 Checking out PR #$PR_NUMBER..."
    gh pr checkout "$PR_NUMBER"

    echo "📁 Analyzing changed files in PR..."
    CHANGED_FILES=$(gh pr view "$PR_NUMBER" --json files --jq '.files[].path' | grep '\.go$' | grep -E '/(v1|v1alpha1|v1beta1)/')
else
    REVIEW_MODE="local"
    echo "🔍 Local review mode: Reviewing local changes against upstream master"

    # Find a remote pointing to openshift/api repository
    OPENSHIFT_REMOTE=""
    for remote in $(git remote); do
        remote_url=$(git remote get-url "$remote" 2>/dev/null || echo "")
        if [[ "$remote_url" =~ github\.com[/:]openshift/api(\.git)?$ ]]; then
            OPENSHIFT_REMOTE="$remote"
            echo "✅ Found OpenShift API remote: '$remote' -> $remote_url"
            break
        fi
    done

    # If no existing remote found, add upstream
    if [ -z "$OPENSHIFT_REMOTE" ]; then
        echo "⚠️  No remote pointing to openshift/api found. Adding upstream remote..."
        git remote add upstream https://github.com/openshift/api.git 2>&1 || true
        OPENSHIFT_REMOTE="upstream"
    fi

    # Local Review: Get changed files compared to openshift remote master
    echo "📁 Analyzing locally changed files compared to $OPENSHIFT_REMOTE/master..."
    CHANGED_FILES=$(git diff --name-only "$OPENSHIFT_REMOTE/master...HEAD" | grep '\.go$' | grep -E '/(v1|v1alpha1|v1beta1)/' || true)

    # Also include staged changes
    STAGED_FILES=$(git diff --cached --name-only | grep '\.go$' | grep -E '/(v1|v1alpha1|v1beta1)/' || true)
    if [ -n "$STAGED_FILES" ]; then
        CHANGED_FILES=$(echo -e "$CHANGED_FILES\n$STAGED_FILES" | sort -u)
    fi
fi

echo "Changed API files:"
echo "$CHANGED_FILES"

if [ -z "$CHANGED_FILES" ]; then
    echo "ℹ️  No API files changed. Nothing to review."
    if [ "$REVIEW_MODE" = "pr" ]; then
        git checkout "$CURRENT_BRANCH"
    fi
    exit 0
fi

# Count the files
FILE_COUNT=$(echo "$CHANGED_FILES" | wc -l)
echo ""
echo "📊 Total API files changed: $FILE_COUNT"

echo "⏳ Running linting checks on changes..."
make lint

if [ $? -ne 0 ]; then
    echo "❌ Linting checks failed. Please fix the issues before proceeding."
    if [ "$REVIEW_MODE" = "pr" ]; then
        echo "🔄 Switching back to original branch: $CURRENT_BRANCH"
        git checkout "$CURRENT_BRANCH"
    fi
    exit 1
fi

echo "✅ Linting checks passed."

# Export variables to restore the environment after review
echo "REVIEW_MODE=$REVIEW_MODE" > /tmp/api_review_vars.sh
echo "CURRENT_BRANCH=$CURRENT_BRANCH" >> /tmp/api_review_vars.sh
