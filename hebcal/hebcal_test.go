package hebcal

import (
	"fmt"
	"testing"
	"time"

	"github.com/hebcal/hebcal-go/hdate"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

func hd2iso(hd hdate.HDate) string {
	year, month, day := hd.Greg()
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return d.Format(time.RFC3339)[:10]
}

func TestHebrewCalendar(t *testing.T) {
	assert := assert.New(t)
	opts := CalOptions{
		Year:         2022,
		Month:        time.April,
		IsHebrewYear: false,
	}
	events, err := HebrewCalendar(&opts)

	assert.Equal(nil, err)
	assert.Equal(15, len(events))
	expected := []string{
		"2022-04-02 Rosh Chodesh Nisan",
		"2022-04-02 Shabbat HaChodesh",
		"2022-04-09 Shabbat HaGadol",
		"2022-04-11 Yom HaAliyah",
		"2022-04-15 Erev Pesach",
		"2022-04-15 Ta'anit Bechorot",
		"2022-04-16 Pesach I",
		"2022-04-17 Pesach II",
		"2022-04-18 Pesach III (CH''M)",
		"2022-04-19 Pesach IV (CH''M)",
		"2022-04-20 Pesach V (CH''M)",
		"2022-04-21 Pesach VI (CH''M)",
		"2022-04-22 Pesach VII",
		"2022-04-23 Pesach VIII",
		"2022-04-28 Yom HaShoah",
	}
	actual := make([]string, 0, len(events))
	for _, ev := range events {
		line := fmt.Sprintf("%s %s", hd2iso(ev.GetDate()), ev.Render("en"))
		actual = append(actual, line)
	}
	assert.Equal(expected, actual)
}

func TestHebrewCalendarSedrotOnly(t *testing.T) {
	assert := assert.New(t)
	opts := CalOptions{
		NoHolidays:   true,
		Sedrot:       true,
		Year:         5783,
		IsHebrewYear: true,
	}
	events, err := HebrewCalendar(&opts)
	assert.Equal(nil, err)
	assert.Equal(47, len(events))
	expected := []string{
		"2022-10-01 Parashat Vayeilech",
		"2022-10-08 Parashat Ha'Azinu",
		"2022-10-22 Parashat Bereshit",
		"2022-10-29 Parashat Noach",
		"2022-11-05 Parashat Lech-Lecha",
		"2022-11-12 Parashat Vayera",
		"2022-11-19 Parashat Chayei Sara",
		"2022-11-26 Parashat Toldot",
		"2022-12-03 Parashat Vayetzei",
		"2022-12-10 Parashat Vayishlach",
		"2022-12-17 Parashat Vayeshev",
		"2022-12-24 Parashat Miketz",
		"2022-12-31 Parashat Vayigash",
		"2023-01-07 Parashat Vayechi",
		"2023-01-14 Parashat Shemot",
		"2023-01-21 Parashat Vaera",
		"2023-01-28 Parashat Bo",
		"2023-02-04 Parashat Beshalach",
		"2023-02-11 Parashat Yitro",
		"2023-02-18 Parashat Mishpatim",
		"2023-02-25 Parashat Terumah",
		"2023-03-04 Parashat Tetzaveh",
		"2023-03-11 Parashat Ki Tisa",
		"2023-03-18 Parashat Vayakhel-Pekudei",
		"2023-03-25 Parashat Vayikra",
		"2023-04-01 Parashat Tzav",
		"2023-04-15 Parashat Shmini",
		"2023-04-22 Parashat Tazria-Metzora",
		"2023-04-29 Parashat Achrei Mot-Kedoshim",
		"2023-05-06 Parashat Emor",
		"2023-05-13 Parashat Behar-Bechukotai",
		"2023-05-20 Parashat Bamidbar",
		"2023-06-03 Parashat Nasso",
		"2023-06-10 Parashat Beha'alotcha",
		"2023-06-17 Parashat Sh'lach",
		"2023-06-24 Parashat Korach",
		"2023-07-01 Parashat Chukat-Balak",
		"2023-07-08 Parashat Pinchas",
		"2023-07-15 Parashat Matot-Masei",
		"2023-07-22 Parashat Devarim",
		"2023-07-29 Parashat Vaetchanan",
		"2023-08-05 Parashat Eikev",
		"2023-08-12 Parashat Re'eh",
		"2023-08-19 Parashat Shoftim",
		"2023-08-26 Parashat Ki Teitzei",
		"2023-09-02 Parashat Ki Tavo",
		"2023-09-09 Parashat Nitzavim-Vayeilech",
	}
	actual := make([]string, 0, len(events))
	for _, ev := range events {
		line := fmt.Sprintf("%s %s", hd2iso(ev.GetDate()), ev.Render("en"))
		actual = append(actual, line)
	}
	assert.Equal(expected, actual)
}

func TestHebrewCalendarCandles(t *testing.T) {
	assert := assert.New(t)
	loc := zmanim.LookupCity("Chicago")
	opts := CalOptions{
		Start:          hdate.New(5782, hdate.Elul, 25),
		End:            hdate.New(5783, hdate.Tishrei, 8),
		CandleLighting: true,
		Location:       loc,
		HavdalahMins:   50,
	}
	events, err := HebrewCalendar(&opts)

	assert.Equal(nil, err)
	assert.Equal(14, len(events))
	expected := []string{
		"2022-09-23 Candle lighting: 6:28",
		"2022-09-24 Havdalah (50 min): 7:35",
		"2022-09-25 Erev Rosh Hashana",
		"2022-09-25 Candle lighting: 6:25",
		"2022-09-26 Rosh Hashana 5783",
		"2022-09-26 Candle lighting: 7:32",
		"2022-09-27 Rosh Hashana II",
		"2022-09-27 Havdalah (50 min): 7:30",
		"2022-09-28 Fast begins: 5:21",
		"2022-09-28 Tzom Gedaliah",
		"2022-09-28 Fast ends: 7:11",
		"2022-09-30 Candle lighting: 6:16",
		"2022-10-01 Shabbat Shuva",
		"2022-10-01 Havdalah (50 min): 7:23",
	}
	actual := make([]string, 0, len(events))
	for _, ev := range events {
		desc := ev.Render("en")
		line := fmt.Sprintf("%s %s", hd2iso(ev.GetDate()), desc)
		actual = append(actual, line)
	}
	assert.Equal(expected, actual)
}

func TestHebrewCalendarChanukahCandles(t *testing.T) {
	assert := assert.New(t)
	loc := zmanim.LookupCity("Jerusalem")
	opts := CalOptions{
		Start:          hdate.New(5783, hdate.Kislev, 24),
		End:            hdate.New(5783, hdate.Tevet, 2),
		CandleLighting: true,
		Location:       loc,
	}
	events, err := HebrewCalendar(&opts)
	assert.Equal(nil, err)
	assert.Equal(13, len(events))
	expected := []string{
		"2022-12-18 Chanukah: 1 Candle: 5:04",
		"2022-12-19 Chanukah: 2 Candles: 5:05",
		"2022-12-20 Chanukah: 3 Candles: 5:05",
		"2022-12-21 Chanukah: 4 Candles: 5:06",
		"2022-12-22 Chanukah: 5 Candles: 5:06",
		"2022-12-23 Chanukah: 6 Candles: 4:00",
		"2022-12-23 Candle lighting: 4:00",
		"2022-12-24 Chanukah: 7 Candles: 5:20",
		"2022-12-24 Rosh Chodesh Tevet",
		"2022-12-24 Havdalah: 5:20",
		"2022-12-25 Chanukah: 8 Candles: 5:08",
		"2022-12-25 Rosh Chodesh Tevet",
		"2022-12-26 Chanukah: 8th Day",
	}
	actual := make([]string, 0, len(events))
	for _, ev := range events {
		desc := ev.Render("en")
		line := fmt.Sprintf("%s %s", hd2iso(ev.GetDate()), desc)
		// fmt.Printf("\"%s\",\n", line)
		actual = append(actual, line)
	}
	assert.Equal(expected, actual)
}

func TestHebrewCalendarMask(t *testing.T) {
	assert := assert.New(t)
	opts := CalOptions{
		Year: 2020,
		Mask: ROSH_CHODESH | SPECIAL_SHABBAT,
	}
	events, err := HebrewCalendar(&opts)
	assert.Equal(nil, err)
	assert.Equal(25, len(events))
	expected := []string{
		"2020-01-27 Rosh Chodesh Sh'vat",
		"2020-02-08 Shabbat Shirah",
		"2020-02-22 Shabbat Shekalim",
		"2020-02-25 Rosh Chodesh Adar",
		"2020-02-26 Rosh Chodesh Adar",
		"2020-03-07 Shabbat Zachor",
		"2020-03-14 Shabbat Parah",
		"2020-03-21 Shabbat HaChodesh",
		"2020-03-26 Rosh Chodesh Nisan",
		"2020-04-04 Shabbat HaGadol",
		"2020-04-24 Rosh Chodesh Iyyar",
		"2020-04-25 Rosh Chodesh Iyyar",
		"2020-05-24 Rosh Chodesh Sivan",
		"2020-06-22 Rosh Chodesh Tamuz",
		"2020-06-23 Rosh Chodesh Tamuz",
		"2020-07-22 Rosh Chodesh Av",
		"2020-07-25 Shabbat Chazon",
		"2020-08-01 Shabbat Nachamu",
		"2020-08-20 Rosh Chodesh Elul",
		"2020-08-21 Rosh Chodesh Elul",
		"2020-09-26 Shabbat Shuva",
		"2020-10-18 Rosh Chodesh Cheshvan",
		"2020-10-19 Rosh Chodesh Cheshvan",
		"2020-11-17 Rosh Chodesh Kislev",
		"2020-12-16 Rosh Chodesh Tevet",
	}
	actual := make([]string, 0, len(events))
	for _, ev := range events {
		desc := ev.Render("en")
		line := fmt.Sprintf("%s %s", hd2iso(ev.GetDate()), desc)
		actual = append(actual, line)
	}
	assert.Equal(expected, actual)
}

func ExampleHebrewCalendar() {
	loc := zmanim.LookupCity("Providence")
	opts := CalOptions{
		Year:           2022,
		Sedrot:         true,
		CandleLighting: true,
		Location:       loc,
		HavdalahMins:   50,
	}
	events, _ := HebrewCalendar(&opts)
	for i := 0; i < 6; i++ {
		ev := events[i]
		dateStr := ev.GetDate().Gregorian().Format("Mon 02-Jan-2006")
		title := ev.Render("en")
		fmt.Println(dateStr, title)
	}
	// Output:
	// Sat 01-Jan-2022 Parashat Vaera
	// Sat 01-Jan-2022 Havdalah (50 min): 5:15
	// Mon 03-Jan-2022 Rosh Chodesh Sh'vat
	// Fri 07-Jan-2022 Candle lighting: 4:12
	// Sat 08-Jan-2022 Parashat Bo
	// Sat 08-Jan-2022 Havdalah (50 min): 5:22
}

func TestHebrewCalendarLocale(t *testing.T) {
	loc := zmanim.LookupCity("Providence")
	opts := CalOptions{
		Year:           2022,
		Sedrot:         true,
		CandleLighting: true,
		Location:       loc,
		HavdalahMins:   50,
	}
	events, _ := HebrewCalendar(&opts)
	actual := make([]string, 6)
	for i := 0; i < 6; i++ {
		ev := events[i]
		dateStr := ev.GetDate().Gregorian().Format("Mon 02-Jan-2006")
		actual[i] = fmt.Sprintf("%s %s", dateStr, ev.Render("es"))
	}
	expected := []string{
		"Sat 01-Jan-2022 Parashá Vaera",
		"Sat 01-Jan-2022 Havdalah (50 min): 5:15",
		"Mon 03-Jan-2022 Rosh Jodesh Sh'vat",
		"Fri 07-Jan-2022 Iluminación de velas: 4:12",
		"Sat 08-Jan-2022 Parashá Bo",
		"Sat 08-Jan-2022 Havdalah (50 min): 5:22",
	}
	assert.Equal(t, expected, actual)
}

func TestHebrewCalendarMishnaYomiOnly(t *testing.T) {
	assert := assert.New(t)
	opts := CalOptions{
		Start:      hdate.New(5782, hdate.Kislev, 23),
		End:        hdate.New(5782, hdate.Kislev, 29),
		MishnaYomi: true,
		NoHolidays: true,
	}
	events, err := HebrewCalendar(&opts)
	assert.Equal(nil, err)
	assert.Equal(7, len(events))
	expected := []string{
		"2021-11-27 Tevul Yom 4:2-3",
		"2021-11-28 Tevul Yom 4:4-5",
		"2021-11-29 Tevul Yom 4:6-7",
		"2021-11-30 Yadayim 1:1-2",
		"2021-12-01 Yadayim 1:3-4",
		"2021-12-02 Yadayim 1:5-2:1",
		"2021-12-03 Yadayim 2:2-3",
	}
	actual := make([]string, 0, len(events))
	for _, ev := range events {
		desc := ev.Render("en")
		line := fmt.Sprintf("%s %s", hd2iso(ev.GetDate()), desc)
		// fmt.Printf("\"%s\",\n", line)
		actual = append(actual, line)
	}
	assert.Equal(expected, actual)
}