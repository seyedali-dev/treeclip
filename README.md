# 👉 **TreeClip**

---

### 🎯 Core Objective

Create a CLI tool that:

* Traverses a directory tree
* Outputs all file contents into a single `.txt` file
* Copies that output to the clipboard
* Optionally opens it in the default text editor
* Deletes the temp file after the editor is closed

---

## 📋 Feature List

### 📦 Core

- [ ] Traverse and read all files in a directory recursively
- [ ] Exclude certain files/patterns (e.g., via `.treeclipignore`)
- [ ] Save combined output to a `.txt` file
- [ ] Copy the result to clipboard
- [ ] Open the output in a default system text editor
- [ ] Delete the file when the editor is closed

### ⚙️ Quality of Life

- [ ] Register CLI tool to open from a system-wide shortcut
- [ ] Support opening in CWD, pinned/starred dirs, add/remove starred dirs
- [ ] Load starred dirs from a config file

### 🧠 Inspired by Ditto

- [ ] Optionally keep a history of last N outputs
- [ ] Optionally search past outputs

---

## 🔧 Technology Choices

| Feature                          | Go Dependencies                  |
|----------------------------------|----------------------------------|
| Cross-platform Clipboard         | `github.com/atotto/clipboard`    |
| File traversal                   | Built-in `filepath.WalkDir`      |
| Ignore files (like `.gitignore`) | `doublestar`, custom logic       |
| Open text editor                 | `os/exec.Command()` + `xdg-open` |
| Temp file handling               | `os.CreateTemp()`                |

✅ Recommendation: **Go** – for fast prototyping, fewer dependencies, and easier system integration (editor launching,
clipboard, etc.).

---

## 🧭 Development Phases

### 🏁 Phase 1 – MVP (Get something working quickly)

1. **CLI scaffold**

    * Set up Go module (`go mod init`)
    * Use Cobra for CLI structure
2. **Directory traversal**

    * Recursively find files in a specified directory
    * Output full file paths and contents in order into a temp `.txt`
3. **Basic exclusion**

    * Support `--exclude` flag (simple wildcard patterns)
4. **Clipboard integration**

    * Use `atotto/clipboard` to copy combined output
5. **Open in editor**

    * Use `skratchdot/open-golang` or `exec.Command` to open in default editor
6. **Cleanup**

    * After editor exits, delete temp file

---

### 🚀 Phase 2 – Enhanced Usability

7. **Ignore file support**

    * Add `.treeclipignore` parsing for file patterns (git‑ignore style)
8. **Config & starred directories**

    * Create `~/.treeclip/config.json`
    * Support commands:

        * `treeclip star <dir>`
        * `treeclip list-stars`
        * `treeclip remove-star <dir>`
9. **Shortcut workflow**

    * Add command `treeclip open`

        * Prompts to choose from CWD or starred directories
        * Executes the traversal + clipboard + editor flow

---

### 🧰 Phase 3 – Advanced Features (Non‑MVP)

10. **History management**

* Store last N outputs (e.g. as files or in a small DB)
* Command `treeclip history` to list past entries
* `treeclip show <id>` to copy and open a past entry

11. **Search within history**

* Support fuzzy search in past outputs (`treeclip search <term>`)

12. **Interactive TUI**

* Use Bubble Tea (Charmbracelet) to navigate starred directories and history

---

## ✅ Summary Table

| Phase       | Task # | Feature                                  | Status |
|-------------|--------|------------------------------------------|--------|
| **MVP**     | 1      | CLI scaffold                             | ☐      |
|             | 2      | Directory traversal & output to txt      | ☐      |
|             | 3      | Basic `--exclude` filtering              | ☐      |
|             | 4      | Copy to clipboard via `atotto/clipboard` | ☐      |
|             | 5      | Open in default editor                   | ☐      |
|             | 6      | Delete temp file after editor closes     | ☐      |
| **Phase 2** | 7      | `.treeclipignore` support                | ☐      |
|             | 8      | Starred dirs config commands             | ☐      |
|             | 9      | Shortcut-style `open` command            | ☐      |
| **Phase 3** | 10     | Store & list history of outputs          | ☐      |
|             | 11     | Search past outputs                      | ☐      |
|             | 12     | Interactive TUI for picks                | ☐      |

---

✅ **Git commit message suggestion**

```bash
git commit -m "chore: define phase‑based task plan for TreeClip MVP → advanced features"
```
