package main

import (
	"fmt"
	"testing"

	"github.com/stiflerGit/fuzzyMe/fuzzy"
)

var (
	err error
	// Temperature
	TEMPERATURE     fuzzy.Universe
	temperatureVar  fuzzy.Set
	Cold, Warm, Hot fuzzy.Set
	// Sunshine
	SUNSHINE                 fuzzy.Universe
	sunshineVar              fuzzy.Set
	Cloudy, PartSunny, Sunny fuzzy.Set
	// Tourists
	TOURISTS          fuzzy.Universe
	touristVar        fuzzy.Set
	Low, Medium, High fuzzy.Set
)

func init() {
	{ // Temperature
		TEMPERATURE = fuzzy.NewUniverse("temperature", 0, 50)
		temperatureVar, _ = TEMPERATURE.NewFuzzySingleton(19)

		Cold, err = TEMPERATURE.NewFuzzyMultipointSet(fuzzy.Points{
			{0, 1}, {17, 1}, {20, 0},
		})
		if err != nil {
			panic(err)
		}
		Warm, err = TEMPERATURE.NewFuzzyMultipointSet(fuzzy.Points{
			{17, 0}, {20, 1}, {26, 1}, {29, 0},
		})
		if err != nil {
			panic(err)
		}
		Hot, err = TEMPERATURE.NewFuzzyMultipointSet(fuzzy.Points{
			{26, 0}, {29, 1}, {50, 1},
		})
		if err != nil {
			panic(err)
		}
	}
	{ // Sunshine
		SUNSHINE = fuzzy.NewUniverse("Sunshine", 0, 100)
		sunshineVar, _ = SUNSHINE.NewFuzzySingleton(60)

		Cloudy, err = SUNSHINE.NewFuzzyMultipointSet(fuzzy.Points{
			{0, 1}, {30, 1}, {50, 0},
		})
		if err != nil {
			panic(err)
		}
		PartSunny, err = SUNSHINE.NewFuzzyMultipointSet(fuzzy.Points{
			{30, 0}, {50, 1}, {100, 0},
		})
		if err != nil {
			panic(err)
		}
		Sunny, err = SUNSHINE.NewFuzzyMultipointSet(fuzzy.Points{
			{50, 0}, {100, 1},
		})
		if err != nil {
			panic(err)
		}
	}
	{ // tourists
		TOURISTS = fuzzy.NewUniverse("Tourists", 0, 100)

		Low, err = TOURISTS.NewFuzzyMultipointSet(fuzzy.Points{
			{0, 1}, {50, 0},
		})
		if err != nil {
			panic(err)
		}
		Medium, err = TOURISTS.NewFuzzyMultipointSet(fuzzy.Points{
			{0, 0}, {50, 1}, {100, 0},
		})
		if err != nil {
			panic(err)
		}
		High, err = TOURISTS.NewFuzzyMultipointSet(fuzzy.Points{
			{50, 0}, {100, 1},
		})
		if err != nil {
			panic(err)
		}
	}
}

func TestDefuzzify(t *testing.T) {
	touristVar, err = TOURISTS.NewFuzzySingleton(40)
	ruleBase := fuzzy.NewRuleBase()

	ruleBase.NewRule().IF(temperatureVar).IS(Hot).OR(sunshineVar).IS(Sunny).THEN(touristVar).IS(High)
	ruleBase.NewRule().IF(temperatureVar).IS(Warm).AND(sunshineVar).IS(PartSunny).THEN(touristVar).IS(Medium)
	ruleBase.NewRule().IF(temperatureVar).IS(Cold).AND(sunshineVar).IS(Cloudy).THEN(touristVar).IS(Low)

	resSet, err := ruleBase.Exec()
	if err != nil {
		panic(err)
	}

	res := Defuzzify(resSet, COG)

	fmt.Print(res)
}

func ExampleDefuzzify() {

}
