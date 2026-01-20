/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	commitType string
	useEmoji   bool
	dryRun     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "commitz",
	Short: "Smart commit message generator",
	Run: func(cmd *cobra.Command, args []string) {
		generateCommitMessage()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&commitType, "type", "t", "", "Commit type (feat, fix, docs, refactor, test, chore)")
	rootCmd.PersistentFlags().BoolVarP(&useEmoji, "emoji", "e", false, "Add emoji to commit message")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview commit message without committing")
}

func generateCommitMessage() {
	// Get staged changes
	diffBytes, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		color.Red("Error getting git diff: %v", err)
		fmt.Println("Make sure you are in a git repository and have staged changes.")
		os.Exit(1)
	}

	diffStr := string(diffBytes)
	if len(diffStr) == 0 {
		color.Yellow("No staged changes found.")
		fmt.Println("Please stage your changes with 'git add' before generating a commit message.")
		os.Exit(0)
	}

	// Determine commit type based on diff content
	suggestedType := detectCommitType(diffStr)
	if commitType != "" {
		suggestedType = commitType
	}

	// Extract scope from branch name
	scope := extractScopeFromBranch()

	// Generate emoji if requested
	emoji := getEmojiForType(suggestedType)

	// Generate summary
	summary := generateSummary(diffStr)

	// Build commit message
	message := buildCommitMessage(emoji, suggestedType, scope, summary)

	// Display suggested message
	displaySuggestedMessage(message)

	// Add optional description
	message = addDescription(message)

	// Handle dry-run
	if dryRun {
		color.Yellow("\n[DRY RUN] Commit not created")
		fmt.Println("\nProposed commit message:")
		fmt.Println(color.CyanString(message))
		return
	}

	// Confirm and commit
	if confirmCommit() {
		executeCommit(message)
	} else {
		color.Yellow("Commit cancelled.")
	}
}

func detectCommitType(diff string) string {
	diffLower := strings.ToLower(diff)

	if strings.Contains(diff, "test/") || strings.Contains(diff, "_test.go") {
		return "test"
	}
	if strings.Contains(diff, "README") || strings.Contains(diff, ".md") || strings.Contains(diff, "docs/") {
		return "docs"
	}
	if strings.Contains(diffLower, "fix") || strings.Contains(diffLower, "bug") {
		return "fix"
	}
	if strings.Contains(diffLower, "feat") || strings.Contains(diffLower, "add ") || strings.Contains(diffLower, "new ") {
		return "feat"
	}

	return "chore"
}

func extractScopeFromBranch() string {
	branchBytes, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return ""
	}

	branchName := strings.TrimSpace(string(branchBytes))
	if strings.Contains(branchName, "/") {
		parts := strings.SplitN(branchName, "/", 2)
		if len(parts) > 1 {
			return strings.TrimSpace(parts[0])
		}
	}

	return ""
}

func getEmojiForType(commitType string) string {
	if !useEmoji {
		return ""
	}

	emojiMap := map[string]string{
		"feat":     "✨ ",
		"fix":      "🐛 ",
		"docs":     "📝 ",
		"refactor": "♻️ ",
		"test":     "✅ ",
		"chore":    "🧹 ",
		"style":    "💄 ",
		"perf":     "⚡ ",
	}

	if emoji, ok := emojiMap[commitType]; ok {
		return emoji
	}

	return ""
}

func generateSummary(diff string) string {
	// TODO: Implement intelligent summary generation
	// For now, return a placeholder
	return "update changes"
}

func buildCommitMessage(emoji, commitType, scope, summary string) string {
	if scope != "" {
		return fmt.Sprintf("%s%s(%s): %s", emoji, commitType, scope, summary)
	}
	return fmt.Sprintf("%s%s: %s", emoji, commitType, summary)
}

func displaySuggestedMessage(message string) {
	fmt.Println()
	color.Green("Suggested commit message:")
	fmt.Printf("  %s\n", color.GreenString(message))
}

func addDescription(message string) string {
	fmt.Println("\nAdd optional description (press Enter twice to finish):")

	scanner := bufio.NewScanner(os.Stdin)
	var bodyLines []string
	emptyLineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			emptyLineCount++
			if emptyLineCount >= 2 || (len(bodyLines) > 0 && emptyLineCount >= 1) {
				break
			}
		} else {
			emptyLineCount = 0
			bodyLines = append(bodyLines, line)
		}
	}

	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))
	if body != "" {
		return message + "\n\n" + body
	}

	return message
}

func confirmCommit() bool {
	fmt.Print("\nProceed with commit? [Y/n]: ")
	var confirm string
	fmt.Scanln(&confirm)

	confirm = strings.ToLower(strings.TrimSpace(confirm))
	return confirm == "y" || confirm == ""
}

func executeCommit(message string) {
	commitCmd := exec.Command("git", "commit", "-F", "-")
	commitCmd.Stdin = strings.NewReader(message)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr

	if err := commitCmd.Run(); err != nil {
		color.Red("Commit failed: %v", err)
		os.Exit(1)
	}

	color.Green("✓ Commit successful! 🎉")
}
