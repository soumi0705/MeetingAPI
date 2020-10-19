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
