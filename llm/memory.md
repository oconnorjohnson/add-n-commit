# add-n-commit Development Memory

## Project Overview

- **Date Started**: July 1, 2025, 12:45 PM PDT
- **Purpose**: Build a more robust CLI app for auto-generating commit messages, replacing a simple bash script
- **Tech Stack**: Go with Bubble Tea TUI framework and OpenAI SDK

## Key Decisions

### Architecture (July 1, 2025, 12:45 PM PDT)

- Chose Go for better performance and distribution
- Used Bubble Tea for beautiful TUI interface
- Modular structure with separate packages for:
  - `internal/app`: Main TUI application logic
  - `internal/config`: Configuration management
  - `internal/git`: Git operations wrapper
  - `internal/openai`: OpenAI API client wrapper
  - `internal/ui`: UI components and styles

### Features Implemented (July 1, 2025, 12:53 PM PDT)

1. **Interactive File Selection**: Users can select specific files to stage (vs. `git add -A`)
2. **Multiple Commit Modes**:
   - All-in-one summary
   - File-by-file summary
   - Custom prompt with additional context
3. **API Key Management**:
   - CLI flags: `--set-key`, `--show-key`, `--delete-key`
   - Interactive config editor: `--config`
   - Secure storage in `~/.config/anc/config.json`
   - Environment variable support
4. **Beautiful TUI**:
   - File selection with checkboxes
   - Mode selection menu
   - Message preview and editing
   - Progress indicators
5. **Configuration Options**:
   - Model selection (default: o4-mini)
   - Temperature control
   - Custom system prompts
   - Default mode preference

### Design Choices

- **State Machine Pattern**: Used explicit states for different UI screens
- **Secure Key Storage**: API keys are masked when displayed, stored in user config
- **Error Handling**: Graceful error states with clear messages
- **Keyboard Navigation**: Vim-style keys (j/k) plus standard arrows

## Testing Instructions

### Build and Install

```bash
# Build the binary
go build -o anc

# Option 1: Run locally
./anc

# Option 2: Install globally
sudo cp anc /usr/local/bin/
# or
cp anc ~/bin/  # if ~/bin is in your PATH
```

### Test Scenarios

1. **First Run (No API Key)**:

   ```bash
   ./anc
   # Should prompt for API key configuration
   ```

2. **API Key Management**:

   ```bash
   # Set key
   ./anc --set-key sk-your-test-key-here

   # View key (should be masked)
   ./anc --show-key

   # Delete key
   ./anc --delete-key
   ```

3. **Configuration Editor**:

   ```bash
   ./anc --config
   # Test navigation with Tab/Shift+Tab
   # Test saving with Ctrl+S
   ```

4. **Main Workflow**:

   ```bash
   # Make some test changes
   echo "test" > test.txt

   # Run the app
   ./anc

   # Test file selection:
   # - Space to toggle files
   # - 'a' to toggle all
   # - Enter to continue

   # Test mode selection
   # Test message generation
   # Test editing (press 'e')
   # Test regeneration (press 'r')
   ```

5. **Edge Cases**:
   - Run outside a git repo (should show error)
   - Run with no changes (should show "No changes detected")
   - Run with only staged files (should show only unstaged files)

### Environment Variables

```bash
# Test with environment variable
export OPENAI_API_KEY=sk-test-key
./anc
```

## Next Steps/Improvements

- Add support for conventional commits format
- Add diff preview in file selection
- Support for amending commits
- Integration with git hooks
- Add more models (GPT-4, Claude, etc.)
- Batch processing for multiple repositories
- Export/import configuration

## 2025-07-01 13:05:44 PDT - Debugging Empty File List Issue

**Issue**: User reported not seeing any files, modes, or commit messages in the app.

**Investigation**:

1. Examined the app.go file to understand the UI flow
2. Checked git.go to verify the GetStatus implementation
3. Ran git status commands to check repository state

**Root Cause**: The repository has no uncommitted changes. The git status showed:

- Working tree is clean
- No staged files
- No modified files
- No untracked files

**Resolution**: The app is functioning correctly. It shows "No changes detected" when there are no files to commit. The files that were shown as staged in the initial git status snapshot were already committed before the user ran the app.

**Key Learning**: The app properly handles the empty state by showing an appropriate message to the user.

## 2025-07-01 13:11:07 PDT - Fixed UI List Rendering and Navigation Issues

**Issues Fixed**:

1. **File list not displaying**: The `NewFileDelegate()` and `NewModeDelegate()` functions were returning `list.DefaultDelegate` instead of the custom delegates with proper render methods.

   - Fixed by changing return type to `list.ItemDelegate` and returning the custom delegate instances

2. **Keyboard navigation not working**: The `updateFileSelection` and `updateModeSelection` methods were not passing through keyboard events to the list components.
   - Fixed by always calling the list's Update method at the end of these functions to handle navigation keys (arrows, j/k)

**Code Changes**:

- Modified `internal/ui/styles.go`: Changed delegate factory functions to return custom delegates
- Modified `internal/app/app.go`: Added list update calls to handle navigation in file and mode selection

**Result**: The app now properly displays files and allows navigation using arrow keys or j/k vim keys.

## 2025-07-01 13:14:27 PDT - Fixed Commit Message Generation

**Issues Fixed**:

1. **Commit message generation not working**: The app was immediately showing empty commit message without calling OpenAI API.

   - Fixed variable shadowing in `generateCommitMessage` function where `err` was being redeclared in nested scopes
   - Added proper error handling with unique variable names (diffErr, filesErr, msgErr)

2. **No loading spinner during generation**: Added proper spinner initialization and tick commands.
   - Modified all paths to generating state to use `tea.Batch` with `m.spinner.Tick`
   - Improved the generating view with a bordered preview box showing the spinner

**Code Changes**:

- Modified `internal/app/app.go`:
  - Fixed commitModeSelectedMsg handler to start spinner tick
  - Fixed updateReviewing "r" handler to start spinner tick
  - Fixed updateEditing custom prompt handler to start spinner tick
  - Improved viewGenerating with a styled preview box
  - Fixed variable shadowing in generateCommitMessage function

**Result**: The app now properly shows a loading spinner during API calls and correctly generates commit messages using the OpenAI API.

## 2025-07-01 13:23:33 PDT - Handled Staged Files Edge Case

**Issue**: When users stage files and quit the app before committing, those files remain staged in Git. On next run, this creates confusion.

**Solution Implemented**:

1. **Detection**: On startup, check for already-staged files using `git diff --cached --name-only`
2. **User Prompt**: If staged files exist, show a prompt with options:
   - `c`: Continue with the already staged files
   - `u`: Unstage all and start fresh
   - `q`: Quit
3. **Visual Indicator**: Files already staged show "(staged)" next to their name in the file list
4. **Cleanup on Quit**: When user quits (via 'q' or Ctrl+C), unstage any files that were staged in this session but not committed

**Implementation Details**:

- Added `stateStagedFilesPrompt` state to handle the prompt
- Added `alreadyStagedFiles` field to track previously staged files
- Added `cleanup()` function that unstages files on early exit
- Modified all quit handlers to call cleanup
- Updated file list rendering to show staged status

**Result**: The app now gracefully handles the edge case of previously staged files, giving users clear options and preventing confusion about git state.
