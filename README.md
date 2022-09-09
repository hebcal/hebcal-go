# hebcal-go

[![Build Status](https://app.travis-ci.com/hebcal/hebcal-go.svg?branch=main)](https://app.travis-ci.com/hebcal/hebcal-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/hebcal/hebcal-go)](https://goreportcard.com/report/github.com/hebcal/hebcal-go)
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

This golang implementation was released in 2022 by Michael J. Radwin.

Many users of this library will utilize the HebrewCalendar and HDate
interfaces.
