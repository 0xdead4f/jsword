package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

func isValidIdentifier(s string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_$][\w$]*$`, s)
	return matched
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: jsword <url>")
		os.Exit(1)
	}
	var jsContent string
	if strings.Contains(os.Args[1], "http://") || strings.Contains(os.Args[1], "https://") {
	   url := os.Args[1]
	   resp, err := http.Get(url)
	   if err != nil {
		   fmt.Fprintf(os.Stderr, "Error fetching URL: %v\n", err)
		   os.Exit(1)
	   }
	   defer resp.Body.Close()
	   body, err := ioutil.ReadAll(resp.Body)
	   if err != nil {
		   fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		   os.Exit(1)
	   }
	   jsContent = string(body)
	} else {
	   body, err := ioutil.ReadFile(os.Args[1])
	   if err != nil {
		   fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		   os.Exit(1)
	   }
	   jsContent = string(body)
	}

	uniqueWords := make(map[string]bool)

	// Extract variables from declarations
	reVar := regexp.MustCompile(`\b(var|let|const)\s+([^;]+)`)
	varMatches := reVar.FindAllStringSubmatch(jsContent, -1)
	for _, match := range varMatches {
		if len(match) < 3 {
			continue
		}
		decl := match[2]
		vars := strings.Split(decl, ",")
		for _, v := range vars {
			v = strings.TrimSpace(v)
			parts := strings.Split(v, "=")
			varName := strings.TrimSpace(parts[0])
			if isValidIdentifier(varName) && len(varName) >= 2 {
				uniqueWords[varName] = true
			}
		}
	}

	// Extract function/class names and object properties
	patterns := []struct {
		pattern string
		checkId bool
	}{
		{`\bfunction\s+([a-zA-Z_$][\w$]*)`, true},
		{`\bclass\s+([a-zA-Z_$][\w$]*)`, true},
		{`(\w+)\s*\([^)]*\)\s*\{`, true},
		{`\.(\w+)\b`, true},
		{`:\s*([a-zA-Z_$][\w$]*)`, true},
		{`const\s+([a-zA-Z_$][\w$]*)\s*=\s*\(?`, true},
		{`id=["'](\w+)["']`, false},  // Get alphanumeric id value with quotes
		{`id=(\w+)`, false},          // Get alphanumeric id without quotes 
		{`name=["'](\w+)["']`, false}, // Get alphanumeric name value with quotes
		{`name=(\w+)`, false},         // Get alphanumeric name without quotes
		{`\{[^{}]*?([a-zA-Z_$][\w$]*)\s*:`, true},           // Object property names
		{`([a-zA-Z_$][\w$]*)\s*:\s*\{`, true},               // Nested object identifiers
		{`([a-zA-Z_$][\w$]*)\s*:\s*["'<]`, true},            // Property-value pairs
		{`<[^>]+(id|name)=["']([^"']+)["']`, false},         // HTML attributes in JSX/templates
		{`{{[\s\n]*(\w+)[\s\n]*}}`, false},      // {{variable}}
		{`{[\s\n]*(\w+)[\s\n]*}`, false},        // {variable}
		{`\${[\s\n]*(\w+)[\s\n]*}`, false},      // ${variable}
		{`@{[\s\n]*(\w+)[\s\n]*}`, false},       // @{variable}
		{`\[\[[\s\n]*(\w+)[\s\n]*\]\]`, false},  // [[variable]]
		{`<%[\s\n]*(\w+)[\s\n]*%>`, false},      // <%variable%>
		{`\(\([\s\n]*(\w+)[\s\n]*\)\)`, false},  // ((variable))
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p.pattern)
		matches := re.FindAllStringSubmatch(jsContent, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			word := match[1]
			if p.checkId && !isValidIdentifier(word) {
				continue
			}
			if len(word) >= 2 {
				uniqueWords[word] = true
			}
		}
	}

	// Process strings
	reString := regexp.MustCompile(`["']([^"']+)["']`)
	splitRegex := regexp.MustCompile(`[\W_]+`)
	stringMatches := reString.FindAllStringSubmatch(jsContent, -1)
	for _, match := range stringMatches {
		if len(match) < 2 {
			continue
		}
		for _, word := range splitRegex.Split(match[1], -1) {
			if len(word) >= 2 {
				uniqueWords[word] = true
			}
		}
	}

	// Prepare and print results
	words := make([]string, 0, len(uniqueWords))
	for word := range uniqueWords {
		words = append(words, word)
	}
	sort.Strings(words)

	for _, word := range words {
		fmt.Println(word)
	}
}