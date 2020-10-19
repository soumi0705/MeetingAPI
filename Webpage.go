// +build ignore

package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"time"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)

type Page struct {
	Title string
	Body  []byte
}

type Meeting struct{
	ID              string		"json: id"
	Title			string		"json: title"
	Participants	string		"json: participants"
	StartTime 		time.Time		"json: starttime"
	EndTime			time.Time		"json: endtime"
	CreationTime 	time.Time  	"json: creationtime"
}
type Participant struct{
	Nme             string		"json: id"
	Email			string		"json: title"
	RSVP			string		"json: participants"
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func write(id string, t string , p string, st time.Time, et time.Time, ct time.Time) {
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := "mongodb://127.0.0.1:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	log.Printf("Successfully connected and pinged.")
	officeDatabase:= client.Database("office")
	meetingCollection := officeDatabase.Collection("meeting")
	meetResult, err := meetingCollection.InsertOne(ctx, bson.D{
		{Key: "id" , Value: id},
		{Key: "title" , Value:t},
		{Key: "participants" , Value: p},
		{Key: "starttime" , Value:st},
		{Key: "endtime" , Value: et},
		{Key: "creationtime" , Value: ct},
	})
	if err!= nil{
		log.Fatal(err)
	}
	fmt.Println(meetResult.InsertedID)

}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("./templates/"+tmpl + ".html")
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	renderTemplate(w, "view", p)
}

func meetCreate(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/meetings" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
	}
	switch r.Method{
	case "GET":
		
		query := r.URL.Query()
		start, present := query["start"] 
		end, present2 := query["end"]
		participant, present3 :=query["participant"]
		uri := "mongodb://127.0.0.1:27017"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}
		// Ping the primary
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			panic(err)
		}
		log.Printf("Successfully connected and pinged.")
		officeDatabase:= client.Database("office")
		meetingCollection := officeDatabase.Collection("meeting")
		if !present || len(start) == 0 {
			fmt.Println("start not present")
		}
		if !present2 || len(end) == 0 {
			fmt.Println("end not present")
		}
		if !present3 || len(participant) ==0{
			fmt.Println("Participant query not present")
		}
		if present3 || len(participant)!=0{
			
			pt := strings.Join(participant,"")
			fmt.Println(pt)
			meetResult, err := meetingCollection.Find(ctx, bson.D{{"participants",primitive.Regex{Pattern: pt , Options: ""}}})
			if err!= nil{
				fmt.Fprint(w, err)
			}
			var meet []bson.M
			if err = meetResult.All(ctx, &meet); err != nil {
				log.Fatal(err)
			}
			
			var meetResult1 Meeting
			for i:=0 ; i<len(meet); i++{
				bsonBytes, _ := bson.Marshal(meet[i])
				bson.Unmarshal(bsonBytes, &meetResult1)
				fmt.Println(meetResult1)
				json.NewEncoder(w).Encode(meetResult1)
			}
		}
		//db.order.find({"OrderDateTime":{ $gte:ISODate("2019-02-10"), $lt:ISODate("2019-02-21") }})
		if (present && present2) || (len(start) != 0 && len(end) != 0){
			layout := "2006-01-02T15:04:05.000Z"
			
			st, err:= time.Parse(layout , strings.Join(start,""))
			en, err1 := time.Parse(layout , strings.Join(end,""))
			if err!= nil{
				fmt.Fprint(w, err)
			}
			if err1!= nil{
				fmt.Fprint(w, err)
			}
			meetResult, err := meetingCollection.Find(ctx, bson.D{{"starttime", bson.D{{"$gt", st},{"$lt",en }}}})
			if err!= nil{
				fmt.Fprint(w, err)
			}
			var meet []bson.M
			if err = meetResult.All(ctx, &meet); err != nil {
				log.Fatal(err)
			}
			
			var meetResult1 Meeting
			for i:=0 ; i<len(meet); i++{
				bsonBytes, _ := bson.Marshal(meet[i])
				bson.Unmarshal(bsonBytes, &meetResult1)
				fmt.Println(meetResult1)
				json.NewEncoder(w).Encode(meetResult1)
			}
		}
		
		
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		var meet Meeting
		var StartTime time.Time
		var EndTime time.Time
		layout := "2006-01-02T15:04:05.000Z"
		meet.ID = r.FormValue("id")
		meet.Title = r.FormValue("title")
		meet.Participants =r.FormValue("participants")
		StartTime, errq:= time.Parse(layout, r.FormValue("starttime"))
		if errq != nil {
			fmt.Println(errq)
		}
		EndTime, err:= time.Parse(layout, r.FormValue("endtime"))
		if err != nil {
			fmt.Println(err)
		}
		meet.StartTime=StartTime
		meet.EndTime=EndTime
		meet.CreationTime = time.Now()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(meet)
		write(meet.ID, meet.Title, meet.Participants, meet.StartTime, meet.EndTime, meet.CreationTime)    
	default:
		fmt.Println("Only Get and Post Supported")
	}
}
func meetView(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/meetings/"):]
	uri := "mongodb://127.0.0.1:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	log.Printf("Successfully connected and pinged.")
	officeDatabase:= client.Database("office")
	meetingCollection := officeDatabase.Collection("meeting")
	var meet bson.M
	err1:= meetingCollection.FindOne(ctx, bson.M{"id": title}).Decode(&meet)
	if err1!= nil{
		fmt.Fprint(w, err1)
	}
	var meetResult Meeting
	bsonBytes, _ := bson.Marshal(meet)
	bson.Unmarshal(bsonBytes, &meetResult)
	json.NewEncoder(w).Encode(meetResult)
}


func main() {
	fmt.Println("http://localhost:8080/view/test")
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/meetings", meetCreate)
	http.HandleFunc("/meetings/", meetView)

	//http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8090", nil))
}
