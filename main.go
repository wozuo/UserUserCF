package main

import (
	"fmt"
	"os"
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"math"
)

type User struct {
	UserID uint32
	Ratings []float32
	SimToUsers []float32
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

// Pearson correlation to compute similarity between users
func pearsonCorrelation(users []User) []User {
	for i := 0; i < len(users); i++ {
		var avg_A, avg_B float32
		var notNull uint32
		for j := 0; j < len(users[i].Ratings); j++ {
			avg_A += users[i].Ratings[j]
			if users[i].Ratings[j] != 0 {
				notNull += 1
			}
		}
		avg_A = avg_A / float32(notNull)
		for j := 0; j < len(users); j++ {
			notNull = 0
			for k := 0; k < len(users[j].Ratings); k++ {
				avg_B += users[j].Ratings[k]
				if users[j].Ratings[k] != 0 {
					notNull += 1
				}
			}
			avg_B = avg_B / float32(notNull)
			var p1, p2, p3 float32
			for k := 0; k < len(users[j].Ratings); k++ {
				if users[j].Ratings[k] != 0 && users[i].Ratings[k] != 0 {
					p1 += (users[i].Ratings[k] - avg_A) * (users[j].Ratings[k] - avg_B)
					p2 += (users[i].Ratings[k] - avg_A) * (users[i].Ratings[k] - avg_A)
					p3 += (users[j].Ratings[k] - avg_B) * (users[j].Ratings[k] - avg_B)
				}
			}
			users[i].SimToUsers = append(users[i].SimToUsers, float32((float64(p1) / (math.Sqrt(float64(p2)) * math.Sqrt(float64(p3))))))
		}
	}
	return users
}

func main() {
	//users, movies := loadCSV()
	//fmt.Println("Done! %v %v", users, movies)
	users, _ := loadCSV()
	users = pearsonCorrelation(users)
	fmt.Println("Similarity values of user 0: ", users[0].SimToUsers)
	fmt.Println("Done!")
}