package main

import (
	"fmt"
	"os"
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"strconv"
)

type User struct {
	UserID uint32
	Ratings []float32
}

func loadCSV() ([]User, []string) {
	csvFile, _ := os.Open("user-movie.csv")
	bufReader := csv.NewReader(bufio.NewReader(csvFile))
	var users []User
	var movies []string
	var lineNbr = 0
	for {
		line, err := bufReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if lineNbr == 0 {
			for i := 1; i < len(line); i++ {
				movies = append(movies, line[i])
			}
			lineNbr++
			continue
		}
		var ratings []float32
		for i := 1; i < len(line); i++ {
			floatRating, _ := strconv.ParseFloat(line[i], 32)
			ratings = append(ratings, float32(floatRating))
		}
		intID, _ := strconv.ParseInt(line[0], 0, 32)
		users = append(users, User{
			UserID: uint32(intID),
			Ratings: ratings,
		})
	}
	return users, movies
}

func main() {
	//users, movies := loadCSV()
	//fmt.Println("Done! %v %v", users, movies)
	fmt.Println("Done!")
}