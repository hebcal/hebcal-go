# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

There is no Makefile. Standard Go tooling:

- Build: `go build -v ./...`
- Test all: `go test ./...` (CI runs `go test -v ./...`)
- Test one package: `go test ./sedra/`
- Test one function: `go test -run TestSedra ./sedra/`
- Coverage: `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`
- Vet / format check: `go vet ./...` and `gofmt -l .` (CI-adjacent hygiene; keep `gofmt -l .` empty)

CI (`.github/workflows/go.yml`) builds and tests on Go 1.17; `go.mod` declares Go 1.18. CodeQL runs on push/PR to `main`.

## Architecture

This is a **library** (no `main`, no `cmd/`), published as `github.com/hebcal/hebcal-go`. It is a Go port of the original C `hebcal` and is kept behavior-compatible with the JS `@hebcal/core`. Recent commit history shows deliberate reconciliation against `@hebcal/core` (URLs, `Erev` flags, candle-lighting), so when changing holiday logic, URLs, or flags, check the JS implementation for parity.

### The event model (`event/`)

Everything the calendar emits satisfies `event.CalEvent`: `GetDate`, `Render(locale)`, `GetFlags`, `GetEmoji`, `Basename`, `GetCategories`. `HolidayFlags` is a bitmask (`event/flags.go`) — `CHAG`, `LIGHT_CANDLES`, `YOM_TOV_ENDS`, `CHUL_ONLY`/`IL_ONLY`, `CHANUKAH_CANDLES`, `PARSHA_HASHAVUA`, `DAF_YOMI`, `EREV`, `CHOL_HAMOED`, etc. Flags drive downstream filtering and rendering; adding a flag means touching `flagNames` (same file, for `String()`) and `getMaskFromOptions` in `hebcal/hebcal.go`.

Concrete event types live in `event/`: `HolidayEvent`, `parshaEvent`, `hebrewDateEvent`, `MevarchimChodeshEvent`, `moladEvent`, user events (yahrzeit/birthday) in `user.go`. URL generation is a separate concern: `event.URL(ev)` dispatches through the optional `URLer` interface (`event/url.go`); only some event types carry a URL, and `gregYearInRange` gates it.

### The calendar engine (`hebcal/`)

`hebcal.HebrewCalendar(opts *CalOptions) ([]event.CalEvent, error)` in `hebcal/hebcal.go` is the primary entry point. `CalOptions` (`hebcal/options.go`) is large and toggles nearly all behavior: date range (`Year`/`IsHebrewYear`/`NumYears`, or explicit `Start`/`End` HDates), `IL` (Israel vs. Diaspora schedule — affects both holidays and Torah readings), suppression flags (`NoHolidays`, `NoMinorFast`, `NoModern`, `NoRoshChodesh`, …), and opt-in event streams (`Sedrot`, `Omer`, `DafYomi`, `YerushalmiYomi`, `MishnaYomi`, `NachYomi`, `ShabbatMevarchim`, `Molad`, `YomKippurKatan`, `CandleLighting`).

- `holidays.go` — `GetHolidaysForYear(year, il)` / `getAllHolidaysForYear` build the fixed holiday set for a Hebrew year, including computed dates (`tzomGedaliahDate`, `taanitEstherDate`, Birkat HaChama).
- `candles.go` — candle-lighting / havdalah / fast start-end / Chanukah / biur chametz. `TimedEvent` is the timed `CalEvent`; times come from `zmanim`. `checkCandleOptions` validates location config.
- Holiday "related" events (candle-lighting, havdalah, mevarchim, chametz) are appended by `appendHolidayAndRelated`.

### Supporting packages

- `sedra/` — `sedra.New(year, il)` then `Lookup(hd)` / `LookupByRD(rd)` for weekly Torah reading; encodes the fixed parsha schedules for all year types.
- `zmanim/` — halachic times from lat/long/elevation; `zmanim/cities.go` + `location.go` hold the small built-in city database with timezones.
- `omer/`, `molad/` — Sefirat HaOmer and molad (new-moon) calculations.
- `dailylearning/` — a **plugin registry** only. `AddCalendar(name, fn, [startDate])` registers a schedule; `Lookup(name, hd, il)` retrieves an event. Names are case-insensitive. The actual schedules (929, Daily Rambam, etc.) live in the separate module `github.com/hebcal/learning`, which calls `AddCalendar` from its `init`; importing that module is what enables `DAILY_LEARNING` events here.

### External hebcal modules this depends on

`github.com/hebcal/hdate` (Hebrew↔Gregorian dates, personal anniversaries), `.../greg` (Gregorian↔R.D. day numbers), `.../locales` (holiday-name translations/transliterations), `.../gematriya`, `.../noaa-go` (sunrise/sunset). Changes to date math or names often belong in those repos, not here.

## Reference originating issues in commits and release notes

When committing or releasing work that addresses a GitHub feature request or bug,
reference the originating issue (full URL) in the commit message and/or release
notes.

**Why:** Keeps the public history and release notes traceable back to the user
request that motivated the change.

**How to apply:** Example — the 929 daily bible chapter schedule in
github.com/hebcal/learning addresses
https://github.com/hebcal/hebcal/issues/300 ("Add events of '929 daily bible
chapter project'"); cite it in the v0.2.0 release notes. Ask the user for the
issue link if a change clearly maps to a request but you don't have the number.
