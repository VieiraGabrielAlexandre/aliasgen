# aliasgen — Automatic Alias Generator for Bash/Zsh/Fish

`aliasgen` observes your shell history (Bash, Zsh, Fish), learns which commands you run most often, and suggests **short aliases** to save keystrokes. It can then apply those suggestions by writing them into a dedicated aliases file that your shell sources automatically.

🚀 Status: **MVP** — functional ingestion (`learn`), suggestion (`suggest`), and application (`apply`). Ready for experimentation and community contributions.

---

## Features

- **Learn**: Ingests shell history (`~/.bash_history`, `~/.zsh_history`, `~/.local/share/fish/fish_history`) and stores normalized commands.
- **Suggest**: Generates aliases based on frequency, typing cost, and recency.
- **Apply**: Writes aliases to `~/.config/aliasgen/aliases.{sh|fish}` and gives you instructions to source them in your shell.
- **List**: Displays applied aliases with scores.
- **Explain**: Shows why a particular alias was suggested.

Planned improvements:
- Blacklist of dangerous commands (`rm -rf`, `mkfs`, …)
- Time-window filtering (e.g., only last 90 days)
- Conflict resolution for duplicate aliases
- Shell completions (`aliasgen completion bash|zsh|fish`)
- Context-based aliases (per project/directory)
- TUI interface for reviewing suggestions (with [Bubbletea](https://github.com/charmbracelet/bubbletea))
- Packaging with [GoReleaser](https://goreleaser.com) (DEB/RPM/AppImage)

---

## Installation (development)

Requirements:
- Go 1.21+
- (Optional) `make`, `goreleaser`, `fpm` for packaging.

Clone and build:

```bash
git clone https://github.com/VieiraGabrielAlexandre/aliasgen
cd aliasgen
go mod tidy
make build
./bin/aliasgen --help
````

Optionally install into your PATH:

```bash
mkdir -p ~/.local/bin
chmod +x ~/.local/bin/aliasgen
cp bin/aliasgen ~/.local/bin/
~/.local/bin/aliasgen learn
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
# (se usar Zsh:)
# echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc
hash -r
```

Then ensure `~/.local/bin` is in your `$PATH`.

---

## Usage

### 1. Learn from your history

```bash
aliasgen learn
```

### 2. See suggestions

```bash
aliasgen suggest --limit 100
# Output example:
# glg         -> git log --oneline --graph --decorate -n 20   (score=25.40)  [freq=12, saved=28, rec=0.85]
```

### 3. Apply aliases

```bash
aliasgen apply --shell auto --top 30
# Creates ~/.config/aliasgen/aliases.sh (or .fish)
# Shows how to add it to your shell RC file
```

### 4. List aliases

```bash
aliasgen list
```

### 5. Explain a suggestion

```bash
aliasgen explain glg
```

---

## Example

```bash
$ aliasgen learn
ingested: 150 records

$ aliasgen suggest --limit 10
glg          -> git log --oneline --graph --decorate -n 20   (score=42.1)  [freq=18, saved=28, rec=0.95]
gst          -> git status                                   (score=30.5)  [freq=12, saved=10, rec=0.87]

$ aliasgen apply --top 2
aliases written to: /home/user/.config/aliasgen/aliases.sh

Add this line to your ~/.zshrc:
[ -f "$HOME/.config/aliasgen/aliases.sh" ] && source "$HOME/.config/aliasgen/aliases.sh"
```

---

## Project Structure

```text
cmd/aliasgen           # main CLI entrypoint (Cobra commands)
internal/shell         # shell detection, history paths, alias file writer
internal/learn         # normalization, scoring, ingestion from history
internal/store         # SQLite persistence (commands + aliases)
internal/config        # config file loader
internal/log           # logger abstraction
pkg/aliasfmt           # helpers for alias formatting per shell
scripts/               # install/uninstall helpers
```
---

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## License

Licensed under the [Apache License 2.0](LICENSE).
