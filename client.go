package appmetrica_push

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type client struct {
	httpClient *http.Client
	oAuthToken string
}

func NewClient(token string) Client {
	return &client{oAuthToken: token, httpClient: &http.Client{}}
}

// CreateGroup is a method to create group
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/post-groups.html
func (c client) CreateGroup(group *Group) (*Group, error) {
	res, err := c.sendRequest(groupEndpoint, http.MethodPost, &request{Group: group})
	if err != nil {
		return nil, err
	}
	return res.Group, nil
}

// GetGroups is a method to get all groups
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/get-groups.html
func (c client) GetGroups(appId int) ([]*Group, error) {
	param := strconv.Itoa(appId)
	res, err := c.sendRequest(groupsEndpoint+"?app_id="+param, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return res.Groups, nil
}

// GetGroup is a method to get group by id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/get-group-id.html
func (c client) GetGroup(id int) (*Group, error) {
	param := strconv.Itoa(id)
	res, err := c.sendRequest(groupEndpoint+param, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return res.Group, err
}

// UpdateGroup is a method to update group by id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/put-group-id.html
func (c client) UpdateGroup(id int, group *Group) (*Group, error) {
	param := strconv.Itoa(id)
	res, err := c.sendRequest(groupEndpoint+param, http.MethodPut, &request{Group: group})
	if err != nil {
		return nil, err
	}
	return res.Group, nil
}

// ArchiveGroup is a method to archive group by id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/delete-group-id.html
func (c client) ArchiveGroup(id int) error {
	param := strconv.Itoa(id)
	_, err := c.sendRequest(groupEndpoint+param, http.MethodDelete, nil)
	return err
}

// RestoreGroup is a method to restore group by id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/post-group-id.html
func (c client) RestoreGroup(id int) error {
	param := strconv.Itoa(id)
	_, err := c.sendRequest(groupEndpoint+param+"/restore", http.MethodPost, nil)
	return err
}

// SendPush is a method to batch send pushes
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/post-send-batch.html
func (c client) SendPush(r *PushBatchRequest) (*PushResponse, error) {
	res, err := c.sendRequest(sendEndpoint, http.MethodPost, &request{PushBatchRequest: r})
	if err != nil {
		return nil, err
	}
	return res.PushResponse, nil
}

// GetStatusByTransferId is a method to get dispatch status by transfer id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/get-status-id.html
func (c client) GetStatusByTransferId(transferId int) (*Transfer, error) {
	param := strconv.Itoa(transferId)
	res, err := c.sendRequest(statusEndpoint+param, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return res.Transfer, err
}

// GetStatusByClientTransferId is a method to get dispatch status by client transfer id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/get-status-group-id.html
func (c client) GetStatusByClientTransferId(groupId int, clientTransferId int64) (*Transfer, error) {
	p1 := strconv.Itoa(groupId)
	p2 := strconv.FormatInt(clientTransferId, 10)
	res, err := c.sendRequest(statusEndpoint+p1+"/"+p2, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return res.Transfer, nil
}

func (c client) sendRequest(endpoint string, method string, req *request) (res *response, err error) {
	var r *http.Request
	url := host + endpoint
	if req != nil {
		payload, err := json.Marshal(req)
		if err != nil {
			return
		}
		r, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	} else {
		r, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "OAuth "+c.oAuthToken)

	resp, err := c.httpClient.Do(r)
	defer resp.Body.Close()

	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return
	}

	if len(res.Errors) > 0 {
		msg := ""
		for _, e := range res.Errors {
			msg += e.Message + ". "
		}
		return nil, errors.New(msg)
	}

	return
}
