package models

import (
	"encoding/json"
	"testing"
)

//TODO: check if struct is empty and throw an error if values are nil

func TestBuildUserStructModel(t *testing.T) {
	// want := &SpotifyUser{ID: "fakeid1234567", Name: "Sergio", Email: "sergiofakeemail@student.fakeschool.edu", Country: "US",
	// 	UserURL: "https://api.spotify.com/v1/usersfakeid1234567",
	// 	Images:  [{"https://i.scdn.co/image/abcd123fakeuri"}],
	// }

	d := []byte(`{"country" : "US","display_name" : "Sergio","email" : "sergiofakeemail@student.fakeschool.edu","explicit_content" : {"filter_enabled" : false,"filter_locked" : false},"external_urls" : {"spotify" : "https://open.spotify.com/user/fakeid1234567"},"followers" : {"href" : null,"total" : 2},"href" : "https://api.spotify.com/v1/usersfakeid1234567","id" : "fakeid1234567","images" : [ {"height" : null,"url" : "https://i.scdn.co/image/abcd123fakeuri","width" : null} ],"product" : "premium","type" : "user","uri" : "spotify:user:fakeid1234567"}`)
	user := &SpotifyUser{}
	err := json.Unmarshal(d, &user)
	if err != nil {
		t.Errorf("Failed to unmarshall data into struct : %v", err)
	}
	t.Logf("User struct: %+v", user)
	t.Logf("type Imaage, %T", user.Images[0])

}
