package main

import (
	"fmt"

	"spiderden.net/go/libgallery"
)

const searchPrintString = `ID: %v
Tags: %v
NSFW: %v
Date: %v 
Description: %v
Uploader: %v
Source: %v
Score: %v
`

func implSearch(impl libgallery.Driver, query string, limit uint64) {
	var i uint64

	loopfunc := func() bool {
		return i >= limit
	}

	if limit == 0 {
		loopfunc = func() bool {
			return true
		}
	}

	var notfirst bool

	for i = 0; loopfunc(); i++ {
		result, err := impl.Search(query, i)
		if err != nil {
			panic(err)
		}

		if len(result) == 0 {
			return
		}

		for _, v := range result {
			if onlyids {
				fmt.Println(v.ID)
			} else {
				if onlyids {
					fmt.Println(v.ID)
				} else {
					str := searchPrintString
					if notfirst {
						fmt.Printf("\n")
					} else {
						notfirst = true
					}
					fmt.Printf(str,
						v.ID, v.Tags, v.NSFW, v.Date, v.Description,
						v.Uploader, sarrayToString(v.Source), v.Score)
				}
			}
		}

	}
}
