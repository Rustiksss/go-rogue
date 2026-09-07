package data

import (
	"encoding/json"
	"io"
	"os"
	"rogue/domain"
)

type Data struct{}

const GameSaveFileName = "rogue_game.save"
const LeaderboardSaveFileName = "rogue_leaderboard.save"

func NewData() *Data {
	return &Data{}
}

func (d *Data) SaveLeaderboard(leaderboard []domain.Statistics) error {
	file, err := os.Create(LeaderboardSaveFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	serializedData, err := json.Marshal(leaderboard)

	if err != nil {
		return err
	}

	file.Write(serializedData)

	return err
}

func (d *Data) LoadLeaderboard() ([]domain.Statistics, error) {
	file, err := os.Open(LeaderboardSaveFileName)
	if os.IsNotExist(err) {
		return []domain.Statistics{}, nil
	}
	if err != nil {
		return []domain.Statistics{}, err
	}
	defer file.Close()

	serializedData, err := io.ReadAll(file)

	if err != nil && !os.IsNotExist(err) {
		return []domain.Statistics{}, err
	}

	var leaderboard []domain.Statistics
	err = json.Unmarshal(serializedData, &leaderboard)

	return leaderboard, err
}

func (d *Data) SaveGame(g domain.GameInfo) error {
	file, err := os.Create(GameSaveFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	serializedData, err := json.MarshalIndent(g, "  ", "")
	if err != nil {
		return err
	}
	_, err = file.Write(serializedData)

	return err
}

func (d *Data) LoadGame() (domain.GameInfo, error) {
	serializedData, err := os.ReadFile(GameSaveFileName)
	if err != nil && !os.IsNotExist(err) {
		return domain.GameInfo{}, err
	}

	var g domain.GameInfo
	err = json.Unmarshal(serializedData, &g)

	return g, err
}

func (d *Data) HasSavedGame() (has bool, err error) {
	has = true
	_, err = os.Stat(GameSaveFileName)
	if os.IsNotExist(err) {
		err = nil
		has = false
	}
	return
}

func (d *Data) DeleteSavedGame() (err error) {
	err = os.Remove(GameSaveFileName)
	if os.IsNotExist(err) {
		err = nil
	}
	return
}
