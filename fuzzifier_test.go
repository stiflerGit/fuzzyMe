package main

import (
	"fmt"
	"fuzzyMe/fuzzy"
	"testing"
)

func TestDefuzzify(t *testing.T) {
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

	touristVar, err = TOURISTS.NewFuzzySingleton(0)
	ruleBase := fuzzy.NewRuleBase()
	//r.IF(temperatureVar.IS(Hot)).THEN(touristVar.IS(High))
	//rule := r.IF(temperatureVar).IS(Warm).THEN(touristVar).IS(Medium)
	rule1 := ruleBase.NewRule().IF(temperatureVar).IS(Hot).OR(sunshineVar).IS(Sunny).THEN(touristVar).IS(High)
	rule2 := ruleBase.NewRule().IF(temperatureVar).IS(Warm).AND(sunshineVar).IS(PartSunny).THEN(touristVar).IS(Medium)
	rule3 := ruleBase.NewRule().IF(temperatureVar).IS(Cold).AND(sunshineVar).IS(Cloudy).THEN(touristVar).IS(Low)
	//rule := r.IS(Hot).THEN(touristVar).IS(High)
	res1 := Defuzzify(rule1.EXEC(), COG)
	res2 := Defuzzify(rule2.EXEC(), COG)
	res3 := Defuzzify(rule3.EXEC(), COG)
	fmt.Print(res1, res2, res3)
	fmt.Println("###############################################################################")
	fmt.Println(temperatureVar, Cold, Warm, Hot, sunshineVar, Cloudy, PartSunny, Sunny, touristVar, Low, Medium, High)
	fmt.Println("###############################################################################")
}
