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