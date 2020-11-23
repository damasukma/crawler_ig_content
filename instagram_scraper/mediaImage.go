package instagram_scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Media struct{
	Graphql struct {
		User struct {
			EdgeOwnerToTimelineMedia struct {
				Edges []struct {
					Node struct {
						DisplayUrl string `json:"display_url"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
		} `json:"user"`
	} `json:"graphql"`
}

const account = "https://www.instagram.com/%s/?__a=1"

func FetchMediaImage(username string, limit int) (*[]string, error){

	var(
		data []string
	)

	//
	resp, err := http.Get(fmt.Sprintf(account,username))

	if err != nil {
		return &[]string{}, err
	}

	defer resp.Body.Close()
	res := &Media{}
	json.NewDecoder(resp.Body).Decode(&res)

	list := res.Graphql.User.EdgeOwnerToTimelineMedia.Edges
	if len(list) > 0 {
		item := list[:limit]

		for _, value := range item{
			data = append(data, value.Node.DisplayUrl)
		}

		return &data,nil
	}

	return &data,nil

}

