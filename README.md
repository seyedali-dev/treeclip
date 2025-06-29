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

### 📁 Phase 1: MVP – Copy folder contents into .txt and clipboard

1. [ ] CLI command to traverse a directory and read all files
2. [ ] Exclude files using a `.treeclipignore` or CLI flags
3. [ ] Save to `treeclip_output.txt` in `/tmp/` or CWD
4. [ ] Copy to clipboard
5. [ ] Open in default editor (e.g., `xdg-open`, `start`, or `open`)
6. [ ] Delete file after editor is closed

### 🧰 Phase 2: Quality-of-Life Enhancements

* Shortcut registration
* Starred directories stored in `~/.treeclip/config.json`
* Command like `treeclip --star ~/Projects/my-repo`

### 🕶️ Phase 3: Nice-to-Haves

* History of past N pastes
* Search through them
* Interactive mode with TUI (e.g., `charmbracelet/bubbletea`)
