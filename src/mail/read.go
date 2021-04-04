package mail

import (
	"io/ioutil"
	"log"
)

func readAddress() string {
	dat, err := ioutil.ReadFile("address.txt")
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}

func readPass() string {
	dat, err := ioutil.ReadFile("addressp.txt")
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}
