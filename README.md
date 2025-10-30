# cmdiet

cmdiet is a terminal app for logging what you eat without leaving your keyboard. It saves entries in a small SQLite file and opens simple on-screen prompts to add meals, review weeks, and glance at trends. The rich terminal UI is powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea), so all flows feel smooth and familiar to Go developers who love text interfaces.

## Features at a glance
- Log breakfast, lunch, dinner, or snacks with a short guided flow.
- Revisit any recent week, jump by day offsets, or open a focused view of today.
- Keep a meal catalog, tune the nutrition numbers, and reuse meals in later logs.
- Draw quick charts of calories, protein, carbs, or fat to spot patterns.
- Spin up a fake dataset for demos or testing without touching your real log.

## What you need
- Go 1.22 or newer.
- A macOS, Linux, or Windows terminal that can run Go programs.

Go will pull every other dependency automatically the first time you run the app.

## Getting started
1. Clone this repo and move into it.
2. Run `go run .` once to download modules and set up the database.
3. A new `diet.db` file appears in the project folder; it holds your real data.

Build a reusable binary if you prefer:

```bash
go build -o cmdiet
./cmdiet --help
```

### Log a meal
```bash
go run . log breakfast
```
Pick an existing meal or type a new one. cmdiet blocks duplicate entries for the same meal type on the same day and timestamps the log for you.

Use `--pdays` (or `-d`) if you need to backfill a past day, for example:

```bash
go run . log lunch --pdays 1
```

### Review your days
```bash
go run . view today
go run . view --offset 14
```
The first command opens today's breakdown. The second shows the last 14 days and lets you move backward or forward with the arrow keys.

### Check trends
```bash
go run . eval calories
```
Swap `calories` for `protein`, `carbs`, or `fat` to focus on a different macro. Press the keys shown on screen to switch between chart styles.

## Work with meals
```bash
go run . meals list-meals
go run . meals add-meal
go run . meals list-components
go run . meals add-component --fixed
```
These screens manage the meals and components you reuse when logging food. Edits update future logs right away.

## Fake data for quick demos
Use the `--fake` flag to point the app at `fake.db`, then populate it with sample rows:

```bash
go run . --fake populate meals
go run . --fake populate diets
```

Switch back to your real log by dropping the `--fake` flag.

## Handy files and folders
- `diet.db`: your everyday log in SQLite form.
- `fake.db`: optional sandbox database created when you use `--fake`.
- `commands/`: code that wires up the terminal commands and flags.
- `ui/`: terminal screen logic for the interactive views.
- `meals/` and `diet/`: data access helpers and validation rules.

## Develop and test
- Run `go test ./...` to make sure the services still behave.
- Run `go run . --help` to see every command and flag.
- Temporary logs from the UI live in `tea.log`; delete it whenever you like.
