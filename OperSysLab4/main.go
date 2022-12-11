package main

import (
	"fmt"
	"main/problems"
)

var problemsDict map[string]problems.Problem

func init() {
	problemsDict = map[string]problems.Problem{
		"p1": &problems.Problem1{},
		"p2": &problems.Problem2{},
		"p3": &problems.Problem3{},
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
}

func main() {
	fmt.Println("Enter problem name to run:")
	for key, pr := range problemsDict {
		fmt.Printf("%s: %s\n", key, pr.Description())
	}
	var choise string
	fmt.Scanf("%s", &choise)
	if pr, ok := problemsDict[choise]; ok {
		pr.Run()
	} else {
		fmt.Printf("Unknown problem name: %s\n", choise)
	}
}
