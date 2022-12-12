package main

import (
	"fmt"
	"main/problems"
	"sort"
)

var problemsDict map[string]problems.Problem

func init() {
	problemsDict = map[string]problems.Problem{
		"p1": &problems.Problem1{},
		"p2": &problems.Problem2{},
		"p3": &problems.Problem3{},
		"p4": &problems.Problem4{},
		"p5": &problems.Problem5{},
	}
	problemsDict["p1"].Init(problems.Problem1Args{
		Files: []string{
			"D:\\Repos\\Operating_systems\\OperSysLab4\\assets\\1\\file_1.txt",
			"D:\\Repos\\Operating_systems\\OperSysLab4\\assets\\1\\file_2.txt",
			"D:\\Repos\\Operating_systems\\OperSysLab4\\assets\\1\\file_3.txt",
		},
	})
	problemsDict["p2"].Init(problems.Problem2Args{
		FiguresCount: 100,
		Width:        600,
		Height:       600,
	})
	problemsDict["p3"].Init(problems.Problem3Args{
		StationCapacity: 4,
		RefillMinSec:    8,
		RefillMaxSec:    12,
		NewCarMinSec:    1,
		NewCarMaxSec:    4,
	})
	problemsDict["p4"].Init("")
	problemsDict["p5"].Init("")
}

func main() {
	fmt.Println("Enter problem name to run:")
	avaliableKeys := make([]string, 0)
	for key := range problemsDict {
		avaliableKeys = append(avaliableKeys, key)
	}
	sort.Slice(avaliableKeys, func(i, j int) bool {
		return avaliableKeys[i] < avaliableKeys[j]
	})
	for _, key := range avaliableKeys {
		fmt.Printf("%s: %s\n", key, problemsDict[key].Description())
	}
	var choise string
	fmt.Scanf("%s", &choise)
	if pr, ok := problemsDict[choise]; ok {
		pr.Run()
	} else {
		fmt.Printf("Unknown problem name: %s\n", choise)
	}
}
