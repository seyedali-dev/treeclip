// Package clipboard - clipboard provides logic for saving the traversed data in clipboard.
package clipboard

import (
	"fmt"
	atottoClip "github.com/atotto/clipboard"
	"github.com/seyedali-dev/treeclip/pkg/utils"
	"os"
	"strings"
)

// HandleClipboardCommandFlag handles clipboardEnabled command flag.
func HandleClipboardCommandFlag(clipboardFlag, clipboardStatsFlag bool, outputFilePath string) error {
	if clipboardFlag {
		fmt.Printf("\n📋  Copying content to clipboard... (ﾉ◕ヮ◕)ﾉ*:･ﾟ✧\n")
		clipboardContent, err := os.ReadFile(outputFilePath)
		if err != nil {
			return fmt.Errorf("failed to read output file for clipboard: %w", err)
		}

		// Copy to clipboard
		err = atottoClip.WriteAll(string(clipboardContent))
		if err != nil {
			fmt.Printf("⚠️  Warning: failed to copy to clipboard: %v\n", err)
			fmt.Printf("💡  Content is still available in: %s\n", outputFilePath)
		} else {
			fmt.Printf("✅  Content copied to clipboard successfully! ヽ(•‿•)ノ\n")

			// Show clipboard statistics if requested
			if clipboardStatsFlag {
				contentStr := string(clipboardContent)
				lines := strings.Split(contentStr, "\n")
				chars := len(contentStr)
				words := len(strings.Fields(contentStr))

				fmt.Printf("📊  Clipboard content stats:\n")
				fmt.Printf("   📝  Characters: %s\n", utils.FormatNumber(chars))
				fmt.Printf("   📄  Lines: %s\n", utils.FormatNumber(len(lines)))
				fmt.Printf("   💬  Words: %s\n", utils.FormatNumber(words))

				// Show size in human-readable format
				fmt.Printf("   💾  Size: %s\n", utils.FormatBytes(int64(chars)))
			}
		}
	} else {
		fmt.Printf("\n📋  Clipboard copy skipped (disabled) (︶︹︶)\n")
	}
	return nil
}
