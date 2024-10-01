package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

// RSS feeds for exploit databases
var rssFeeds = []string{
	"https://www.exploit-db.com/rss.xml",
	"https://rss.packetstormsecurity.com/files/tags/exploit/",
}

// Keywords to filter for RCE-related exploits
var rceKeywords = []string{"RCE", "Remote Code Execution"}

// Base URLs for handling relative paths
const exploitDBBaseURL = "https://www.exploit-db.com"
const packetStormBaseURL = "https://packetstormsecurity.com"

// Ensure the "exploits" directory exists
func createDirs() {
	// Create the exploits directory if it doesn't exist
	if _, err := os.Stat("exploits"); os.IsNotExist(err) {
		err := os.Mkdir("exploits", os.ModePerm)
		if err != nil {
			fmt.Println("Error creating 'exploits' directory:", err)
		}
	}
}

// Clean the title to use as a file name (remove special characters and unnecessary text)
func cleanTitleForFileName(title string) string {
	// Remove unwanted parts from the title, like "Remote Code Execution" or "(RCE)"
	title = strings.ReplaceAll(title, "Remote Code Execution", "")
	title = strings.ReplaceAll(title, "(RCE)", "")

	// Replace special characters with underscores or remove them
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	cleanTitle := reg.ReplaceAllString(title, "_")

	// Trim any leading or trailing underscores
	cleanTitle = strings.Trim(cleanTitle, "_")

	return cleanTitle
}

// function to check if the entry explicitly mentions RCE
func isRCEEntry(entry *gofeed.Item) bool {
	// Check if the title or description contains RCE-related keywords
	for _, keyword := range rceKeywords {
		if strings.Contains(strings.ToLower(entry.Title), strings.ToLower(keyword)) || strings.Contains(strings.ToLower(entry.Description), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// function to extract exploit download link from ExploitDB pages
func extractExploitDBLink(url string) []string {
	var links []string

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return links
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return links
	}

	// Search for the download link (based on the button in the UI)
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if exists && strings.Contains(href, "/download/") {
			// If it's a relative path, build the full URL
			fullURL := exploitDBBaseURL + href
			links = append(links, fullURL)
		}
	})

	if len(links) == 0 {
		fmt.Println("No exploit links found on page:", url)
	}
	return links
}

// function to extract exploit download link from Packet Storm pages
func extractPacketStormLink(url string) []string {
	var links []string

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return links
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return links
	}

	// Search for the download link
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if exists && strings.Contains(href, "/files/download/") {
			// If it's a relative path, build the full URL
			fullURL := packetStormBaseURL + href
			links = append(links, fullURL)
		}
	})

	if len(links) == 0 {
		fmt.Println("No exploit links found on page:", url)
	}
	return links
}

// function to determine file extension based on the file content
func determineFileExtension(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ".txt" // default extension
	}

	fileContent := string(content)

	// Detect language based on content
	if strings.Contains(fileContent, "import ") && strings.Contains(fileContent, "def ") {
		return ".py" // Python
	} else if strings.Contains(fileContent, "#include") && strings.Contains(fileContent, "int main") {
		return ".c" // C
	} else if strings.Contains(fileContent, "package main") || strings.Contains(fileContent, "func main") {
		return ".go" // GoLang
	} else if strings.Contains(fileContent, "#!/bin/bash") || strings.Contains(fileContent, "bash") {
		return ".sh" // Bash
	} else if strings.Contains(fileContent, "#!/usr/bin/perl") {
		return ".pl" // Perl
	}

	// Default to text if unable to detect language
	return ".txt"
}

// function to download file from URL into the "exploits" directory and rename it with correct extension
func downloadFile(url string, cleanedTitle string) {
	// Download the file temporarily
	tempFilePath := filepath.Join("exploits", cleanedTitle+".temp")

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(tempFilePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	// Ensure the file is closed before renaming
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error saving file:", err)
	}
	out.Close()

	// Determine the file extension based on content
	ext := determineFileExtension(tempFilePath)

	// Rename the file with the correct extension
	finalFilePath := filepath.Join("exploits", cleanedTitle+ext)
	err = os.Rename(tempFilePath, finalFilePath)
	if err != nil {
		fmt.Println("Error renaming file:", err)
	} else {
		fmt.Println("Downloaded file:", finalFilePath)
	}
}

// function to check if the item was published within the past 12 months
func isFromLastTwelveMonths(entry *gofeed.Item) bool {
	// Get the current time and the time 12 months ago
	now := time.Now()
	twelveMonthsAgo := now.AddDate(0, -12, 0)

	// Parse the item's publish date
	publishDate, err := time.Parse(time.RFC1123Z, entry.Published)
	if err != nil {
		// Try parsing in case another time format is used
		publishDate, err = time.Parse(time.RFC1123, entry.Published)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return false
		}
	}

	// Check if the item's publish date is within the past 12 months
	return publishDate.After(twelveMonthsAgo) && publishDate.Before(now)
}

func main() {
	// Ensure the exploits directory exists
	createDirs()

	fp := gofeed.NewParser()

	// Loop through RSS feeds
	for _, feedURL := range rssFeeds {
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			fmt.Println("Error parsing RSS feed:", err)
			continue
		}

		// Filter for RCE exploits within the past 12 months
		for _, item := range feed.Items {
			if isRCEEntry(item) && isFromLastTwelveMonths(item) {
				fmt.Println("Title:", item.Title)
				fmt.Println("Link:", item.Link)
				fmt.Println("Published:", item.Published)

				// Clean the title to use as the file name
				cleanedTitle := cleanTitleForFileName(item.Title)

				var exploitLinks []string
				// Detect which site we are scraping from and extract the correct download link
				if strings.Contains(item.Link, "exploit-db.com") {
					exploitLinks = extractExploitDBLink(item.Link)
				} else if strings.Contains(item.Link, "packetstormsecurity.com") {
					exploitLinks = extractPacketStormLink(item.Link)
				}

				for _, link := range exploitLinks {
					// Download and rename the file based on its content
					downloadFile(link, cleanedTitle)
				}
			}
		}
	}
}
