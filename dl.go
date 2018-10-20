package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Base summon chances out of 1,000,000
const (
	// Regular
	baseA = 5000  // 5* Adventurer
	baseD = 5000  // 5* Dragon
	baseW = 10000 // 5* Wyrmprint

	// Featured
	baseAF = 5000  // 5* Adventurer
	baseDF = 5000  // 5* Dragon
	baseWF = 10000 // 5* Wyrmprint
)

// Bonus summon chances
const (
	numSummonPerBonus = 10

	bonusA = 625
	bonusD = 625
	bonusW = 1250

	bonusAF = 625
	bonusDF = 625
	bonusWF = 1250
)

// Pity mechanics
const (
	numSummonPerPity = 100
)

// Sample size
const sampleSize = 1000000

var rateAF, rate5 int64
var clearChance = 0

// Summon counters
var bonusCounter = 0
var pityCounter = 0
var maxNumNotHit = 0
var numNotHit = 0

// RNG
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Result struct {
	rng int64
}

func (res *Result) IsAF() bool {
	return res.rng < rateAF
}

func (res *Result) Is5Star() bool {
	return res.rng < rate5
}

func main() {
	averageSummonPerAF(sampleSize * 10)
	averageAFPerSummon(sampleSize * 10)
	rateAFInFixedSummons(10)
	rateAFInFixedSummons(20)
	rateAFInFixedSummons(30)
	rateAFInFixedSummons(50)
	rateAFInFixedSummons(100)
	rateAFInFixedSummons(200)
	rateAFInFixedSummons(300)
	rateAFInFixedSummons(500)
	rateAFInFixedSummons(1000)
}

// averageSummonPerAF prints how many summons it takes
// to get featured adventurer on average.
func averageSummonPerAF(numAF int) {
	resetAll()
	count := 0

	// No need to resetAll() again, since resetAll()
	// is executed when a 5* is pulled.
	for i := 0; i < numAF; i++ {
		isAF := false
		if maxNumNotHit < numNotHit {
			maxNumNotHit = numNotHit
		}
		numNotHit = 0
		for !isAF {
			count += 1
			numNotHit += 1
			summon := single()
			isAF = summon.IsAF()
		}
	}

	ave := float64(count) / float64(numAF)
	aveStr := strconv.FormatFloat(ave, 'f', 2, 64)

	fmt.Println("===== Average Summons For Halloween Eli =====")
	fmt.Println(fmt.Sprintf("No. of Eli pulled: %d", numAF))
	fmt.Println(fmt.Sprintf("Average summons per Eli: %s", aveStr))
	fmt.Println(fmt.Sprintf("Max summons to get Eli: %d", maxNumNotHit))
}

// averageAFPerSummon prints how many featured adventurer
// you get per summon.
func averageAFPerSummon(numSummon int) {
	resetAll()
	count := 0

	for i := 0; i < numSummon; i++ {
		summon := single()
		if summon.IsAF() {
			count += 1
		}
	}

	ave := 100.0 * float64(count) / float64(numSummon)
	aveStr := strconv.FormatFloat(ave, 'f', 4, 64)

	fmt.Println("===== Average Halloween Eli per pull =====")
	fmt.Println(fmt.Sprintf("No. of pulls: %d", numSummon))
	fmt.Println(fmt.Sprintf("Average no. of Eli per 100 summons: %s", aveStr))
}

// rateAFInFixedSummons prints the rate of getting the featured adventurer
// in the given number of summons.
// Standard sample size is used.
func rateAFInFixedSummons(numSummon int) {
	resetAll()
	count := 0

	for i := 0; i < sampleSize; i++ {
		// Previous round of summons may not have yielded
		// a 5*, so we must reset and start afresh.
		resetAll()
		for j := 0; j < numSummon; j++ {
			summon := single()
			if summon.IsAF() {
				count += 1
				break
			}
		}
	}

	ave := 100.0 * float64(count) / float64(sampleSize)
	aveStr := strconv.FormatFloat(ave, 'f', 2, 64)

	fmt.Println("===== Rate of getting Halloween Eli in X pulls =====")
	fmt.Println(fmt.Sprintf("Sample size: %d", sampleSize))
	fmt.Println(fmt.Sprintf("Chance of getting Eli in %d pulls: %s%%", numSummon, aveStr))
}

// single performs a single summon, taking into account pity pulls
// and adjusting the summon rates where required.
func single_old() *Result {
	var res *Result

	addPityCounter()
	if isPitySummon() {
		res = &Result{rng: r.Int63n(rate5)}
	} else {
		res = &Result{rng: r.Int63n(1000000)}
	}

	if res.Is5Star() {
		resetAll()
	} else {
		addBonusCounter()
		if isBonusReached() {
			resetBonusCounter()
			increaseRates()
		}
	}

	return res
}

func single() *Result {
	var res *Result

	addPityCounter()
	if isPitySummon() {
		res = &Result{rng: r.Int63n(rate5)}
	} else {
		res = &Result{rng: r.Int63n(1000000)}
	}

	if res.Is5Star() {
		setResetFlag()
	}
	addBonusCounter()
	if isBonusReached() {
		resetBonusCounter()
		increaseRates()
		if isClearChance() {
			resetAll()
		}
	}

	return res
}

func isPitySummon() bool {
	return pityCounter >= numSummonPerPity
}

func isClearChance() bool {
	return clearChance >= 1
}

func isBonusReached() bool {
	return bonusCounter >= numSummonPerBonus
}

func addPityCounter() {
	pityCounter += 1
}

func resetPityCounter() {
	pityCounter = 0
}

func setResetFlag() {
	clearChance = 1;
}

func addBonusCounter() {
	bonusCounter += 1
}

func resetBonusCounter() {
	bonusCounter = 0
}

func resetClearChance() {
	clearChance = 0
}

func resetRates() {
	rateAF = baseAF
	rate5 = baseA + baseD + baseW + baseAF + baseDF + baseWF
}

func increaseRates() {
	rateAF += bonusAF
	rate5 += bonusA + bonusD + bonusW + bonusAF + bonusDF + bonusWF
}

func resetAll() {
	resetRates()
	resetPityCounter()
	resetBonusCounter()
	resetClearChance()
}
