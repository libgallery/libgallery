package main

import (
	"fmt"
	"log"
	"os"

	"github.com/integrii/flaggy"
	"spiderden.net/go/libgallery"
)

var driversel string
var query string
var limit uint64
var onlyids bool

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flaggy.DefaultParser.DisableShowVersionWithVersion()
	// Search
	search := flaggy.NewSubcommand("search")
	search.AddPositionalValue(&driversel, "driver", 1, true, "sets the module to use with an action.")
	search.AddPositionalValue(&query, "query", 2, true, "tags which will be sent to the module to look for, snake-case and space-separated.")
	search.UInt64(&limit, "l", "limit", "maximum number of pages to search in, zero sets it to unlimited")
	search.Bool(&onlyids, "n", "ids", "only print IDs")

	search.Description = "Search for posts and print metadata."
	flaggy.AttachSubcommand(search, 1)

	// Archive
	archive := flaggy.NewSubcommand("archive")
	archive.AddPositionalValue(&driversel, "driver", 1, true, "sets the module to use with an action.")
	archive.AddPositionalValue(&query, "query", 2, true, "tags which will be sent to the module to look for, snake-case and space-separated.")
	archive.UInt64(&limit, "l", "limit", "maximum number of posts to search for, anything below 1 sets it to unlimited")
	flaggy.AttachSubcommand(archive, 1)

	// List
	list := flaggy.NewSubcommand("list")
	flaggy.AttachSubcommand(list, 1)
	flaggy.Parse()

	if !archive.Used && !search.Used && !list.Used {
		flaggy.ShowHelp("Please select an operation.")
		os.Exit(1)
	}

	var driver libgallery.Driver
	if !list.Used {
		var ok bool
		driver, ok = libgallery.Registry[driversel]
		if !ok {
			fmt.Fprint(os.Stderr, "Driver does not exist. Use the list subcommand to view avaliable subcommands.\n")
			os.Exit(1)
		}
	}

	switch {
	case archive.Used:
		implArchive(driversel, driver, query, limit)
	case search.Used:
		implSearch(driver, query, limit)
	case list.Used:
		for k := range libgallery.Registry {
			fmt.Println(k)
		}
	}

}
