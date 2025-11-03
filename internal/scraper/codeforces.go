package scraper

import (
	"coding-profile-service/pkg/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Base URL for the Codeforces API
const codeforcesBase = "https://codeforces.com/api/"

// CFUserInfo represents response from user.info endpoint
type CFUserInfo struct {
	Handle     string `json:"handle"`
	Rating     int    `json:"rating"`
	MaxRating  int    `json:"maxRating"`
	Rank       string `json:"rank"`
	MaxRank    string `json:"maxRank"`
	Contribution int  `json:"contribution"`
	Country    string `json:"country"`
	Organization string `json:"organization"`
}

// CFRatingChange represents response item from user.rating endpoint
type CFRatingChange struct {
	ContestId   int    `json:"contestId"`
	ContestName string `json:"contestName"`
	Rank        int    `json:"rank"`
	OldRating   int    `json:"oldRating"`
	NewRating   int    `json:"newRating"`
}

// CFSubmission represents response item from user.status endpoint
type CFSubmission struct {
	Id        int64 `json:"id"`
	Problem   struct {
		Index      string `json:"index"`
		Name       string `json:"name"`
		Rating     int    `json:"rating"`
		Tags       []string `json:"tags"`
	} `json:"problem"`
	Verdict    string `json:"verdict"`
	ProgrammingLanguage string `json:"programmingLanguage"`
	CreationTimeSeconds int64 `json:"creationTimeSeconds"`
}

// FetchCodeforces fetches data for a given Codeforces handle
func FetchCodeforces(username string) (model.StatsResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// --- user.info ---
	infoURL := fmt.Sprintf("%suser.info?handles=%s", codeforcesBase, username)
	infoResp, err := client.Get(infoURL)
	if err != nil {
		return errorResp("codeforces", username, err), err
	}
	defer infoResp.Body.Close()

	var infoResult struct {
		Status string       `json:"status"`
		Result []CFUserInfo `json:"result"`
	}
	if err := json.NewDecoder(infoResp.Body).Decode(&infoResult); err != nil {
		return errorResp("codeforces", username, err), err
	}
	if infoResult.Status != "OK" || len(infoResult.Result) == 0 {
		return errorResp("codeforces", username, fmt.Errorf("user not found")), fmt.Errorf("user not found")
	}
	info := infoResult.Result[0]

	// --- user.rating ---
	ratingURL := fmt.Sprintf("%suser.rating?handle=%s", codeforcesBase, username)
	ratingResp, err := client.Get(ratingURL)
	if err != nil {
		return errorResp("codeforces", username, err), err
	}
	defer ratingResp.Body.Close()

	var ratingResult struct {
		Status string            `json:"status"`
		Result []CFRatingChange  `json:"result"`
	}
	if err := json.NewDecoder(ratingResp.Body).Decode(&ratingResult); err != nil {
		return errorResp("codeforces", username, err), err
	}

	contestsParticipated := len(ratingResult.Result)

	// --- user.status ---
	statusURL := fmt.Sprintf("%suser.status?handle=%s&from=1&count=100", codeforcesBase, username)
	statusResp, err := client.Get(statusURL)
	if err != nil {
		return errorResp("codeforces", username, err), err
	}
	defer statusResp.Body.Close()

	var statusResult struct {
		Status string          `json:"status"`
		Result []CFSubmission  `json:"result"`
	}
	if err := json.NewDecoder(statusResp.Body).Decode(&statusResult); err != nil {
		return errorResp("codeforces", username, err), err
	}

	// --- Count solved problems by difficulty (if rating info present) ---
	questionsByType := make(map[string]int)
	uniqueSolved := make(map[string]bool)
	for _, s := range statusResult.Result {
		if s.Verdict == "OK" {
			key := s.Problem.Name
			if !uniqueSolved[key] {
				uniqueSolved[key] = true
				switch {
				case s.Problem.Rating < 1300:
					questionsByType["easy"]++
				case s.Problem.Rating < 1800:
					questionsByType["medium"]++
				default:
					questionsByType["hard"]++
				}
			}
		}
	}

	stats := model.StatsResponse{
		Platform:             "codeforces",
		Username:             username,
		Rating:               info.Rating,
		MaxRating:            info.MaxRating,
		ContestsParticipated: contestsParticipated,
		TotalSolved:          len(uniqueSolved),
		QuestionsByType:      questionsByType,
		GlobalRank:           0, // Codeforces doesn't expose global rank directly
		CountryRank:          0, // not exposed via API
	}

	return stats, nil
}

func errorResp(platform, username string, err error) model.StatsResponse {
	return model.StatsResponse{
		Platform: platform,
		Username: username,
		Error:    fmt.Sprintf("could not fetch %s data for user: %s, error: %v", platform, username, err),
	}
}
