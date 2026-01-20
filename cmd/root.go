/*
Copyright © 2026 NAME HERE <barisdilekci@outlook.com.tr>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	commitType  string
	useEmoji    bool
	dryRun      bool
	interactive bool
	commitScope string
)

type CommitType struct {
	Type        string
	Emoji       string
	Description string
}

var commitTypes = []CommitType{
	{"feat", "✨", "A new feature"},
	{"fix", "🐛", "A bug fix"},
	{"docs", "📝", "Documentation only changes"},
	{"style", "💄", "Changes that don't affect code meaning"},
	{"refactor", "♻️", "Code change that neither fixes a bug nor adds a feature"},
	{"perf", "⚡", "Performance improvements"},
	{"test", "✅", "Adding or correcting tests"},
	{"build", "🔨", "Changes to build system or dependencies"},
	{"ci", "👷", "Changes to CI configuration"},
	{"chore", "🧹", "Other changes that don't modify src or test files"},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "commitz",
	Short: "Smart commit message generator",
	Long: `Commitz helps you create well-formatted conventional commits.
It can auto-detect commit types or guide you through an interactive process.`,
	Run: func(cmd *cobra.Command, args []string) {
		generateCommitMessage()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&commitType,
		"type",
		"t",
		"",
		"Commit type (feat, fix, docs, refactor, test, chore)",
	)

	rootCmd.PersistentFlags().StringVarP(
		&commitScope,
		"scope",
		"s",
		"",
		"Commit scope",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&useEmoji,
		"emoji",
		"e",
		false,
		"Add emoji to commit message",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&dryRun,
		"dry-run",
		"d",
		false,
		"Preview commit message without committing",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&interactive,
		"interactive",
		"i",
		false,
		"Enable interactive commit mode",
	)
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

	var selectedType string
	var selectedScope string
	var selectedEmoji string

	// Interactive mode
	if interactive {
		selectedType, selectedEmoji = selectCommitTypeInteractive()
		selectedScope = selectScopeInteractive()
	} else {
		// Auto-detect or use provided flags
		selectedType = detectCommitType(diffStr)
		if commitType != "" {
			selectedType = commitType
		}

		selectedScope = extractScopeFromBranch()
		if commitScope != "" {
			selectedScope = commitScope
		}

		selectedEmoji = getEmojiForType(selectedType)
	}

	// Generate summary with smart suggestion
	summary := generateSummaryInteractive(interactive, diffStr, selectedType)

	// Build commit message
	message := buildCommitMessage(selectedEmoji, selectedType, selectedScope, summary)

	// Display suggested message
	displaySuggestedMessage(message)

	// Add optional description
	message = addDescriptionInteractive(message, interactive)

	// Handle dry-run
	if dryRun {
		color.Yellow("\n[DRY RUN] Commit not created")
		fmt.Println("\nProposed commit message:")
		fmt.Println(color.CyanString(message))
		return
	}

	// Confirm and commit
	if confirmCommitInteractive(interactive) {
		executeCommit(message)
	} else {
		color.Yellow("Commit cancelled.")
	}
}

func selectCommitTypeInteractive() (string, string) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "▸ {{ .Emoji }} {{ .Type | cyan }} - {{ .Description }}",
		Inactive: "  {{ .Emoji }} {{ .Type | cyan }} - {{ .Description }}",
		Selected: "{{ .Emoji }} {{ .Type | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select commit type",
		Items:     commitTypes,
		Templates: templates,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		color.Red("Selection cancelled")
		os.Exit(0)
	}

	selected := commitTypes[idx]
	emoji := ""
	if useEmoji {
		emoji = selected.Emoji + " "
	}

	return selected.Type, emoji
}

func selectScopeInteractive() string {
	// Try to extract scope from branch first
	branchScope := extractScopeFromBranch()

	// Get common scopes from project structure
	commonScopes := getCommonScopes()

	// Add branch scope if available
	if branchScope != "" {
		commonScopes = append([]string{branchScope + " (from branch)"}, commonScopes...)
	}

	// Add "no scope" option
	commonScopes = append(commonScopes, "Skip (no scope)")

	prompt := promptui.Select{
		Label: "Select scope (optional)",
		Items: commonScopes,
		Size:  8,
	}

	_, result, err := prompt.Run()
	if err != nil || strings.Contains(result, "Skip") {
		return ""
	}

	// Remove "(from branch)" suffix if present
	result = strings.TrimSuffix(result, " (from branch)")
	return result
}

func getCommonScopes() []string {
	scopes := []string{}

	// Check for common directories
	dirs := []string{"cmd", "pkg", "internal", "api", "web", "docs", "test", "config", "auth", "db", "ui"}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			scopes = append(scopes, dir)
		}
	}

	// Add generic options
	scopes = append(scopes, "core", "deps", "ci")

	return scopes
}

func generateSmartSummary(diff string, commitType string) string {
	diffLower := strings.ToLower(diff)

	// Extract file names from diff
	lines := strings.Split(diff, "\n")
	var modifiedFiles []string
	var addedContent []string

	for _, line := range lines {
		// Check for file changes
		if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
			parts := strings.Fields(line)
			if len(parts) >= 2 && parts[1] != "/dev/null" {
				fileName := strings.TrimPrefix(parts[1], "b/")
				fileName = strings.TrimPrefix(fileName, "a/")
				if fileName != "" && !contains(modifiedFiles, fileName) {
					modifiedFiles = append(modifiedFiles, fileName)
				}
			}
		}

		// Look for added lines with meaningful content
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			content := strings.TrimPrefix(line, "+")
			content = strings.TrimSpace(content)
			if len(content) > 10 && !strings.HasPrefix(content, "//") && !strings.HasPrefix(content, "/*") {
				addedContent = append(addedContent, content)
			}
		}
	}

	// Generate smart summary based on commit type and changes
	switch commitType {
	case "feat":
		if strings.Contains(diffLower, "interactive") {
			return "add interactive mode"
		}
		if strings.Contains(diffLower, "api") {
			return "add API endpoints"
		}
		if len(modifiedFiles) > 0 {
			baseName := getBaseName(modifiedFiles[0])
			return fmt.Sprintf("add %s functionality", baseName)
		}
		return "add new feature"

	case "fix":
		if strings.Contains(diffLower, "bug") || strings.Contains(diffLower, "error") {
			return "fix bug in error handling"
		}
		if len(modifiedFiles) > 0 {
			baseName := getBaseName(modifiedFiles[0])
			return fmt.Sprintf("fix issue in %s", baseName)
		}
		return "fix bug"

	case "docs":
		if strings.Contains(diffLower, "readme") {
			return "update README documentation"
		}
		return "update documentation"

	case "refactor":
		if len(modifiedFiles) > 0 {
			baseName := getBaseName(modifiedFiles[0])
			return fmt.Sprintf("refactor %s", baseName)
		}
		return "refactor code structure"

	case "test":
		return "add/update tests"

	case "style":
		return "improve code formatting"

	case "perf":
		return "improve performance"

	case "build":
		if strings.Contains(diffLower, "go.mod") || strings.Contains(diffLower, "go.sum") {
			return "update dependencies"
		}
		return "update build configuration"

	case "ci":
		return "update CI configuration"

	case "chore":
		if strings.Contains(diffLower, "cleanup") {
			return "cleanup code"
		}
		return "update project files"
	}

	return "update changes"
}

func getBaseName(filePath string) string {
	// Remove extension and get base name
	parts := strings.Split(filePath, "/")
	fileName := parts[len(parts)-1]
	fileName = strings.TrimSuffix(fileName, ".go")
	fileName = strings.TrimSuffix(fileName, ".js")
	fileName = strings.TrimSuffix(fileName, ".ts")
	fileName = strings.TrimSuffix(fileName, ".md")
	return fileName
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func generateSummaryInteractive(interactive bool, diff string, commitType string) string {
	// Generate smart suggestion
	suggestion := generateSmartSummary(diff, commitType)

	if !interactive {
		return suggestion
	}

	validate := func(input string) error {
		if len(input) < 3 {
			return fmt.Errorf("summary must be at least 3 characters")
		}
		if len(input) > 72 {
			return fmt.Errorf("summary should be under 72 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("Commit summary (suggestion: %s)", color.CyanString(suggestion)),
		Default:  suggestion,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		color.Red("Input cancelled")
		os.Exit(0)
	}

	return strings.TrimSpace(result)
}

func addDescriptionInteractive(message string, interactive bool) string {
	if interactive {
		prompt := promptui.Prompt{
			Label:     "Add detailed description? (y/N)",
			IsConfirm: true,
		}

		result, err := prompt.Run()
		if err != nil || strings.ToLower(result) != "y" {
			return message
		}
	}

	fmt.Println("\n" + color.CyanString("Enter description (press Enter twice to finish):"))

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

func confirmCommitInteractive(interactive bool) bool {
	if !interactive {
		fmt.Print("\nProceed with commit? [Y/n]: ")
		var confirm string
		fmt.Scanln(&confirm)
		confirm = strings.ToLower(strings.TrimSpace(confirm))
		return confirm == "y" || confirm == ""
	}

	prompt := promptui.Prompt{
		Label:     "Proceed with commit",
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return false
	}

	return strings.ToLower(result) == "y" || result == ""
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

	for _, ct := range commitTypes {
		if ct.Type == commitType {
			return ct.Emoji + " "
		}
	}

	return ""
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
