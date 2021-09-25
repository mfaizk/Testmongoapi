package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfaizk/mongoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString="mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000"
const dbName="netflix"
const colName="watchlist"


//MOST IMPORTANT 

var collection *mongo.Collection

//connect with mongodb

func init(){
	//client option

	clientOption:=options.Client().ApplyURI(connectionString)

	//connect to monngodb
	client,err:= mongo.Connect(context.TODO(),clientOption)

	if err!=nil {
		log.Fatal(err)
	}

	fmt.Println("MonDB Connection success")

	collection=client.Database(dbName).Collection(colName)

	//collection instance 

	fmt.Println("Collection instance is ready")


}

// MongoDB helper -file

//insert 1 record

func insertOneMovie(movie model.Netflix){
     
   inserted,err:=collection.InsertOne(context.Background(),movie)

   if err!=nil {
	   log.Fatal(err)
   }
   fmt.Println("Inserted one movie in db with id",inserted.InsertedID)
}

// update 1 record

func updateOneMovie(movieId string){

   id,_:=primitive.ObjectIDFromHex(movieId)
   
   filter:=bson.M{"_id":id}
   update:=bson.M{"$set":bson.M{"watched":true}}

  result,err:= collection.UpdateOne(context.Background(),filter,update)

  if err!=nil {
	  log.Fatal(err)
  }

  fmt.Println("Modified count: ",result.ModifiedCount)

}

//delete 1 record

func deleteOneMovie(movieId string){
	id,_:= primitive.ObjectIDFromHex(movieId)
    filter:=bson.M{"_id":id}

  deleteCount,err:=	collection.DeleteOne(context.Background(),filter)
  if err!=nil {
	  log.Fatal(err)
  }
  fmt.Println("Movie got deleted with delete count: ",deleteCount)

}

//delete all record from mongo db

func deleteAllMovie()int64{

  deleteResult,err:=	collection.DeleteMany(context.Background(),bson.D{{}},nil)

  if err!=nil {
	  log.Fatal(err)
  }

  fmt.Println("No of movies deleted",deleteResult.DeletedCount)

    return deleteResult.DeletedCount
}

//get all movies from database

func getAllMovies() []primitive.M{
	cur,err:= collection.Find(context.Background(),bson.D{{}})


	if err!=nil {
		log.Fatal(err)

	}

	var movies []primitive.M

	for cur.Next(context.Background()){
		var movie bson.M
		err:=cur.Decode(&movie)
		if err!=nil {
			log.Fatal(err)
		}
		movies=append(movies, movie)	

	}
	defer cur.Close(context.Background())
	return movies
}

//Actual controller -file

func GetMyAllMovies(w http.ResponseWriter,r *http.Request){
 
	w.Header().Set("Content-type","application/x-www-form-urlencode")

	allMovies:=getAllMovies()

	json.NewEncoder(w).Encode(allMovies)


}


func CreateMovie(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-type","application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Alllow-Methods","POST")

	var movie model.Netflix

	_=json.NewDecoder(r.Body).Decode(&movie)
  
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)


}

func MarkAsWatched(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-type","application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Alllow-Methods","PUT")

	params:=mux.Vars(r)

	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}


func DeleteAMovie(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-type","application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Alllow-Methods","DELETE")

	params:=mux.Vars(r)

	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAllMovies(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-type","application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Alllow-Methods","DELETE")
   count:= deleteAllMovie()
	json.NewEncoder(w).Encode(count)
	
}