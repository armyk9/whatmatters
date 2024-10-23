package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// Vulnerability represents an individual exploit item with title, description, and link
type Vulnerability struct {
	Title       string
	Description template.HTML
	Link        string
}

func main() {
	// Define RSS feed URLs
	rssURLs := []string{
		"https://www.exploit-db.com/rss.xml",
		"https://rss.packetstormsecurity.com/files/tags/exploit/",
	}

	// Get RCE exploits from the past 12 months
	vulnerabilities := []Vulnerability{}
	for _, rssURL := range rssURLs {
		exploits := fetchExploits(rssURL)
		vulnerabilities = append(vulnerabilities, exploits...)
	}

	// Generate HTML report
	err := generateHTMLReport(vulnerabilities)
	if err != nil {
		log.Fatalf("Error generating HTML report: %v", err)
	}
	fmt.Println("Vulnerabilities report generated: vulnerabilities_report.html")
}

// fetchExploits scrapes the RSS feed and filters RCE exploits from the past 12 months
func fetchExploits(feedURL string) []Vulnerability {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		log.Fatalf("Error fetching RSS feed: %v", err)
	}

	var vulnerabilities []Vulnerability
	now := time.Now()
	oneYearAgo := now.AddDate(-1, 0, 0)

	for _, item := range feed.Items {
		if item.PublishedParsed != nil && item.PublishedParsed.After(oneYearAgo) {
			if strings.Contains(strings.ToLower(item.Title), "remote code execution") ||
				strings.Contains(strings.ToLower(item.Description), "rce") {
				// Combine the title and description with extra details for more context
				description := getDetailedDescription(item)
				vuln := Vulnerability{
					Title:       item.Title,
					Description: template.HTML(description), // Marking it as safe HTML
					Link:        item.Link,
				}
				vulnerabilities = append(vulnerabilities, vuln)
			}
		}
	}
	return vulnerabilities
}

// getDetailedDescription provides a more detailed description with clickable link
func getDetailedDescription(item *gofeed.Item) string {
	// Start with the item description
	description := item.Description

	// Add clickable "Read more" link
	if item.Link != "" {
		description += fmt.Sprintf(`<br><a href="%s" target="_blank">Read more</a>`, item.Link)
	}

	// If the description is too brief, add a fallback or more context
	if len(description) < 50 {
		description = "No detailed description available. For more information, visit the full exploit entry."
	}

	return description
}

// generateHTMLReport creates an HTML page with the vulnerabilities data
func generateHTMLReport(vulnerabilities []Vulnerability) error {
	// Define HTML template
	const tpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vulnerabilities Report</title>
    <style>
        body {
            background-color: #1e2a38;
            color: #ffffff;
            font-family: Arial, sans-serif;
        }
        .container {
            width: 80%;
            margin: 0 auto;
        }
        h1 {
            color: #ffcc00;
            text-align: center;
        }
        table {
            width: 100%;
            border-collapse: collapse;
        }
        table, th, td {
            border: 1px solid #cccccc;
        }
        th, td {
            padding: 12px;
            text-align: left;
        }
        th {
            background-color: #2f3e4f;
        }
        tr:nth-child(even) {
            background-color: #293542;
        }
        a {
            color: #ffcc00;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Remote Code Execution Vulnerabilities Report</h1>
        <table>
            <thead>
                <tr>
                    <th>Exploit Title</th>
                    <th>Description</th>
                </tr>
            </thead>
            <tbody>
                {{ range . }}
                <tr>
                    <td>{{ .Title }}</td>
                    <td>{{ .Description }}</td>
                </tr>
                {{ end }}
            </tbody>
        </table>
    </div>
</body>
</html>`

	// Parse the template
	t, err := template.New("report").Parse(tpl)
	if err != nil {
		return err
	}

	// Create HTML file
	file, err := os.Create("vulnerabilities_report.html")
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute the template with the vulnerabilities data
	err = t.Execute(file, vulnerabilities)
	if err != nil {
		return err
	}

	return nil
}
