package spotifyapi

import "testing"

func TestBuildAPIEndpoints(t *testing.T) {
	api, err := NewAPI()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("api endpoints: %+v\n", api)

}
