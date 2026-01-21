package hook

import (
	"os"
	"path/filepath"

	"github.com/omgitsads/gh-pair/internal/git"
)

const hookScript = `#!/bin/sh
# gh-pair: Adds co-author trailers to commits
#
# This hook is managed by gh-pair. Do not edit manually.
# To update, run: gh pair init
#
# Arguments:
#   $1 - Path to the temporary file containing the commit message
#   $2 - Source of the commit message (message, template, merge, squash, commit)
#   $3 - SHA1 of commit (only when source is "commit")

COMMIT_MSG_FILE="$1"
COMMIT_SOURCE="$2"

# Only add co-authors for regular commits (not merges, squashes, or amends)
case "$COMMIT_SOURCE" in
merge | squash | commit)
  exit 0
  ;;
esac

# Path to the pairs config file
GIT_DIR=$(git rev-parse --git-dir 2>/dev/null)
CONFIG_FILE="$GIT_DIR/gh-pair/pairs.json"

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
  exit 0
fi

# Check if there are any pairs configured
if ! grep -q '"username"' "$CONFIG_FILE" 2>/dev/null; then
  exit 0
fi

# Check if co-authors are already present in the commit message
if grep -q "^Co-Authored-By:" "$COMMIT_MSG_FILE" 2>/dev/null; then
  exit 0
fi

# Extract co-author lines from JSON and prepend to commit message
# Uses simple text processing to avoid requiring jq
COAUTHORS=$(grep -E '"(name|email)"' "$CONFIG_FILE" | \
  sed 's/.*"name": *"\([^"]*\)".*/\1/' | \
  paste - - | \
  while read name email; do
    # Clean up the email extraction
    clean_email=$(echo "$email" | sed 's/.*"email": *"\([^"]*\)".*/\1/')
    echo "Co-Authored-By: $name <$clean_email>"
  done)

# If no co-authors found, exit
if [ -z "$COAUTHORS" ]; then
  exit 0
fi

# Read original message
ORIGINAL_MSG=$(cat "$COMMIT_MSG_FILE")

# Write co-authors followed by blank line and original message
{
  echo "$COAUTHORS"
  echo ""
  echo "$ORIGINAL_MSG"
} > "$COMMIT_MSG_FILE"

exit 0
`

const hookMarker = "# gh-pair: Adds co-author trailers to commits"

// Install installs the prepare-commit-msg hook.
func Install() error {
	hooksDir, err := git.HooksDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return err
	}

	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")

	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		// Read existing hook to check if it's ours
		content, err := os.ReadFile(hookPath)
		if err != nil {
			return err
		}

		// If it's our hook, just update it
		if !isOurHook(string(content)) {
			// Backup existing hook
			backupPath := hookPath + ".gh-pair-backup"
			if err := os.Rename(hookPath, backupPath); err != nil {
				return err
			}
		}
	}

	// Write the hook
	if err := os.WriteFile(hookPath, []byte(hookScript), 0755); err != nil {
		return err
	}

	return nil
}

// Uninstall removes the prepare-commit-msg hook if it's ours.
func Uninstall() error {
	hooksDir, err := git.HooksDir()
	if err != nil {
		return err
	}

	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")

	// Check if hook exists
	content, err := os.ReadFile(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Nothing to uninstall
		}
		return err
	}

	// Only remove if it's our hook
	if !isOurHook(string(content)) {
		return nil
	}

	// Remove the hook
	if err := os.Remove(hookPath); err != nil {
		return err
	}

	// Restore backup if exists
	backupPath := hookPath + ".gh-pair-backup"
	if _, err := os.Stat(backupPath); err == nil {
		if err := os.Rename(backupPath, hookPath); err != nil {
			return err
		}
	}

	return nil
}

// IsInstalled checks if the gh-pair hook is installed.
func IsInstalled() bool {
	hooksDir, err := git.HooksDir()
	if err != nil {
		return false
	}

	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")
	content, err := os.ReadFile(hookPath)
	if err != nil {
		return false
	}

	return isOurHook(string(content))
}

func isOurHook(content string) bool {
	return len(content) > len(hookMarker) && content[0:len(hookMarker)] == hookMarker ||
		(len(content) > 10 && content[0:10] == "#!/bin/sh\n" && len(content) > len("#!/bin/sh\n"+hookMarker) &&
			content[10:10+len(hookMarker)] == hookMarker)
}
