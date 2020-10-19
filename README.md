# MeetingAPI
Packages Used:
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
  
  Unit Tests:
  Schedule a meeting
   1   Should be a POST request
   2   Use JSON request body
   3   URL should be ‘/meetings’
   4   Must return the meeting in JSON format
   
   
