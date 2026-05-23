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

type gfgUserData struct {
	Score                      int    `json:"score"`
	MonthlyScore               int    `json:"monthly_score"`
	TotalProblemsSolved        int    `json:"total_problems_solved"`
	InstituteRank              string `json:"institute_rank"`
	PodSolvedLongestStreak     int    `json:"pod_solved_longest_streak"`
	PodSolvedCurrentStreak     int    `json:"pod_solved_current_streak"`
	PodCorrectSubmissionsCount int    `json:"pod_correct_submissions_count"`
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

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	profile := &GFGProfile{}

	// Collect all script contents
	var allScripts []string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		allScripts = append(allScripts, s.Text())
	})

	// Parse main stats from Next.js payload
	for _, script := range allScripts {
		if strings.Contains(script, "total_problems_solved") {
			if parseNextJSPayload(script, profile) {
				break
			}
		}
	}

	// Parse difficulty from the HTML directly
	// The snippet shows: School (0)  Basic (16)  Easy (205)  Medium (275)  Hard (33)
	// These are inside the raw HTML not script tags
	parseDifficultyFromHTML(htmlStr, profile)

	return profile, nil
}

func parseNextJSPayload(script string, profile *GFGProfile) bool {
	re := regexp.MustCompile(`"data":\{([^{}]*"total_problems_solved":[^{}]*)\}`)

	// Try on raw script first
	matches := re.FindStringSubmatch(script)
	if len(matches) >= 2 {
		jsonStr := `{"data":{` + matches[1] + `}}`
		var wrapper struct {
			Data gfgUserData `json:"data"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &wrapper); err == nil && wrapper.Data.TotalProblemsSolved > 0 {
			applyProfile(wrapper.Data, profile)
			return true
		}
	}

	// Try after unescaping Next.js payload
	unescaped := unescapeNextJS(script)
	if unescaped == "" {
		return false
	}

	matches2 := re.FindStringSubmatch(unescaped)
	if len(matches2) >= 2 {
		jsonStr := `{"data":{` + matches2[1] + `}}`
		var wrapper struct {
			Data gfgUserData `json:"data"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &wrapper); err == nil && wrapper.Data.TotalProblemsSolved > 0 {
			applyProfile(wrapper.Data, profile)
			return true
		}
	}

	return false
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

// parseDifficultyFromHTML searches raw HTML for patterns like "School (0)", "Easy (205)"
// These appear in the doughnut chart legend in the raw HTML body
func parseDifficultyFromHTML(html string, profile *GFGProfile) {
	patterns := map[string]*int{
		`School\s*\((\d+)\)`: nil,
		`Basic\s*\((\d+)\)`:  nil,
		`Easy\s*\((\d+)\)`:   nil,
		`Medium\s*\((\d+)\)`: &profile.MediumSolved,
		`Hard\s*\((\d+)\)`:   &profile.HardSolved,
	}

	school, basic, easy := 0, 0, 0

	for pattern, target := range patterns {
		re := regexp.MustCompile(pattern)
		m := re.FindStringSubmatch(html)
		if len(m) < 2 {
			continue
		}
		val, _ := strconv.Atoi(m[1])
		switch {
		case strings.Contains(pattern, "School"):
			school = val
		case strings.Contains(pattern, "Basic"):
			basic = val
		case strings.Contains(pattern, "Easy"):
			easy = val
		default:
			if target != nil {
				*target = val
			}
		}
	}

	profile.EasySolved = school + basic + easy
}

func unescapeNextJS(s string) string {
	re := regexp.MustCompile(`self\.__next_f\.push\(\[1,"((?:[^"\\]|\\.)*)"\]\)`)
	matches := re.FindAllStringSubmatch(s, -1)

	var parts []string
	for _, m := range matches {
		if len(m) >= 2 {
			jsonBytes := []byte(`"` + m[1] + `"`)
			var unescaped string
			if err := json.Unmarshal(jsonBytes, &unescaped); err == nil {
				parts = append(parts, unescaped)
			}
		}
	}
	return strings.Join(parts, "")
}