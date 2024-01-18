package hebcal_test

import (
	"fmt"
	"testing"

	"github.com/hebcal/hdate"
	"github.com/MaxBGreenberg/hebcal-go/event"
	"github.com/MaxBGreenberg/hebcal-go/hebcal"
	"github.com/stretchr/testify/assert"
)

func TestGetHolidaysForYearArrayDiaspora(t *testing.T) {
	events := hebcal.GetHolidaysForYear(5771, false)
	assert.Equal(t, 110, len(events))

	expected := []string{
		"2010-09-09 Rosh Hashana 5771",
		"2010-09-10 Rosh Hashana II",
		"2010-09-11 Shabbat Shuva",
		"2010-09-12 Tzom Gedaliah",
		"2010-09-17 Erev Yom Kippur",
		"2010-09-18 Yom Kippur",
		"2010-09-22 Erev Sukkot",
		"2010-09-23 Sukkot I",
		"2010-09-24 Sukkot II",
		"2010-09-25 Sukkot III (CH''M)",
		"2010-09-26 Sukkot IV (CH''M)",
		"2010-09-27 Sukkot V (CH''M)",
		"2010-09-28 Sukkot VI (CH''M)",
		"2010-09-29 Sukkot VII (Hoshana Raba)",
		"2010-09-30 Shmini Atzeret",
		"2010-10-01 Simchat Torah",
		"2010-10-02 Shabbat Mevarchim Chodesh Cheshvan",
		"2010-10-08 Rosh Chodesh Cheshvan",
		"2010-10-09 Rosh Chodesh Cheshvan",
		"2010-11-04 Yom Kippur Katan Kislev",
		"2010-11-06 Shabbat Mevarchim Chodesh Kislev",
		"2010-11-06 Sigd",
		"2010-11-07 Rosh Chodesh Kislev",
		"2010-11-08 Rosh Chodesh Kislev",
		"2010-12-01 Chanukah: 1 Candle",
		"2010-12-02 Chanukah: 2 Candles",
		"2010-12-03 Chanukah: 3 Candles",
		"2010-12-04 Chanukah: 4 Candles",
		"2010-12-04 Shabbat Mevarchim Chodesh Tevet",
		"2010-12-05 Chanukah: 5 Candles",
		"2010-12-06 Chanukah: 6 Candles",
		"2010-12-07 Chag HaBanot",
		"2010-12-07 Chanukah: 7 Candles",
		"2010-12-07 Rosh Chodesh Tevet",
		"2010-12-08 Chanukah: 8 Candles",
		"2010-12-08 Rosh Chodesh Tevet",
		"2010-12-09 Chanukah: 8th Day",
		"2010-12-17 Asara B'Tevet",
		"2011-01-01 Shabbat Mevarchim Chodesh Sh'vat",
		"2011-01-05 Yom Kippur Katan Sh'vat",
		"2011-01-06 Rosh Chodesh Sh'vat",
		"2011-01-15 Shabbat Shirah",
		"2011-01-20 Tu BiShvat",
		"2011-01-29 Shabbat Mevarchim Chodesh Adar I",
		"2011-02-03 Yom Kippur Katan Adar I",
		"2011-02-04 Rosh Chodesh Adar I",
		"2011-02-05 Rosh Chodesh Adar I",
		"2011-02-18 Purim Katan",
		"2011-02-19 Shushan Purim Katan",
		"2011-03-03 Yom Kippur Katan Adar II",
		"2011-03-05 Shabbat Mevarchim Chodesh Adar II",
		"2011-03-05 Shabbat Shekalim",
		"2011-03-06 Rosh Chodesh Adar II",
		"2011-03-07 Rosh Chodesh Adar II",
		"2011-03-17 Ta'anit Esther",
		"2011-03-19 Erev Purim",
		"2011-03-19 Shabbat Zachor",
		"2011-03-20 Purim",
		"2011-03-21 Shushan Purim",
		"2011-03-26 Shabbat Parah",
		"2011-04-02 Shabbat HaChodesh",
		"2011-04-02 Shabbat Mevarchim Chodesh Nisan",
		"2011-04-04 Yom Kippur Katan Nisan",
		"2011-04-05 Rosh Chodesh Nisan",
		"2011-04-16 Shabbat HaGadol",
		"2011-04-18 Erev Pesach",
		"2011-04-18 Ta'anit Bechorot",
		"2011-04-19 Pesach I",
		"2011-04-20 Pesach II",
		"2011-04-21 Pesach III (CH''M)",
		"2011-04-22 Pesach IV (CH''M)",
		"2011-04-23 Pesach V (CH''M)",
		"2011-04-24 Pesach VI (CH''M)",
		"2011-04-25 Pesach VII",
		"2011-04-26 Pesach VIII",
		"2011-04-30 Shabbat Mevarchim Chodesh Iyyar",
		"2011-05-02 Yom HaShoah",
		"2011-05-04 Rosh Chodesh Iyyar",
		"2011-05-05 Rosh Chodesh Iyyar",
		"2011-05-09 Yom HaZikaron",
		"2011-05-10 Yom HaAtzma'ut",
		"2011-05-18 Pesach Sheni",
		"2011-05-22 Lag BaOmer",
		"2011-05-28 Shabbat Mevarchim Chodesh Sivan",
		"2011-06-01 Yom Yerushalayim",
		"2011-06-02 Yom Kippur Katan Sivan",
		"2011-06-03 Rosh Chodesh Sivan",
		"2011-06-07 Erev Shavuot",
		"2011-06-08 Shavuot I",
		"2011-06-09 Shavuot II",
		"2011-06-25 Shabbat Mevarchim Chodesh Tamuz",
		"2011-06-30 Yom Kippur Katan Tamuz",
		"2011-07-02 Rosh Chodesh Tamuz",
		"2011-07-03 Rosh Chodesh Tamuz",
		"2011-07-19 Tzom Tammuz",
		"2011-07-30 Shabbat Mevarchim Chodesh Av",
		"2011-07-31 Yom Kippur Katan Av",
		"2011-08-01 Rosh Chodesh Av",
		"2011-08-06 Shabbat Chazon",
		"2011-08-08 Erev Tish'a B'Av",
		"2011-08-09 Tish'a B'Av",
		"2011-08-13 Shabbat Nachamu",
		"2011-08-15 Tu B'Av",
		"2011-08-27 Shabbat Mevarchim Chodesh Elul",
		"2011-08-29 Yom Kippur Katan Elul",
		"2011-08-30 Rosh Chodesh Elul",
		"2011-08-31 Rosh Chodesh Elul",
		"2011-08-31 Rosh Hashana LaBehemot",
		"2011-09-24 Leil Selichot",
		"2011-09-28 Erev Rosh Hashana",
	}

	actual := make([]string, 0, len(events))
	for _, ev := range events {
		line := fmt.Sprintf("%s %s", hd2iso(ev.Date), ev.Desc)
		actual = append(actual, line)
	}

	assert.Equal(t, expected, actual)
}

func TestGetHolidaysForYearArrayIL(t *testing.T) {
	events := hebcal.GetHolidaysForYear(5720, true)
	assert.Equal(t, 99, len(events))

	expected := []string{
		"1959-10-03 Rosh Hashana 5720",
		"1959-10-04 Rosh Hashana II",
		"1959-10-05 Tzom Gedaliah",
		"1959-10-10 Shabbat Shuva",
		"1959-10-11 Erev Yom Kippur",
		"1959-10-12 Yom Kippur",
		"1959-10-16 Erev Sukkot",
		"1959-10-17 Sukkot I",
		"1959-10-18 Sukkot II (CH''M)",
		"1959-10-19 Sukkot III (CH''M)",
		"1959-10-20 Sukkot IV (CH''M)",
		"1959-10-21 Sukkot V (CH''M)",
		"1959-10-22 Sukkot VI (CH''M)",
		"1959-10-23 Sukkot VII (Hoshana Raba)",
		"1959-10-24 Shmini Atzeret",
		"1959-10-31 Shabbat Mevarchim Chodesh Cheshvan",
		"1959-11-01 Rosh Chodesh Cheshvan",
		"1959-11-02 Rosh Chodesh Cheshvan",
		"1959-11-28 Shabbat Mevarchim Chodesh Kislev",
		"1959-11-30 Yom Kippur Katan Kislev",
		"1959-12-01 Rosh Chodesh Kislev",
		"1959-12-02 Rosh Chodesh Kislev",
		"1959-12-25 Chanukah: 1 Candle",
		"1959-12-26 Chanukah: 2 Candles",
		"1959-12-26 Shabbat Mevarchim Chodesh Tevet",
		"1959-12-27 Chanukah: 3 Candles",
		"1959-12-28 Chanukah: 4 Candles",
		"1959-12-29 Chanukah: 5 Candles",
		"1959-12-30 Chanukah: 6 Candles",
		"1959-12-31 Chag HaBanot",
		"1959-12-31 Chanukah: 7 Candles",
		"1959-12-31 Rosh Chodesh Tevet",
		"1960-01-01 Chanukah: 8 Candles",
		"1960-01-01 Rosh Chodesh Tevet",
		"1960-01-02 Chanukah: 8th Day",
		"1960-01-10 Asara B'Tevet",
		"1960-01-23 Shabbat Mevarchim Chodesh Sh'vat",
		"1960-01-28 Yom Kippur Katan Sh'vat",
		"1960-01-30 Rosh Chodesh Sh'vat",
		"1960-02-13 Shabbat Shirah",
		"1960-02-13 Tu BiShvat",
		"1960-02-25 Yom Kippur Katan Adar",
		"1960-02-27 Shabbat Mevarchim Chodesh Adar",
		"1960-02-27 Shabbat Shekalim",
		"1960-02-28 Rosh Chodesh Adar",
		"1960-02-29 Rosh Chodesh Adar",
		"1960-03-10 Ta'anit Esther",
		"1960-03-12 Erev Purim",
		"1960-03-12 Shabbat Zachor",
		"1960-03-13 Purim",
		"1960-03-14 Shushan Purim",
		"1960-03-19 Shabbat Parah",
		"1960-03-26 Shabbat HaChodesh",
		"1960-03-26 Shabbat Mevarchim Chodesh Nisan",
		"1960-03-28 Yom Kippur Katan Nisan",
		"1960-03-29 Rosh Chodesh Nisan",
		"1960-04-09 Shabbat HaGadol",
		"1960-04-11 Erev Pesach",
		"1960-04-11 Ta'anit Bechorot",
		"1960-04-12 Pesach I",
		"1960-04-13 Pesach II (CH''M)",
		"1960-04-14 Pesach III (CH''M)",
		"1960-04-15 Pesach IV (CH''M)",
		"1960-04-16 Pesach V (CH''M)",
		"1960-04-17 Pesach VI (CH''M)",
		"1960-04-18 Pesach VII",
		"1960-04-23 Shabbat Mevarchim Chodesh Iyyar",
		"1960-04-25 Yom HaShoah",
		"1960-04-27 Rosh Chodesh Iyyar",
		"1960-04-28 Rosh Chodesh Iyyar",
		"1960-05-01 Yom HaZikaron",
		"1960-05-02 Yom HaAtzma'ut",
		"1960-05-11 Pesach Sheni",
		"1960-05-15 Lag BaOmer",
		"1960-05-21 Shabbat Mevarchim Chodesh Sivan",
		"1960-05-26 Yom Kippur Katan Sivan",
		"1960-05-27 Rosh Chodesh Sivan",
		"1960-05-31 Erev Shavuot",
		"1960-06-01 Shavuot",
		"1960-06-18 Shabbat Mevarchim Chodesh Tamuz",
		"1960-06-23 Yom Kippur Katan Tamuz",
		"1960-06-25 Rosh Chodesh Tamuz",
		"1960-06-26 Rosh Chodesh Tamuz",
		"1960-07-12 Tzom Tammuz",
		"1960-07-23 Shabbat Mevarchim Chodesh Av",
		"1960-07-24 Yom Kippur Katan Av",
		"1960-07-25 Rosh Chodesh Av",
		"1960-07-30 Shabbat Chazon",
		"1960-08-01 Erev Tish'a B'Av",
		"1960-08-02 Tish'a B'Av",
		"1960-08-06 Shabbat Nachamu",
		"1960-08-08 Tu B'Av",
		"1960-08-20 Shabbat Mevarchim Chodesh Elul",
		"1960-08-22 Yom Kippur Katan Elul",
		"1960-08-23 Rosh Chodesh Elul",
		"1960-08-24 Rosh Chodesh Elul",
		"1960-08-24 Rosh Hashana LaBehemot",
		"1960-09-17 Leil Selichot",
		"1960-09-21 Erev Rosh Hashana",
	}

	actual := make([]string, 0, len(events))
	for _, ev := range events {
		line := fmt.Sprintf("%s %s", hd2iso(ev.Date), ev.Desc)
		actual = append(actual, line)
	}

	assert.Equal(t, expected, actual)
}

func TestModernILHolidays(t *testing.T) {
	events0Israel := hebcal.GetHolidaysForYear(5783, true)
	eventsIsrael := make([]event.HolidayEvent, 0, 12)
	for _, ev := range events0Israel {
		if (ev.Flags & event.MODERN_HOLIDAY) != 0 {
			eventsIsrael = append(eventsIsrael, ev)
		}
	}
	actualIsrael := make([]string, 0, len(eventsIsrael))
	for _, ev := range eventsIsrael {
		line := fmt.Sprintf("%s %s", hd2iso(ev.Date), ev.Desc)
		// fmt.Printf("\"%s\",\n", line)
		actualIsrael = append(actualIsrael, line)
	}
	expectedIsrael := []string{
		"2022-11-01 Yom HaAliyah School Observance",
		"2022-11-06 Yitzhak Rabin Memorial Day",
		"2022-11-23 Sigd",
		"2022-11-30 Ben-Gurion Day",
		"2023-02-21 Family Day",
		"2023-04-01 Yom HaAliyah",
		"2023-04-18 Yom HaShoah",
		"2023-04-25 Yom HaZikaron",
		"2023-04-26 Yom HaAtzma'ut",
		"2023-05-01 Herzl Day",
		"2023-05-19 Yom Yerushalayim",
		"2023-07-18 Jabotinsky Day",
	}
	assert.Equal(t, expectedIsrael, actualIsrael)

	events0Diaspora := hebcal.GetHolidaysForYear(5783, false)
	eventsDiaspora := make([]event.HolidayEvent, 0, 12)
	for _, ev := range events0Diaspora {
		if (ev.Flags & event.MODERN_HOLIDAY) != 0 {
			eventsDiaspora = append(eventsDiaspora, ev)
		}
	}
	actualDiaspora := make([]string, 0, len(eventsDiaspora))
	for _, ev := range eventsDiaspora {
		line := fmt.Sprintf("%s %s", hd2iso(ev.Date), ev.Desc)
		// fmt.Printf("\"%s\",\n", line)
		actualDiaspora = append(actualDiaspora, line)
	}
	expectedDiaspora := []string{
		"2022-11-23 Sigd",
		"2023-04-01 Yom HaAliyah",
		"2023-04-18 Yom HaShoah",
		"2023-04-25 Yom HaZikaron",
		"2023-04-26 Yom HaAtzma'ut",
		"2023-05-19 Yom Yerushalayim",
	}
	assert.Equal(t, expectedDiaspora, actualDiaspora)
}

func TestModernFriSatMovetoThu(t *testing.T) {
	events := hebcal.GetHolidaysForYear(5781, true)
	var rabinDay event.HolidayEvent
	for _, ev := range events {
		if ev.Desc == "Yitzhak Rabin Memorial Day" {
			rabinDay = ev
		}
	}
	assert.Equal(t, "2020-10-29", hd2iso(rabinDay.Date))
}

func TestBirkatHachamah(t *testing.T) {
	actual := make([]int, 0, 10)
	for year := 5650; year < 5920; year++ {
		events := hebcal.GetHolidaysForYear(year, false)
		for _, ev := range events {
			if ev.Desc == "Birkat Hachamah" {
				actual = append(actual, year)
			}
		}
	}
	expected := []int{5657, 5685, 5713, 5741, 5769, 5797, 5825, 5853, 5881, 5909}
	assert.Equal(t, expected, actual)

	events := hebcal.GetHolidaysForYear(5965, false)
	var hd hdate.HDate
	for _, ev := range events {
		if ev.Desc == "Birkat Hachamah" {
			hd = ev.Date
		}
	}
	assert.Equal(t, "19 Nisan 5965", hd.String())

	hd = hdate.HDate{}
	events = hebcal.GetHolidaysForYear(5993, false)
	for _, ev := range events {
		if ev.Desc == "Birkat Hachamah" {
			hd = ev.Date
		}
	}
	assert.Equal(t, "29 Adar II 5993", hd.String())
}

func TestPurimMeshulash(t *testing.T) {
	actual := make([]string, 0, 10)
	events := hebcal.GetHolidaysForYear(5781, true)
	for _, ev := range events {
		if ev.Date.Month() == hdate.Adar1 && ev.Date.Day() >= 13 && ev.Date.Day() <= 17 {
			line := fmt.Sprintf("%s / %s / %s", hd2iso(ev.Date), ev.Date.String(), ev.Desc)
			actual = append(actual, line)
		}
	}
	expected := []string{
		"2021-02-25 / 13 Adar 5781 / Erev Purim",
		"2021-02-25 / 13 Adar 5781 / Ta'anit Esther",
		"2021-02-26 / 14 Adar 5781 / Purim",
		"2021-02-27 / 15 Adar 5781 / Shushan Purim",
		"2021-02-28 / 16 Adar 5781 / Purim Meshulash",
	}
	assert.Equal(t, expected, actual)
}

func TestHolidayEmoji(t *testing.T) {
	var expectedEmoji = map[string]string{
		"Asara B'Tevet":                   "✡️",
		"Chanukah: 1 Candle":              "🕎",
		"Chanukah: 3 Candles":             "🕎",
		"Chanukah: 8 Candles":             "🕎",
		"Chanukah: 8th Day":               "🕎",
		"Lag BaOmer":                      "🔥",
		"Leil Selichot":                   "🕍",
		"Pesach Sheni":                    "✡️",
		"Erev Pesach":                     "🫓🍷",
		"Pesach I":                        "🫓🍷",
		"Pesach":                          "🫓",
		"Purim Katan":                     "🎭️",
		"Purim":                           "🎭️📜",
		"Rosh Chodesh Nisan":              "🌒",
		"Rosh Chodesh Iyyar":              "🌒",
		"Rosh Chodesh Sivan":              "🌒",
		"Rosh Chodesh Tamuz":              "🌒",
		"Rosh Chodesh Av":                 "🌒",
		"Rosh Chodesh Elul":               "🌒",
		"Rosh Chodesh Cheshvan":           "🌒",
		"Rosh Chodesh Kislev":             "🌒",
		"Rosh Chodesh Tevet":              "🌒",
		"Rosh Chodesh Sh'vat":             "🌒",
		"Rosh Chodesh Adar":               "🌒",
		"Rosh Chodesh Adar I":             "🌒",
		"Rosh Chodesh Adar II":            "🌒",
		"Rosh Hashana":                    "🍏🍯",
		"Rosh Hashana LaBehemot":          "🐑",
		"Shabbat Chazon":                  "🕍",
		"Shabbat HaChodesh":               "🕍",
		"Shabbat HaGadol":                 "🕍",
		"Shabbat Machar Chodesh":          "🕍",
		"Shabbat Nachamu":                 "🕍",
		"Shabbat Parah":                   "🕍",
		"Shabbat Rosh Chodesh":            "🕍",
		"Shabbat Shekalim":                "🕍",
		"Shabbat Shirah":                  "🕍",
		"Shabbat Shuva":                   "🕍",
		"Shabbat Zachor":                  "🕍",
		"Shavuot":                         "⛰️🌸",
		"Shmini Atzeret":                  "✡️",
		"Shushan Purim":                   "🎭️📜",
		"Purim Meshulash":                 "✡️",
		"Sigd":                            "✡️",
		"Simchat Torah":                   "✡️",
		"Sukkot":                          "🌿🍋",
		"Ta'anit Bechorot":                "✡️",
		"Ta'anit Esther":                  "✡️",
		"Tish'a B'Av":                     "✡️",
		"Tu B'Av":                         "❤️",
		"Tu BiShvat":                      "🌳",
		"Tzom Gedaliah":                   "✡️",
		"Tzom Tammuz":                     "✡️",
		"Yom HaAliyah":                    "🇮🇱",
		"Yom HaAtzma'ut":                  "🇮🇱",
		"Yom HaShoah":                     "✡️",
		"Yom HaZikaron":                   "🇮🇱",
		"Yom Kippur":                      "✡️",
		"Yom Yerushalayim":                "🇮🇱",
		"Yom Kippur Katan Sh'vat":         "",
		"Yom Kippur Katan Kislev":         "",
		"Shabbat Mevarchim Chodesh Iyyar": "",
	}
	events := hebcal.GetHolidaysForYear(5765, false)
	for _, ev := range events {
		actual := ev.GetEmoji()
		expected, ok := expectedEmoji[ev.Desc]
		if ok {
			assert.Equal(t, expected, actual, ev.Desc)
		} else {
			expected, ok := expectedEmoji[ev.Basename()]
			if ok {
				assert.Equal(t, expected, actual, ev.Desc)
			}
		}
	}
}

func TestHolidaysEarlyYears(t *testing.T) {
	events := hebcal.GetHolidaysForYear(3763, false)
	assert.Equal(t, 98, len(events))
	events = hebcal.GetHolidaysForYear(3762, false)
	assert.Equal(t, 104, len(events))
	events = hebcal.GetHolidaysForYear(3761, false)
	assert.Equal(t, 99, len(events))
	events = hebcal.GetHolidaysForYear(3760, false)
	assert.Equal(t, 105, len(events))
	events = hebcal.GetHolidaysForYear(2, false)
	assert.Equal(t, 99, len(events))
	events = hebcal.GetHolidaysForYear(1, false)
	assert.Equal(t, 99, len(events))
}
