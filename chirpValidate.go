package main

import (
	"fmt"
	"strings"
)

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", fmt.Errorf("chirp is too long")
	}
	cleanedBody := cleanChirpProfanity(body)
	return cleanedBody, nil
}

func cleanChirpProfanity(body string) string {
	prohibitedWords := map[string]interface{}{
		"kerfuffle": struct{}{},
		"sharbert":  struct{}{},
		"fornax":    struct{}{},
	}

	words := strings.Split(body, " ")
	result := []string{}
	for _, word := range words {
		lowerWord := strings.ToLower(word)
		_, found := prohibitedWords[lowerWord]
		if found {
			result = append(result, "****")
		} else {
			result = append(result, word)
		}
	}
	newBody := strings.Join(result, " ")
	return newBody
}
