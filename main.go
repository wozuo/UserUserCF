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
	"sort"
)

type User struct {
	UserID int
	Ratings []float64
	SimToUsers []float64
}

type KeyValue struct {
	Key int
	Value float64
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
		var ratings []float64
		for i := 1; i < len(line); i++ {
			floatRating, _ := strconv.ParseFloat(line[i], 64)
			ratings = append(ratings, floatRating)
		}
		intID, _ := strconv.ParseInt(line[0], 0, 64)
		users = append(users, User{
			UserID: int(intID),
			Ratings: ratings,
		})
	}
	return users, movies
}

// Pearson correlation to compute similarity between users
func pearsonCorrelation(users []User) []User {
	for i := 0; i < len(users); i++ {
		var avg_A, avg_B float64
		var notNull int
		for j := 0; j < len(users[i].Ratings); j++ {
			avg_A += users[i].Ratings[j]
			if users[i].Ratings[j] != 0 {
				notNull += 1
			}
		}
		avg_A = avg_A / float64(notNull)
		for j := 0; j < len(users); j++ {
			notNull = 0
			for k := 0; k < len(users[j].Ratings); k++ {
				avg_B += users[j].Ratings[k]
				if users[j].Ratings[k] != 0 {
					notNull += 1
				}
			}
			avg_B = avg_B / float64(notNull)
			var p1, p2, p3 float64
			for k := 0; k < len(users[j].Ratings); k++ {
				if users[j].Ratings[k] != 0 && users[i].Ratings[k] != 0 {
					p1 += (users[i].Ratings[k] - avg_A) * (users[j].Ratings[k] - avg_B)
					p2 += (users[i].Ratings[k] - avg_A) * (users[i].Ratings[k] - avg_A)
					p3 += (users[j].Ratings[k] - avg_B) * (users[j].Ratings[k] - avg_B)
				}
			}
			users[i].SimToUsers = append(users[i].SimToUsers, (p1 / (math.Sqrt(p2) * math.Sqrt(p3))))
		}
	}
	return users
}

func getTopSimNeighbors(ownIndex int, n int, simUsers []float64) []KeyValue {
	var topNeighbors []KeyValue
	for i := 0; i < len(simUsers); i++ {
		if i != ownIndex {
			topNeighbors = append(topNeighbors, KeyValue{i, simUsers[i]})
		}
	}
	sort.Slice(topNeighbors, func(i, j int) bool {
		return topNeighbors[i].Value > topNeighbors[j].Value
	})
	return topNeighbors[:n]
}

func main() {
	//users, movies := loadCSV()
	//fmt.Println("Done! %v %v", users, movies)
	users, _ := loadCSV()
	users = pearsonCorrelation(users)
	fmt.Println("Similarity values of user 0: ", users[0].SimToUsers)
	topNeighbors := getTopSimNeighbors(0, 5, users[0].SimToUsers)
	fmt.Println("Top neighbors: ", topNeighbors)
	fmt.Println("Done!")
}