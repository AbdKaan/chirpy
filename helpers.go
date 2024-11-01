package main

import "strings"

var profaneWords = map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}

func cencorProfane(words string) string {
	splitWords := strings.Split(words, " ")

	for i, word := range splitWords {
		if profaneWords[strings.ToLower(word)] {
			splitWords[i] = "****"
		}
	}

	return strings.Join(splitWords, " ")
}