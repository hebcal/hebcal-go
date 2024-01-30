package hebcal_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/MaxBGreenberg/hebcal-go/event"
	"github.com/MaxBGreenberg/hebcal-go/hebcal"
	"github.com/MaxBGreenberg/hebcal-go/yerushalmi"
	"github.com/MaxBGreenberg/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

func hd2iso(hd hdate.HDate) string {
	year, month, day := hd.Greg()
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return d.Format(time.RFC3339)[:10]
}

func TestHebrewCalendar(t *testing.T) {
	assert := assert.New(t)
	opts := hebcal.CalOptions{
		Year:         2022,
		Month:        time.April,
		IsHebrewYear: false,
	}
	events, err := hebcal.HebrewCalendar(&opts)

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
	opts := hebcal.CalOptions{
		NoHolidays:   true,
		Sedrot:       true,
		Year:         5783,
		IsHebrewYear: true,
	}
	events, err := hebcal.HebrewCalendar(&opts)
	assert.Equal(nil, err)
	assert.Equal(47, len(events))
	expected := []string{
		"2022-10-01 Parashat Vayeilech",
		"2022-10-08 Parashat Ha'azinu",
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
	opts := hebcal.CalOptions{
		Start:          hdate.New(5782, hdate.Elul, 25),
		End:            hdate.New(5783, hdate.Tishrei, 8),
		CandleLighting: true,
		Location:       loc,
		HavdalahMins:   50,
	}
	events, err := hebcal.HebrewCalendar(&opts)

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
	opts := hebcal.CalOptions{
		Start:          hdate.New(5783, hdate.Kislev, 24),
		End:            hdate.New(5783, hdate.Tevet, 2),
		CandleLighting: true,
		Location:       loc,
	}
	events, err := hebcal.HebrewCalendar(&opts)
	assert.Equal(nil, err)
	assert.Equal(14, len(events))
	expected := []string{
		"2022-12-18 Chanukah: 1 Candle: 5:04",
		"2022-12-19 Chanukah: 2 Candles: 5:05",
		"2022-12-20 Chanukah: 3 Candles: 5:05",
		"2022-12-21 Chanukah: 4 Candles: 5:06",
		"2022-12-22 Chanukah: 5 Candles: 5:06",
		"2022-12-23 Chanukah: 6 Candles: 4:00",
		"2022-12-23 Candle lighting: 4:00",
		"2022-12-24 Chag HaBanot",
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
	opts := hebcal.CalOptions{
		Year: 2020,
		Mask: event.ROSH_CHODESH | event.SPECIAL_SHABBAT,
	}
	events, err := hebcal.HebrewCalendar(&opts)
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
	opts := hebcal.CalOptions{
		Year:           2022,
		Sedrot:         true,
		CandleLighting: true,
		Location:       loc,
		HavdalahMins:   50,
	}
	events, _ := hebcal.HebrewCalendar(&opts)
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
	opts := hebcal.CalOptions{
		Year:           2022,
		Sedrot:         true,
		CandleLighting: true,
		Location:       loc,
		HavdalahMins:   50,
	}
	events, _ := hebcal.HebrewCalendar(&opts)
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
	opts := hebcal.CalOptions{
		Start:      hdate.New(5782, hdate.Kislev, 23),
		End:        hdate.New(5782, hdate.Kislev, 29),
		MishnaYomi: true,
		NoHolidays: true,
	}
	events, err := hebcal.HebrewCalendar(&opts)
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

func TestNoModern(t *testing.T) {
	opts := hebcal.CalOptions{
		Year:             2022,
		IL:               true,
		NoMinorFast:      true,
		NoModern:         true,
		NoRoshChodesh:    true,
		NoSpecialShabbat: true,
	}
	events, _ := hebcal.HebrewCalendar(&opts)
	actual := make([]string, 0, len(events))
	for _, ev := range events {
		desc := ev.Render("en")
		line := fmt.Sprintf("%s %s", hd2iso(ev.GetDate()), desc)
		// fmt.Printf("\"%s\",\n", line)
		actual = append(actual, line)
	}
	expected := []string{
		"2022-01-17 Tu BiShvat",
		"2022-02-15 Purim Katan",
		"2022-02-16 Shushan Purim Katan",
		"2022-03-16 Erev Purim",
		"2022-03-17 Purim",
		"2022-03-18 Shushan Purim",
		"2022-04-15 Erev Pesach",
		"2022-04-16 Pesach I",
		"2022-04-17 Pesach II (CH''M)",
		"2022-04-18 Pesach III (CH''M)",
		"2022-04-19 Pesach IV (CH''M)",
		"2022-04-20 Pesach V (CH''M)",
		"2022-04-21 Pesach VI (CH''M)",
		"2022-04-22 Pesach VII",
		"2022-05-15 Pesach Sheni",
		"2022-05-19 Lag BaOmer",
		"2022-06-04 Erev Shavuot",
		"2022-06-05 Shavuot",
		"2022-08-06 Erev Tish'a B'Av",
		"2022-08-07 Tish'a B'Av (observed)",
		"2022-08-12 Tu B'Av",
		"2022-08-28 Rosh Hashana LaBehemot",
		"2022-09-17 Leil Selichot",
		"2022-09-25 Erev Rosh Hashana",
		"2022-09-26 Rosh Hashana 5783",
		"2022-09-27 Rosh Hashana II",
		"2022-10-04 Erev Yom Kippur",
		"2022-10-05 Yom Kippur",
		"2022-10-09 Erev Sukkot",
		"2022-10-10 Sukkot I",
		"2022-10-11 Sukkot II (CH''M)",
		"2022-10-12 Sukkot III (CH''M)",
		"2022-10-13 Sukkot IV (CH''M)",
		"2022-10-14 Sukkot V (CH''M)",
		"2022-10-15 Sukkot VI (CH''M)",
		"2022-10-16 Sukkot VII (Hoshana Raba)",
		"2022-10-17 Shmini Atzeret",
		"2022-12-18 Chanukah: 1 Candle",
		"2022-12-19 Chanukah: 2 Candles",
		"2022-12-20 Chanukah: 3 Candles",
		"2022-12-21 Chanukah: 4 Candles",
		"2022-12-22 Chanukah: 5 Candles",
		"2022-12-23 Chanukah: 6 Candles",
		"2022-12-24 Chag HaBanot",
		"2022-12-24 Chanukah: 7 Candles",
		"2022-12-25 Chanukah: 8 Candles",
		"2022-12-26 Chanukah: 8th Day",
	}
	assert.Equal(t, expected, actual)
}

func TestDailyZemanim(t *testing.T) {
	hd := hdate.New(5782, hdate.Kislev, 23)
	loc := zmanim.LookupCity("Providence")
	opts := hebcal.CalOptions{
		Start:       hd,
		End:         hd,
		NoHolidays:  true,
		DailyZmanim: true,
		Location:    loc,
		Hour24:      true,
	}
	events, _ := hebcal.HebrewCalendar(&opts)
	actual := make([]string, 0, len(events))
	for _, ev := range events {
		desc := ev.Render("en")
		line := fmt.Sprintf("%s %s", hd2iso(ev.GetDate()), desc)
		actual = append(actual, line)
	}
	expected := []string{
		"2021-11-27 Alot haShachar: 05:21",
		"2021-11-27 Misheyakir: 05:47",
		"2021-11-27 Misheyakir Machmir: 05:54",
		"2021-11-27 Dawn: 06:18",
		"2021-11-27 Sunrise: 06:49",
		"2021-11-27 Kriat Shema, sof zeman (MGA): 08:35",
		"2021-11-27 Kriat Shema, sof zeman (GRA): 09:11",
		"2021-11-27 Tefilah, sof zeman (MGA): 09:34",
		"2021-11-27 Tefilah, sof zeman (GRA): 09:58",
		"2021-11-27 Chatzot hayom: 11:33",
		"2021-11-27 Mincha Gedolah: 11:57",
		"2021-11-27 Mincha Ketanah: 14:19",
		"2021-11-27 Plag HaMincha: 15:18",
		"2021-11-27 Sunset: 16:17",
		"2021-11-27 Dusk: 16:48",
		"2021-11-27 Tzeit HaKochavim: 17:02",
	}
	assert.Equal(t, expected, actual)
}

func TestHebrewCalendarYYomi(t *testing.T) {
	opts := hebcal.CalOptions{
		NoHolidays:     true,
		YerushalmiYomi: true,
		Start:          hdate.New(5783, hdate.Cheshvan, 18),
		End:            hdate.New(5783, hdate.Cheshvan, 23),
	}
	events, err := hebcal.HebrewCalendar(&opts)
	assert.Equal(t, nil, err)
	assert.Equal(t, 6, len(events))
	expected := []string{
		"2022-11-12 Yerushalmi Niddah 12",
		"2022-11-13 Yerushalmi Niddah 13",
		"2022-11-14 Yerushalmi Berakhot 1",
		"2022-11-15 Yerushalmi Berakhot 2",
		"2022-11-16 Yerushalmi Berakhot 3",
		"2022-11-17 Yerushalmi Berakhot 4",
	}
	actual := make([]string, 0, len(events))
	for _, ev := range events {
		desc := ev.Render("en")
		line := fmt.Sprintf("%s %s", hd2iso(ev.GetDate()), desc)
		actual = append(actual, line)
	}
	assert.Equal(t, expected, actual)
}

func TestHebrewCalendarSchottenstein(t *testing.T) {
	opts := hebcal.CalOptions{
		NoHolidays:        true,
		YerushalmiYomi:    true,
		YerushalmiEdition: yerushalmi.Schottenstein,
		Start:             hdate.FromGregorian(2022, time.November, 14),
		End:               hdate.FromGregorian(2028, time.August, 7),
	}
	events, err := hebcal.HebrewCalendar(&opts)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2094, len(events))
	actual := make([]string, 0, 40)
	for _, ev := range events {
		desc := ev.Render("en")
		if strings.HasSuffix(desc, " 1") {
			line := fmt.Sprintf("%s %s", hd2iso(ev.GetDate()), desc)
			actual = append(actual, line)
		}
	}
	expected := []string{
		"2022-11-14 Yerushalmi Berakhot 1",
		"2023-02-16 Yerushalmi Peah 1",
		"2023-04-30 Yerushalmi Demai 1",
		"2023-07-16 Yerushalmi Kilayim 1",
		"2023-10-08 Yerushalmi Sheviit 1",
		"2024-01-03 Yerushalmi Terumot 1",
		"2024-04-19 Yerushalmi Maasrot 1",
		"2024-06-04 Yerushalmi Maaser Sheni 1",
		"2024-08-02 Yerushalmi Challah 1",
		"2024-09-20 Yerushalmi Orlah 1",
		"2024-11-01 Yerushalmi Bikkurim 1",
		"2024-11-27 Yerushalmi Shabbat 1",
		"2025-03-20 Yerushalmi Eruvin 1",
		"2025-05-30 Yerushalmi Pesachim 1",
		"2025-08-24 Yerushalmi Shekalim 1",
		"2025-10-24 Yerushalmi Yoma 1",
		"2025-12-20 Yerushalmi Sukkah 1",
		"2026-01-22 Yerushalmi Beitzah 1",
		"2026-03-12 Yerushalmi Rosh Hashanah 1",
		"2026-04-08 Yerushalmi Taanit 1",
		"2026-05-09 Yerushalmi Megillah 1",
		"2026-06-19 Yerushalmi Chagigah 1",
		"2026-07-17 Yerushalmi Moed Katan 1",
		"2026-08-09 Yerushalmi Yevamot 1",
		"2026-11-05 Yerushalmi Ketubot 1",
		"2027-01-21 Yerushalmi Nedarim 1",
		"2027-03-04 Yerushalmi Nazir 1",
		"2027-04-26 Yerushalmi Sotah 1",
		"2027-06-17 Yerushalmi Gittin 1",
		"2027-08-09 Yerushalmi Kiddushin 1",
		"2027-10-01 Yerushalmi Bava Kamma 1",
		"2027-11-10 Yerushalmi Bava Metzia 1",
		"2027-12-15 Yerushalmi Bava Batra 1",
		"2028-01-23 Yerushalmi Sanhedrin 1",
		"2028-04-07 Yerushalmi Shevuot 1",
		"2028-05-26 Yerushalmi Avodah Zarah 1",
		"2028-06-29 Yerushalmi Makkot 1",
		"2028-07-10 Yerushalmi Horayot 1",
		"2028-07-28 Yerushalmi Niddah 1",
	}
	assert.Equal(t, expected, actual)
}

func TestYear2(t *testing.T) {
	opts := hebcal.CalOptions{
		Year:         2,
		IsHebrewYear: true,
		Sedrot:       true,
	}
	events, err := hebcal.HebrewCalendar(&opts)
	assert.Equal(t, nil, err)
	assert.Equal(t, 127, len(events))
}

func TestYear1(t *testing.T) {
	opts := hebcal.CalOptions{
		Year:         1,
		IsHebrewYear: true,
	}
	events, err := hebcal.HebrewCalendar(&opts)
	assert.Equal(t, nil, err)
	assert.Equal(t, 79, len(events))
}

func TestHebrewCalendarZmanimOnly(t *testing.T) {
	assert := assert.New(t)
	loc := zmanim.LookupCity("Amsterdam")
	opts := hebcal.CalOptions{
		Start:       hdate.New(5783, hdate.Sivan, 9),
		End:         hdate.New(5783, hdate.Sivan, 9),
		Location:    loc,
		DailyZmanim: true,
	}
	events, err := hebcal.HebrewCalendar(&opts)
	assert.Equal(nil, err)
	assert.Equal(15, len(events)) // not 16 (no Alot HaShachar)
}
