package main

import (
	"context"
	"fmt"
	"os"
)

const successCode = 0
const errorCode = 1

func printError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

func printAuthHelp(project string) {
	fmt.Fprintf(
		os.Stderr,
		"Make sure you are authenticated with Google Cloud by running "+
			"following command:\n\n",
	)
	fmt.Fprintf(
		os.Stderr,
		"  gcloud auth application-default set-quota-project %s\n",
		project,
	)
	fmt.Fprintf(
		os.Stderr,
		"  gcloud auth application-default login\n\n",
	)
}

func reviewCode(ctx context.Context, cfg *Config) int {
	patch, err := fetchPatch(ctx, cfg.BaseCommit, cfg.ContextLines)
	if err != nil {
		printError(err)
		return errorCode
	}

	reviewPrompt := createReviewPrompt(patch)

	if cfg.Verbose {
		fmt.Printf(
			"=== REVIEW PROMPT START ===\n\n%s\n\n=== REVIEW PROMPT END ===\n\n",
			reviewPrompt,
		)
	}

	client := GeminiClient{Project: cfg.GCPProject, Location: cfg.Location, Model: cfg.Model}
	reviewText, err := client.generate(ctx, reviewPrompt)
	if err != nil {
		printError(err)
		printAuthHelp(cfg.GCPProject)
		return errorCode
	}

	fmt.Println(reviewText)

	untracked, err := untrackedFiles(ctx)
	if err == nil && len(untracked) > 0 {
		fmt.Println("\nWarning: The following untracked files were ignored:")
		for _, file := range untracked {
			fmt.Printf("  - %s\n", file)
		}
	}

	return successCode
}

func main() {
	cfg, err := Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(errorCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	os.Exit(reviewCode(ctx, cfg))
}
