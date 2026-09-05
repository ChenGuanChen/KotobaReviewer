package main

import (
	"fmt"
	"kotobaReviewer/cli"
	"kotobaReviewer/db"
	"os"
)

func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "import":
			if len(os.Args) < 3 {
				fmt.Println("usage: kotobaReviewer import <file>")
				return
			}
			cli.RunImport(os.Args[2])
			return
		case "add":
			cli.RunAdd()
			return
		case "review":
			cli.RunReview(os.Args[2:])
			return
		case "study":
			cli.RunStudy(os.Args[2:])
			return
		case "stats":
			cli.RunStats()
			return
		case "reclassify":
			cli.RunReclassify()
			return
		case "help", "-h", "--help":
			cli.PrintHelp()
			return
		default:
			fmt.Printf("unknown command %q\n\n", os.Args[1])
			cli.PrintHelp()
			return
		}
	}

	database, err := db.Init()
	if err != nil {
		fmt.Println("Error loading:", err)
		return
	}
	fmt.Printf("You have %d entries. Run 'kotobatrainer help' to see everything you can do.\n", len(database.Database))
}
