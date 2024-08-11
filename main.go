package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	resty "github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)


const prefix string = "!bidenbot"


func getGiphyGif(apiKey string, query string) (string, error) {
	client := resty.New()
	var randNum = rand.IntN(100)
	var result map[string]interface{}
	_, err := client.R().SetQueryParams(map[string]string{
		"api_key": apiKey,
		"q": `Joe Biden ${query}`,
		"limit:": "100",
		"offset": strconv.Itoa(randNum),
	}).SetResult(&result).Get("https://api.giphy.com/v1/gifs/search")

	if err != nil {
		return "", err
	}	
	
	

	if data, ok := result["data"].([]interface{}); ok && len(data) > 0 {
        if gif, ok := data[0].(map[string]interface{}); ok {
            if images, ok := gif["images"].(map[string]interface{}); ok {
                if original, ok := images["original"].(map[string]interface{}); ok {
                    if url, ok := original["url"].(string); ok {
                        return url, nil
                    }
                }
            }
        }
    }

	return "", nil
}


func main(){
	godotenv.Load(".env.dev")

	token := os.Getenv("BOT_TOKEN")
	gifApiKey := os.Getenv("GIPHY_API")
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate){
		// If the message in the channel was sent by the biden bot then return
		if m.Author.ID == s.State.User.ID {
			return
		}
		

		//server logic
		args := strings.Split(m.Content, " ")

		if args[0] != prefix {
			return
		}	

		if args[1] == "gif" {		
			query := strings.Join(args[2:], " ")
			gifURL, gifErr := getGiphyGif(gifApiKey, query)
			if gifErr != nil {
				return
			}			
		
	
		s.ChannelMessageSend(m.ChannelID, gifURL)
	}


	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println("Biden is in office")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
