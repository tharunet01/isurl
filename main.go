package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func main() {
	// File paths for File A and File B
	fileAPath := "fileA.txt" // Replace with the actual path to File A
	fileBPath := "fileB.txt" // Replace with the actual path to File B

	// Read File B into a map of domains for quick lookup
	domainMap, err := readDomains(fileBPath)
	if err != nil {
		fmt.Printf("Error reading File B: %v\n", err)
		return
	}

	// Filter URLs from File A based on the domains in File B
	filteredURLs, err := filterURLs(fileAPath, domainMap)
	if err != nil {
		fmt.Printf("Error reading File A: %v\n", err)
		return
	}

	// Output the filtered URLs
	fmt.Println("Filtered URLs:")
	for _, url := range filteredURLs {
		fmt.Println(url)
	}
}

// readDomains reads domains from a file into a map for quick lookup
func readDomains(filePath string) (map[string]struct{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	domains := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			domains[line] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}

// filterURLs filters URLs from File A based on the domains in the provided map
func filterURLs(filePath string, domainMap map[string]struct{}) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var filteredURLs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rawURL := strings.TrimSpace(scanner.Text())
		if rawURL == "" {
			continue
		}

		// Parse the URL to extract the domain
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			fmt.Printf("Skipping invalid URL: %s\n", rawURL)
			continue
		}

		host := parsedURL.Hostname()
		// Check if the domain or subdomain matches
		if isDomainMatch(host, domainMap) {
			filteredURLs = append(filteredURLs, rawURL)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return filteredURLs, nil
}

// isDomainMatch checks if a given host matches any domain in the map
func isDomainMatch(host string, domainMap map[string]struct{}) bool {
	for domain := range domainMap {
		if strings.HasSuffix(host, domain) {
			return true
		}
	}
	return false
}
