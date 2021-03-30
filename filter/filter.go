package filter

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var hosts []string

func init() {
	readSafeList()
}

// Read safe domains
func readSafeList() {
	f, err := os.OpenFile("safe.txt", os.O_RDONLY, 0644)
	if err != nil {
		log.Printf("Failed to read domain list %s", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hosts = append(hosts, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner failed to read %s", err)
	}
}

func Hosts() []string {
	return hosts
}

func IsSafe(host string) bool {
	for _, h := range hosts {
		if strings.EqualFold(h, host) {
			return true
		}
	}
	return false
}

func Add(host string) {
	if !IsSafe(host) {
		hosts = append(hosts, host)
		addToFile(host)
		log.Printf("Adding %s to whitelist", host)
	}
}

func Remove(host string) {
	var final []string
	for _, h := range hosts {
		if strings.EqualFold(h, host) {
			continue
		}
		final = append(final, h)
	}
	hosts = final
	updateFile()
}

func updateFile() {
	f, err := os.OpenFile("safe.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(strings.Join(hosts, "\n") + "\n"); err != nil {
		log.Println(err)
	}
}

func addToFile(from string) {
	f, err := os.OpenFile("safe.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(from + "\n"); err != nil {
		log.Println(err)
	}
}
