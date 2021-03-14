package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type entry struct {
	t   time.Time
	val int
}

func checkSum(file string) {
	p := regexp.MustCompile(`Received (\d+) at (.*), \w+ (\d+(.\d+)?)\.`)

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	entries := []entry{}

	scanner := bufio.NewScanner(f)
	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()
		match := p.FindStringSubmatch(line)
		num, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatalln(fmt.Sprintf("At line %d, Expect number to be an Integer, but got %s", i, match[1]))
		}
		curTime, err := time.Parse(time.RFC3339Nano, match[2])
		if err != nil {
			log.Fatalln(fmt.Sprintf("At line %d, Expect time to be a time.RFC3339Nano time, but got %s", i, match[2]))
		}
		sum, err := strconv.Atoi(match[3])
		if err != nil {
			log.Fatalln(fmt.Sprintf("At line %d, Expect sum to be an Integer, but got %s", i, match[3]))
		}

		// remove outdated entry
		for len(entries) > 0 && curTime.Sub(entries[0].t) > 5*time.Second {
			entries = entries[1:]
		}
		entries = append(entries, entry{curTime, num})
		expect := calcSum(entries)
		if expect != sum {
			log.Fatal(fmt.Sprintf("At line %d, Expect to have sum %d, but got %d\n", i, expect, sum))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("All sliding window sum are correct.")
}

func calcSum(entries []entry) int {
	sum := 0
	for _, e := range entries {
		sum += e.val
	}
	return sum
}

func checkMedian(file string) {
	p := regexp.MustCompile(`Received (\d+) at (.*), \w+ (\d+(.\d+)?)\.`)

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	entries := []entry{}

	scanner := bufio.NewScanner(f)
	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()
		match := p.FindStringSubmatch(line)
		num, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatalln(fmt.Sprintf("At line %d, Expect number to be an Integer, but got %s", i, match[1]))
		}
		curTime, err := time.Parse(time.RFC3339Nano, match[2])
		if err != nil {
			log.Fatalln(fmt.Sprintf("At line %d, Expect time to be a time.RFC3339Nano time, but got %s", i, match[2]))
		}
		median, err := strconv.ParseFloat(match[3], 64)
		if err != nil {
			log.Fatalln(fmt.Sprintf("At line %d, Expect median to be an Integer, but got %s", i, match[3]))
		}
		for len(entries) > 0 && curTime.Sub(entries[0].t) > 5*time.Second {
			entries = entries[1:]
		}
		entries = append(entries, entry{curTime, num})
		expect := calcMedian(entries)
		if expect != median {
			log.Fatal(fmt.Sprintf("At line %d, Expect to have median %g, but got %g.\n", i, expect, median))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("All sliding window median are correct.")
}

func calcMedian(entries []entry) float64 {
	nums := []int{}
	for _, e := range entries {
		nums = append(nums, e.val)
	}
	sort.Ints(nums)

	n := len(nums)
	if n%2 == 1 {
		return float64(nums[n/2])
	}
	return float64(nums[n/2-1]+nums[n/2]) / 2
}

func main() {
	args := os.Args
	if len(args) != 3 {
		log.Fatalln("Usage: ./validate file_name task(sum or median)")
	}

	file := os.Args[1]
	task := os.Args[2]

	if task == "sum" {
		checkSum(file)
	} else if task == "median" {
		checkMedian(file)
	} else {
		log.Fatalln("Task can only be 'sum' or 'median'.")
	}

}
