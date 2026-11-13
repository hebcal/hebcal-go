# hebcal-go

[![Build Status](https://github.com/hebcal/hebcal-go/actions/workflows/go.yml/badge.svg)](https://github.com/hebcal/hebcal-go/actions/workflows/go.yml)
[![GoDoc](https://pkg.go.dev/badge/github.com/hebcal/hebcal-go?status.svg)](https://pkg.go.dev/github.com/hebcal/hebcal-go)

Hebcal is a perpetual Jewish Calendar. This library converts between
Hebrew and Gregorian dates, and generates lists of Jewish holidays for
any year (past, present or future). Shabbat and holiday candle lighting
and havdalah times are approximated based on location. Torah readings
(Parashat HaShavua), Daf Yomi, and counting of the Omer can also be
specified. Hebcal also includes algorithms to calculate yahrzeits,
birthdays and anniversaries.

Hebcal was created in 1992 by Danny Sadinoff as a Unix/Linux program
written in C, inspired by similar functionality written in Emacs Lisp.
This Golang implementation was released in 2022 by Michael J. Radwin.

Many users of this library will utilize the HebrewCalendar and HDate
interfaces.

## Upgrading to v0.50

v0.50 changes the API in ways that break existing callers. The calendar's
output is unchanged; the Go API around it is cleaner:

  - `event.CalEvent` gained `Desc()` and `URL()`, and its getters lost their
    `Get` prefix: `Date()`, `Flags()`, `Emoji()`, `Categories()`.
  - `event.HolidayFlags` is `event.Flags`, with MixedCaps constants
    (`event.Chag`, `event.LightCandles`), which `String()` prints too.
  - Several event fields became methods, and events are built with
    constructors such as `event.NewHolidayEvent`.
  - Timed events are found through the `event.Timed` interface and told apart
    with `hebcal.TimedKind`.
  - The 12/24-hour options became `event.ClockFormat`.
  - `hebcal.CalOptions` is `hebcal.Options`, and `HebrewCalendar` no longer
    modifies the options it is given.
  - `zmanim.New` and `zmanim.NewLocation` return an error instead of
    panicking, and `HebrewCalendar`'s option errors are exported values.

See [MIGRATION.md](MIGRATION.md) for every change, with before-and-after code
and commands that do most of the renaming. Read its first section even if your
code compiles after the upgrade: a few old field reads still compile and are
wrong.

Hebcal incorporates and uses several related packages:

  - dailylearning: a plugin registry for daily learning schedules
    (Daf Yomi, Mishna Yomi, etc.). The schedules themselves live in a
    separate module, github.com/hebcal/learning, which registers them
    here; import that module to enable daily learning events.
  - event: an interface for calendar events.
  - greg: converts between Gregorian dates and R.D. (Rata Die)
    day numbers.
  - hdate: converts between Hebrew and Gregorian dates.
    Also includes functions for calculating personal anniversaries
    (Yahrzeit, Birthday) according to the Hebrew calendar.
  - hebcal: provides functionality for calculating Jewish holidays,
    candle-lighting and havdalah times, and fast start/end times.
  - locales: translations and transliterations of Jewish holiday
    names into several languages.
  - molad: calculates the time at which the New Moon is born.
  - omer: calculates the Sefirat HaOmer.
  - sedra: weekly Torah reading (Parashat HaShavua).
  - zmanim: calculates halachic times.
