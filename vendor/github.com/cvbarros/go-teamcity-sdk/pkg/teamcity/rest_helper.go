package teamcity

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

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

func (r *restHelper) getCustom(path string, out interface{}, resourceDescription string, reader responseReadFunc) error {
	request, _ := r.sling.New().Get(path).Request()
	response, err := r.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		err = reader(bodyBytes, out)
		if err != nil {
			return err
		}
		return nil
	}

	return r.handleRestError(response, "GET", resourceDescription)
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

func (r *restHelper) putTextPlain(path string, data string, resourceDescription string) (string, error) {
	req, err := r.sling.New().Put(path).
		BodyProvider(textPlainBodyProvider{payload: data}).
		Add("Accept", "text/plain").
		Request()

	if err != nil {
		return "", err
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		return string(bodyBytes), nil
	}

	return "", r.handleRestError(resp, "PUT", resourceDescription)
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

func (r *restHelper) put(path string, data interface{}, out interface{}, resourceDescription string) error {
	request, _ := r.sling.New().Put(path).BodyJSON(data).Request()
	response, err := r.httpClient.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 201 || response.StatusCode == 200 {
		json.NewDecoder(response.Body).Decode(out)
		return nil
	}

	return r.handleRestError(response, "PUT", resourceDescription)
}

func (r *restHelper) delete(path string, resourceDescription string) error {
	return r.deleteByIDWithSling(r.sling, path, resourceDescription)
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

func replaceValue(i, v interface{}) {
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr {
		panic("not a pointer")
	}
	val = val.Elem()

	newVal := reflect.Indirect(reflect.ValueOf(v))
	if !val.Type().AssignableTo(newVal.Type()) {
		panic("mismatched types")
	}

	val.Set(newVal)
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

type textPlainBodyProvider struct {
	payload interface{}
}

func (p textPlainBodyProvider) ContentType() string {
	return "text/plain; charset=utf-8"
}

func (p textPlainBodyProvider) Body() (io.Reader, error) {
	return strings.NewReader(p.payload.(string)), nil
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Printf("[ERROR] %s\n\n", err)
	}
}
