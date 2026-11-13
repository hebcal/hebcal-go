# Migrating to hebcal-go v0.50.0

v0.50.0 changes the public API of `github.com/hebcal/hebcal-go` in ways that
break existing callers. The version number jumps from v0.27 to say so. This
module is still v0, so the import path does not change.

Nothing about the calendar itself changed. Every event has the same date,
description, flags, categories, URL, rendered text and time as before, and the
events come in the same order. Only the Go API around them is different.

Most of the upgrade is renaming. [Section 12](#12-quick-reference) is a table
of every rename, and [section 11](#11-doing-the-renames-mechanically) has
commands that do most of them for you.

## Contents

1. [Read this first: code that compiles but is wrong](#1-read-this-first-code-that-compiles-but-is-wrong)
2. [The `CalEvent` interface](#2-the-calevent-interface)
3. [Flags: `event.Flags` and MixedCaps constants](#3-flags-eventflags-and-mixedcaps-constants)
4. [Event fields became methods and constructors](#4-event-fields-became-methods-and-constructors)
5. [Timed events: `event.Timed`, `event.Linked`, `TimedKind`](#5-timed-events-eventtimed-eventlinked-timedkind)
6. [Clock format](#6-clock-format)
7. [`HebrewCalendar` no longer modifies your options](#7-hebrewcalendar-no-longer-modifies-your-options)
8. [Errors instead of panics](#8-errors-instead-of-panics)
9. [Other renames](#9-other-renames)
10. [If you implement `event.CalEvent` or register a schedule](#10-if-you-implement-eventcalevent-or-register-a-schedule)
11. [Doing the renames mechanically](#11-doing-the-renames-mechanically)
12. [Quick reference](#12-quick-reference)

## 1. Read this first: code that compiles but is wrong

Several exported fields are now methods with the same name. Most uses stop
compiling, which is what you want. Two kinds of use keep compiling and do the
wrong thing, because a method name without parentheses is a legal Go value:

```go
// Prints a function address instead of the description.
// go vet reports this one; plain `go build` does not.
fmt.Printf("%s\n", ev.Desc)

// Always true. The next line then calls Render on a nil event and panics.
if timed.LinkedEvent != nil {
    fmt.Println(timed.LinkedEvent().Render("en"))
}

// Compiles when the parameter is interface{} (testify, loggers, templates).
assert.Equal(t, "Havdalah", ev.Desc)
```

After upgrading, run `go vet ./...`, and search your code for these names
without a following parenthesis:

```bash
grep -rnE '\.(Desc|LinkedEvent|Date|Flags|Emoji)\b[^(]' --include='*.go' .
```

That search also finds fields of your own types with those names; only the
ones on hebcal-go events matter.

## 2. The `CalEvent` interface

`event.CalEvent` has eight methods. `Desc` and `URL` are new, and the four
getters lost their `Get` prefix:

```go
type CalEvent interface {
    Date() hdate.HDate           // was GetDate
    Desc() string                // new: untranslated description, the stable identifier
    Basename() string
    Render(locale string) string
    Flags() Flags                // was GetFlags() HolidayFlags
    Emoji() string               // was GetEmoji
    Categories() []string        // was GetCategories
    URL() string                 // new: "" when the event has no page
}
```

`Desc` and `URL` replace the optional interfaces `event.Descer` and
`event.URLer` and the helper functions that went with them, all of which are
removed:

| Before | After |
|---|---|
| `event.Desc(ev)` | `ev.Desc()` |
| `event.URL(ev)` | `ev.URL()` |
| `ev.(event.Descer)`, then `.GetDesc()` | `ev.Desc()` |
| `ev.(event.URLer)`, then `.URL()` | `ev.URL()` |
| `x.GetDesc()` on any event type | `x.Desc()` |
| `ev.GetDate()` | `ev.Date()` |
| `ev.GetFlags()` | `ev.Flags()` |
| `ev.GetEmoji()` | `ev.Emoji()` |
| `ev.GetCategories()` | `ev.Categories()` |

Each event returns exactly the strings it returned before. In particular the
daily-learning events of `github.com/hebcal/learning` return their
`Basename()` from `Desc()`, which is what `event.Desc` fell back to.

`event.Memoer` (`Memo`) and `event.BriefRenderer` (`RenderBrief`) are still
optional interfaces, and `event.Memo(ev, locale)` and
`event.RenderBrief(ev, locale)` work as before.

## 3. Flags: `event.Flags` and MixedCaps constants

`event.HolidayFlags` is now `event.Flags`, and its constants follow Go
naming:

| Before | After | | Before | After |
|---|---|---|---|---|
| `CHAG` | `Chag` | | `MOLAD` | `Molad` |
| `LIGHT_CANDLES` | `LightCandles` | | `USER_EVENT` | `User` |
| `YOM_TOV_ENDS` | `YomTovEnds` | | `HEBREW_DATE` | `HebrewDate` |
| `CHUL_ONLY` | `ChulOnly` | | `MINOR_HOLIDAY` | `MinorHoliday` |
| `IL_ONLY` | `ILOnly` | | `EREV` | `Erev` |
| `LIGHT_CANDLES_TZEIS` | `LightCandlesTzeis` | | `CHOL_HAMOED` | `CholHamoed` |
| `CHANUKAH_CANDLES` | `ChanukahCandles` | | `MISHNA_YOMI` | `MishnaYomi` |
| `ROSH_CHODESH` | `RoshChodesh` | | `YOM_KIPPUR_KATAN` | `YomKippurKatan` |
| `MINOR_FAST` | `MinorFast` | | `ZMANIM` | `Zmanim` |
| `SPECIAL_SHABBAT` | `SpecialShabbat` | | `YERUSHALMI_YOMI` | `YerushalmiYomi` |
| `PARSHA_HASHAVUA` | `ParshaHashavua` | | `NACH_YOMI` | `NachYomi` |
| `DAF_YOMI` | `DafYomi` | | `DAILY_LEARNING` | `DailyLearning` |
| `OMER_COUNT` | `OmerCount` | | `BEHAB` | `BeHaB` |
| `MODERN_HOLIDAY` | `ModernHoliday` | | `KIDDUSH_LEVANA` | `KiddushLevana` |
| `MAJOR_FAST` | `MajorFast` | | `YIZKOR` | `Yizkor` |
| `SHABBAT_MEVARCHIM` | `ShabbatMevarchim` | | | |

Two to notice: `USER_EVENT` is `event.User` (`event.UserEvent` is the event
type), and `IL_ONLY` is `event.ILOnly`.

**`Flags.String()` did not change.** It still returns the old names, such as
`"CHAG|LIGHT_CANDLES"` and `"USER_EVENT"`, so anything you log, store or
compare keeps working. The bit values are the same too.

`Has`, `HasAny`, `With` and `Without` are unchanged.

## 4. Event fields became methods and constructors

A type cannot have a field and a method of the same name, so the fields that
the interface reads are now unexported. Read them through the methods, and
build the event with its constructor.

### `event.HolidayEvent`

```go
// Before
ev := event.HolidayEvent{Date: hd, Desc: "Purim", Flags: event.MINOR_HOLIDAY, Emoji: "🎭"}
if ev.Desc == "Purim" && ev.Flags.Has(event.MINOR_HOLIDAY) { ... }
day := ev.Date

// After
ev := event.NewHolidayEvent(hd, "Purim", event.MinorHoliday, "🎭")
if ev.Desc() == "Purim" && ev.Flags().Has(event.MinorHoliday) { ... }
day := ev.Date()
```

Pass `""` as the emoji to get the default that `Emoji()` chooses from the
flags. `CholHaMoedDay` and `ChanukahDay` are still exported fields; set them
after constructing the event.

`Emoji()` returns the emoji the event was built with when there is one, and
the flag-based default otherwise, as `GetEmoji()` always did. The raw field
is no longer readable.

### `event.UserEvent`

```go
// Before
ev := event.UserEvent{Date: hd, Desc: "Birthday"}
// After
ev := event.NewUserEvent(hd, "Birthday")
```

### `Date` on the other event types

`MoladEvent.Date`, `MevarchimChodeshEvent.Date` and `omer.OmerEvent.Date` are
now the method `Date()`. Their other fields (`Molad`, `MonthName`,
`CountryCode`, `OmerDay`, ...) are unchanged. Build them with `NewMoladEvent`,
`NewMevarchimChodeshEvent` and `omer.NewOmerEvent`.

`omer.OmerEvent.GetWeeks()` is now `Weeks()`.

### Constructors return concrete types

| Function | Returned before | Returns now |
|---|---|---|
| `event.NewParshaEvent` | `event.CalEvent` | `event.ParshaEvent` |
| `event.NewHebrewDateEvent` | `event.CalEvent` | `event.HebrewDateEvent` |
| `event.NewMoladEvent` | `event.CalEvent` | `event.MoladEvent` |

`ParshaEvent` and `HebrewDateEvent` are newly exported, so you can assert them
and read `ParshaEvent.Parsha` and `ParshaEvent.IL`:

```go
if pe, ok := ev.(event.ParshaEvent); ok {
    fmt.Println(pe.Parsha.Name, pe.IL)
}
```

Code that assigned the result to a `CalEvent` variable still compiles. Code
that declared a variable with `:=` and later assigned a different event type
to it needs an explicit `event.CalEvent` type.

## 5. Timed events: `event.Timed`, `event.Linked`, `TimedKind`

`hebcal.TimedEvent` no longer embeds `event.HolidayEvent`, and
`hebcal.KiddushLevanaEvent` no longer embeds `TimedEvent`. Both are found
through an interface:

```go
type Timed interface {   // package event
    CalEvent
    Time() time.Time
}

type Linked interface {  // package event
    LinkedEvent() CalEvent
}
```

### Finding timed events

```go
// Before: silently skipped "Latest Kiddush Levana"
if te, ok := ev.(hebcal.TimedEvent); ok {
    fmt.Println(te.EventTime)
}
// or
if te, ok := hebcal.AsTimedEvent(ev); ok {
    fmt.Println(te.EventTime)
}

// After: every timed event, Kiddush Levana included
if te, ok := ev.(event.Timed); ok {
    fmt.Println(te.Time())
}
```

`hebcal.AsTimedEvent` is removed.

If your code asserted `hebcal.TimedEvent` and you generate Kiddush Levana
events (`CalOptions.KiddushLevanaMaharil`), switching to `event.Timed` makes
that event take your timed path for the first time. That is usually the fix
you want, but check it.

### Field renames on `TimedEvent`

| Before | After |
|---|---|
| `te.EventTime` | `te.Time()` |
| `te.LinkedEvent` | `te.LinkedEvent()` (may return nil) |
| `te.Desc` | `te.Desc()` |
| `te.Date` | `te.Date()` |
| `te.Flags` | `te.Flags()` |
| `te.Emoji` | `te.Emoji()` |
| `te.ChanukahDay` | unchanged |
| `te.CholHaMoedDay` | removed; it was always 0 |

`KiddushLevanaEvent` has the same methods, plus its `Molad` and `AdjustedTo`
fields as before. `YizkorEvent.Date` and `YizkorEvent.LinkedEvent` are the
methods `Date()` and `LinkedEvent()`.

### Linked events

`TimedEvent` and `YizkorEvent` satisfy `event.Linked`:

```go
if l, ok := ev.(event.Linked); ok && l.LinkedEvent() != nil {
    fmt.Println(ev.Desc(), "belongs to", l.LinkedEvent().Desc())
}
```

### Which timed event is it?

Switch on the kind instead of comparing description strings:

```go
// Before
if te.Desc == "Candle lighting" || te.Desc == "Havdalah" { ... }

// After
switch kind, _ := hebcal.TimedKindOf(ev); kind {
case hebcal.CandleLighting, hebcal.Havdalah:
    ...
}
```

The kinds are `CandleLighting`, `Havdalah`, `ChanukahCandles`, `FastBegins`,
`FastEnds`, `FinishChametz`, `BiurChametz`, `KiddushLevana`, `Zman` (the daily
zmanim of `CalOptions.DailyZmanim`) and the zero value `TimedOther`.
`TimedEvent.Kind()` and `KiddushLevanaEvent.Kind()` return it directly;
`hebcal.TimedKindOf(ev)` reads it from any `CalEvent` and reports false for an
event that is not timed.

Do not confuse the kinds with the flags of the same name: `hebcal.KiddushLevana`
is a `TimedKind`, `event.KiddushLevana` is a flag.

`Desc()` is unchanged and remains the translation key and the stable
identifier. Comparing it still works; the kind is simply the safer thing to
switch on.

One case to know about: Chanukah candle-lighting has the kind
`ChanukahCandles` and the holiday's description ("Chanukah: 3 Candles"). With
`NoHolidays` set, the Friday and Saturday candle-lighting during Chanukah is
an ordinary `CandleLighting` whose flags include `event.ChanukahCandles`. That
is how it behaved before too.

### `NewTimedEvent`

The signature is unchanged. Its `*CalOptions` argument now only chooses the
clock format; the event does not keep the pointer. It may be nil.

## 6. Clock format

Three overlapping settings became one type, `event.ClockFormat`, with the
values `event.ClockAuto` (follow the country, the zero value), `event.Clock12`
and `event.Clock24`.

| Before | After |
|---|---|
| `CalOptions{Hour24: true}` | `CalOptions{Clock: event.Clock24}` |
| `CalOptions{Hour12: true}` | `CalOptions{Clock: event.Clock12}` |
| neither set | `Clock` left at its zero value |
| `event.ReformatTimeStr(s, suffix, cc, nil)` | `event.ReformatTimeStr(s, suffix, cc, event.ClockAuto)` |
| `event.ReformatTimeStr(s, suffix, cc, &yes)` | `event.ReformatTimeStr(s, suffix, cc, event.Clock12)` |
| `event.ReformatTimeStr(s, suffix, cc, &no)` | `event.ReformatTimeStr(s, suffix, cc, event.Clock24)` |
| `event.RenderMolad(m, month, locale, cc, hour12)` | same, with a `ClockFormat` |
| `MoladEvent.Hour12`, `MevarchimChodeshEvent.Hour12` (`*bool`) | `.Clock` (`ClockFormat`) |
| `event.NewMoladEvent(hd, m, month, cc)` | `event.NewMoladEvent(hd, m, month, cc, event.ClockAuto)` |
| `event.NewMoladEventHour12(hd, m, month, cc, hour12)` | `event.NewMoladEvent(hd, m, month, cc, clock)` |

If you had a `*bool` to convert:

```go
func clockFormat(hour12 *bool) event.ClockFormat {
    switch {
    case hour12 == nil:
        return event.ClockAuto
    case *hour12:
        return event.Clock12
    }
    return event.Clock24
}
```

`hebcal.ReformatTimeStr(timeStr, suffix, opts)` keeps its signature and reads
`opts.Clock`.

## 7. `HebrewCalendar` no longer modifies your options

`HebrewCalendar(opts)` used to write its normalizations back into `*opts`, and
every timed event kept a pointer to that struct. It now works on a private
copy, and events keep only the clock format they were generated with. A nil
`opts` is accepted and means the zero `CalOptions`.

This is safer (one options value can be shared by concurrent calls), but two
habits no longer work.

**Reading normalized values back.** After the call, your struct is exactly as
you left it. These fields used to change and no longer do:

| Field | What `HebrewCalendar` used to write |
|---|---|
| `IL` | `true` when `Location.CountryCode == "IL"` |
| `Mask` | the mask derived from the boolean options |
| `CandleLightingMins` | negated, with the default 18 and the Israeli city customs applied |
| `HavdalahDeg` | 8.5 when neither Havdalah option was set |
| `HavdalahMins`, the fast offsets | their absolute values |
| `KiddushLevanaMaharil`, `Sedrot`, `Omer`, `DafYomi`, ... | `true` when `Mask` asked for them |
| `NoRoshChodesh`, `NoModern`, `NoMinorFast`, `NoSpecialShabbat` | `false` when `Mask` asked for them |
| `YerushalmiEdition` | `Vilna` when unset |
| `CandleLighting` | `false` for years before 1753 (Hebrew 5514) |

If you need one of these, compute it yourself. For the candle-lighting offset
of an Israeli location, `hebcal.IsraelCandleLightingMins(loc)` is exported.

**Changing the options to change rendering.** This no longer has any effect:

```go
events, _ := hebcal.HebrewCalendar(&opts)
opts.Hour24 = true           // before: every TimedEvent now rendered 24-hour
fmt.Println(events[0].Render("en"))
```

Set `opts.Clock` before calling `HebrewCalendar` instead.

## 8. Errors instead of panics

### `zmanim.New` and `zmanim.NewLocation`

Both used to panic on a coordinate out of range, and `New` on an unknown time
zone. Both now return an error:

```go
// Before
loc := zmanim.NewLocation("Home", "US", lat, long, 0, tzid)   // might panic
z := zmanim.New(&loc, date)                                   // might panic

// After
loc, err := zmanim.NewLocation("Home", "US", lat, long, 0, tzid)
if err != nil {
    return err
}
z, err := zmanim.New(&loc, date)
if err != nil {
    return err
}
```

The errors are `zmanim.ErrLatitude`, `zmanim.ErrLongitude` and (from `New`,
wrapped with the identifier) `zmanim.ErrTimeZone`; test with `errors.Is`.

If the location is known to be good (a `zmanim.LookupCity` result, a constant
in a test), `zmanim.MustNewLocation` and `zmanim.MustNew` keep the old
one-value form and panic on failure.

### `HebrewCalendar`

`HebrewCalendar` checks the location before generating anything, when the
calendar needs times from it, and returns the zmanim error. An unusable
location used to panic part-way through.

Its own option errors are now exported values, so you can test for them
instead of comparing message text. The messages themselves are unchanged.

| Error | When | Also matches |
|---|---|---|
| `ErrCandleLightingLocation` | `CandleLighting` without `Location` | `ErrLocationRequired` |
| `ErrKiddushLevanaLocation` | `KiddushLevanaMaharil` without `Location` | `ErrLocationRequired` |
| `ErrSunriseSunsetLocation` | `SunriseSunset` without `Location` | `ErrLocationRequired` |
| `ErrDailyZmanimLocation` | `DailyZmanim` without `Location` | `ErrLocationRequired` |
| `ErrHavdalahMinsAndDeg` | both Havdalah offsets set | `ErrMutuallyExclusive` |
| `ErrFastStartMinsAndDeg` | both fast-start offsets set | `ErrMutuallyExclusive` |
| `ErrFastEndMinsAndDeg` | both fast-end offsets set | `ErrMutuallyExclusive` |
| `ErrTishaBavEndMinsAndDeg` | both Tish'a B'Av offsets set | `ErrMutuallyExclusive` |
| `ErrStartRequiresEnd` | only one of `Start` and `End` set | |
| `ErrInvalidHebrewYear` | `IsHebrewYear` with `Year` below 1 | |

```go
_, err := hebcal.HebrewCalendar(&opts)
if errors.Is(err, hebcal.ErrLocationRequired) {
    // ask the user for a location
}
```

`hebcal.IsAssurBemelacha` (see below) also returns an error for an unusable
location instead of panicking.

## 9. Other renames

| Before | After | Notes |
|---|---|---|
| `zmanim.Location.TimeZoneId` | `TimeZoneID` | field rename; also in struct literals |
| `hebcal.IsAssurBemlacha` | `hebcal.IsAssurBemelacha` | spelling |
| `hebcal.UserEvent` | `hebcal.UserEventSpec` | the *input* in `CalOptions.UserEvents`; `event.UserEvent` is the event it produces. The `UserEvents` field keeps its name. |
| `hebcal.TimeOfDay(ev) int` | `hebcal.DayPartOf(ev) hebcal.DayPart` | typed result |
| `hebcal.TimeOfDayDawn` ... `TimeOfDayNight` | `hebcal.DayPartDawn` ... `DayPartNight` | same values, same order |
| `dailylearning.GetCalendars()` | `dailylearning.Calendars()` | |
| `dailylearning.GetStartDate(name)` | `dailylearning.StartDate(name)` | |
| `dailylearning.AddCalendar(name, fn, startDate)` | `dailylearning.AddCalendarFrom(name, fn, startDate)` | see section 10 |

If you stored a `TimeOfDay` result in an `int` field or variable, change its
type to `hebcal.DayPart` (or convert with `int(...)`). A sentinel such as `-1`
for "before everything" still works as an untyped constant.

`zmanim.AllCities()` now returns a copy of the built-in list. If you relied on
changing the list in place to change what `LookupCity` returns, that no longer
works.

`hebcal.GetHolidaysForYear` and `hebcal.GetHolidaysOnDate` keep their names.

## 10. If you implement `event.CalEvent` or register a schedule

### Event types

Rename your four getters and add the two new methods:

```go
func (ev myEvent) Date() hdate.HDate    { return ev.date }   // was GetDate
func (ev myEvent) Flags() event.Flags   { return event.DailyLearning } // was GetFlags
func (ev myEvent) Emoji() string        { return "" }        // was GetEmoji
func (ev myEvent) Categories() []string { return []string{"mySchedule"} } // was GetCategories

// Desc returns the untranslated description. If you had no GetDesc method,
// return what Basename returns: that is what event.Desc(ev) gave callers.
func (ev myEvent) Desc() string { return ev.Basename() }

// URL returns the event's page, or "" if it has none.
func (ev myEvent) URL() string { return "" }
```

If your struct has a field named `Date`, `Flags`, `Emoji` or `Categories`,
rename the field (for example to `date`): it now collides with the method.

If your type had a `GetDesc()` method, rename it to `Desc()`. If it embeds
`event.HolidayEvent` to inherit methods, note that it inherits `Desc()` and
`URL()` too, and the URL will be a hebcal.com holiday page built from your
description. Override `URL()` if that is not what you want. We stopped
embedding event types inside this module for that reason.

### Registering a schedule

`AddCalendar` no longer takes an optional start date:

```go
// Before
dailylearning.AddCalendar("mySchedule", fn)
dailylearning.AddCalendar("mySchedule", fn, startDate)

// After
dailylearning.AddCalendar("mySchedule", fn)
dailylearning.AddCalendarFrom("mySchedule", fn, startDate)
```

The registry is now safe to use from several goroutines, including
registering while calendars are being generated.

`github.com/hebcal/learning` has been updated; use a release of it made for
hebcal-go v0.50.0 or later.

## 11. Doing the renames mechanically

`gofmt -r` can do the method and constant renames. Run these from the root of
your module, with a clean working tree, and review the diff. They rewrite
every identifier of that name, so check that you have no unrelated methods
called `GetDate`, `GetFlags`, `GetEmoji` or `GetCategories` (protobuf-generated
getters are a common example; exclude those files or revert them).

```bash
gofmt -l -w -r 'a.GetDate() -> a.Date()' .
```

```bash
gofmt -l -w -r 'a.GetFlags() -> a.Flags()' .
```

```bash
gofmt -l -w -r 'a.GetEmoji() -> a.Emoji()' .
```

```bash
gofmt -l -w -r 'a.GetCategories() -> a.Categories()' .
```

```bash
gofmt -l -w -r 'event.Desc(a) -> a.Desc()' .
```

```bash
gofmt -l -w -r 'event.URL(a) -> a.URL()' .
```

```bash
gofmt -l -w -r 'event.HolidayFlags -> event.Flags' .
```

The flag constants, one command for all of them:

```bash
for pair in CHAG:Chag LIGHT_CANDLES:LightCandles YOM_TOV_ENDS:YomTovEnds \
  CHUL_ONLY:ChulOnly IL_ONLY:ILOnly LIGHT_CANDLES_TZEIS:LightCandlesTzeis \
  CHANUKAH_CANDLES:ChanukahCandles ROSH_CHODESH:RoshChodesh MINOR_FAST:MinorFast \
  SPECIAL_SHABBAT:SpecialShabbat PARSHA_HASHAVUA:ParshaHashavua DAF_YOMI:DafYomi \
  OMER_COUNT:OmerCount MODERN_HOLIDAY:ModernHoliday MAJOR_FAST:MajorFast \
  SHABBAT_MEVARCHIM:ShabbatMevarchim MOLAD:Molad USER_EVENT:User \
  HEBREW_DATE:HebrewDate MINOR_HOLIDAY:MinorHoliday EREV:Erev \
  CHOL_HAMOED:CholHamoed MISHNA_YOMI:MishnaYomi YOM_KIPPUR_KATAN:YomKippurKatan \
  ZMANIM:Zmanim YERUSHALMI_YOMI:YerushalmiYomi NACH_YOMI:NachYomi \
  DAILY_LEARNING:DailyLearning BEHAB:BeHaB KIDDUSH_LEVANA:KiddushLevana YIZKOR:Yizkor
do
  gofmt -l -w -r "event.${pair%%:*} -> event.${pair##*:}" .
done
```

`gofmt -r` only rewrites expressions. It does not touch method declarations
(`func (ev myEvent) GetDate()`), struct literals with removed fields, or
comments; the compiler will point at the first two.

Then fix what remains by following the compiler: field reads that became
methods (sections 4 and 5), `zmanim.New` and `NewLocation` now returning two
values (section 8), and the `Hour12`/`Hour24` options (section 6). Finish with
`go vet ./...` and the search in section 1.

## 12. Quick reference

| Before | After |
|---|---|
| `ev.GetDate()`, `ev.GetFlags()`, `ev.GetEmoji()`, `ev.GetCategories()` | `ev.Date()`, `ev.Flags()`, `ev.Emoji()`, `ev.Categories()` |
| `event.Desc(ev)` / `event.URL(ev)` | `ev.Desc()` / `ev.URL()` |
| `x.GetDesc()` | `x.Desc()` |
| `event.Descer`, `event.URLer` | removed; use `event.CalEvent` |
| `event.HolidayFlags` | `event.Flags` |
| `event.CHAG`, `event.LIGHT_CANDLES`, ... | `event.Chag`, `event.LightCandles`, ... (table in section 3) |
| `event.USER_EVENT` | `event.User` |
| `event.HolidayEvent{Date: d, Desc: s, Flags: f, Emoji: e}` | `event.NewHolidayEvent(d, s, f, e)` |
| `event.UserEvent{Date: d, Desc: s}` | `event.NewUserEvent(d, s)` |
| `he.Date`, `he.Desc`, `he.Flags`, `he.Emoji` | `he.Date()`, `he.Desc()`, `he.Flags()`, `he.Emoji()` |
| `omerEv.GetWeeks()` | `omerEv.Weeks()` |
| `ev.(hebcal.TimedEvent)`, `hebcal.AsTimedEvent(ev)` | `ev.(event.Timed)` |
| `te.EventTime` | `te.Time()` |
| `te.LinkedEvent`, `yizkor.LinkedEvent` | `ev.(event.Linked)`, then `.LinkedEvent()` |
| `te.Desc == "Havdalah"` | `kind, _ := hebcal.TimedKindOf(ev); kind == hebcal.Havdalah` |
| `CalOptions.Hour24` / `.Hour12` | `CalOptions.Clock` |
| `*bool` hour12 arguments and fields | `event.ClockFormat` |
| `event.NewMoladEventHour12(...)` | `event.NewMoladEvent(..., clock)` |
| reading `opts.IL`, `opts.Mask`, ... after `HebrewCalendar` | no longer written; compute them yourself |
| `zmanim.NewLocation(...)` (one result) | `loc, err := zmanim.NewLocation(...)` or `zmanim.MustNewLocation(...)` |
| `zmanim.New(loc, date)` (one result) | `z, err := zmanim.New(loc, date)` or `zmanim.MustNew(loc, date)` |
| comparing `err.Error()` from `HebrewCalendar` | `errors.Is(err, hebcal.ErrLocationRequired)` and the others in section 8 |
| `loc.TimeZoneId` | `loc.TimeZoneID` |
| `hebcal.IsAssurBemlacha` | `hebcal.IsAssurBemelacha` |
| `hebcal.UserEvent` | `hebcal.UserEventSpec` |
| `hebcal.TimeOfDay(ev)`, `hebcal.TimeOfDayDawn` ... | `hebcal.DayPartOf(ev)`, `hebcal.DayPartDawn` ... |
| `dailylearning.GetCalendars()` / `GetStartDate(n)` | `dailylearning.Calendars()` / `StartDate(n)` |
| `dailylearning.AddCalendar(n, fn, start)` | `dailylearning.AddCalendarFrom(n, fn, start)` |
