package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	TrainNo   string `bson:"trainNo   string"`
	TrainName string `bson:"trainName string"`
	SEQ       string `bson:"seq       string"`
	Code      string `bson:"code      string"`
	StName    string `bson:"stName    string"`
	ATime     string `bson:"atime     string"`
	DTime     string `bson:"dtime     string"`
	Distance  string `bson:"distance  string"`
	SS        string `bson:"ss        string"`
	SSname    string `bson:"ssname    string"`
	Ds        string `bson:"ds        string"`
	DsName    string `bson:"dsName    string"`
}

func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func dbConn() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("test").Collection("railway1")
	fmt.Println("Connected to MongoDB!")

	return collection

}
func getallTrains(w http.ResponseWriter, r *http.Request) {

	collection := dbConn()
	cursor, err := collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}
	var trains []Data
	if err = cursor.All(context.TODO(), &trains); err != nil {
		log.Fatal(err)
	}
	bytedata, err := json.MarshalIndent(trains, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytedata)
}
func main() {
	collection := dbConn()

	csvData, err := ReadCsv("Indian_railway1.csv")
	if err != nil {
		panic(err)
	}

	csvData = csvData[1:501]

	for _, line := range csvData {
		data := Data{
			TrainNo:   line[0],
			TrainName: line[1],
			SEQ:       line[2],
			Code:      line[3],
			StName:    line[4],
			ATime:     line[5],
			DTime:     line[6],
			Distance:  line[7],
			SS:        line[8],
			SSname:    line[9],
			Ds:        line[10],
			DsName:    line[11],
		}

		_, err := collection.InsertOne(context.TODO(), data)
		if err != nil {
			panic(err)
		}
		//break

	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/Trains", getallTrains)

	http.ListenAndServe(":8080", nil)
}
