package files

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/franchesko/assembly-labyrinth/src/internal/emu"
)

type LevelInfo struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Layout      []nodeLayout `json:"layout"`
	Streams     []emu.Stream `json:"streams"`
}

type nodeLayout struct {
	Index uint8        `json:"index"`
	Type  emu.NodeType `json:"type"`
}

func LoadLevels(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	levels := make([]string, 0)
	for _, file := range files {
		levels = append(levels, strings.Split(file.Name(), ".")[0])
	}

	return levels, nil
}

func LoadLevelInfo(filePath string) (LevelInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return LevelInfo{}, err
	}
	defer file.Close()

	var levelInfo LevelInfo
	if err = json.NewDecoder(file).Decode(&levelInfo); err != nil {
		return LevelInfo{}, err
	}

	return levelInfo, nil
}

func LoadNodesCode(filePath string) ([]emu.NodeCode, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nodesCodes []emu.NodeCode
	if err = json.NewDecoder(file).Decode(&nodesCodes); err != nil {
		return nil, err
	}

	return nodesCodes, nil
}
