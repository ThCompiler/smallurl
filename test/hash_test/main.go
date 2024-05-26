package main

import (
	"bufio"
	"crypto/md5" // nolint: gosec // Нет потребности в криптобезопастности хэша
	"encoding/binary"
	"log"
	"os"
	"smallurl/test/hash_test/simple"
	"smallurl/test/hash_test/sponge"

	"golang.org/x/exp/rand"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	countGeneratedNumber        = 100000
	countBytesInUint64          = 8
	countGeneratedBase63Numbers = 4
)

func histPlot(values plotter.Values, resFileName string) error {
	p := plot.New()
	p.Title.Text = "histogram plot"

	hist, err := plotter.NewHist(values, 2<<14)
	if err != nil {
		return err
	}

	p.Add(hist)

	return p.Save(3*vg.Inch, 3*vg.Inch, "./test/hash_test/images/"+resFileName+".png")
}

// nolint: revive // Функция для тестирования поэтому размер не так важен
func main() {
	fl, err := os.OpenFile("./test/hash_test/1000-most-common-words.txt", os.O_RDONLY, 0o777)
	if err != nil {
		log.Fatal(err)
	}
	defer fl.Close()

	scanner := bufio.NewScanner(fl)

	var valuesWords plotter.Values

	for scanner.Scan() {
		// nolint: gosec // Нет потребности в криптобезопастности хэша
		hasher := simple.NewHash63(md5.New(), countGeneratedBase63Numbers)

		word := scanner.Text()
		_, _ = hasher.Write([]byte(word)) // nolint: errcheck // использует md5, который не возвращает ошибку
		hash := hasher.Sum(nil)

		uintHash := uint64(hash[0])*sponge.BaseOfResult*sponge.BaseOfResult*sponge.BaseOfResult +
			uint64(hash[1])*sponge.BaseOfResult*sponge.BaseOfResult +
			uint64(hash[2])*sponge.BaseOfResult +
			uint64(hash[3])

		valuesWords = append(valuesWords, float64(uintHash))
	}

	var valuesNumber plotter.Values

	for i := 0; i < countGeneratedNumber; i++ {
		// nolint: gosec // Нет потребности в криптобезопастности хэша
		hasher := simple.NewHash63(md5.New(), countGeneratedBase63Numbers)

		buf := make([]byte, countBytesInUint64)
		binary.LittleEndian.PutUint64(buf, uint64(i))

		_, _ = hasher.Write(buf) // nolint: errcheck // использует md5, который не возвращает ошибку
		hash := hasher.Sum(nil)

		uintHash := uint64(hash[0])*sponge.BaseOfResult*sponge.BaseOfResult*sponge.BaseOfResult +
			uint64(hash[1])*sponge.BaseOfResult*sponge.BaseOfResult +
			uint64(hash[2])*sponge.BaseOfResult +
			uint64(hash[3])

		valuesNumber = append(valuesNumber, float64(uintHash))
	}

	var valuesRandom plotter.Values

	for i := 0; i < countGeneratedNumber; i++ {
		// nolint: gosec // Нет потребности в криптобезопастности хэша
		hasher := simple.NewHash63(md5.New(), countGeneratedBase63Numbers)

		buf := make([]byte, countBytesInUint64)
		binary.LittleEndian.PutUint64(buf, uint64(rand.Int63()))

		_, _ = hasher.Write(buf) // nolint: errcheck // использует md5, который не возвращает ошибку
		hash := hasher.Sum(nil)

		uintHash := uint64(hash[0])*sponge.BaseOfResult*sponge.BaseOfResult*sponge.BaseOfResult +
			uint64(hash[1])*sponge.BaseOfResult*sponge.BaseOfResult +
			uint64(hash[2])*sponge.BaseOfResult +
			uint64(hash[3])

		valuesRandom = append(valuesRandom, float64(uintHash))
	}

	if err := histPlot(valuesWords, "words"); err != nil {
		log.Println(err)

		return
	}

	if err := histPlot(valuesNumber, "number"); err != nil {
		log.Println(err)

		return
	}

	if err := histPlot(valuesRandom, "random"); err != nil {
		log.Println(err)
	}
}
