package dota

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables/stream"

	model "tsuwari/parser/internal/models"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/lib/pq"
	"github.com/samber/lo"
)

type MatchPlayer struct {
	AccountId  int `json:"account_id"`
	TeamNumber int `json:"team_number"`
	HeroId     int `json:"hero_id"`
	Kills      int `json:"kills"`
	Deaths     int `json:"deaths"`
	Assists    int `json:"assists"`
}

type MatchResult struct {
	Error      *string       `json:"error"`
	MatchID    int           `json:"match_id"`
	Players    []MatchPlayer `json:"players"`
	RadiantWin bool          `json:"radiant_win"`
	GameMode   int32         `json:"game_mode"`
}

type MatchResponse struct {
	Result MatchResult `json:"result"`
}

var WlCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "wl",
		Description: lo.ToPtr("Score for played games on stream"),
		Permission:  "VIEWER",
		Visible:     true,
		Module:      lo.ToPtr("DOTA"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		streamString, err := ctx.Services.Redis.Get(context.TODO(), "streams:"+ctx.ChannelId).
			Result()
		streamData := stream.HelixStream{}

		if err != nil || streamString == "" {
			return []string{"Stream not found"}
		}

		err = json.Unmarshal([]byte(streamString), &streamData)

		if err != nil {
			fmt.Println(err)
			return []string{"Something went wrong on getting stream."}
		}

		accounts := GetAccountsByChannelId(ctx.Services.Db, ctx.ChannelId)

		if accounts == nil || len(*accounts) == 0 {
			return []string{NO_ACCOUNTS}
		}

		dbGames := []model.DotaMatchWithRelation{}
		intAccounts := lo.Map(*accounts, func(a string, _ int) int {
			acc, _ := strconv.Atoi(a)
			return acc
		})

		err = ctx.Services.Db.
			Table("dota_matches").
			Where(
				`"dota_matches"."players" && ?
				AND "startedAt" >= ? 
				AND "lobby_type" IN ?
				`,
				pq.Array(intAccounts),
				streamData.StartedAt.Add(-time.Minute*10),
				[2]int{0, 7},
			).
			Order(`"startedAt" DESC`).
			Joins("GameMode").
			Joins("Result").
			Find(&dbGames).Error

		if err != nil {
			fmt.Println(err)
			return []string{"Something went wrong on fetching games."}
		}

		dbGames = lo.Filter(dbGames, func(g model.DotaMatchWithRelation, _ int) bool {
			return !lo.SomeBy(g.PlayersHeroes, func(id int64) bool {
				return id == 0
			})
		})

		if len(dbGames) == 0 {
			return []string{"Games not played on the stream."}
		}

		matchesData := []MatchResult{}

		for _, match := range dbGames {
			if match.Result != nil {
				matchId, _ := strconv.Atoi(match.MatchID)
				players := []MatchPlayer{}
				_ = json.Unmarshal([]byte(match.Result.Players), &players)
				result := MatchResult{
					Error:      nil,
					MatchID:    matchId,
					RadiantWin: match.Result.RadiantWin,
					GameMode:   match.Result.GameMode,
					Players:    players,
				}
				matchesData = append(matchesData, result)
			}
		}

		gamesForRequest := lo.Filter(dbGames, func(g model.DotaMatchWithRelation, _ int) bool {
			return g.Result == nil
		})

		gamesForRequestWg := sync.WaitGroup{}
		for _, game := range gamesForRequest {
			gamesForRequestWg.Add(1)

			go func(game model.DotaMatchWithRelation) {
				defer gamesForRequestWg.Done()
				data := MatchResponse{}
				r, err := ApiInstance.
					R().
					SetQueryParams(map[string]string{
						"match_id": game.MatchID,
						"key":      "2B5C2069282D28E79B60B494489E31C5",
					}).
					SetResult(&data).
					Get("http://api.steampowered.com/IDOTA2Match_570/GetMatchDetails/v1")

				if err != nil || r.StatusCode != 200 || data.Result.Error != nil {
					fmt.Println(r.StatusCode, err, *data.Result.Error, game)
					return
				}

				dataForPush := MatchResult{
					MatchID:    data.Result.MatchID,
					Players:    data.Result.Players,
					RadiantWin: data.Result.RadiantWin,
					GameMode:   data.Result.GameMode,
				}

				matchesData = append(matchesData, dataForPush)
				players, _ := json.Marshal(data.Result.Players)
				err = ctx.Services.Db.
					Table("dota_matches_results").
					Create(map[string]interface{}{
						"match_id":    strconv.Itoa(data.Result.MatchID),
						"players":     string(players),
						"radiant_win": data.Result.RadiantWin,
						"game_mode":   data.Result.GameMode,
					}).
					Error
				if err != nil {
					fmt.Println(game.MatchID, err)
				}
				ctx.Services.Db.
					Model(&model.DotaMatches{}).
					Where("match_id = ?", game.MatchID).
					Update("finished", true)
			}(game)
		}

		gamesForRequestWg.Wait()

		matchesByGameMode := make(map[int32]MatchByGameMode)

		for _, mode := range DotaGameModes {
			matchesByGameMode[int32(mode.ID)] = MatchByGameMode{
				Matches: []Match{},
			}
		}

		for _, account := range *accounts {
			for _, match := range matchesData {
				dbMatch, ok := lo.Find(dbGames, func(m model.DotaMatchWithRelation) bool {
					return m.MatchID == strconv.Itoa(match.MatchID)
				})

				if !ok {
					continue
				}

				_, playerIndex, ok := lo.FindIndexOf(dbMatch.Players, func(p int64) bool {
					stringedP := strconv.FormatInt(p, 10)
					return stringedP == account
				})

				if !ok {
					continue
				}

				player := match.Players[playerIndex]
				isPlayerRadiant := player.TeamNumber == 0
				isWinner := lo.
					If(isPlayerRadiant && match.RadiantWin, true).
					ElseIf(!isPlayerRadiant && !match.RadiantWin, true).
					Else(false)
				hero, _ := lo.Find(DotaHeroes, func(h Hero) bool {
					return h.ID == player.HeroId
				})

				shortHeroName := lo.FromPtr(hero.ShortName)
				heroName := lo.
					If(shortHeroName != "", shortHeroName).
					Else(hero.LocalizedName)

				if entry, ok := matchesByGameMode[match.GameMode]; ok {
					entry.Matches = append(
						entry.Matches,
						Match{
							IsWinner: isWinner,
							Hero:     heroName,
							Kills:    player.Kills,
							Deaths:   player.Deaths,
							Assists:  player.Assists,
						},
					)

					matchesByGameMode[match.GameMode] = entry
				}
			}
		}

		result := []string{}

		for _, entry := range lo.Entries(matchesByGameMode) {
			if len(entry.Value.Matches) == 0 {
				continue
			}
			wins := lo.Reduce(entry.Value.Matches, func(acc int, m Match, _ int) int {
				if m.IsWinner {
					return acc + 1
				} else {
					return acc
				}
			}, 0)
			mode, _ := lo.Find(DotaGameModes, func(m GameMode) bool {
				return int(entry.Key) == m.ID
			})

			heroesResult := lo.Map(entry.Value.Matches, func(m Match, _ int) string {
				return fmt.Sprintf(
					"%s(%s) [%v/%v/%v]",
					m.Hero,
					lo.If(m.IsWinner, "W").Else("L"),
					m.Kills,
					m.Deaths,
					m.Assists,
				)
			})

			msg := fmt.Sprintf(
				"%s W %v — %v: %s",
				mode.Name,
				wins,
				len(entry.Value.Matches)-wins,
				strings.Join(heroesResult, ", "),
			)

			result = append(result, msg)
		}

		if len(result) == 0 {
			return []string{"W 0 — L 0"}
		}

		return []string{strings.Join(result, "")}
	},
}

type Match struct {
	IsWinner bool
	Hero     string
	Kills    int
	Deaths   int
	Assists  int
}

type MatchByGameMode struct {
	Matches []Match
}
