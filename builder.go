package appmetrica_push

func NewCreateGroupRequest(appId int, name string) *Group {
	return &Group{AppId: appId, Name: name}
}

func NewUpdateGroupRequest(name string) *Group {
	return &Group{Name: name}
}

func NewPushBatchRequestBody(groupId int, tag string) *PushBatchRequest {
	return &PushBatchRequest{
		GroupID: groupId,
		Tag:     tag,
		Batch:   make([]*Batch, 0),
	}
}

func NewBatch() *Batch {
	return &Batch{
		Messages: &Message{},
		Devices:  make([]*Device, 0),
	}
}

func NewAndroidMessage(title string, text string, silent bool) *AndroidMessage {
	return &AndroidMessage{
		Silent:  silent,
		Content: &AndroidContent{Title: title, Text: text},
	}
}

func NewAndroidOpenAction(deeplink string) *AndroidAction {
	return &AndroidAction{Deeplink: deeplink}
}

func NewIOSMessage(title string, text string, silent bool) *IOSMessage {
	return &IOSMessage{
		Silent:  silent,
		Content: &IOSContent{Title: title, Text: text, Attachments: make([]*Attachment, 0)},
	}
}

func NewIOSOpenAction(url string) *IOSAction {
	return &IOSAction{URL: url}
}

func NewDevice(idType string, idValues ...string) *Device {
	return &Device{
		IDType:   idType,
		IDValues: idValues,
	}
}
