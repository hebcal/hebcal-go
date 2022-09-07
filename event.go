package hebcal

type HolidayFlags uint32

const (
	// Chag, yontiff, yom tov
	CHAG HolidayFlags = 1 << iota
	// Light candles 18 minutes before sundown
	LIGHT_CANDLES
	// End of holiday (end of Yom Tov)
	YOM_TOV_ENDS
	// Observed only in the Diaspora (chutz l'aretz)
	CHUL_ONLY
	// Observed only in Israel
	IL_ONLY
	// Light candles in the evening at Tzeit time (3 small stars)
	LIGHT_CANDLES_TZEIS
	// Candle-lighting for Chanukah
	CHANUKAH_CANDLES
	// Rosh Chodesh, beginning of a new Hebrew month
	ROSH_CHODESH
	// Minor fasts like Tzom Tammuz, Ta'anit Esther, ...
	MINOR_FAST
	// Shabbat Shekalim, Zachor, ...
	SPECIAL_SHABBAT
	// Weekly sedrot on Saturdays
	PARSHA_HASHAVUA
	// Daily page of Talmud
	DAF_YOMI
	// Days of the Omer
	OMER_COUNT
	// Yom HaShoah, Yom HaAtzma'ut, ...
	MODERN_HOLIDAY
	// Yom Kippur and Tish'a B'Av
	MAJOR_FAST
	// On the Saturday before Rosh Chodesh
	SHABBAT_MEVARCHIM
	// Molad
	MOLAD
	// Yahrzeit or Hebrew Anniversary
	USER_EVENT
	// Daily Hebrew date ("11th of Sivan, 5780")
	HEBREW_DATE
	// A holiday that's not major, modern, rosh chodesh, or a fast day
	MINOR_HOLIDAY
	// Evening before a major or minor holiday
	EREV
	// Chol haMoed, intermediate days of Pesach or Sukkot
	CHOL_HAMOED
	// Mishna Yomi
	MISHNA_YOMI
	// Yom Kippur Katan, minor day of atonement on the day preceding each Rosh Chodesh
	YOM_KIPPUR_KATAN
)

type HEvent interface {
	GetDate() HDate         // Holiday date of occurrence
	Render() string         // Description (e.g. "Pesach III (CH''M)")
	GetFlags() HolidayFlags // Event flag bitmask
	GetEmoji() string       // Holiday-specific emoji
}

type HolidayEvent struct {
	Date          HDate        // Holiday date of occurrence
	Desc          string       // Description (e.g. "Pesach III (CH''M)")
	Flags         HolidayFlags // Event flag bitmask
	Emoji         string       // Holiday-specific emoji
	CholHaMoedDay int          // used only for Pesach and Sukkot
	ChanukahDay   int          // used only for Chanukah
}

func (ev HolidayEvent) GetDate() HDate {
	return ev.Date
}

func (ev HolidayEvent) Render() string {
	return ev.Desc
}

func (ev HolidayEvent) GetFlags() HolidayFlags {
	return ev.Flags
}

func (ev HolidayEvent) GetEmoji() string {
	return ev.Emoji
}
