package searchranker

import (
	"memesearch/internal/models"
	"sort"
	"strings"
)

type ScroredMeme struct {
	Score float64
	Meme  models.Meme
}
type Ranker interface {
	Rank(mems []models.Meme, req map[string]string) ([]ScroredMeme, error)
}

var _ Ranker = &DefaultRanker{}

type DefaultRanker struct {
}

// Rank implements Ranker.
func (d *DefaultRanker) Rank(memes []models.Meme, req map[string]string) ([]ScroredMeme, error) {
	res := make([]ScroredMeme, 0, len(memes))
	for _, m := range memes {
		s := d.score(m.Description, req)
		if s < 0.01 {
			continue
		}
		res = append(res, ScroredMeme{Score: s, Meme: m})
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Score > res[j].Score
	})
	return res, nil

}

func (d *DefaultRanker) score(dsc map[string]string, req map[string]string) float64 {
	dscCnt := map[string]int{}
	reqCnt := map[string]int{}
	for _, v := range dsc {
		words := strings.Split(v, " ")
		for _, word := range words {
			dscCnt[word] += 1
		}
	}
	totalReq := 0

	for _, v := range req {
		words := strings.Split(v, " ")
		for _, word := range words {
			reqCnt[word] += 1
			totalReq += 1
		}
	}

	totalScore := 0
	for k := range reqCnt {
		totalScore += min(reqCnt[k], dscCnt[k])
	}
	return float64(totalScore) / float64(totalReq)
}
