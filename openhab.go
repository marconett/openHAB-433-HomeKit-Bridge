package main

import "encoding/json"
import "errors"
import "io/ioutil"
import "net/http"
import "github.com/Jeffail/gabs"
import "github.com/boltdb/bolt"
import "log"
import "regexp"
import "strings"

var (
	defaultTherms = []string{"marco", "henrik", "buero"}
)

func querySitemap(host string, port string, sitemap string) (value *gabs.Container, err error) {
	resp, err := http.Get("http://" + host + ":" + port + "/rest/sitemaps/" + sitemap + "?type=json")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	jsonParsed, err := gabs.ParseJSON(body)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return jsonParsed.Path("homepage"), nil
}

func getThermostats(value *gabs.Container) {

	// open bolt db
	db, err := bolt.Open("hk.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	children, _ := value.S("widget").Children()
	for _, child := range children {
		name, okName := child.Path("item.name").Data().(string)
		link, okLink := child.Path("item.link").Data().(string)
		state, okState := child.Path("item.state").Data().(string)

		if okName && okLink && okState {
			matchHum := regexp.MustCompile("(?i)currenthum")
			matchSetTemp := regexp.MustCompile("(?i)settemp")
			matchTemp := regexp.MustCompile("(?i)currenttemp")
			matchMode := regexp.MustCompile("(?i)setmode")

			switch {
			case matchHum.MatchString(name):
				db.Update(func(tx *bolt.Tx) error {
					b, err := tx.CreateBucketIfNotExists([]byte(strings.Split(name, "_")[0] + "Thermostat"))
					if err != nil {
						return err
					}
					encoded, err := json.Marshal(OpenHabItem{name, link, state})
					if err != nil {
						return err
					}

					return b.Put([]byte("hum"), encoded)
				})
			case matchSetTemp.MatchString(name):
				db.Update(func(tx *bolt.Tx) error {
					b, err := tx.CreateBucketIfNotExists([]byte(strings.Split(name, "_")[0] + "Thermostat"))
					if err != nil {
						return err
					}
					encoded, err := json.Marshal(OpenHabItem{name, link, state})
					if err != nil {
						return err
					}
					return b.Put([]byte("settemp"), encoded)
				})
			case matchTemp.MatchString(name):
				db.Update(func(tx *bolt.Tx) error {
					b, err := tx.CreateBucketIfNotExists([]byte(strings.Split(name, "_")[0] + "Thermostat"))
					if err != nil {
						return err
					}
					encoded, err := json.Marshal(OpenHabItem{name, link, state})
					if err != nil {
						return err
					}
					return b.Put([]byte("temp"), encoded)
				})
			case matchMode.MatchString(name):
				db.Update(func(tx *bolt.Tx) error {
					b, err := tx.CreateBucketIfNotExists([]byte(strings.Split(name, "_")[0] + "Thermostat"))
					if err != nil {
						return err
					}
					encoded, err := json.Marshal(OpenHabItem{name, link, state})
					if err != nil {
						return err
					}
					return b.Put([]byte("mode"), encoded)
				})
			default:
				err = errors.New("Couldn't match string " + name)
				log.Fatal(err)
			}
		}
	}
}

func setValue(host string, port string, itemName string, value string) (status string, err error) {
	resp, err := http.Get("http://" + host + ":" + port + "/CMD?" + itemName + "=" + value)
	if err != nil {
		log.Fatal(err)
		return resp.Status, err
	}

	defer resp.Body.Close()

	return resp.Status, nil
}
