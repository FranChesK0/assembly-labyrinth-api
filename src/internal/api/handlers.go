package api

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/franchesko/assembly-labyrinth/src/internal/config"
	"github.com/franchesko/assembly-labyrinth/src/internal/emu"
	"github.com/franchesko/assembly-labyrinth/src/internal/emu/program"
	"github.com/franchesko/assembly-labyrinth/src/internal/files"
	"github.com/gorilla/mux"
)

type LevelsResponse struct {
	Levels []string `json:"levels"`
}

type LevelInfoResponse struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Layout      []nodeLayoutResponse `json:"layout"`
	In          []ioeStreamResponse  `json:"in"`
	Expected    []ioeStreamResponse  `json:"expected"`
}

type RunLevelResponse struct {
	CodeValidation bool                `json:"code_validation"`
	CheckStatus    bool                `json:"check_status"`
	In             []ioeStreamResponse `json:"in"`
	Expected       []ioeStreamResponse `json:"expected"`
	Out            []ioeStreamResponse `json:"out"`
}

type nodeLayoutResponse struct {
	Index uint8        `json:"index"`
	Type  emu.NodeType `json:"type"`
}

type ioeStreamResponse struct {
	Index  uint8   `json:"index"`
	Name   string  `json:"name,omitempty"`
	Values []int16 `json:"values"`
}

type runLevelRequest struct {
	Nodes    []emu.NodeCode      `json:"nodes"`
	In       []ioeStreamResponse `json:"in"`
	Expected []ioeStreamResponse `json:"expected"`
}

func GetLevelsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		levels, err := files.LoadLevels(cfg.LevelPath)
		if err != nil {
			http.Error(w, "Unable to load levels", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(LevelsResponse{Levels: levels})
	}
}

func GetLevelInfoHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		levelInfo, err := files.LoadLevelInfo(cfg.LevelPath + "/" + params["level"] + ".json")
		if err != nil {
			http.Error(w, "Unable to load level", http.StatusInternalServerError)
			return
		}
		layout := make([]nodeLayoutResponse, 0)
		for _, nl := range levelInfo.Layout {
			layout = append(layout, nodeLayoutResponse{
				Index: nl.Index,
				Type:  nl.Type,
			})
		}

		in := make([]ioeStreamResponse, 0)
		out := make([]ioeStreamResponse, 0)
		for i := range levelInfo.Streams {
			if levelInfo.Streams[i].Type == emu.IN {
				levelInfo.Streams[i].Values = generateValues(int(levelInfo.Streams[i].MinValue), int(levelInfo.Streams[i].MaxValue))
				in = append(in, ioeStreamResponse{
					Index:  levelInfo.Streams[i].Index,
					Name:   levelInfo.Streams[i].Name,
					Values: levelInfo.Streams[i].Values,
				})
			} else {
				out = append(out, ioeStreamResponse{
					Index: levelInfo.Streams[i].Index,
					Name:  levelInfo.Streams[i].Name,
				})
			}
		}

		code, err := files.LoadNodesCode(cfg.CodePath + "/" + params["level"] + ".json")
		if err != nil {
			http.Error(w, "Unable to load level code", http.StatusInternalServerError)
			return
		}
		expected, err := program.Run(levelInfo.Streams, code)
		if err != nil {
			http.Error(w, "Unable to get expected values", http.StatusInternalServerError)
			return
		}

		expRes := make([]ioeStreamResponse, 0)
		for i, expStream := range expected {
			expRes = append(expRes, ioeStreamResponse{
				Index:  out[i].Index,
				Name:   out[i].Name,
				Values: expStream.Values,
			})
		}

		json.NewEncoder(w).Encode(LevelInfoResponse{
			Title:       levelInfo.Title,
			Description: levelInfo.Description,
			Layout:      layout,
			In:          in,
			Expected:    expRes,
		})
	}
}

func GetRunLevelHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		var runLevel runLevelRequest
		if err := json.NewDecoder(r.Body).Decode(&runLevel); err != nil {
			http.Error(w, "Wrong request format", http.StatusBadRequest)
			return
		}

		codeValidation := true
		status := true
		levelInfo, err := files.LoadLevelInfo(cfg.LevelPath + "/" + params["level"] + ".json")
		if err != nil {
			http.Error(w, "Unable to load level", http.StatusInternalServerError)
			return
		}
		for i := range levelInfo.Streams {
			if levelInfo.Streams[i].Type == emu.IN {
				levelInfo.Streams[i].Values = runLevel.In[i].Values
			}
		}
		out, err := program.Run(levelInfo.Streams, runLevel.Nodes)
		if err != nil {
			codeValidation = false
			status = false
		}
		outResp := make([]ioeStreamResponse, 0)
		if codeValidation {
			for _, outStream := range out {
				outResp = append(outResp, ioeStreamResponse{
					Index:  outStream.Index,
					Values: outStream.Values,
				})
			}
			status = checkResult(runLevel.Expected, outResp)
		}

		json.NewEncoder(w).Encode(RunLevelResponse{
			CodeValidation: codeValidation,
			CheckStatus:    status,
			In:             runLevel.In,
			Expected:       runLevel.Expected,
			Out:            outResp,
		})
	}
}

func generateValues(minValue, maxValue int) []int16 {
	out := make([]int16, 0)
	for range emu.StreamLength {
		out = append(out, int16(rand.Intn(maxValue-minValue)+minValue))
	}
	return out
}

func checkResult(expected, out []ioeStreamResponse) bool {
	if len(expected) != len(out) {
		return false
	}

	for i := range expected {
		if len(expected[i].Values) != len(out[i].Values) {
			return false
		}

		for j := range expected[i].Values {
			if expected[i].Values[j] != out[i].Values[j] {
				return false
			}
		}
	}

	return true
}
