package fuzzy

import (
	"testing"
)

func TestRule_IF(t *testing.T) {

}

func TestRule_IS(t *testing.T) {

}

func TestRule_THEN(t *testing.T) {
	var (
		err error
		r   Rule
		// Temperature
		TEMPERATURE     Universe
		temperatureVar  base
		Cold, Warm, Hot Set
		// Sunshine
		SUNSHINE                 Universe
		sunshineVar              base
		Cloudy, PartSunny, Sunny Set
		// Tourists
		TOURISTS          Universe
		touristVar        base
		Low, Medium, High Set
	)
	{ // Temperature
		TEMPERATURE = NewUniverse("temperature", 0, 50)
		temperatureVar = TEMPERATURE.NewVar(19)

		Cold, err = NewFuzzySet(&TEMPERATURE, points{
			{0, 1}, {17, 1}, {20, 0},
		})
		if err != nil {
			panic(err)
		}
		Warm, err = NewFuzzySet(&TEMPERATURE, points{
			{17, 0}, {20, 1}, {26, 1}, {29, 0},
		})
		if err != nil {
			panic(err)
		}
		Hot, err = NewFuzzySet(&TEMPERATURE, points{
			{26, 0}, {29, 1}, {50, 1},
		})
		if err != nil {
			panic(err)
		}
	}
	{ // Sunshine
		SUNSHINE = NewUniverse("Sunshine", 0, 100)
		sunshineVar = SUNSHINE.NewVar(60)

		Cloudy, err = NewFuzzySet(&SUNSHINE, points{
			{0, 1}, {30, 1}, {50, 0},
		})
		if err != nil {
			panic(err)
		}
		PartSunny, err = NewFuzzySet(&SUNSHINE, points{
			{30, 0}, {50, 1}, {100, 0},
		})
		if err != nil {
			panic(err)
		}
		Sunny, err = NewFuzzySet(&SUNSHINE, points{
			{50, 0}, {100, 1},
		})
		if err != nil {
			panic(err)
		}
	}
	{ // tourists
		TOURISTS = NewUniverse("Tourists", 0, 100)

		Low, err = NewFuzzySet(&TOURISTS, points{
			{0, 1}, {50, 0},
		})
		if err != nil {
			panic(err)
		}
		Medium, err = NewFuzzySet(&TOURISTS, points{
			{0, 0}, {50, 1}, {100, 0},
		})
		if err != nil {
			panic(err)
		}
		High, err = NewFuzzySet(&TOURISTS, points{
			{50, 0}, {100, 1},
		})
		if err != nil {
			panic(err)
		}
	}

	r.IF(temperatureVar.IS(Hot)).THEN(touristVar.IS(High))
	//fmt.Println(r, temperatureVar, Cold, Warm, Hot, sunshineVar, Cloudy, PartSunny, Sunny, touristVar, Low, Medium, High)
}
