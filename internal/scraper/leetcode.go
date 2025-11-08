package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"coding-profile-service/pkg/model"
)

func FetchLeetCode(username string) (model.StatsResponse, error) {
	query := fmt.Sprintf(`
    {
      matchedUser(username: "%s") {
        username
        submitStatsGlobal {
          acSubmissionNum {
            difficulty
            count
          }
        }
        profile {
          ranking
          userAvatar
          realName
          reputation
        }
      }
      userContestRanking(username: "%s") {
        attendedContestsCount
        rating
        globalRanking
        topPercentage
      }
      userContestRankingHistory(username: "%s") {
        attended
        rating
        ranking
        contest {
          title
          startTime
        }
      }
    }`, username, username, username)

	body, _ := json.Marshal(map[string]string{"query": query})
	req, _ := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return model.StatsResponse{}, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return model.StatsResponse{}, err
	}

	data := result["data"].(map[string]interface{})["matchedUser"]
	if data == nil {
		return model.StatsResponse{}, fmt.Errorf("user not found")
	}

	stats := data.(map[string]interface{})["submitStatsGlobal"].(map[string]interface{})["acSubmissionNum"].([]interface{})
	
	respModel := model.StatsResponse{
		Platform: "leetcode",
		Username: username,
	}

	// Parse submission statistics
	for _, s := range stats {
		item := s.(map[string]interface{})
		switch item["difficulty"].(string) {
		case "All":
			respModel.TotalSolved = int(item["count"].(float64))
		case "Easy":
			respModel.EasySolved = int(item["count"].(float64))
		case "Medium":
			respModel.MediumSolved = int(item["count"].(float64))
		case "Hard":
			respModel.HardSolved = int(item["count"].(float64))
		}
	}

	// Parse contest ranking data
	contestRanking := result["data"].(map[string]interface{})["userContestRanking"]
	if contestRanking != nil {
		contestData := contestRanking.(map[string]interface{})
		
		if attended, ok := contestData["attendedContestsCount"]; ok && attended != nil {
			respModel.ContestsParticipated = int(attended.(float64))
		}
		
		if rating, ok := contestData["rating"]; ok && rating != nil {
			respModel.Rating = int(rating.(float64))
		}
		
		if globalRank, ok := contestData["globalRanking"]; ok && globalRank != nil {
			respModel.GlobalRank = int(globalRank.(float64))
		}
	}

	// Parse contest history to find max rating
	contestHistory := result["data"].(map[string]interface{})["userContestRankingHistory"]
	if contestHistory != nil {
		historyArray := contestHistory.([]interface{})
		maxRating := 0
		
		for _, contest := range historyArray {
			contestData := contest.(map[string]interface{})
			if rating, ok := contestData["rating"]; ok && rating != nil {
				currentRating := int(rating.(float64))
				if currentRating > maxRating {
					maxRating = currentRating
				}
			}
		}
		
		if maxRating > 0 {
			respModel.MaxRating = maxRating
		}
	}

	return respModel, nil
}