package handler

import (
	"chirpy/core/config"
	"chirpy/dto"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

func NewChirp(apiCfg *config.ApiConfig) *Chirp {
	return &Chirp{
		ApiCfg: apiCfg,
	}
}

type Chirp struct {
	ApiCfg *config.ApiConfig
}

func (c *Chirp) ValidateChirp(w http.ResponseWriter, r *http.Request) {
	invalidWords := []string{"kerfuffle", "sharbert", "fornax"}
	replacement := "****"

	decoder := json.NewDecoder(r.Body)
	var chirp dto.Chirp
	err := decoder.Decode(&chirp)
	if err != nil {
		fmt.Printf("error decoding request: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isValid := len(chirp.Body) <= 140

	w.Header().Set("Content-Type", "application/json")

	if isValid {
		w.WriteHeader(http.StatusOK)

		result := chirp.Body
		for _, word := range invalidWords {
			re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(word))
			result = re.ReplaceAllString(result, replacement)
		}

		marshal, err := json.Marshal(dto.Response{
			CleanedBody: result,
		})
		if err != nil {
			fmt.Printf("error marshaling response: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(marshal)
		if err != nil {
			fmt.Printf("error writing response: %s\n", err)
			return
		}

		return
	} else {
		w.WriteHeader(http.StatusBadRequest)

		marshal, err := json.Marshal(dto.Response{
			Error: "Chirp is too long",
		})
		if err != nil {
			fmt.Printf("error marshaling response: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(marshal)
		if err != nil {
			fmt.Printf("error writing response: %s\n", err)
			return
		}

		return
	}
}
