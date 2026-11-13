# hebcal-go: interface cleanup and idiomatic-Go review

Reviewed: `hebcal-go` at HEAD (`0d44bdb`), against the three consumers in this
session: `hebcal-api-go` (pins hebcal-go v0.27.0, Go 1.26), `learning` (v0.26.3,
Go 1.18), and the `hebcal` CLI (v0.24.0, Go 1.18, vendored). `go vet`, `gofmt -l`
and `go test ./...` are all clean today. Nothing was changed.

---

## 1. Is breaking the API allowed?

**Yes.** Go modules follow semantic versioning, and the Go docs say plainly that
a `v0.x.y` module "is still in development and unstable" and "carries no
backward compatibility or stability guarantees" (go.dev/doc/modules/version-numbers).
The *import compatibility rule* (a new version with the same import path must
stay compatible) is the promise you make **at v1**. Before that, a minor bump is
the accepted way to ship a breaking change.

What that means in practice:

- Ship the break as a **minor** bump (e.g. v0.27 → v0.28), never a patch. No
  `/v2` import path is needed below v2.
- Nothing breaks for anyone until they upgrade: each consumer is pinned by
  `go.mod`/`go.sum`, and the CLI also vendors.
- Courtesy that makes migration cheap: for one release, keep the old names as
  `// Deprecated:` aliases (`const CHAG = Chag`, `type Descer = …`) so
  consumers can upgrade first and fix call sites second. Then delete them.
- Since you would be asking every consumer to touch their code anyway, this is
  the right time to settle the API you would be willing to call **v1.0.0**.
  Bundle the breaks into one or two releases, not ten.
- The order matters because the dependency runs one way (`learning → hebcal-go`):
  release hebcal-go, then `learning`, then the API service and the CLI.

---

## 2. Can URLer and Descer go away?

### What they are for today

`event.CalEvent` has six methods. Four more abilities are layered on top as
optional interfaces, each with a helper that falls back when the interface is
missing:

| Interface | Method | Helper | Fallback | Implemented by |
|---|---|---|---|---|
| `URLer` | `URL() string` | `event.URL(ev)` | `""` | Holiday, parsha, Hebrew date, Omer, TimedEvent (overrides one it inherits), **all 11 learning events** |
| `Descer` | `GetDesc() string` | `event.Desc(ev)` | `Basename()` | Holiday, User, Hebrew date, Mevarchim, parsha, Molad, Omer, TimedEvent, Yizkor. **No learning event.** |
| `Memoer` | `Memo(locale)` | `event.Memo(ev, l)` | `""` | Holiday (YKK), Mevarchim, TimedEvent, KiddushLevana, Tanakh Yomi |
| `BriefRenderer` | `RenderBrief(locale)` | `event.RenderBrief(ev, l)` | `Render` | Holiday, Hebrew date, Mevarchim, Omer, TimedEvent, 929, Yerushalmi, Tanakh Yomi |

There is also a fifth, hidden one: `hebcal.AsTimedEvent` relies on an
unexported `timedEvent()` method so that it can find a `TimedEvent` inside a
`KiddushLevanaEvent`.

The doc comments explain why: `CalEvent` is implemented outside the package
(by `learning`), and adding a method would break those implementations. That
reasoning is correct **under a compatibility promise**. Without one, the choice
is just a design question.

### Optional interfaces are idiomatic Go in themselves

`io.WriterTo`, `http.Flusher` and `fmt.Stringer` are all optional interfaces
that callers find with a type assertion. So the pattern is not wrong. The test
is whether the fallback gives a **correct** answer:

- **`Descer` fails that test. Remove it and make `Desc()` part of the core
  interface.** Its fallback, `Basename()`, is a different concept: Basename
  merges the days of one holiday ("Sukkot"), and Desc keeps them apart ("Sukkot
  III (CH''M)"). The fallback silently gives a *wrong* identifier. The learning
  events get away with it only because their Basename happens to be the right
  Desc. Every event has an untranslated description, so the method should be
  required.
  Evidence from the consumers: `hebcal-api-go`'s `download.UntranslatedDesc`
  runs a four-step chain (assert `HolidayEvent` → `AsTimedEvent` → `Descer` →
  `Render("en")`) to answer "what is this event's description". With a required
  `Desc()` that becomes a single method call.
- **`URLer` passes the test (`""` really does mean "no page"), but I would still
  fold it in.** All three consumers call `event.URL` on every event, every
  learning event already implements it, and `hebcal-api-go` wraps hebcal-go
  events in its own `model.CalEv` interface mainly to get `Desc()` + `URL()`
  onto one value. If both are required, that adapter layer mostly disappears.
  Keeping `URLer` optional is defensible. Making it required is more
  consistent, and in v0 it costs only the 11 learning types plus a few
  hebcal-go types. (An embeddable zero type can give them the method for free;
  see §3.4.)
- **Keep `Memoer` and `BriefRenderer` optional.** Their fallbacks are genuinely
  correct ("no memo"; "brief form = full form"), only a minority of events have
  them, and these are real optional capabilities. Keep the helpers.
- **Replace the hidden `timedEvent()` marker with a public interface**, e.g.
  `event.Timed { CalEvent; Time() time.Time }` (and maybe `event.Linked {
  LinkedEvent() CalEvent }`). See §3.2 for why this matters.

**Recommended core interface** (named with the getter rename from §3.1. If you
skip that rename, keep the `Get` prefixes):

```go
type CalEvent interface {
    Date() hdate.HDate
    Desc() string              // untranslated identifier ("Erev Pesach")
    Basename() string          // groups related events ("Pesach")
    Render(locale string) string
    Flags() HolidayFlags
    Emoji() string
    Categories() []string
    URL() string               // "" when there is no page
}
// optional: RenderBrief, Memo, Timed, Linked
```

Eight methods is on the heavy side for Go ("the bigger the interface, the
weaker the abstraction"). But this interface describes *what a calendar event
is* rather than a behaviour, and the JS ancestor is a base class. An
embeddable default type (§3.4) keeps the cost of implementing it low.

---

## 3. Other recommendations, ranked by value per unit of churn

### 3.1 High value

**1. `HebrewCalendar` changes the caller's `*CalOptions`.** It writes
`opts.Mask`, `opts.IL`, `opts.KiddushLevanaMaharil`, flips
`CandleLightingMins` negative, fills in `HavdalahDeg`, sets
`YerushalmiEdition`, clears `NoRoshChodesh`/`NoModern`/… when a mask bit is
set, and turns `CandleLighting` off for early years. Every `TimedEvent` also
keeps a pointer to that same struct for later rendering. So:
- two goroutines sharing one options value is a data race;
- if the caller changes the options after the call, events it already received
  render differently;
- callers have to know which fields come back "normalized".

Fix: copy on entry (`o := *opts`), resolve it into an unexported
`resolvedOptions` struct, and have events keep only what they need (e.g. a clock
format, see item 5), not a pointer to the options.

**2. Embedding is used as inheritance, and it leaks.**
- `TimedEvent` embeds `event.HolidayEvent`, so it silently inherits
  `URL`, `Memo`, `RenderBrief`, `GetDesc`, `CholHaMoedDay` and `ChanukahDay`.
  It has to override `URL()` so that it doesn't produce `/holidays/havdalah-2026`
  (the comment on `TimedEvent.URL` says so). Every new method on
  `HolidayEvent` is a new latent bug in `TimedEvent`.
- `KiddushLevanaEvent` embeds `TimedEvent`, so `ev.(hebcal.TimedEvent)` does
  **not** match it, and `AsTimedEvent` exists to work around that. Consumers
  still use the plain assertion in four places that will miss Kiddush Levana:
  `hebcal-api-go/internal/service/shabbat/shabbat.go:273` and `:483`,
  `internal/service/mcp/tools.go:308`, and `hebcal/main.go:606`. The api-go
  CLAUDE.md records the PDF hitting exactly this bug.

Fix: give `TimedEvent` named fields (`Date`, `Desc`, `Flags`, `Emoji`, `Time`,
`Linked`) instead of an embedded holiday, and expose the `event.Timed`
interface. Then `if t, ok := ev.(event.Timed)` works for every timed event,
including ones in other modules.

**3. Behaviour is chosen by comparing description strings.**
`TimedEvent` checks `ev.Desc == "Candle lighting" / "Havdalah" / "Fast begins" /
"Finish eating chametz"` in `Memo`, `RenderBrief`, `GetCategories`, `TimeOfDay`
and `HebrewCalendar` (`candlesEv.Desc == "Havdalah"`). The engine checks
`event.Desc(ev) == "Erev Pesach"` / `"Yom Kippur"`, and the consumers do the same
(9 places in api-go and the CLI). `HolidayEvent.render` slices `ev.Desc[13:]`
and `[17:]`. Add a `Kind` enum to `TimedEvent` (`CandleLighting`, `Havdalah`,
`FastBegins`, `FastEnds`, `FinishChametz`, `BiurChametz`, `KiddushLevana`,
`Zman`, …) and switch on it. `Desc` stays as the translation key and stable id.

**4. Constructors that return interfaces, and unexported event types.**
`NewParshaEvent`, `NewHebrewDateEvent`, `NewMoladEvent` and
`NewMoladEventHour12` return `CalEvent`, while `NewMevarchimChodeshEvent` and
`omer.NewOmerEvent` return concrete types. `parshaEvent` and `hebrewDateEvent`
are unexported, so nobody can reach `.Parsha` or `.IL`, which is why api-go has
its own `model.ParshaEv`. Go guidance is "accept interfaces, return structs".
Export `ParshaEvent` and `HebrewDateEvent`, return concrete types, and merge the
two Molad constructors into one (or just use a struct literal).

### 3.2 Medium value

**5. Clock format is three overlapping knobs.** `CalOptions.Hour12`,
`CalOptions.Hour24` ("Hour24 wins") and `*bool hour12` on the events and on
`ReformatTimeStr`. Replace them with one enum: `type ClockFormat int`
(`ClockAuto`, `Clock12`, `Clock24`). It reads better than a tri-state `*bool`,
and it is the value events should keep instead of `*CalOptions`.

**6. Zero means "default", so sentinel booleans were added.**
`CandleLightingAtSunset` and `SuppressHavdalah` exist only because `0` already
means "use the default". That works, and it is documented. If you redesign the
options, consider `*int` for the minute offsets, or a small
`type Offset struct{ Mins int; Deg float64; Set bool }` that also absorbs the
four "Mins and Deg are mutually exclusive" pairs (Havdalah, FastStart, FastEnd,
TishaBavEnd), which are currently checked at runtime.

**7. Two sources of truth for what is included.** `Mask` and the ~20 booleans
both select event streams, and a non-zero `Mask` rewrites the booleans (item 1).
Pick one as the public API (the booleans read better) and make the mask an
internal detail, or document that `Mask` replaces all the booleans.

**8. Panics in library code.** `zmanim.NewLocation` panics on out-of-range
coordinates, `zmanim.New` on an unknown timezone, and `Sedra.LookupByRD` on a
zero-value `Sedra` or an early date. api-go's MCP tool validates beforehand only
to avoid the panic. Return `(Location, error)` and `(Zmanim, error)`, or add
`Must…` variants next to error-returning ones.

**9. Naming fixes** (each one is mechanical):
- ALL_CAPS flag constants (`CHAG`, `LIGHT_CANDLES`, …) → `Chag`,
  `LightCandles`. This touches 92 references in the consumers. A strong
  candidate for one release of deprecated aliases.
- Getters: `GetDate/GetFlags/GetEmoji/GetCategories/GetDesc` → drop `Get`
  (Effective Go). This **collides with exported fields** (`HolidayEvent.Date`,
  `.Desc`, `.Flags`, `.Emoji`), so it means either unexporting those fields
  (api-go reads `he.Desc`, `he.Flags` and `he.ChanukahDay` directly) or renaming
  them. This is the single largest piece of churn. Do it only if you want a
  clean v1, and do it together with §2.
- `Location.TimeZoneId` → `TimeZoneID`; `IsAssurBemlacha` → `IsAssurBemelacha`
  (the unexported helper already spells it that way).
- `hebcal.UserEvent` (an input spec) and `event.UserEvent` (an output event)
  share a name. Rename the input to `UserEventSpec` or `Anniversary`.
- `TimeOfDay(ev) int` returns the untyped `TimeOfDayDawn…` constants. Use
  `type DayPart int` with a `func DayPartOf(ev) DayPart`.

**10. The `dailylearning` registry.** The global map has no lock. That is fine
while every `AddCalendar` call comes from an `init`, but
`hebcal/timeofday_test.go` registers at run time, so a parallel test could race
with `HebrewCalendar`. Add a `sync.RWMutex` (cheap) or document "init only".
`AddCalendar(name, fn, startDate ...hdate.HDate)` uses a variadic parameter to
fake an optional argument. Prefer two functions or a `Calendar` struct.
`GetCalendars`/`GetStartDate` → `Calendars`/`StartDate`.

### 3.3 Low value / housekeeping

11. **Doc comments.** `options.go` documents fields with `/* … */` and fragment
    sentences. Go style is `// FieldName does …`. Package docs should start
    "Package event …" (now "Hebcal's event package …"; `dailylearning` is the
    same). `hebcal/doc.go` lists a `dafyomi` package that now lives in
    `learning`. The comment "Returns a slice of holidays for the year." sits on
    `holidaysYearKey` (`holidays.go:755`) instead of `GetHolidaysForYear`. Several
    exported methods and constructors (`NewParshaEvent`, the `UserEvent`
    methods, `AllCities`) have no doc comment at all.
12. **Returning internal slices.** `zmanim.AllCities()` returns the package's
    `classicCities` slice, so a caller can change it. `sedra.Parshiot()` already
    copies. Do the same.
13. **Locale handling.** `strings.ToLower(locale)` and `locale == "he" ||
    locale == "he-x-nonikud"` are repeated across `event`, `sedra`, `omer` and
    `learning` (`isHebrewLocale` again). One exported helper, ideally in
    `github.com/hebcal/locales`, would replace them.
14. **Sentinel errors.** The `errors.New("opts.CandleLighting requires
    opts.Location")` family could be exported variables (`ErrNoLocation`, …) so
    that callers can use `errors.Is`.
15. **CLI-only options in the library.** `SunriseSunset`, `DailyZmanim`,
    `WeeklyAbbreviated`, `DailySedra`, `Yahrzeits` and `UserEvents` exist for the
    CLI. Moving them out is a bigger project. At least group them in an embedded
    `CLIOptions` struct so that the library surface reads clearly.
16. **Go version.** hebcal-go, `learning` and the CLI declare `go 1.18`, and
    api-go is on 1.26. If you are breaking things anyway, moving to `go 1.21`
    gives you `slices`, `maps`, `min`/`max` and `sort` → `slices.SortStableFunc`.
    That requires bumping the CLI's CI, so it is your call.
17. **Stale notes.** `learning/CLAUDE.md` says "Latest releases … hebcal-go
    v0.14.0" and "`go 1.17` in go.mod". Both are out of date (v0.27, 1.18).

### 3.4 One design idea that ties §2 and §3 together

An embeddable zero-value base removes the boilerplate from every event type,
and makes it cheap to grow the core interface later:

```go
// Base supplies the zero answers for the optional parts of an event.
type Base struct{}
func (Base) URL() string     { return "" }
func (Base) Emoji() string   { return "" }
```

Use it carefully. This is the same embedding mechanism that bit `TimedEvent`
(§3.1 item 2). It is safe here only because `Base` has no fields and only
"empty" answers. Never embed a concrete event type in another one.

---

## 4. What I would *not* change

- `HolidayFlags` as a bitmask with `Has`, `HasAny`, `With`, `Without` and
  `String` is good, idiomatic Go.
- Keep the `dailylearning` plugin registry and its one-way dependency
  (`learning → hebcal-go`). It is the right shape.
- The `GetHolidaysForYear` memo (RWMutex, capacity clamped against aliasing) is
  well done.
- Keep `Memoer` and `BriefRenderer` as optional interfaces (§2).
- The deliberate divergences from `@hebcal/core` (gematriya year, Tammuz) are
  unaffected.

---

## 5. Suggested release plan

1. **v0.28.0: interface and correctness.** Required `Desc()` + `URL()`;
   `event.Timed`; de-embed `TimedEvent`/`KiddushLevanaEvent`; `Kind` enum;
   copy `CalOptions` on entry; `ClockFormat`; export `ParshaEvent` and
   `HebrewDateEvent` with concrete constructors. Keep `Descer`, `URLer` and
   `AsTimedEvent` as deprecated aliases for this one release.
2. **learning vNext**, then **hebcal-api-go** and the **CLI**: adopt it, and
   delete api-go's adapter chain (`UntranslatedDesc`, most of `model.CalEv`).
3. **v0.29.0: naming.** MixedCaps flags, getter renames (with field
   unexporting), `TimeZoneID`, error-returning constructors, docs. Delete the
   deprecated aliases from v0.28.
4. Then consider tagging **v1.0.0**.

---

## 6. Prompt for a separate Claude Code session

> Copy everything inside the fence. It assumes a session with `hebcal/hebcal-go`,
> `hebcal/learning`, `hebcal/hebcal-api-go` and `hebcal/hebcal` checked out side
> by side.

````text
You are refactoring github.com/hebcal/hebcal-go into a cleaner, more idiomatic
Go API. It is a v0 module, so breaking changes are allowed and will ship as a
MINOR version bump. The sibling checkouts github.com/hebcal/learning,
github.com/hebcal/hebcal-api-go and github.com/hebcal/hebcal (the CLI) are
consumers. Update them in the same session so that the whole graph builds and
passes its tests. Read every repository's CLAUDE.md first and follow it
(parity with @hebcal/core, the deliberate divergences, GPL headers, gofmt, the
Go version in each go.mod, port provenance kept out of doc comments in api-go,
byte-identical golden tests in api-go).

SETUP
- Create a go.work in the parent directory that `use`s all four modules, so the
  consumers compile against your local hebcal-go without a release. Do not
  commit go.work. The CLI vendors its dependencies: build it with
  `-mod=mod` under the workspace, and re-vendor at the end.
- Record a baseline first: `go vet ./... && go test ./...` in each repository
  (api-go: `make test`). Every test that passes now must still pass, and
  api-go's golden tests (yahrzeit, staticcal, zmanim feed, CSV) must stay
  byte-identical. UIDs depend on event descriptions, so be especially careful
  that each event's Desc() keeps returning exactly what event.Desc(ev) returns
  today. In particular, learning events must return their current Basename().

PHASE 1: interface and correctness (one hebcal-go minor release)
1. In package event, add Desc() string and URL() string to CalEvent. Implement
   them on every event type in hebcal-go and learning (11 types). Desc()
   returns exactly what event.Desc(ev) returns today. URL() returns exactly what
   event.URL(ev) returns today ("" for no page). Keep Descer, URLer, event.Desc
   and event.URL for this release only, as `// Deprecated:` aliases/wrappers.
   Keep Memoer, BriefRenderer and their helpers unchanged (they stay optional).
2. Add exported interfaces in package event:
     type Timed interface { CalEvent; Time() time.Time }
     type Linked interface { LinkedEvent() CalEvent }
   Stop embedding event.HolidayEvent in hebcal.TimedEvent: give it explicit
   fields and methods. Stop embedding TimedEvent in KiddushLevanaEvent in a way
   that hides it from callers: both must satisfy event.Timed. Keep
   hebcal.AsTimedEvent as a deprecated wrapper. Every behaviour TimedEvent used
   to inherit implicitly (URL, Memo, RenderBrief, GetDesc, Basename,
   GetCategories) must be explicit, with today's exact output.
3. Add `type TimedKind int` (CandleLighting, Havdalah, FastBegins, FastEnds,
   FinishChametz, BiurChametz, KiddushLevana, Zman, SunriseSunset …) on
   TimedEvent. Replace the string comparisons on Desc inside hebcal-go with
   switches on Kind. Desc stays as the translation key.
4. HebrewCalendar must not modify *CalOptions. Copy it on entry, normalize into
   an unexported struct, and stop storing *CalOptions in events: store only
   what rendering needs (the clock format). Add a test that runs HebrewCalendar
   concurrently on one shared options value under -race, and a test that the
   caller's struct is unchanged after the call.
5. Replace CalOptions.Hour12/Hour24 and the *bool hour12 parameters with
   `type ClockFormat int` (ClockAuto, Clock12, Clock24) wherever they are used
   (CalOptions, MoladEvent, MevarchimChodeshEvent, ReformatTimeStr). Merge
   NewMoladEvent and NewMoladEventHour12 into one constructor.
6. Export ParshaEvent and HebrewDateEvent. Make NewParshaEvent,
   NewHebrewDateEvent and NewMoladEvent return their concrete types.
7. Migrate the consumers: learning (implement Desc/URL), hebcal-api-go (replace
   type assertions on hebcal.TimedEvent with event.Timed; this fixes the missed
   KiddushLevanaEvent at internal/service/shabbat/shabbat.go and
   internal/service/mcp/tools.go; collapse download.UntranslatedDesc and the
   Desc/URL parts of model.CalEv into the new methods; replace Desc-string
   matching with Kind where it is a TimedEvent), and the CLI (main.go's
   TimedEvent assertion).
8. Update the CLAUDE.md files and READMEs that describe these APIs (hebcal-go's
   event-model section, api-go's notes on AsTimedEvent and the Kiddush Levana
   pitfall, learning's recipe).

STOP after phase 1. Show me the diff summary per repository and the test
results, and wait for approval before phase 2.

PHASE 2: naming and robustness (the next minor release)
- Flag constants to MixedCaps (CHAG → Chag, LIGHT_CANDLES → LightCandles, …),
  keeping flagNames/String() output UNCHANGED (it may be logged or compared).
- Drop the Get prefix: Date(), Flags(), Emoji(), Categories(). Where this
  collides with exported struct fields (HolidayEvent.Date/Desc/Flags/Emoji,
  UserEvent, MoladEvent, …), propose a field strategy to me BEFORE doing it
  (unexport and add a constructor, or rename the fields). It is the largest
  change for consumers.
- zmanim.Location.TimeZoneId → TimeZoneID; IsAssurBemlacha → IsAssurBemelacha;
  hebcal.UserEvent (the input spec) → a non-colliding name; TimeOfDay → a typed
  DayPart; dailylearning GetCalendars/GetStartDate → Calendars/StartDate, with a
  sync.RWMutex on the registry, and no variadic optional argument in
  AddCalendar.
- Error-returning constructors: zmanim.NewLocation and zmanim.New return an
  error instead of panicking (keep Must… variants if that helps tests).
  Exported sentinel errors for HebrewCalendar's option validation.
- zmanim.AllCities returns a copy. Convert CalOptions field comments to Go doc
  style. Fix the package docs ("Package event …"), remove the stale dafyomi
  mention in hebcal/doc.go, and move the misplaced GetHolidaysForYear comment.
- Delete every deprecated alias added in phase 1.

RULES
- Do not change observable output: rendered strings, URLs, categories, flags'
  String(), event order, times. If a test's expected value would have to change,
  stop and ask me.
- Keep hebcal-go buildable on the Go version its go.mod declares (currently
  1.18) unless I approve a bump.
- Do not change the dependency direction (learning → hebcal-go).
- Commit each logical step separately, on the branch this session is told to
  use. Do not tag or release. Do not open PRs unless I ask.
- Before each commit, run gofmt -l, go vet ./... and go test ./... (api-go:
  make test) in every affected repository.
````
