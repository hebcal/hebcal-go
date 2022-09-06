package hebcal

import (
	"errors"
	"strconv"
	"time"
)

type YearType int

const (
	Incomplete YearType = 1 + iota
	Regular
	Complete
)

var parshiot = []string{
	"Bereshit",
	"Noach",
	"Lech-Lecha",
	"Vayera",
	"Chayei Sara",
	"Toldot",
	"Vayetzei",
	"Vayishlach",
	"Vayeshev",
	"Miketz",
	"Vayigash",
	"Vayechi",
	"Shemot",
	"Vaera",
	"Bo",
	"Beshalach",
	"Yitro",
	"Mishpatim",
	"Terumah",
	"Tetzaveh",
	"Ki Tisa",
	"Vayakhel",
	"Pekudei",
	"Vayikra",
	"Tzav",
	"Shmini",
	"Tazria",
	"Metzora",
	"Achrei Mot",
	"Kedoshim",
	"Emor",
	"Behar",
	"Bechukotai",
	"Bamidbar",
	"Nasso",
	"Beha'alotcha",
	"Sh'lach",
	"Korach",
	"Chukat",
	"Balak",
	"Pinchas",
	"Matot",
	"Masei",
	"Devarim",
	"Vaetchanan",
	"Eikev",
	"Re'eh",
	"Shoftim",
	"Ki Teitzei",
	"Ki Tavo",
	"Nitzavim",
	"Vayeilech",
	"Ha'Azinu",
}

/* parsha doubler */
func D(n int) int {
	return -n
}

/* parsha undoubler */
func U(n int) int {
	return -n
}

func isValidDouble(n int) bool {
	switch n {
	case -21, -26, -28, -31, -38, -41, -50:
		return true
	default:
		return false
	}
}

/*
 * These indices were originally included in the emacs 19 distribution.
 * These arrays determine the correct indices into the parsha names
 * -1 means no parsha that week.
 */
var Sat_short = []int{
	-1, 52, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	17, 18, 19, 20, D(21), 23, 24, -1, 25, D(26), D(28), 30, D(31), 33, 34, 35, 36, 37, 38, 39, 40, D(41), 43, 44, 45, 46, 47,
	48, 49, 50}

var Sat_long = []int{
	-1, 52, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	17, 18, 19, 20, D(21), 23, 24, -1, 25, D(26), D(28), 30, D(31), 33, 34, 35, 36, 37, 38, 39, 40, D(41), 43, 44, 45, 46, 47,
	48, 49, D(50)}

var Mon_short = []int{
	51, 52, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17,
	18, 19, 20, D(21), 23, 24, -1, 25, D(26), D(28), 30, D(31), 33, 34, 35, 36, 37, 38, 39, 40, D(41), 43, 44, 45, 46, 47, 48,
	49, D(50)}

var Mon_long = []int{
	51, 52, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, D(21), 23, 24, -1, 25, D(26), D(28),
	30, D(31), 33, -1, 34, 35, 36, 37, D(38), 40, D(41), 43, 44, 45, 46, 47, 48, 49, D(50)}

var Thu_normal = []int{
	52, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17,
	18, 19, 20, D(21), 23, 24, -1, -1, 25, D(26), D(28), 30, D(31), 33, 34, 35, 36, 37, 38, 39, 40, D(41), 43, 44, 45, 46, 47,
	48, 49, 50}
var Thu_normal_Israel = []int{
	52, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	16, 17, 18, 19, 20, D(21), 23, 24, -1, 25, D(26), D(28), 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, D(41), 43, 44, 45,
	46, 47, 48, 49, 50}

var Thu_long = []int{
	52, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17,
	18, 19, 20, 21, 22, 23, 24, -1, 25, D(26), D(28), 30, D(31), 33, 34, 35, 36, 37, 38, 39, 40, D(41), 43, 44, 45, 46, 47,
	48, 49, 50}

var Sat_short_leap = []int{
	-1, 52, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, -1, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, D(41),
	43, 44, 45, 46, 47, 48, 49, D(50)}

var Sat_long_leap = []int{
	-1, 52, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, -1, 28, 29, 30, 31, 32, 33, -1, 34, 35, 36, 37, D(38), 40, D(41),
	43, 44, 45, 46, 47, 48, 49, D(50)}

var Mon_short_leap = []int{
	51, 52, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, -1, 28, 29, 30, 31, 32, 33, -1, 34, 35, 36, 37, D(38), 40, D(41), 43,
	44, 45, 46, 47, 48, 49, D(50)}

var Mon_short_leap_Israel = []int{
	51, 52, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, -1, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
	D(41), 43, 44, 45, 46, 47, 48, 49, D(50)}

var Mon_long_leap = []int{
	51, 52, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, -1, -1, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, D(41),
	43, 44, 45, 46, 47, 48, 49, 50}

var Mon_long_leap_Israel = []int{
	51, 52, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, -1, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
	41, 42, 43, 44, 45, 46, 47, 48, 49, 50}

var Thu_short_leap = []int{
	52, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, -1, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42,
	43, 44, 45, 46, 47, 48, 49, 50}

var Thu_long_leap = []int{
	52, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, -1, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42,
	43, 44, 45, 46, 47, 48, 49, D(50)}

func getSedraArray(leap bool, rhDay time.Weekday, yearType YearType, il bool) []int {
	if !leap {
		switch rhDay {
		case time.Saturday:
			if yearType == Incomplete {
				return Sat_short
			} else if yearType == Complete {
				return Sat_long
			}
		case time.Monday:
			if yearType == Incomplete {
				return Mon_short
			} else if yearType == Complete {
				if il {
					return Mon_short
				} else {
					return Mon_long
				}
			}
		case time.Tuesday:
			if yearType == Regular {
				if il {
					return Mon_short
				} else {
					return Mon_long
				}
			}
		case time.Thursday:
			if yearType == Regular {
				if il {
					return Thu_normal_Israel
				} else {
					return Thu_normal
				}
			} else if yearType == Complete {
				return Thu_long
			}
		}
	} else {
		/* leap year */
		switch rhDay {
		case time.Saturday:
			if yearType == Incomplete {
				return Sat_short_leap
			} else if yearType == Complete {
				if il {
					return Sat_short_leap
				} else {
					return Sat_long_leap
				}
			}
		case time.Monday:
			if yearType == Incomplete {
				if il {
					return Mon_short_leap_Israel
				} else {
					return Mon_short_leap
				}
			} else if yearType == Complete {
				if il {
					return Mon_long_leap_Israel
				} else {
					return Mon_long_leap
				}
			}
		case time.Tuesday:
			if yearType == Regular {
				if il {
					return Mon_long_leap_Israel
				} else {
					return Mon_long_leap
				}
			}
		case time.Thursday:
			if yearType == Incomplete {
				return Thu_short_leap
			} else if yearType == Complete {
				return Thu_long_leap
			}
		}
	}
	panic("improper sedra year type calculated")
}

type Sedra struct {
	Year          int
	IL            bool
	firstSaturday int
	theSedraArray []int
}

type Parsha struct {
	Name []string
	Num  []int
	Chag bool
}

func NewSedra(year int, il bool) Sedra {
	longC := LongCheshvan(year)
	shortK := ShortKislev(year)
	var yearType YearType
	if longC && !shortK {
		yearType = Complete
	} else if !longC && shortK {
		yearType = Incomplete
	} else {
		yearType = Regular
	}
	rh := HebrewToRD(year, Tishrei, 1)
	rhDay := time.Weekday(rh % 7)
	leap := IsHebLeapYear(year)
	firstSaturday := dayOnOrBefore(time.Saturday, rh+6)
	theSedraArray := getSedraArray(leap, rhDay, yearType, il)
	return Sedra{Year: year, IL: il, firstSaturday: firstSaturday, theSedraArray: theSedraArray}
}

func (sedra *Sedra) LookupByRD(rataDie int) Parsha {
	abs := dayOnOrBefore(time.Saturday, rataDie+6)
	if abs < sedra.firstSaturday {
		panic("lookup date " + strconv.Itoa(abs) + " is earlier than start of year " + strconv.Itoa(sedra.firstSaturday))
	}
	weekNum := ((abs - sedra.firstSaturday) / 7)
	if weekNum >= len(sedra.theSedraArray) {
		nextYear := NewSedra(sedra.Year+1, sedra.IL)
		return nextYear.LookupByRD(rataDie)
	}
	idx := sedra.theSedraArray[weekNum]
	if idx >= 0 {
		name := parshiot[idx]
		return Parsha{Name: []string{name}, Num: []int{idx + 1}, Chag: false}
	} else if idx == -1 {
		return Parsha{Chag: true}
	} else {
		// undouble
		p1 := U(idx)
		p2 := p1 + 1
		n1 := parshiot[p1]
		n2 := parshiot[p2]
		return Parsha{Name: []string{n1, n2}, Num: []int{p1 + 1, p2 + 1}, Chag: false}
	}
}

func (sedra *Sedra) Lookup(hd HDate) Parsha {
	return sedra.LookupByRD(hd.Abs())
}

func (sedra *Sedra) FindParshaNum(num int) (HDate, error) {
	parshaNum := num - 1
	if parshaNum > 53 || (parshaNum < 0 && !isValidDouble(parshaNum)) {
		return HDate{}, errors.New("invalid parsha number " + strconv.Itoa(num))
	}
	for idx, candidate := range sedra.theSedraArray {
		if candidate == parshaNum {
			return NewHDateFromRD(sedra.firstSaturday + (idx * 7)), nil
		}
	}
	panic("not found parsha num " + strconv.Itoa(num))
}
