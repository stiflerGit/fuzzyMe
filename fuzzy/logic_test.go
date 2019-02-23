package fuzzy

import (
	"fmt"
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
		temperatureVar  Set
		Cold, Warm, Hot Set
		// Sunshine
		SUNSHINE                 Universe
		sunshineVar              Set
		Cloudy, PartSunny, Sunny Set
		// Tourists
		TOURISTS          Universe
		touristVar        Set
		Low, Medium, High Set
	)
	{ // Temperature
		TEMPERATURE = NewUniverse("temperature", 0, 50)
		temperatureVar, _ = TEMPERATURE.NewFuzzySingleton(19)

		Cold, err = newFuzzySet(&TEMPERATURE, points{
			{0, 1}, {17, 1}, {20, 0},
		})
		if err != nil {
			panic(err)
		}
		Warm, err = newFuzzySet(&TEMPERATURE, points{
			{17, 0}, {20, 1}, {26, 1}, {29, 0},
		})
		if err != nil {
			panic(err)
		}
		Hot, err = newFuzzySet(&TEMPERATURE, points{
			{26, 0}, {29, 1}, {50, 1},
		})
		if err != nil {
			panic(err)
		}
	}
	{ // Sunshine
		SUNSHINE = NewUniverse("Sunshine", 0, 100)
		sunshineVar, _ = SUNSHINE.NewFuzzySingleton(60)

		Cloudy, err = newFuzzySet(&SUNSHINE, points{
			{0, 1}, {30, 1}, {50, 0},
		})
		if err != nil {
			panic(err)
		}
		PartSunny, err = newFuzzySet(&SUNSHINE, points{
			{30, 0}, {50, 1}, {100, 0},
		})
		if err != nil {
			panic(err)
		}
		Sunny, err = newFuzzySet(&SUNSHINE, points{
			{50, 0}, {100, 1},
		})
		if err != nil {
			panic(err)
		}
	}
	{ // tourists
		TOURISTS = NewUniverse("Tourists", 0, 100)

		Low, err = newFuzzySet(&TOURISTS, points{
			{0, 1}, {50, 0},
		})
		if err != nil {
			panic(err)
		}
		Medium, err = newFuzzySet(&TOURISTS, points{
			{0, 0}, {50, 1}, {100, 0},
		})
		if err != nil {
			panic(err)
		}
		High, err = newFuzzySet(&TOURISTS, points{
			{50, 0}, {100, 1},
		})
		if err != nil {
			panic(err)
		}
	}

	r.IF(temperatureVar.IS(Hot)).THEN(touristVar.IS(High))
	//r.IF(temperatureVar).IS(Hot).THEN(touristVar).IS(High)
	fmt.Println(r, temperatureVar, Cold, Warm, Hot, sunshineVar, Cloudy, PartSunny, Sunny, touristVar, Low, Medium, High)
}
