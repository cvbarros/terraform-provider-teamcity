package teamcity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

type restHelper struct {
	httpClient *http.Client
}

func newRestHelper(httpClient *http.Client) *restHelper {
	return &restHelper{httpClient: httpClient}
}

func (r *restHelper) postJSONWithSling(path string, sling *sling.Sling, data interface{}, out interface{}, resourceDescription string) error {
	request, _ := sling.New().Post(path).BodyJSON(data).Request()
	response, err := r.httpClient.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 201 || response.StatusCode == 200 {
		json.NewDecoder(response.Body).Decode(out)
		return nil
	}

	respData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("Error '%d' when posting %s: %s", response.StatusCode, resourceDescription, string(respData))
}

func (r *restHelper) deleteByIDWithSling(sling *sling.Sling, resourceID string, resourceDescription string) error {
	request, _ := sling.New().Delete(resourceID).Request()
	response, err := r.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode == 204 {
		return nil
	}

	if response.StatusCode != 200 && response.StatusCode != 204 {
		respData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error '%d' when deleting %s: %s", response.StatusCode, resourceDescription, string(respData))
	}

	return nil
}
