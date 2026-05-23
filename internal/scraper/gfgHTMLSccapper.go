package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GFGProfile struct {
	TotalSolved          int
	Streak               int
	EasySolved           int
	MediumSolved         int
	HardSolved           int
	ContestsParticipated int
	MaxRating            int
	CodingScore          int
	GlobalRank           int
	CountryRank          int
}

func FetchGFGProfile(username string) (*GFGProfile, error) {
	url := fmt.Sprintf("https://www.geeksforgeeks.org/profile/%s?tab=activity", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	htmlStr := string(body)

	// DEBUG: check if key fields exist in raw HTML
	fmt.Println("=== DEBUG ===")
	fmt.Println("HTML length:", len(htmlStr))
	fmt.Println("Contains 'total_problems_solved':", strings.Contains(htmlStr, "total_problems_solved"))
	fmt.Println("Contains 'score':", strings.Contains(htmlStr, "score"))
	fmt.Println("Contains '__next_f':", strings.Contains(htmlStr, "__next_f"))
	fmt.Println("Contains 'mearjuntripathi':", strings.Contains(htmlStr, username))

	// Print a snippet around total_problems_solved if it exists
	if idx := strings.Index(htmlStr, "total_problems_solved"); idx != -1 {
		start := idx - 100
		if start < 0 {
			start = 0
		}
		end := idx + 200
		if end > len(htmlStr) {
			end = len(htmlStr)
		}
		fmt.Println("Snippet around total_problems_solved:")
		fmt.Println(htmlStr[start:end])
	}
	fmt.Println("=== END DEBUG ===")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	profile := &GFGProfile{}

	// Collect all script contents and log them
	var allScripts []string
	scriptCount := 0
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		allScripts = append(allScripts, text)
		scriptCount++

		// Log scripts that might have our data
		if strings.Contains(text, "total_problems_solved") ||
			strings.Contains(text, "score") ||
			strings.Contains(text, "__next_f") {
			fmt.Printf("Script[%d] length=%d contains_total=%v contains_next_f=%v\n",
				i, len(text),
				strings.Contains(text, "total_problems_solved"),
				strings.Contains(text, "__next_f"),
			)
		}
	})
	fmt.Println("Total scripts found:", scriptCount)

	// Try parsing Next.js payload
	parsed := false
	for _, script := range allScripts {
		if strings.Contains(script, "total_problems_solved") {
			fmt.Println("Found total_problems_solved in script, attempting parse...")
			if parseNextJSPayload(script, profile) {
				fmt.Println("parseNextJSPayload succeeded")
				parsed = true
				break
			}
			fmt.Println("parseNextJSPayload failed")
		}
	}

	// Difficulty data
	for _, script := range allScripts {
		if strings.Contains(script, "School") && strings.Contains(script, "Medium") {
			parseDifficultyFromScript(script, profile)
			break
		}
	}

	if !parsed {
		fmt.Println("Falling back to HTML parsing...")
		parseFromHTML(doc, profile)
	}

	fmt.Printf("Final profile: %+v\n", profile)
	return profile, nil
}

func parseNextJSPayload(script string, profile *GFGProfile) bool {
	// Removed local GFGUserData struct - using package-level gfgUserData instead

	re := regexp.MustCompile(`"data":\{([^{}]*"total_problems_solved":[^{}]*)\}`)
	matches := re.FindStringSubmatch(script)

	if len(matches) >= 2 {
		fmt.Println("Regex matched in raw script")
		jsonStr := `{"data":{` + matches[1] + `}}`
		var wrapper struct {
			Data gfgUserData `json:"data"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &wrapper); err == nil {
			applyProfile(wrapper.Data, profile)
			return profile.TotalSolved > 0
		}
		fmt.Println("Direct unmarshal failed, trying unescaped...")
	}

	// Try unescaping the Next.js payload
	unescaped := unescapeNextJS(script)
	fmt.Println("Unescaped length:", len(unescaped))
	fmt.Println("Unescaped contains total_problems_solved:", strings.Contains(unescaped, "total_problems_solved"))

	if len(unescaped) > 0 {
		matches2 := re.FindStringSubmatch(unescaped)
		if len(matches2) >= 2 {
			fmt.Println("Regex matched in unescaped content")
			jsonStr := `{"data":{` + matches2[1] + `}}`
			var wrapper struct {
				Data gfgUserData `json:"data"`
			}
			if err := json.Unmarshal([]byte(jsonStr), &wrapper); err == nil {
				applyProfile(wrapper.Data, profile)
				return profile.TotalSolved > 0
			}
		}
	}

	return false
}

type gfgUserData struct {
	Score                      int    `json:"score"`
	MonthlyScore               int    `json:"monthly_score"`
	TotalProblemsSolved        int    `json:"total_problems_solved"`
	InstituteRank              string `json:"institute_rank"`
	PodSolvedLongestStreak     int    `json:"pod_solved_longest_streak"`
	PodSolvedCurrentStreak     int    `json:"pod_solved_current_streak"`
	PodCorrectSubmissionsCount int    `json:"pod_correct_submissions_count"`
}

func applyProfile(d gfgUserData, profile *GFGProfile) {
	profile.CodingScore = d.Score
	profile.MaxRating = d.Score
	profile.TotalSolved = d.TotalProblemsSolved
	profile.Streak = d.PodSolvedLongestStreak
	profile.ContestsParticipated = d.PodCorrectSubmissionsCount
	if d.InstituteRank != "" {
		rank, _ := strconv.Atoi(d.InstituteRank)
		profile.GlobalRank = rank
	}
}

func parseDifficultyFromScript(script string, profile *GFGProfile) {
	difficulties := map[string]int{}
	patterns := []string{"School", "Basic", "Easy", "Medium", "Hard"}
	for _, name := range patterns {
		re := regexp.MustCompile(name + `\s*\((\d+)\)`)
		m := re.FindStringSubmatch(script)
		if len(m) >= 2 {
			val, _ := strconv.Atoi(m[1])
			difficulties[name] = val
		}
	}
	profile.EasySolved = difficulties["School"] + difficulties["Basic"] + difficulties["Easy"]
	profile.MediumSolved = difficulties["Medium"]
	profile.HardSolved = difficulties["Hard"]
}

func unescapeNextJS(s string) string {
	re := regexp.MustCompile(`self\.__next_f\.push\(\[1,"((?:[^"\\]|\\.)*)"\]\)`)
	matches := re.FindAllStringSubmatch(s, -1)

	fmt.Println("unescapeNextJS: found", len(matches), "push matches")

	var parts []string
	for _, m := range matches {
		if len(m) >= 2 {
			jsonBytes := []byte(`"` + m[1] + `"`)
			var unescaped string
			if err := json.Unmarshal(jsonBytes, &unescaped); err == nil {
				parts = append(parts, unescaped)
			} else {
				fmt.Println("Unescape error:", err)
			}
		}
	}
	return strings.Join(parts, "")
}

func parseFromHTML(doc *goquery.Document, profile *GFGProfile) {
	doc.Find(".ScoreContainer_score-row___bfdI").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Find(".ScoreContainer_label__aVpLE").Text())
		value := strings.TrimSpace(s.Find(".ScoreContainer_value__7yy7h").Text())
		val, _ := strconv.Atoi(value)
		fmt.Printf("HTML score card: label=%q value=%q\n", label, value)
		switch label {
		case "Coding Score":
			profile.CodingScore = val
			profile.MaxRating = val
		case "Problems Solved":
			profile.TotalSolved = val
		case "Institute Rank":
			profile.GlobalRank = val
		}
	})

	streakText := strings.TrimSpace(doc.Find(".PotdContainer_streakText__oNgWh").Text())
	fmt.Println("Streak text:", streakText)
	if m := regexp.MustCompile(`^(\d+)`).FindStringSubmatch(streakText); len(m) >= 2 {
		profile.Streak, _ = strconv.Atoi(m[1])
	}

	doc.Find(".DoughnutChart_legendText__tQ2hK").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		fmt.Println("Legend text:", text)
		m := regexp.MustCompile(`(\w+)\s*\((\d+)\)`).FindStringSubmatch(text)
		if len(m) < 3 {
			return
		}
		count, _ := strconv.Atoi(m[2])
		switch m[1] {
		case "School", "Basic", "Easy":
			profile.EasySolved += count
		case "Medium":
			profile.MediumSolved = count
		case "Hard":
			profile.HardSolved = count
		}
	})
}