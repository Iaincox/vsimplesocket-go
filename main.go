package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
)

var Tpl *template.Template
/* *********************************************************************************************************************
	Init
********************************************************************************************************************* */
func init() {
	/**********************
		Validate Templates
	********************* */
	cwd, _ := os.Getwd()
	cwd = cwd + "/templates/*.gohtml"
	log.Printf("Current Directory =%s\n", cwd)

	Tpl = template.Must(template.ParseGlob(cwd))

	/*****************************
		Pickup Environment data
	**************************** */
	/*e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}
	WebPort = WebPort + os.Getenv("web_port")*/
	/*
		e := godotenv.Load() //Load .env file
		if e != nil {
			fmt.Print(e)
		}

		username := os.Getenv("db_user")
		password := os.Getenv("db_pass")
		dbName := os.Getenv("db_name")
		dbHost := os.Getenv("db_host")


		dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
		fmt.Println(dbUri)

		conn, err := gorm.Open("postgres", dbUri)
		if err != nil {
			fmt.Print(err)
		}

		db = conn
		db.Debug().AutoMigrate(&Account{}, &Contact{}) //Database migration
	*/
}


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func reader (conn *websocket.Conn)  {
	for{
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil{
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, []byte("Tula Puppy likes Dinner and Walks")); err != nil{
			log.Println(err)
			return
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("Home Page")

	if err := Tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("Page Rendered")
}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request)  {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade( w,r, nil)
	if err != nil{
		log.Println(err)
	}
	log.Println("Client Successfully Conected...")
	reader(ws)
}

func setupRoutes()  {
	http.HandleFunc( "/", homePage)
	http.HandleFunc("/ws", wsEndpoint )
}

func main()  {
	fmt.Println("Go WebSockets")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil ))
}

