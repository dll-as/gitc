package utils

import (
	"fmt"
	"strings"

	"github.com/dll-as/gitc/internal/ai"
)

// PrintDryRun displays the exact prompt and configuration sent to the AI model.
// Used exclusively in --dry-run mode to help debugging prompt engineering and cost estimation.
// No network request is made.
func PrintDryRun(diff string, cfg *ai.Config) {
	prompt := GetPromptForSingleCommit(diff, cfg.CommitType, cfg.CustomConvention, cfg.Language)

	fmt.Println("Prompt sent to model:")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Println(prompt)
	fmt.Println(strings.Repeat("─", 70))

	fmt.Printf(
		"\nProvider : \033[1m%s\033[0m  |  Model : \033[1m%s\033[0m  |  Lang : \033[1m%s\033[0m  |  Timeout : %s\n\n",
		strings.ToUpper(cfg.Provider), cfg.Model, cfg.Language, cfg.Timeout,
	)
}
