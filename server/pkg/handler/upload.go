package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	postgresql "github.com/trannguyen33398/google-search-data-scraping/pkg/database"
	"github.com/trannguyen33398/google-search-data-scraping/pkg/model"
	scrapingService "github.com/trannguyen33398/google-search-data-scraping/pkg/service"
	"golang.org/x/sync/errgroup"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Upload struct {
	l *log.Logger
}

func NewUpload(l *log.Logger) *Upload {
	return &Upload{l}
}

type ConnectionHandlerInterface interface {
	SendMessageToClients(message []byte)
}

var connectionHandler ConnectionHandlerInterface

func SetConnectionHandler(handler ConnectionHandlerInterface) {
	connectionHandler = handler
}

func (p *Upload) UploadFile(w http.ResponseWriter, r *http.Request) {
	
	sqlDDL := "insert into scraping_items(file_id,user_id,key_word,total_advertised,total_link,total_search,html) values"
	// Parse the multipart form from the request
	err := r.ParseMultipartForm(10 << 20) // 10MB is the maximum form size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// ReadAll reads all the records from the CSV file
	// and Returns them as slice of slices of string
	// and an error if any
	records, err := reader.ReadAll()

	// Checks for the error
	if err != nil {
		fmt.Println("Error reading records")
	}

	// Loop to iterate through
	// and print each of the string slice
	g, ctx := errgroup.WithContext(context.Background())
	channelOutput := make(chan scrapingService.TScrapingData,1)
	db, _ := postgresql.InitConnection()

	uuid := uuid.New()
	if err != nil {
		log.Fatal(err)
	}

	//totalRecords := len(records)
	g.Go(func() error {
		defer close(channelOutput)

		for _, row := range records {
			data1 := scrapingService.ScrapingData(row[0])
			channelOutput <- data1

		}
		return nil
	})

	g.Go(func() error {
	
		for data := range channelOutput {
			p.l.Println(data.Keyword)
			sqlDDL += fmt.Sprintf("('%s','%s','%s', %d, %d,'%s','%s'), \n",
				uuid,
				"abc",
				data.Keyword,
				data.TotalAdvertised,
				data.TotalLink,
				data.TotalSearch,
				strings.Replace(data.Html, `'`, `"`, -1),
			)

		
			if err != nil {
				panic(err)
			}
		}
		select {
		default:
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Panic(err)
	}
	
	_, err = db.Exec(sqlDDL[:len(sqlDDL)-3])
	if err != nil {
		panic(err)
	}
	p.l.Println(connectionHandler)
	connectionHandler.SendMessageToClients([]byte(uuid.String()))
	return
}

func (p *Upload) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Upload) GetHistoryUpload(w http.ResponseWriter, r *http.Request){
	db, _ := postgresql.InitConnection()
	rows, queryError := db.Query("select  total_advertised,total_link,total_search,html from scraping_items")

	if queryError != nil {
		p.l.Println("Error: ",queryError.Error())
	}

	defer rows.Close()
	var scrapings [] model.TScraping
	if rows == nil {
		return 
	}
	
	for rows.Next() {
		scraping := model.TScraping{}
		
		err := rows.Scan(
			&scraping.TotalAdvertised,
			&scraping.TotalLink,
			&scraping.TotalSearch,
			&scraping.Html,
		)
		if err == nil {
			scrapings = append(scrapings, scraping)
		} else {
			fmt.Printf("[ERROR] error parse data %v\n", err)
		}
	}
	// serialize the list to JSON
	e := json.NewEncoder(w)
	e.Encode(scrapings)
	
}