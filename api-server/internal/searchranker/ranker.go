package searchranker

import (
	"context"
	"log/slog"
	"memesearch/internal/models"
	"sort"
	"strings"
)

type ScroredMeme struct {
	Score float64
	Meme  models.Meme
}
type Ranker interface {
	Rank(ctx context.Context, mems []models.Meme, req map[string]string) ([]ScroredMeme, error)
}

var _ Ranker = &DefaultRanker{}

type DefaultRanker struct {
}

// Rank implements Ranker.
func (dr *DefaultRanker) Rank(ctx context.Context, memes []models.Meme, req map[string]string) ([]ScroredMeme, error) {
	res := make([]ScroredMeme, 0, len(memes))
	for _, m := range memes {
		s := dr.score(ctx, m.Description, req)
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

func (dr *DefaultRanker) score(ctx context.Context, dsc map[string]string, req map[string]string) float64 {
	r, ok := req["general"]
	if !ok || len(r) == 0 {
		slog.DebugContext(ctx, "Can't find 'general' field in request")
		return -1
	}
	d, ok := dsc["general"]
	if !ok || len(r) == 0 {
		slog.DebugContext(ctx, "Can't find 'general' field in description")
		return -1
	}

	rs := strings.Split(strings.ToLower(r), " ")
	ds := strings.Split(strings.ToLower(d), " ")
	n := 0
	totalScore := 0.0
	for _, i := range rs {
		if commonlyUsed(i) {
			continue
		}
		n += 1
		mx := 0.0
		for _, j := range ds {
			if commonlyUsed(j) {
				continue
			}
			score := 1 - normlizedLevenstainDist(i, j)
			mx = max(mx, score)
		}
		if mx > 0.5 {
			totalScore += mx
		}
	}
	return totalScore / float64(n)
}

func normlizedLevenstainDist(as, bs string) float64 {
	a := []rune(as)
	b := []rune(bs)

	la, lb := len(a), len(b)
	matrix := make([][]float64, la+1)
	for i := range matrix {
		matrix[i] = make([]float64, lb+1)
	}

	for i := 0; i <= la; i++ {
		matrix[i][0] = float64(i)
	}
	for j := 0; j <= lb; j++ {
		matrix[0][j] = float64(j)
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			matrix[i][j] = min(
				matrix[i-1][j]+1,                      // Удаление
				matrix[i][j-1]+1,                      // Вставка
				matrix[i-1][j-1]+cost(a[i-1], b[j-1]), // Замена
			)
		}
	}

	return matrix[la][lb] / float64(min(len(a), len(b)))
}

var keyboardLayout = map[rune]struct{ x, y int }{
	'q': {0, 0}, 'w': {1, 0}, 'e': {2, 0}, 'r': {3, 0}, 't': {4, 0}, 'y': {5, 0}, 'u': {6, 0}, 'i': {7, 0}, 'o': {8, 0}, 'p': {9, 0},
	'a': {0, 1}, 's': {1, 1}, 'd': {2, 1}, 'f': {3, 1}, 'g': {4, 1}, 'h': {5, 1}, 'j': {6, 1}, 'k': {7, 1}, 'l': {8, 1},
	'z': {0, 2}, 'x': {1, 2}, 'c': {2, 2}, 'v': {3, 2}, 'b': {4, 2}, 'n': {5, 2}, 'm': {6, 2},

	'й': {0, 0}, 'ц': {1, 0}, 'у': {2, 0}, 'к': {3, 0}, 'е': {4, 0}, 'н': {5, 0}, 'г': {6, 0}, 'ш': {7, 0}, 'щ': {8, 0}, 'з': {9, 0}, 'х': {10, 0}, 'ъ': {11, 0},
	'ф': {0, 1}, 'ы': {1, 1}, 'в': {2, 1}, 'а': {3, 1}, 'п': {4, 1}, 'р': {5, 1}, 'о': {6, 1}, 'л': {7, 1}, 'д': {8, 1}, 'ж': {9, 1}, 'э': {10, 1},
	'я': {0, 2}, 'ч': {1, 2}, 'с': {2, 2}, 'м': {3, 2}, 'и': {4, 2}, 'т': {5, 2}, 'ь': {6, 2}, 'б': {7, 2}, 'ю': {8, 2},
}

func cost(a, b rune) float64 {
	pa, ok := keyboardLayout[a]
	if !ok {
		return 100
	}
	pb, ok := keyboardLayout[b]
	if !ok {
		return 100
	}
	dst := abs(pa.x-pb.x) + abs(pa.y-pb.y)

	switch dst {
	case 0:
		return 0
	case 1:
		return 0.2
	case 2:
		return 1
	default:
		return 2
	}
}

func abs(i int) int {
	if i > 0 {
		return i
	}
	return -i
}

// var commonly map[string]struct{} = map[string]struct{}{
// 	"или": struct{}{},
// }

func commonlyUsed(a string) bool {
	if len([]rune(a)) <= 2 {
		return true
	}
	return false
}
