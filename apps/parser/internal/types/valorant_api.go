package types

type ValorantProfileImages struct {
	Small        string `json:"small"`
	Large        string `json:"large"`
	TriangleDown string `json:"triangle_down"`
	TriangleUp   string `json:"triangle_up"`
}

type ValorantProfile struct {
	Status int `json:"status"`
	Data   struct {
		Name        string `json:"name"`
		Tag         string `json:"tag"`
		CurrentData struct {
			Currenttier        int    `json:"currenttier"`
			Currenttierpatched string `json:"currenttierpatched"`
			Images             struct {
				Small        string `json:"small"`
				Large        string `json:"large"`
				TriangleDown string `json:"triangle_down"`
				TriangleUp   string `json:"triangle_up"`
			} `json:"images"`
			RankingInTier        int  `json:"ranking_in_tier"`
			MmrChangeToLastGame  int  `json:"mmr_change_to_last_game"`
			Elo                  int  `json:"elo"`
			GamesNeededForRating int  `json:"games_needed_for_rating"`
			Old                  bool `json:"old"`
		} `json:"current_data"`
		HighestRank struct {
			Old         bool   `json:"old"`
			Tier        int    `json:"tier"`
			PatchedTier string `json:"patched_tier"`
			Season      string `json:"season"`
			Converted   int    `json:"converted"`
		} `json:"highest_rank"`
		BySeason struct {
			E1A1 struct {
				Wins             int    `json:"wins"`
				NumberOfGames    int    `json:"number_of_games"`
				FinalRank        int    `json:"final_rank"`
				FinalRankPatched string `json:"final_rank_patched"`
				ActRankWins      []struct {
					PatchedTier string `json:"patched_tier"`
					Tier        int    `json:"tier"`
				} `json:"act_rank_wins"`
				Old bool `json:"old"`
			} `json:"e1a1"`
			E1A2 struct {
				Error string `json:"error"`
			} `json:"e1a2"`
			E1A3 struct {
				Error string `json:"error"`
			} `json:"e1a3"`
			E2A1 struct {
				Wins             int    `json:"wins"`
				NumberOfGames    int    `json:"number_of_games"`
				FinalRank        int    `json:"final_rank"`
				FinalRankPatched string `json:"final_rank_patched"`
				ActRankWins      []struct {
					PatchedTier string `json:"patched_tier"`
					Tier        int    `json:"tier"`
				} `json:"act_rank_wins"`
				Old bool `json:"old"`
			} `json:"e2a1"`
			E2A2 struct {
				Error string `json:"error"`
			} `json:"e2a2"`
			E2A3 struct {
				Error string `json:"error"`
			} `json:"e2a3"`
			E3A1 struct {
				Error string `json:"error"`
			} `json:"e3a1"`
			E3A2 struct {
				Error string `json:"error"`
			} `json:"e3a2"`
			E3A3 struct {
				Error string `json:"error"`
			} `json:"e3a3"`
			E4A1 struct {
				Error string `json:"error"`
			} `json:"e4a1"`
			E4A2 struct {
				Error string `json:"error"`
			} `json:"e4a2"`
			E4A3 struct {
				Wins             int    `json:"wins"`
				NumberOfGames    int    `json:"number_of_games"`
				FinalRank        int    `json:"final_rank"`
				FinalRankPatched string `json:"final_rank_patched"`
				ActRankWins      []struct {
					PatchedTier string `json:"patched_tier"`
					Tier        int    `json:"tier"`
				} `json:"act_rank_wins"`
				Old bool `json:"old"`
			} `json:"e4a3"`
			E5A1 struct {
				Error string `json:"error"`
			} `json:"e5a1"`
			E5A2 struct {
				Error string `json:"error"`
			} `json:"e5a2"`
			E5A3 struct {
				Error string `json:"error"`
			} `json:"e5a3"`
			E6A1 struct {
				Error string `json:"error"`
			} `json:"e6a1"`
			E6A2 struct {
				Error string `json:"error"`
			} `json:"e6a2"`
			E6A3 struct {
				Error string `json:"error"`
			} `json:"e6a3"`
			E7A1 struct {
				Error string `json:"error"`
			} `json:"e7a1"`
			E7A2 struct {
				Error string `json:"error"`
			} `json:"e7a2"`
			E7A3 struct {
				Error string `json:"error"`
			} `json:"e7a3"`
			E8A1 struct {
				Error string `json:"error"`
			} `json:"e8a1"`
			E8A2 struct {
				Error string `json:"error"`
			} `json:"e8a2"`
			E8A3 struct {
				Error string `json:"error"`
			} `json:"e8a3"`
		} `json:"by_season"`
	} `json:"data"`
}

type ValorantMatchPlayer struct {
	Puuid              string `json:"puuid"`
	Name               string `json:"name"`
	Tag                string `json:"tag"`
	Team               string `json:"team"`
	Level              int    `json:"level"`
	Character          string `json:"character"`
	CurrentTier        int    `json:"currenttier"`
	CurrentTierPatched string `json:"currenttier_patched"`
	Behavior           struct {
		AfkRounds    float64 `json:"afk_rounds"`
		FriendlyFire struct {
			Incoming int `json:"incoming"`
			Outgoing int `json:"outgoing"`
		} `json:"friendly_fire"`
		RoundsInSpawn float64 `json:"rounds_in_spawn"`
	} `json:"behavior"`
	Stats struct {
		Score     int `json:"score"`
		Kills     int `json:"kills"`
		Deaths    int `json:"deaths"`
		Assists   int `json:"assists"`
		Bodyshots int `json:"bodyshots"`
		Headshots int `json:"headshots"`
		Legshots  int `json:"legshots"`
	} `json:"stats"`
	Economy struct {
		Spent struct {
			Overall int `json:"overall"`
			Average int `json:"average"`
		} `json:"spent"`
		LoadoutValue struct {
			Overall int `json:"overall"`
			Average int `json:"average"`
		} `json:"loadout_value"`
	} `json:"economy"`
	DamageMade     int `json:"damage_made"`
	DamageReceived int `json:"damage_received"`
}

type ValorantMatchPlayers struct {
	AllPlayers []ValorantMatchPlayer `json:"all_players"`
}

type ValorantMatchesResponse struct {
	Data []ValorantMatch `json:"data"`
}

type ValorantMatch struct {
	MetaData struct {
		Map              string `json:"map"`
		GameVersion      string `json:"game_version"`
		GameLength       int    `json:"game_length"`
		GameStart        int    `json:"game_start"`
		GameStartPatched string `json:"game_start_patched"`
		RoundsPlayed     int    `json:"rounds_played"`
		Mode             string `json:"mode"`
		Queue            string `json:"queue"`
		SeasonID         string `json:"season_id"`
		Platform         string `json:"platform"`
		MatchID          string `json:"match_id"`
		Region           string `json:"region"`
		Cluster          string `json:"cluster"`
	}
	Players ValorantMatchPlayers `json:"players"`
	Teams   map[string]struct {
		HasWon     bool `json:"has_won"`
		RoundsWon  int  `json:"rounds_won"`
		RoundsLost int  `json:"rounds_lost"`
	}
}

type ValorantMmrHistoryMatch struct {
	Currenttier        int    `json:"currenttier"`
	CurrenttierPatched string `json:"currenttier_patched"`
	Images             struct {
		Small        string `json:"small"`
		Large        string `json:"large"`
		TriangleDown string `json:"triangle_down"`
		TriangleUp   string `json:"triangle_up"`
	} `json:"images"`
	MatchId string `json:"match_id"`
	Map     struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"map"`
	SeasonId            string `json:"season_id"`
	RankingInTier       int    `json:"ranking_in_tier"`
	MmrChangeToLastGame int    `json:"mmr_change_to_last_game"`
	Elo                 int    `json:"elo"`
	Date                string `json:"date"`
	DateRaw             int    `json:"date_raw"`
}

type ValorantMmrHistoryResponse struct {
	Data []*ValorantMmrHistoryMatch `json:"data"`
}
