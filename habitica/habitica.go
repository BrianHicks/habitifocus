package habitica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HabiticaTODO struct {
	Alias     string `json:"alias"`
	Text      string `json:"text"`
	Type      string `json:"type"` // should always be "todo"
	Completed bool   `json:"completed,omitempty"`
}

func (ht *HabiticaTODO) String() string {
	return fmt.Sprintf("%s (done: %t)", ht.Text, ht.Completed)
}

type HabiticaTODOWrapper struct {
	Data []*HabiticaTODO `json:"data"`
}

type Client struct {
	UserID string
	APIKey string
}

func (c *Client) get(url string) (response []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-user", c.UserID)
	req.Header.Add("x-api-key", c.APIKey)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Error! Error code %s, Content: %s", resp.Status, string(content))
	}

	return content, err
}

func (c *Client) post(url string, data []byte) (response []byte, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-user", c.UserID)
	req.Header.Add("x-api-key", c.APIKey)
	req.Header.Add("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Error! Error code %s, Content: %s", resp.Status, string(content))
	}

	return content, err
}

func (c *Client) delete(url string) (response []byte, err error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-user", c.UserID)
	req.Header.Add("x-api-key", c.APIKey)
	req.Header.Add("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Error! Error code %s, Content: %s", resp.Status, string(content))
	}

	return content, err
}

func (c *Client) List() (map[string]*HabiticaTODO, error) {
	resp, err := c.get("https://habitica.com/api/v3/tasks/user?type=todos")
	if err != nil {
		return nil, err
	}

	wrapper := new(HabiticaTODOWrapper)
	err = json.Unmarshal(resp, wrapper)
	if err != nil {
		return nil, err
	}

	out := map[string]*HabiticaTODO{}
	for _, todo := range wrapper.Data {
		out[todo.Alias] = todo
	}

	return out, nil
}

func (c *Client) Create(todo *HabiticaTODO) error {
	todoBytes, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	_, err = c.post("https://habitica.com/api/v3/tasks/user", todoBytes)
	return err
}

func (c *Client) Complete(todo *HabiticaTODO) error {
	_, err := c.post(
		fmt.Sprintf("https://habitica.com/api/v3/tasks/%s/score/up", todo.Alias),
		nil,
	)
	return err
}

func (c *Client) Delete(todo *HabiticaTODO) error {
	_, err := c.delete("https://habitica.com/api/v3/tasks/" + todo.Alias)
	return err
}
