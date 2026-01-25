package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// GFGProfile holds the scraped data temporarily
type GFGProfile struct {
	TotalSolved          int
	Streak               int
	EasySolved           int
	MediumSolved         int
	HardSolved           int
	ContestsParticipated int
	MaxRating            int
	GlobalRank           int
	CountryRank          int
}

// FetchGFGHTML scrapes the GFG user profile using headless Chrome
func FetchGFGHTML(username string) (*GFGProfile, error) {
	// Configure Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-extensions", true),
	)

	// Use chromium from environment or default path
	chromePath := os.Getenv("CHROME_PATH")
	if chromePath == "" {
		chromePath = "/usr/bin/chromium-browser"
	}
	
	// Check if chromium exists, if so set the path
	if _, err := os.Stat(chromePath); err == nil {
		opts = append(opts, chromedp.ExecPath(chromePath))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://www.geeksforgeeks.org/profile/%s?tab=activity", username)
	
	var result string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.ScoreContainer_score-card__zI4vG`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		
		// Extract all data in one JavaScript call
		chromedp.Evaluate(`
			JSON.stringify({
				codingScore: (() => {
					const cards = document.querySelectorAll('.ScoreContainer_score-card__zI4vG');
					for (let card of cards) {
						const label = card.querySelector('.ScoreContainer_label__aVpLE')?.textContent?.trim();
						if (label === 'Coding Score') {
							return card.querySelector('.ScoreContainer_value__7yy7h')?.textContent?.trim() || '0';
						}
					}
					return '0';
				})(),
				problemsSolved: (() => {
					const cards = document.querySelectorAll('.ScoreContainer_score-card__zI4vG');
					for (let card of cards) {
						const label = card.querySelector('.ScoreContainer_label__aVpLE')?.textContent?.trim();
						if (label === 'Problems Solved') {
							return card.querySelector('.ScoreContainer_value__7yy7h')?.textContent?.trim() || '0';
						}
					}
					return '0';
				})(),
				instituteRank: (() => {
					const cards = document.querySelectorAll('.ScoreContainer_score-card__zI4vG');
					for (let card of cards) {
						const label = card.querySelector('.ScoreContainer_label__aVpLE')?.textContent?.trim();
						if (label === 'Institute Rank') {
							const val = card.querySelector('.ScoreContainer_value__7yy7h')?.textContent?.trim();
							return val === '__' ? '0' : (val || '0');
						}
					}
					return '0';
				})(),
				streak: (() => {
					const el = document.querySelector('.PotdContainer_streakText__oNgWh');
					if (el) {
						const match = el.textContent.trim().match(/^(\d+)/);
						return match ? match[1] : '0';
					}
					return '0';
				})(),
				school: (() => {
					const items = document.querySelectorAll('.ProblemNavbar_head_nav--text__7u4wN');
					for (let item of items) {
						const match = item.textContent.trim().match(/SCHOOL\s*\((\d+)\)/i);
						if (match) return match[1];
					}
					return '0';
				})(),
				basic: (() => {
					const items = document.querySelectorAll('.ProblemNavbar_head_nav--text__7u4wN');
					for (let item of items) {
						const match = item.textContent.trim().match(/BASIC\s*\((\d+)\)/i);
						if (match) return match[1];
					}
					return '0';
				})(),
				easy: (() => {
					const items = document.querySelectorAll('.ProblemNavbar_head_nav--text__7u4wN');
					for (let item of items) {
						const match = item.textContent.trim().match(/EASY\s*\((\d+)\)/i);
						if (match) return match[1];
					}
					return '0';
				})(),
				medium: (() => {
					const items = document.querySelectorAll('.ProblemNavbar_head_nav--text__7u4wN');
					for (let item of items) {
						const match = item.textContent.trim().match(/MEDIUM\s*\((\d+)\)/i);
						if (match) return match[1];
					}
					return '0';
				})(),
				hard: (() => {
					const items = document.querySelectorAll('.ProblemNavbar_head_nav--text__7u4wN');
					for (let item of items) {
						const match = item.textContent.trim().match(/HARD\s*\((\d+)\)/i);
						if (match) return match[1];
					}
					return '0';
				})()
			})
		`, &result),
	)

	if err != nil {
		return nil, fmt.Errorf("chromedp error: %v", err)
	}

	// Parse JSON result
	var data map[string]string
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return nil, fmt.Errorf("failed to parse result: %v", err)
	}

	profile := &GFGProfile{}
	
	// Helper to convert string to int
	toInt := func(s string) int {
		var val int
		fmt.Sscanf(s, "%d", &val)
		return val
	}

	profile.MaxRating = toInt(data["codingScore"])
	profile.TotalSolved = toInt(data["problemsSolved"])
	profile.GlobalRank = toInt(data["instituteRank"])
	profile.Streak = toInt(data["streak"])
	
	// Combine School + Basic + Easy
	profile.EasySolved = toInt(data["school"]) + toInt(data["basic"]) + toInt(data["easy"])
	profile.MediumSolved = toInt(data["medium"])
	profile.HardSolved = toInt(data["hard"])

	return profile, nil
}