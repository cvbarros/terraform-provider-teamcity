package teamcity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

type responseReadFunc func([]byte, interface{}) error

type restHelper struct {
	httpClient *http.Client
	sling      *sling.Sling
}

func newRestHelper(httpClient *http.Client) *restHelper {
	return newRestHelperWithSling(httpClient, nil)
}

func newRestHelperWithSling(httpClient *http.Client, s *sling.Sling) *restHelper {
	return &restHelper{
		httpClient: httpClient,
		sling:      s,
	}
}

func (r *restHelper) get(path string, out interface{}, resourceDescription string) error {
	request, _ := r.sling.New().Get(path).Request()
	response, err := r.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode == 200 {
		json.NewDecoder(response.Body).Decode(out)
		return nil
	}

	return r.handleRestError(response, "GET", resourceDescription)
}

func (r *restHelper) postCustom(path string, data interface{}, out interface{}, resourceDescription string, reader responseReadFunc) error {
	request, _ := r.sling.New().Post(path).BodyJSON(data).Request()
	response, err := r.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode == 201 || response.StatusCode == 200 {
		err := reader(bodyBytes, out)
		if err != nil {
			return err
		}
		return nil
	}

	return r.handleRestError(response, "POST", resourceDescription)
}

func (r *restHelper) post(path string, data interface{}, out interface{}, resourceDescription string) error {
	request, _ := r.sling.New().Post(path).BodyJSON(data).Request()
	response, err := r.httpClient.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 201 || response.StatusCode == 200 {
		json.NewDecoder(response.Body).Decode(out)
		return nil
	}

	return r.handleRestError(response, "POST", resourceDescription)
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
		return r.handleRestError(response, "DELETE", resourceDescription)
	}

	return nil
}

func (r *restHelper) handleRestError(resp *http.Response, op string, res string) error {
	dt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("Error '%d' when performing '%s' operation - %s: %s", resp.StatusCode, op, res, string(dt))
}
