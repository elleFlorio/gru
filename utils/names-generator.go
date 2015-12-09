package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	source rand.Source
	gen    *rand.Rand

	left = [...]string{
		"admiring",
		"adoring",
		"agitated",
		"angry",
		"backstabbing",
		"berserk",
		"boring",
		"clever",
		"cocky",
		"compassionate",
		"condescending",
		"cranky",
		"desperate",
		"determined",
		"distracted",
		"dreamy",
		"drunk",
		"ecstatic",
		"elated",
		"elegant",
		"evil",
		"fervent",
		"focused",
		"furious",
		"gloomy",
		"goofy",
		"grave",
		"happy",
		"high",
		"hopeful",
		"hungry",
		"insane",
		"jolly",
		"jovial",
		"kickass",
		"lonely",
		"loving",
		"mad",
		"modest",
		"naughty",
		"nostalgic",
		"pensive",
		"prickly",
		"reverent",
		"romantic",
		"sad",
		"serene",
		"sharp",
		"sick",
		"silly",
		"sleepy",
		"stoic",
		"stupefied",
		"suspicious",
		"tender",
		"thirsty",
		"trusting",
	}
	right = [...]string{
		"sirius",
		"canopus",
		"rigil",
		"arcturus",
		"vega",
		"capella",
		"rigel",
		"procyon",
		"betelgeuse",
		"achernar",
		"hadar",
		"altair",
		"acrux",
		"aldebaran",
		"spica",
		"antares",
		"pollux",
		"fomalhaut",
		"deneb",
		"mimosa",
		"regulus",
		"adhara",
		"castor",
		"gacrux",
		"shaula",
		"bellatrix",
		"elnath",
		"miaplacidus",
		"alnilam",
		"alnair",
		"alnitak",
		"regor",
		"alioth",
		"kaus",
		"mirfak",
		"dubhe",
		"wezen",
		"alkaid",
		"sargas",
		"avior",
		"menkalinan",
		"atria",
		"koo",
		"alhena",
		"peacock",
		"polaris",
		"mirzam",
		"alphard",
		"algieba",
		"hamal",
	}
)

func init() {
	source = rand.NewSource(time.Now().UnixNano())
	gen = rand.New(source)
}

func GetRandomName(retry int) string {
	name := fmt.Sprintf("%s_%s", left[gen.Intn(len(left))], right[gen.Intn(len(right))])
	if retry > 0 {
		name = fmt.Sprintf("%s%d", name, gen.Intn(10))
	}
	return name
}
