package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"spiderden.net/go/libgallery"
)

func implArchive(dname string, driver libgallery.Driver, query string, limit uint64) {
	dir := dname
	err := os.Mkdir(dir, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var i uint64
	loopfunc := func() bool {
		return i >= limit
	}

	if limit == 0 {
		loopfunc = func() bool {
			return true
		}
	}

	for i = 0; loopfunc(); {
		result, err := driver.Search(query, i)
		if err != nil {
			panic(err)
		}

		if len(result) == 0 {
			return
		}

		for _, v := range result {
			log.Printf("[%s] Downloading ID: %s\n", driver.Name(), v.ID)
			b, err := json.MarshalIndent(v, "", "	")
			if err != nil {
				log.Fatal(err)
			}

			os.WriteFile(fmt.Sprintf("%s.json", v.ID), b, 0700)

			files, err := driver.File(v.ID)
			if err != nil {
				log.Fatal(err)
			}

			defer files.Close()

			for i, j := range files {
				filename := fmt.Sprintf("%s_%v", v.ID, i)
				file, err := os.Create(filename)
				if err != nil {
					log.Fatal(err)
				}

				_, err = io.Copy(file, j)
				if err != nil {
					log.Fatal(err)
				}

				defer file.Close()

				mime, err := mimetype.DetectFile(filename)
				if err != nil {
					log.Fatal(err)
				}

				err = os.Rename(filename, filename+mime.Extension())
				if err != nil {
					log.Fatal(err)
				}

			}

			files.Close()
		}

		i++
	}
}
