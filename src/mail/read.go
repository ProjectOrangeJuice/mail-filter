package mail

import (
	"io/ioutil"
	"log"
)

func readAddress() string {
	dat, err := ioutil.ReadFile("/data/address.txt")
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}

func readPass() string {
	dat, err := ioutil.ReadFile("/data/addressp.txt")
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}
