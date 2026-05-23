// package scraper

// import (
// 	"coding-profile-service/pkg/model"
// 	"fmt"
// )

// // FetchGFG wraps FetchGFGHTML and returns StatsResponse
// func FetchGFG(username string) (model.StatsResponse, error) {
// 	stats, err := FetchGFGHTML(username)
// 	if err != nil {
// 		return model.StatsResponse{
// 			Platform: "gfg",
// 			Username: username,
// 			Error:    fmt.Sprintf("could not fetch GFG data for user: %s (%v)", username, err),
// 		}, err
// 	}

// 	return model.StatsResponse{
// 		Platform:             "gfg",
// 		Username:             username,
// 		TotalSolved:          stats.TotalSolved,
// 		EasySolved:           stats.EasySolved,
// 		MediumSolved:         stats.MediumSolved,
// 		HardSolved:           stats.HardSolved,
// 		Streak:               stats.Streak,
// 		MaxRating:            stats.MaxRating,
// 		ContestsParticipated: stats.ContestsParticipated,
// 		GlobalRank:           stats.GlobalRank,
// 		CountryRank:          stats.CountryRank,
// 	}, nil
// }

package scraper

import (
	"coding-profile-service/pkg/model"
	"fmt"
)

func FetchGFG(username string) (model.StatsResponse, error) {
	profile, err := FetchGFGProfile(username)
	if err != nil {
		return model.StatsResponse{
			Platform: "gfg",
			Username: username,
			Error:    fmt.Sprintf("could not fetch GFG data for user: %s (%v)", username, err),
		}, err
	}

	return model.StatsResponse{
		Platform:             "gfg",
		Username:             username,
		TotalSolved:          profile.TotalSolved,
		EasySolved:           profile.EasySolved,
		MediumSolved:         profile.MediumSolved,
		HardSolved:           profile.HardSolved,
		Streak:               profile.Streak,
		MaxRating:            profile.MaxRating,
		CodingScore:          profile.CodingScore,
		GlobalRank:           profile.GlobalRank,
		ContestsParticipated: profile.ContestsParticipated,
	}, nil
}