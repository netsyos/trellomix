package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/draganshadow/trello"
	"github.com/gorilla/mux"
)

type Config struct {
	AppKey                 string `json:"appKey"`
	Token                  string `json:"token"`
	Member                 string `json:"member"`
	OverviewBoardShortLink string `json:"overviewBoardShortLink"`
	OnGoingColumnName      string `json:"onGoingColumnName"`
}

func readConfig() Config {
	configFile := os.Getenv("TRELLOMIX_CONFIG")
	if "" == configFile {
		configFile = "config.json"
	}
	fmt.Println("TRELLOMIX_CONFIG: ", configFile)
	// Open our jsonFile
	jsonFile, err := os.Open(configFile)

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(byteValue, &config)
	fmt.Printf("config : %+v\n", config)
	return config
}

func main() {
	config := readConfig()
	r := mux.NewRouter()
	r.HandleFunc("/status/{item}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		item := vars["item"]

		fmt.Fprintf(w, "You've requested the status of %s\n", item)
	})

	// http.ListenAndServe(":80", r)

	client := trello.NewClient(config.AppKey, config.Token)
	member, err := client.GetMember(config.Member, trello.Defaults())
	if err != nil {
		// Handle error
	}
	overviewBoard, err := client.GetBoard(config.OverviewBoardShortLink, trello.Defaults())
	if err != nil {
		// Handle error
	}
	fmt.Println("Overview Board", overviewBoard.Name)
	overviewBoardLists, err := overviewBoard.GetLists(trello.Defaults())
	if err != nil {
		// Handle error
	}
	// cards, err := overviewBoard.GetCards(trello.Defaults())
	// if err != nil {
	// 	// Handle error
	// }
	// for _, card := range cards {

	// 	fmt.Println("Overview Board Card", card.Name)
	// 	fmt.Printf("%+v\n", card)
	// 	err := card.Delete(trello.Defaults())
	// 	if err != nil {
	// 		// Handle error
	// 	}
	// 	break
	// }

	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {
		// Handle error
	}
	if boards != nil {
		// Handle error
	}
	for _, board := range boards {
		if !board.Closed {
			if board.Url != overviewBoard.Url {
				lists, err := board.GetLists(trello.Defaults())
				if err != nil {
					// Handle error
				}
				for _, list := range lists {
					fmt.Printf("L '%s'\n", list.Name)
					if list.Name == config.OnGoingColumnName {
						fmt.Println("Board", board.Name)
						// createBoardList := true
						for _, oblist := range overviewBoardLists {
							if oblist.Name == board.Name {
								// createBoardList = false
								fmt.Println("Board list exists lets sync it")

								obcards, err := oblist.GetCards(trello.Defaults())
								if err != nil {
									// Handle error
								}

								for _, card := range obcards {
									err := card.Delete(trello.Defaults())
									if err != nil {
										// Handle error
									}
								}

								cards, err := list.GetCards(trello.Defaults())
								if err != nil {
									// Handle error
								}

								for _, card := range cards {
									obCard, err := card.CopyToList(oblist.ID, trello.Defaults())
									if err != nil {
										// Handle error
									}
									fmt.Println(obCard.Name)
									// fmt.Printf("%+v\n", obCard)
								}

								break
							}
						}
						// if createBoardList {
						// }
					}
				}
				// fmt.Printf("%+v\n", board)
			}
		}
	}
}
