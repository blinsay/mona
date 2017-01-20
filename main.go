package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	MonaLisaPath = "data/mona_lisa.jpg"
)

// A function that generates a score for an image. Lower scores are better.
type ScoreFunc func(*image.Gray) float64

// A func that mutates an image from one generation to the next. Must be safe to
// call from multiple goroutines.
type MutationFunc func(*image.Gray)

type Experiment struct {
	BaselineImage *image.Gray
	BaselineScore float64
	Candidates    []Candidate
}

type Candidate struct {
	sync.Mutex
	data image.Gray
}

func (c *Candidate) Unwrap() *image.Gray {
	c.Lock()
	defer c.Unlock()

	return &c.data
}

func (c *Candidate) Score(s ScoreFunc) float64 {
	c.Lock()
	defer c.Unlock()

	return s(&c.data)
}

func (c *Candidate) Do(f func(*image.Gray)) {
	c.Lock()
	defer c.Unlock()

	f(&c.data)
}

// FIXME: doc + test
func NewExperiment(size int, src *image.Gray, initialScore float64) *Experiment {
	copySrc := func(img *image.Gray) {
		CopyGray(img, src)
	}

	candidates := make([]Candidate, size)
	for i := range candidates {
		candidates[i].Do(copySrc)
	}

	return &Experiment{
		BaselineImage: src,
		BaselineScore: initialScore,
		Candidates:    candidates,
	}
}

// FIXME: doc + test
func Winner(experiment *Experiment, score ScoreFunc) (*image.Gray, float64) {
	winner, minScore := experiment.BaselineImage, experiment.BaselineScore

	for i := range experiment.Candidates {
		c := &experiment.Candidates[i]
		if s := c.Score(score); s < minScore {
			winner, minScore = c.Unwrap(), s
		}
	}

	return winner, minScore
}

func ParallelDo(candidates []Candidate, f func(*image.Gray)) {
	wg := sync.WaitGroup{}
	wg.Add(len(candidates))

	for i := range candidates {
		go func(c *Candidate) {
			c.Do(f)
			wg.Done()
		}(&candidates[i])
	}

	wg.Wait()
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	mona, err := ReadJpeg(MonaLisaPath)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	baseImage := NewImage(mona.Bounds(), color.White)
	distanceFunc := EuclideanDistanceTo(ConvertToGray(mona))
	rectangleGenerator := GenerateRects(50, mona.Bounds())
	expt := NewExperiment(100, baseImage, distanceFunc(baseImage))

	lastScore := expt.BaselineScore
	for n := 0; n < 100000; n++ {
		ParallelDo(expt.Candidates, rectangleGenerator.Apply)
		// FIXME: make the score calculation parallel as well. do it in the same func as apply?
		winner, score := Winner(expt, distanceFunc)

		log.Printf("%d: winning score: %.06f", n, score)

		if lastScore != score {
			outFile, err := os.Create(fmt.Sprintf("data/out/test_gen_%09d.jpg", n))
			if err != nil {
				log.Fatalf("error: %s", err)
			}
			jpeg.Encode(outFile, winner, nil)
			lastScore = score
		}

		expt = NewExperiment(10, winner, score)
	}
}
