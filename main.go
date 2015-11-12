package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/benmanns/goworker"
	"github.com/dukex/mixpanel"
)

var mp *mixpanel.Mixpanel

func mixpanelTrack(queue string, args ...interface{}) error {
	if len(args) < 3 {
		ex := fmt.Sprintf("Not enough arguments for %s: %d present", queue, len(args))
		log.Print(ex)
		return errors.New(ex)
	}

	distinctID, ok := args[0].(string)
	if !ok {
		ex := fmt.Sprintf("Can't convert distinctID %s to string", args[0])
		log.Print(ex)
		return errors.New(ex)
	}

	eventName, ok := args[1].(string)
	if !ok {
		ex := fmt.Sprintf("Can't convert event name %v to string", args[1])
		log.Print(ex)
		return errors.New(ex)
	}

	props, ok := args[2].(map[string]interface{})
	if !ok {
		ex := fmt.Sprintf("Can't convert properties %v to mixpanel.Properties", args[2])
		log.Print(ex)
		return errors.New(ex)
	}

	res, err := mp.Track(distinctID, eventName, props)

	if err != nil {
		log.Print(err)
		return err
	}

	log.Printf("Tracked event %s for %s with %s, response code %d", eventName, distinctID, props, res.StatusCode)

	return nil
}

func init() {
	mixpanelToken := os.Getenv("MIXPANEL_TOKEN")

	if mixpanelToken == "" {
		log.Fatal("MIXPANEL_TOKEN must be defined")
	}

	mp = mixpanel.NewMixpanel(mixpanelToken)

	goworker.Register("MixpanelTrackJob", mixpanelTrack)
}

func main() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}

// func (m *Mixpanel) Track(distinctID string, event string, props Properties) error {
