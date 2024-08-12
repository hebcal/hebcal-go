package omer_test

import (
	"fmt"
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/omer"
	"github.com/stretchr/testify/assert"
)

func TestSefira(t *testing.T) {
	assert := assert.New(t)
	omer := omer.NewOmerEvent(hdate.New(5770, hdate.Sivan, 2), 46)
	assert.Equal("Eternity within Majesty", omer.Sefira("en"))
	assert.Equal("נֶּֽצַח שֶׁבְּמַּלְכוּת", omer.Sefira("he"))
	assert.Equal("Netzach sheb'Malkhut", omer.Sefira("translit"))
}

func ExampleOmerEvent_Sefira() {
	ev := omer.NewOmerEvent(hdate.New(5770, hdate.Sivan, 2), 46)
	fmt.Println(ev.Sefira("en"))
	fmt.Println(ev.Sefira("he"))
	fmt.Println(ev.Sefira("translit"))
	// Output:
	// Eternity within Majesty
	// נֶּֽצַח שֶׁבְּמַּלְכוּת
	// Netzach sheb'Malkhut
}

func TestTodayIsEn(t *testing.T) {
	assert := assert.New(t)

	var ev = omer.NewOmerEvent(hdate.New(5770, hdate.Nisan, 16), 1)
	assert.Equal("Today is 1 day of the Omer", ev.TodayIs("en"))

	ev = omer.NewOmerEvent(hdate.New(5770, hdate.Nisan, 17), 2)
	assert.Equal("Today is 2 days of the Omer", ev.TodayIs("en"))

	ev = omer.NewOmerEvent(hdate.New(5770, hdate.Nisan, 22), 7)
	assert.Equal("Today is 7 days, which is 1 week of the Omer", ev.TodayIs("en"))

	ev = omer.NewOmerEvent(hdate.New(5770, hdate.Nisan, 23), 8)
	assert.Equal("Today is 8 days, which is 1 week and 1 day of the Omer", ev.TodayIs("en"))

	ev = omer.NewOmerEvent(hdate.New(5770, hdate.Nisan, 28), 13)
	assert.Equal("Today is 13 days, which is 1 week and 6 days of the Omer", ev.TodayIs("en"))

	ev = omer.NewOmerEvent(hdate.New(5770, hdate.Nisan, 29), 14)
	assert.Equal("Today is 14 days, which is 2 weeks of the Omer", ev.TodayIs("en"))

	ev = omer.NewOmerEvent(hdate.New(5770, hdate.Iyyar, 26), 41)
	assert.Equal("Today is 41 days, which is 5 weeks and 6 days of the Omer", ev.TodayIs("en"))

	ev = omer.NewOmerEvent(hdate.New(5770, hdate.Sivan, 2), 46)
	assert.Equal("Today is 46 days, which is 6 weeks and 4 days of the Omer", ev.TodayIs("en"))

}

func TestTodayIsHe(t *testing.T) {
	expected := []string{
		"",
		"הַיוֹם יוֹם אֶחָד לָעוֹמֶר",
		"הַיוֹם שְׁנֵי יָמִים לָעוֹמֶר",
		"הַיוֹם שְׁלוֹשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם אַרְבָּעָה יָמִים לָעוֹמֶר",
		"הַיוֹם חֲמִשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם שִׁשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם שִׁבְעָה יָמִים, שְׁהֵם שָׁבוּעַ אֶחָד לָעוֹמֶר",
		"הַיוֹם שְׁמוֹנָה יָמִים, שְׁהֵם שָׁבוּעַ אֶחָד וְיוֹם אֶחָד לָעוֹמֶר",
		"הַיוֹם תִּשְׁעָה יָמִים, שְׁהֵם שָׁבוּעַ אֶחָד וְשְׁנֵי יָמִים לָעוֹמֶר",
		"הַיוֹם עֲשָׂרָה יָמִים, שְׁהֵם שָׁבוּעַ אֶחָד וְשְׁלוֹשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם אֶחָד עָשָׂר יוֹם, שְׁהֵם שָׁבוּעַ אֶחָד וְאַרְבָּעָה יָמִים לָעוֹמֶר",
		"הַיוֹם שְׁנַיִם עָשָׂר יוֹם, שְׁהֵם שָׁבוּעַ אֶחָד וְחֲמִשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם שְׁלוֹשָׁה עָשָׂר יוֹם, שְׁהֵם שָׁבוּעַ אֶחָד וְשִׁשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם אַרְבָּעָה עָשָׂר יוֹם, שְׁהֵם שְׁנֵי שָׁבוּעוֹת לָעוֹמֶר",
		"הַיוֹם חֲמִשָׁה עָשָׂר יוֹם, שְׁהֵם שְׁנֵי שָׁבוּעוֹת וְיוֹם אֶחָד לָעוֹמֶר",
		"הַיוֹם שִׁשָׁה עָשָׂר יוֹם, שְׁהֵם שְׁנֵי שָׁבוּעוֹת וְשְׁנֵי יָמִים לָעוֹמֶר",
		"הַיוֹם שִׁבְעָה עָשָׂר יוֹם, שְׁהֵם שְׁנֵי שָׁבוּעוֹת וְשְׁלוֹשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם שְׁמוֹנָה עָשָׂר יוֹם, שְׁהֵם שְׁנֵי שָׁבוּעוֹת וְאַרְבָּעָה יָמִים לָעוֹמֶר",
		"הַיוֹם תִּשְׁעָה עָשָׂר יוֹם, שְׁהֵם שְׁנֵי שָׁבוּעוֹת וְחֲמִשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם עֶשְׂרִים יוֹם, שְׁהֵם שְׁנֵי שָׁבוּעוֹת וְשִׁשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם אֶחָד וְעֶשְׂרִים יוֹם, שְׁהֵם שְׁלוֹשָׁה שָׁבוּעוֹת לָעוֹמֶר",
		"הַיוֹם שְׁנַיִם וְעֶשְׂרִים יוֹם, שְׁהֵם שְׁלוֹשָׁה שָׁבוּעוֹת וְיוֹם אֶחָד לָעוֹמֶר",
		"הַיוֹם שְׁלוֹשָׁה וְעֶשְׂרִים יוֹם, שְׁהֵם שְׁלוֹשָׁה שָׁבוּעוֹת וְשְׁנֵי יָמִים לָעוֹמֶר",
		"הַיוֹם אַרְבָּעָה וְעֶשְׂרִים יוֹם, שְׁהֵם שְׁלוֹשָׁה שָׁבוּעוֹת וְשְׁלוֹשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם חֲמִשָׁה וְעֶשְׂרִים יוֹם, שְׁהֵם שְׁלוֹשָׁה שָׁבוּעוֹת וְאַרְבָּעָה יָמִים לָעוֹמֶר",
		"הַיוֹם שִׁשָׁה וְעֶשְׂרִים יוֹם, שְׁהֵם שְׁלוֹשָׁה שָׁבוּעוֹת וְחֲמִשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם שִׁבְעָה וְעֶשְׂרִים יוֹם, שְׁהֵם שְׁלוֹשָׁה שָׁבוּעוֹת וְשִׁשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם שְׁמוֹנָה וְעֶשְׂרִים יוֹם, שְׁהֵם אַרְבָּעָה שָׁבוּעוֹת לָעוֹמֶר",
		"הַיוֹם תִּשְׁעָה וְעֶשְׂרִים יוֹם, שְׁהֵם אַרְבָּעָה שָׁבוּעוֹת וְיוֹם אֶחָד לָעוֹמֶר",
		"הַיוֹם שְׁלוֹשִׁים יוֹם, שְׁהֵם אַרְבָּעָה שָׁבוּעוֹת וְשְׁנֵי יָמִים לָעוֹמֶר",
		"הַיוֹם אֶחָד וְשְׁלוֹשִׁים יוֹם, שְׁהֵם אַרְבָּעָה שָׁבוּעוֹת וְשְׁלוֹשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם שְׁנַיִם וְשְׁלוֹשִׁים יוֹם, שְׁהֵם אַרְבָּעָה שָׁבוּעוֹת וְאַרְבָּעָה יָמִים לָעוֹמֶר",
		"הַיוֹם שְׁלוֹשָׁה וְשְׁלוֹשִׁים יוֹם, שְׁהֵם אַרְבָּעָה שָׁבוּעוֹת וְחֲמִשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם אַרְבָּעָה וְשְׁלוֹשִׁים יוֹם, שְׁהֵם אַרְבָּעָה שָׁבוּעוֹת וְשִׁשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם חֲמִשָׁה וְשְׁלוֹשִׁים יוֹם, שְׁהֵם חֲמִשָׁה שָׁבוּעוֹת לָעוֹמֶר",
		"הַיוֹם שִׁשָׁה וְשְׁלוֹשִׁים יוֹם, שְׁהֵם חֲמִשָׁה שָׁבוּעוֹת וְיוֹם אֶחָד לָעוֹמֶר",
		"הַיוֹם שִׁבְעָה וְשְׁלוֹשִׁים יוֹם, שְׁהֵם חֲמִשָׁה שָׁבוּעוֹת וְשְׁנֵי יָמִים לָעוֹמֶר",
		"הַיוֹם שְׁמוֹנָה וְשְׁלוֹשִׁים יוֹם, שְׁהֵם חֲמִשָׁה שָׁבוּעוֹת וְשְׁלוֹשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם תִּשְׁעָה וְשְׁלוֹשִׁים יוֹם, שְׁהֵם חֲמִשָׁה שָׁבוּעוֹת וְאַרְבָּעָה יָמִים לָעוֹמֶר",
		"הַיוֹם אַרְבָּעִים יוֹם, שְׁהֵם חֲמִשָׁה שָׁבוּעוֹת וְחֲמִשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם אֶחָד וְאַרְבָּעִים יוֹם, שְׁהֵם חֲמִשָׁה שָׁבוּעוֹת וְשִׁשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם שְׁנַיִם וְאַרְבָּעִים יוֹם, שְׁהֵם שִׁשָׁה שָׁבוּעוֹת לָעוֹמֶר",
		"הַיוֹם שְׁלוֹשָׁה וְאַרְבָּעִים יוֹם, שְׁהֵם שִׁשָׁה שָׁבוּעוֹת וְיוֹם אֶחָד לָעוֹמֶר",
		"הַיוֹם אַרְבָּעָה וְאַרְבָּעִים יוֹם, שְׁהֵם שִׁשָׁה שָׁבוּעוֹת וְשְׁנֵי יָמִים לָעוֹמֶר",
		"הַיוֹם חֲמִשָׁה וְאַרְבָּעִים יוֹם, שְׁהֵם שִׁשָׁה שָׁבוּעוֹת וְשְׁלוֹשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם שִׁשָׁה וְאַרְבָּעִים יוֹם, שְׁהֵם שִׁשָׁה שָׁבוּעוֹת וְאַרְבָּעָה יָמִים לָעוֹמֶר",
		"הַיוֹם שִׁבְעָה וְאַרְבָּעִים יוֹם, שְׁהֵם שִׁשָׁה שָׁבוּעוֹת וְחֲמִשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם שְׁמוֹנָה וְאַרְבָּעִים יוֹם, שְׁהֵם שִׁשָׁה שָׁבוּעוֹת וְשִׁשָׁה יָמִים לָעוֹמֶר",
		"הַיוֹם תִּשְׁעָה וְאַרְבָּעִים יוֹם, שְׁהֵם שִׁבְעָה שָׁבוּעוֹת לָעוֹמֶר",
	}
	actual := make([]string, 50)
	start := hdate.New(5782, hdate.Nisan, 16)
	startAbs := start.Abs()
	for i := 1; i <= 49; i++ {
		abs := startAbs + int64(i) - 1
		ev := omer.NewOmerEvent(hdate.FromRD(abs), i)
		actual[i] = ev.TodayIs("he")
	}
	assert.Equal(t, expected, actual)
}

func TestEmoji(t *testing.T) {
	expected := []string{
		"",
		"①", "②", "③", "④", "⑤", "⑥", "⑦",
		"⑧", "⑨", "⑩", "⑪", "⑫", "⑬", "⑭",
		"⑮", "⑯", "⑰", "⑱", "⑲", "⑳", "㉑",
		"㉒", "㉓", "㉔", "㉕", "㉖", "㉗", "㉘",
		"㉙", "㉚", "㉛", "㉜", "㉝", "㉞", "㉟",
		"㊱", "㊲", "㊳", "㊴", "㊵", "㊶", "㊷",
		"㊸", "㊹", "㊺", "㊻", "㊼", "㊽", "㊾",
	}
	actual := make([]string, 50)
	start := hdate.New(5782, hdate.Nisan, 16)
	startAbs := start.Abs()
	for i := 1; i <= 49; i++ {
		abs := startAbs + int64(i) - 1
		ev := omer.NewOmerEvent(hdate.FromRD(abs), i)
		actual[i] = ev.GetEmoji()
	}
	assert.Equal(t, expected, actual)
}

func ExampleOmerEvent_GetEmoji() {
	omer := omer.NewOmerEvent(hdate.New(5770, hdate.Nisan, 28), 13)
	fmt.Println(omer.GetEmoji())
	// Output: ⑬
}

func ExampleOmerEvent_TodayIs() {
	omer := omer.NewOmerEvent(hdate.New(5770, hdate.Nisan, 28), 13)
	fmt.Println(omer.TodayIs("en"))
	fmt.Println(omer.TodayIs("he"))
	fmt.Println(omer.TodayIs("he-x-NoNikud"))
	// Output:
	// Today is 13 days, which is 1 week and 6 days of the Omer
	// הַיוֹם שְׁלוֹשָׁה עָשָׂר יוֹם, שְׁהֵם שָׁבוּעַ אֶחָד וְשִׁשָׁה יָמִים לָעוֹמֶר
	// היום שלושה עשר יום, שהם שבוע אחד וששה ימים לעומר
}
