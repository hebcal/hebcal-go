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

CI (`.github/workflows/go.yml`) builds and tests with the Go version `go.mod` declares (`go-version-file`), currently 1.18 -- keep code buildable on that version. CodeQL runs on push/PR to `main`.

## Architecture

This is a **library** (no `main`, no `cmd/`), published as `github.com/hebcal/hebcal-go`. It is a Go port of the original C `hebcal` and is kept behavior-compatible with the JS `@hebcal/core`. Recent commit history shows deliberate reconciliation against `@hebcal/core` (URLs, `Erev` flags, candle-lighting), so when changing holiday logic, URLs, or flags, check the JS implementation for parity.

### The event model (`event/`)

Everything the calendar emits satisfies `event.CalEvent`, which has eight methods: `GetDate`, `Desc`, `Basename`, `Render(locale)`, `GetFlags`, `GetEmoji`, `GetCategories`, `URL`. `Desc()` is the untranslated description and the event's stable identifier ("Sukkot III (CH''M)"; consumers hash it into iCalendar UIDs, so never change what an event returns). `Basename()` groups related events ("Sukkot"). `URL()` is `""` for an event with no page. `CalEvent` is also implemented outside this module (the schedules in `github.com/hebcal/learning`), so adding a method to it is a breaking change for them.

Four abilities are optional interfaces, found with a type assertion or a helper:

| Interface | Method | Helper | Who has it |
|---|---|---|---|
| `event.BriefRenderer` | `RenderBrief(locale)` | `event.RenderBrief(ev, l)` falls back to `Render` | holidays, Hebrew dates, Mevarchim, Omer, timed events |
| `event.Memoer` | `Memo(locale)` | `event.Memo(ev, l)` falls back to `""` | Yom Kippur Katan, Mevarchim, timed events, Kiddush Levana |
| `event.Timed` | `Time() time.Time` | — | `hebcal.TimedEvent`, `hebcal.KiddushLevanaEvent` |
| `event.Linked` | `LinkedEvent() CalEvent` | — | `hebcal.TimedEvent`, `hebcal.YizkorEvent` |

To find the events with a clock time, assert `event.Timed`, never the concrete `hebcal.TimedEvent`: that assertion misses Kiddush Levana.

`HolidayFlags` is a bitmask (`event/flags.go`) — `CHAG`, `LIGHT_CANDLES`, `YOM_TOV_ENDS`, `CHUL_ONLY`/`IL_ONLY`, `CHANUKAH_CANDLES`, `PARSHA_HASHAVUA`, `DAF_YOMI`, `EREV`, `CHOL_HAMOED`, etc. Flags drive downstream filtering and rendering; adding a flag means touching `flagNames` (same file, for `String()`) and `getMaskFromOptions` in `hebcal/hebcal.go`.

Concrete event types live in `event/`: `HolidayEvent`, `ParshaEvent`, `HebrewDateEvent`, `MevarchimChodeshEvent`, `MoladEvent`, `UserEvent` (yahrzeit/birthday). Conventions for an event type:

- What the interface reads (date, description, flags, emoji) is an **unexported field** behind the interface method, and the type is built with a constructor that returns the concrete type (`NewHolidayEvent`, `NewUserEvent`, `NewParshaEvent`, `NewHebrewDateEvent`, `NewMoladEvent`, `NewMevarchimChodeshEvent`, `omer.NewOmerEvent`). A `Desc()` method cannot coexist with a `Desc` field. Fields with no method of the same name stay exported (`HolidayEvent.ChanukahDay`, `ParshaEvent.Parsha`, `MoladEvent.Molad`).
- **Never embed one concrete event type in another.** Embedding promotes every method, present and future: `TimedEvent` once embedded `HolidayEvent` and had to override `URL()` to avoid `/holidays/havdalah-2026`. Hold the other event in an unexported field and write each method out, as `KiddushLevanaEvent` does with its `TimedEvent`.
- Clock times are rendered through `event.ClockFormat` (`ClockAuto`, `Clock12`, `Clock24`) and a country code; see `event.ReformatTimeStr`.

URL generation lives in `event/url.go` (holidays, parsha, Hebrew date) and `omer/url.go`; Hebrew years outside 5513-6260 have no page.

`event.Descer`, `event.URLer`, `event.Desc`, `event.URL` and `hebcal.AsTimedEvent` are deprecated leftovers of the optional-interface design, kept for one release.

### The calendar engine (`hebcal/`)

`hebcal.HebrewCalendar(opts *CalOptions) ([]event.CalEvent, error)` in `hebcal/hebcal.go` is the primary entry point. `CalOptions` (`hebcal/options.go`) is large and toggles nearly all behavior: date range (`Year`/`IsHebrewYear`/`NumYears`, or explicit `Start`/`End` HDates), `IL` (Israel vs. Diaspora schedule — affects both holidays and Torah readings), suppression flags (`NoHolidays`, `NoMinorFast`, `NoModern`, `NoRoshChodesh`, …), opt-in event streams (`Sedrot`, `Omer`, `DafYomi`, `YerushalmiYomi`, `MishnaYomi`, `NachYomi`, `ShabbatMevarchim`, `Molad`, `YomKippurKatan`, `CandleLighting`), and `Clock` (an `event.ClockFormat`).

**`HebrewCalendar` never writes to the caller's `CalOptions`.** `newCalConfig` copies it and normalizes the copy into the unexported `calConfig` (defaults, the Israel schedule implied by an Israeli location, the resolved mask, a negated candle-lighting offset), and the internal helpers take `*calConfig`. Events keep only the resolved `clock` (country code and `ClockFormat`), never a pointer to the options, so one options value can be shared by concurrent calls. `options_immutable_test.go` guards this, including under `-race`; a new normalization must go on the copy.

- `holidays.go` — `GetHolidaysForYear(year, il)` / `getAllHolidaysForYear` build the fixed holiday set for a Hebrew year, including computed dates (`tzomGedaliahDate`, `taanitEstherDate`, Birkat HaChama).
- `candles.go` — candle-lighting / havdalah / fast start-end / Chanukah / biur chametz. `TimedEvent` is the timed `CalEvent`; times come from `zmanim`. Its `Kind()` is a `TimedKind` (`CandleLighting`, `Havdalah`, `ChanukahCandles`, `FastBegins`, `FastEnds`, `FinishChametz`, `BiurChametz`, `KiddushLevana`, `Zman`): switch on the kind, not on the description string, which is only the translation key and identifier. `checkCandleOptions` validates location config.
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

## Deliberate divergences from @hebcal/core

Everything else should match `@hebcal/core`; these differences are on purpose:

- **Rosh Hashana's year is gematriya in Hebrew** ("רֹאשׁ הַשָּׁנָה תשפ״ז"), a
  number in every other locale. `@hebcal/core` writes a number in Hebrew too.
- **Tammuz is spelled with two m's** in event descriptions (from `hdate`),
  where `@hebcal/core` and hebcal.com write "Tamuz". Translations and URLs
  already handle both spellings.
