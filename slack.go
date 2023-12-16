package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func StartSlack(settings SlackSetting) (chan string, chan bool) {
	slackAPI := slack.New(settings.Token, slack.OptionAppLevelToken(settings.AppLevelToken))

	scm := socketmode.New(slackAPI)

	go func() {
		err := scm.Run()
		if err != nil {
			fmt.Println(err)
		}
	}()

	botinfo, _ := slackAPI.AuthTest()

	var output = make(chan string)
	var donechan = make(chan bool)

	var useridRegexp = regexp.MustCompile(`<@(\S+)>`)

	go func() {
		for ev := range scm.Events {
			switch ev.Type {
			case socketmode.EventTypeConnected:
				fmt.Printf("Start websocket connection with Slack\n")
			case socketmode.EventTypeEventsAPI:
				scm.Ack(*ev.Request)

				evp, ok := ev.Data.(slackevents.EventsAPIEvent)
				if !ok {
					continue
				}

				switch evp.Type {
				case slackevents.CallbackEvent:
					switch evi := evp.InnerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						text := strings.ReplaceAll(evi.Text, fmt.Sprintf("<@%s>", botinfo.BotID), "")

						matchstrings := useridRegexp.FindAllStringSubmatch(text, -1)

						for _, m := range matchstrings {
							info, err := slackAPI.GetUserInfo(m[1])
							if err != nil {
								fmt.Println("Failed to get user details: ", err)
								continue
							}

							text = strings.ReplaceAll(text, fmt.Sprintf("<@%s>", info.ID), info.Name)
						}

						output <- text
						<-donechan

						slackAPI.SendMessage(
							evi.Channel,
							slack.MsgOptionAsUser(false),
							slack.MsgOptionIconEmoji(settings.Icon),
							slack.MsgOptionText("Message was successfully sent.", false),
						)
					}
				}
			}
		}
	}()

	return output, donechan
}
