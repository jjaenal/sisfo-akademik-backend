package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/todoexec"
)

func main() {
	root, _ := os.Getwd()
	path := filepath.Join(root, "docs", "TODO.md")
	tasks, lines, err := todoexec.Parse(path)
	if err != nil {
		fmt.Println("error:", err.Error())
		os.Exit(1)
	}
	var pending []todoexec.Task
	for _, t := range tasks {
		if !t.Completed {
			pending = append(pending, t)
		}
	}
	sorted := todoexec.SortByPriority(pending)
	var results []todoexec.ExecResult
	for _, t := range sorted {
		r := todoexec.Execute(root, t)
		results = append(results, r)
	}
	b, err := todoexec.UpdateStatuses(lines, results)
	if err != nil {
		fmt.Println("error:", err.Error())
		os.Exit(1)
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		fmt.Println("error:", err.Error())
		os.Exit(1)
	}
	var ok, fail, skip int
	for _, r := range results {
		if r.Executed && r.Succeeded {
			ok++
		} else if r.Executed && !r.Succeeded {
			fail++
		} else {
			skip++
		}
	}
	fmt.Println("executed:", ok)
	fmt.Println("failed:", fail)
	fmt.Println("skipped:", skip)
	for _, r := range results {
		if r.Executed {
			if r.Succeeded {
				fmt.Println("ok:", r.Task.Title)
			} else {
				fmt.Println("fail:", r.Task.Title, r.Error)
			}
		} else {
			fmt.Println("skip:", r.Task.Title)
		}
	}
}
