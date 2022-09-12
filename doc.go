/*
Package hebcal provides functionality for converting between Hebrew
and Gregorian dates.

It also generates lists of Jewish holidays for any year (past,
present or future).

Shabbat and holiday candle lighting and havdalah times are
approximated based on location.

Torah readings (Parashat HaShavua), Daf Yomi, and counting of the Omer
can also be specified.

Hebcal also includes algorithms to calculate yahrzeits, birthdays and
anniversaries.

Hebcal includes several sub-packages:

  - hdate: converts between Hebrew and Gregorian dates.
    Also includes functions for calculating personal anniversaries
    (Yahrzeit, Birthday) according to the Hebrew calendar.

  - dafyomi: Daf Yomi, a daily regimen of learning the Talmud.

  - greg: converts between Gregorian dates and R.D. (Rata Die)
    day numbers.

  - sedra: weekly Torah reading (Parashat HaShavua).

  - zmanim: calculates halachic times.
*/
package hebcal
