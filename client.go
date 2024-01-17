package appmetrica_push

import (
	"bytes"
	"encoding/json"
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
func (c client) CreateGroup(group *Group) *Group {
	res := c.sendRequest(groupEndpoint, http.MethodPost, &request{Group: group})
	return res.Group
}

// GetGroups is a method to get all groups
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/get-groups.html
func (c client) GetGroups() []*Group {
	res := c.sendRequest(groupsEndpoint, http.MethodGet, nil)
	return res.Groups
}

// GetGroup is a method to get group by id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/get-group-id.html
func (c client) GetGroup(id int) *Group {
	param := strconv.Itoa(id)
	res := c.sendRequest(groupEndpoint+param, http.MethodGet, nil)
	return res.Group
}

// UpdateGroup is a method to update group by id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/put-group-id.html
func (c client) UpdateGroup(id int, group *Group) *Group {
	param := strconv.Itoa(id)
	res := c.sendRequest(groupEndpoint+param, http.MethodPut, &request{Group: group})
	return res.Group
}

// ArchiveGroup is a method to archive group by id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/delete-group-id.html
func (c client) ArchiveGroup(id int) {
	param := strconv.Itoa(id)
	c.sendRequest(groupEndpoint+param, http.MethodDelete, nil)
}

// RestoreGroup is a method to restore group by id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/post-group-id.html
func (c client) RestoreGroup(id int) {
	param := strconv.Itoa(id)
	c.sendRequest(groupEndpoint+param+"/restore", http.MethodPost, nil)
}

// SendPush is a method to batch send pushes
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/post-send-batch.html
func (c client) SendPush(r *PushBatchRequest) *PushResponse {
	res := c.sendRequest(sendEndpoint, http.MethodPost, &request{PushBatchRequest: r})
	return res.PushResponse
}

// GetStatusByTransferId is a method to get dispatch status by transfer id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/get-status-id.html
func (c client) GetStatusByTransferId(transferId int) *Transfer {
	param := strconv.Itoa(transferId)
	res := c.sendRequest(statusEndpoint+param, http.MethodGet, nil)
	return res.Transfer
}

// GetStatusByClientTransferId is a method to get dispatch status by client transfer id
// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/get-status-group-id.html
func (c client) GetStatusByClientTransferId(groupId int, clientTransferId int64) *Transfer {
	p1 := strconv.Itoa(groupId)
	p2 := strconv.FormatInt(clientTransferId, 10)
	res := c.sendRequest(statusEndpoint+p1+"/"+p2, http.MethodGet, nil)
	return res.Transfer
}

func (c client) sendRequest(endpoint string, method string, req *request) (res *response) {
	var r *http.Request
	var err error
	url := host + endpoint
	if req != nil {
		payload, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}
		r, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	} else {
		r, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "OAuth "+c.oAuthToken)

	resp, err := c.httpClient.Do(r)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		panic(err)
	}

	return
}
