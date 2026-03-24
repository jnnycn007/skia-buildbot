package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReviewPrompt(t *testing.T) {
	patch := "diff --git a/file.txt b/file.txt\n+added"
	prompt := createReviewPrompt(patch)
	assert.Contains(t, prompt, patch)
}
