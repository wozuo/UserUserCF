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
	AverageRating float64
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
		users[i].AverageRating = avg_A
		for j := 0; j < len(users); j++ {
			notNull = 0
			avg_B = 0
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
					p1 += ((users[i].Ratings[k] - avg_A) * (users[j].Ratings[k] - avg_B))
					p2 += ((users[i].Ratings[k] - avg_A) * (users[i].Ratings[k] - avg_A))
					p3 += ((users[j].Ratings[k] - avg_B) * (users[j].Ratings[k] - avg_B))
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

// n: number of predictions we want to get
func notNormalizedPrediction(ownIndex int, n int, topNeighbors []KeyValue, users []User, movies []string) []string {
	var topMovies []KeyValue
	for i := 0; i < len(movies); i++ {
		var p1, p2 float64
		for j := 0; j < len(topNeighbors); j++ {
			if users[topNeighbors[j].Key].Ratings[i] != 0 {
				p1 += (users[topNeighbors[j].Key].Ratings[i] * users[ownIndex].SimToUsers[topNeighbors[j].Key])
				p2 += users[ownIndex].SimToUsers[topNeighbors[j].Key]
			}
		}
		var prediction = p1 / p2
		topMovies = append(topMovies, KeyValue{i, prediction})
	}
	sort.Slice(topMovies, func(i, j int) bool {
		return topMovies[i].Value > topMovies[j].Value
	})
	var topMovieNames []string
	for i := 0; i < n; i++ {
		topMovieNames = append(topMovieNames, movies[topMovies[i].Key])
	}
	return topMovieNames
}

// n: number of predictions we want to get
func normalizedPrediction(ownIndex int, n int, topNeighbors []KeyValue, users []User, movies []string) []string {
	var topMovies []KeyValue
	for i := 0; i < len(movies); i++ {
		var p1, p2 float64
		for j := 0; j < len(topNeighbors); j++ {
			if users[topNeighbors[j].Key].Ratings[i] != 0 {
				p1 += ((users[topNeighbors[j].Key].Ratings[i] - users[topNeighbors[j].Key].AverageRating) * users[ownIndex].SimToUsers[topNeighbors[j].Key])
				p2 += users[ownIndex].SimToUsers[topNeighbors[j].Key]
			}
		}
		var prediction = users[ownIndex].AverageRating + (p1 / p2)
		topMovies = append(topMovies, KeyValue{i, prediction})
	}
	sort.Slice(topMovies, func(i, j int) bool {
		return topMovies[i].Value > topMovies[j].Value
	})
	var topMovieNames []string
	for i := 0; i < n; i++ {
		topMovieNames = append(topMovieNames, movies[topMovies[i].Key])
	}
	return topMovieNames
}

func main() {
	users, movies := loadCSV()
	users = pearsonCorrelation(users)
	fmt.Println("Similarity values of user 4: ", users[4].SimToUsers)
	topNeighbors := getTopSimNeighbors(4, 5, users[4].SimToUsers)
	fmt.Println("Top neighbors: ", topNeighbors)
	topMovieNames := notNormalizedPrediction(4, 6, topNeighbors, users, movies)
	fmt.Println("(Not normalized) Top movies for user 4: ", topMovieNames)
	topMovieNames = normalizedPrediction(4, 6, topNeighbors, users, movies)
	fmt.Println("(Normalized) Top movies for user 4: ", topMovieNames)
}