# Migrating to hebcal-go v0.50.0

v0.50.0 changes the public API of `github.com/hebcal/hebcal-go` in ways that
break existing callers. The version number jumps from v0.27 to say so. This
module is still v0, so the import path does not change.

Nothing about the calendar itself changed. Every event has the same date,
description, flags, categories, URL, rendered text and time as before, and the
events come in the same order. Only the Go API around them is different.

If you only call `hebcal.HebrewCalendar` and print `ev.Render(locale)`, the one
thing you may have to change is the 12/24-hour option ([section 5](#5-clock-format)).

## Contents

1. [Read this first: code that compiles but is wrong](#1-read-this-first-code-that-compiles-but-is-wrong)
2. [`Desc()` and `URL()` are part of `CalEvent`](#2-desc-and-url-are-part-of-calevent)
3. [Event fields became methods and constructors](#3-event-fields-became-methods-and-constructors)
4. [Timed events: `event.Timed`, `event.Linked`, `TimedKind`](#4-timed-events-eventtimed-eventlinked-timedkind)
5. [Clock format](#5-clock-format)
6. [`HebrewCalendar` no longer modifies your options](#6-hebrewcalendar-no-longer-modifies-your-options)
7. [Constructors return concrete types](#7-constructors-return-concrete-types)
8. [If you implement `event.CalEvent` yourself](#8-if-you-implement-eventcalevent-yourself)
9. [Deprecated, still present](#9-deprecated-still-present)
10. [Quick reference](#10-quick-reference)

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
grep -rnE '\.(Desc|LinkedEvent)\b[^(]' --include='*.go' .
```

## 2. `Desc()` and `URL()` are part of `CalEvent`

`event.CalEvent` now has eight methods. `Desc` and `URL` are new:

```go
type CalEvent interface {
    GetDate() hdate.HDate
    Desc() string                // untranslated description, the stable identifier
    Basename() string
    Render(locale string) string
    GetFlags() HolidayFlags
    GetEmoji() string
    GetCategories() []string
    URL() string                 // "" when the event has no page
}
```

They replace the optional interfaces `event.Descer` and `event.URLer` and the
helper functions that went with them:

| Before | After |
|---|---|
| `event.Desc(ev)` | `ev.Desc()` |
| `event.URL(ev)` | `ev.URL()` |
| `ev.(event.Descer)`, then `.GetDesc()` | `ev.Desc()` |
| `ev.(event.URLer)`, then `.URL()` | `ev.URL()` |
| `x.GetDesc()` on any event type | `x.Desc()` |

Each event returns exactly the string the helper returned before. In
particular the daily-learning events of `github.com/hebcal/learning` return
their `Basename()`, which is what `event.Desc` fell back to.

`event.Memoer` (`Memo`) and `event.BriefRenderer` (`RenderBrief`) are still
optional interfaces, and `event.Memo(ev, locale)` and
`event.RenderBrief(ev, locale)` work as before.

## 3. Event fields became methods and constructors

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
ev := event.NewHolidayEvent(hd, "Purim", event.MINOR_HOLIDAY, "🎭")
if ev.Desc() == "Purim" && ev.GetFlags().Has(event.MINOR_HOLIDAY) { ... }
day := ev.GetDate()
```

Pass `""` as the emoji to get the default that `GetEmoji()` chooses from the
flags. `CholHaMoedDay` and `ChanukahDay` are still exported fields; set them
after constructing the event.

There was no exported way to read the raw `Emoji` field apart from the field
itself. `GetEmoji()` returns it when it is set, and the flag-based default
otherwise, as it always did.

### `event.UserEvent`

```go
// Before
ev := event.UserEvent{Date: hd, Desc: "Birthday"}
// After
ev := event.NewUserEvent(hd, "Birthday")
```

### `Date` on the other event types

`MoladEvent.Date`, `MevarchimChodeshEvent.Date` and `omer.OmerEvent.Date` are
now read with `GetDate()`. Their other fields (`Molad`, `MonthName`,
`CountryCode`, `OmerDay`, ...) are unchanged. Build them with `NewMoladEvent`,
`NewMevarchimChodeshEvent` and `omer.NewOmerEvent`.

### `hebcal.YizkorEvent`

`YizkorEvent.Date` is `GetDate()` and `YizkorEvent.LinkedEvent` is the method
`LinkedEvent()`.

## 4. Timed events: `event.Timed`, `event.Linked`, `TimedKind`

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
| `te.Date` | `te.GetDate()` |
| `te.Flags` | `te.GetFlags()` |
| `te.Emoji` | `te.GetEmoji()` |
| `te.ChanukahDay` | unchanged |
| `te.CholHaMoedDay` | removed; it was always 0 |

`KiddushLevanaEvent` has the same methods, plus its `Molad` and `AdjustedTo`
fields as before.

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

`Desc()` is unchanged and remains the translation key and the stable
identifier. Comparing it still works; the kind is simply the safer thing to
switch on.

One case to know about: Chanukah candle-lighting has the kind
`ChanukahCandles` and the holiday's description ("Chanukah: 3 Candles"). With
`NoHolidays` set, the Friday and Saturday candle-lighting during Chanukah is
an ordinary `CandleLighting` whose flags include `CHANUKAH_CANDLES`. That is
how it behaved before too.

### `NewTimedEvent`

The signature is unchanged. Its `*CalOptions` argument now only chooses the
clock format; the event does not keep the pointer. It may be nil.

## 5. Clock format

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

## 6. `HebrewCalendar` no longer modifies your options

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

## 7. Constructors return concrete types

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
that declared a slice or variable with `:=` and later assigned a different
event type to it needs an explicit `event.CalEvent` type.

## 8. If you implement `event.CalEvent` yourself

Add the two new methods to each of your event types:

```go
// Desc returns the untranslated description. If you had no GetDesc method,
// return what Basename returns: that is what event.Desc(ev) gave callers.
func (ev myEvent) Desc() string { return ev.Basename() }

// URL returns the event's page, or "" if it has none.
func (ev myEvent) URL() string { return "" }
```

If your type had a `GetDesc()` method, rename it to `Desc()`. If it embeds
`event.HolidayEvent` to inherit methods, it keeps compiling, but note that it
inherits `Desc()` and `URL()` too, and the URL will be a hebcal.com holiday
page built from your description. Override `URL()` if that is not what you
want. We stopped embedding event types inside this module for that reason.

`github.com/hebcal/learning` has been updated; use a release of it made for
hebcal-go v0.50.0 or later.

## 9. Deprecated, still present

These still work in v0.50.0 and will be removed in a later release:

| Deprecated | Use instead |
|---|---|
| `event.Desc(ev)` | `ev.Desc()` |
| `event.URL(ev)` | `ev.URL()` |
| `event.Descer` | `event.CalEvent` (note: its method is now `Desc()`, not `GetDesc()`) |
| `event.URLer` | `event.CalEvent` |
| `hebcal.AsTimedEvent(ev)` | `ev.(event.Timed)` |

## 10. Quick reference

| Before | After |
|---|---|
| `event.Desc(ev)` / `event.URL(ev)` | `ev.Desc()` / `ev.URL()` |
| `x.GetDesc()` | `x.Desc()` |
| `event.HolidayEvent{Date: d, Desc: s, Flags: f, Emoji: e}` | `event.NewHolidayEvent(d, s, f, e)` |
| `event.UserEvent{Date: d, Desc: s}` | `event.NewUserEvent(d, s)` |
| `he.Date`, `he.Desc`, `he.Flags`, `he.Emoji` | `he.GetDate()`, `he.Desc()`, `he.GetFlags()`, `he.GetEmoji()` |
| `ev.(hebcal.TimedEvent)`, `hebcal.AsTimedEvent(ev)` | `ev.(event.Timed)` |
| `te.EventTime` | `te.Time()` |
| `te.LinkedEvent`, `yizkor.LinkedEvent` | `ev.(event.Linked)`, then `.LinkedEvent()` |
| `te.Desc == "Havdalah"` | `kind, _ := hebcal.TimedKindOf(ev); kind == hebcal.Havdalah` |
| `CalOptions.Hour24` / `.Hour12` | `CalOptions.Clock` |
| `*bool` hour12 arguments and fields | `event.ClockFormat` |
| `event.NewMoladEventHour12(...)` | `event.NewMoladEvent(..., clock)` |
| reading `opts.IL`, `opts.Mask`, ... after `HebrewCalendar` | no longer written; compute them yourself |
