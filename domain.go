package appmetrica_push

type Client interface {
	CreateGroup(group *Group) (*Group, error)
	GetGroups(appId int) ([]*Group, error)
	GetGroup(id int) (*Group, error)
	UpdateGroup(id int, group *Group) (*Group, error)
	ArchiveGroup(id int) error
	RestoreGroup(id int) error
	SendPush(r *PushBatchRequest) (*PushResponse, error)
	GetStatusByTransferId(transferId int) (*Transfer, error)
	GetStatusByClientTransferId(groupId int, clientTransferId int64) (*Transfer, error)
}

const (
	TransferStatusFailed     = "failed"
	TransferStatusInProgress = "in_progress"
	TransferStatusPending    = "pending"
	TransferStatusSent       = "sent"
)

type (
	// response is unified struct for all appmetrica API responses
	response struct {
		Group        *Group        `json:"group,omitempty"`
		Groups       []*Group      `json:"groups,omitempty"`
		PushResponse *PushResponse `json:"push_response,omitempty"`
		Transfer     *Transfer     `json:"transfer,omitempty"`
		Errors       []*Error      `json:"errors,omitempty"`
	}

	// request is unified struct for all appmetrica API requests
	request struct {
		PushBatchRequest *PushBatchRequest `json:"push_batch_request,omitempty"`
		Group            *Group            `json:"group,omitempty"`
	}

	// Group struct represents the response for push group API.
	// AppId and Name is required.
	// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/post-groups.html
	Group struct {
		ID       int    `json:"id,omitempty"`        // ID of group
		AppId    int    `json:"app_id,omitempty"`    // AppId is an id of the app that the group was created for.
		Name     string `json:"name,omitempty"`      // Name of the group, should be unique
		SendRate int    `json:"send_rate,omitempty"` // SendRate is the dispatch speed limit for push messages (number per second). Default is 5000 (max), min is 100
	}

	// PushResponse represents the response for send-batch push request
	// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/post-send-batch.html
	PushResponse struct {
		TransferId       int   `json:"transfer_id"`        // TransferId is used for querying API for dispatch status
		ClientTransferId int64 `json:"client_transfer_id"` // ClientTransferId is specified by user and used for querying API for dispatch status. Should be unique in one group
	}

	// Transfer - Information about the dispatch.
	// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/get-status-id.html
	// https://appmetrica.yandex.com/docs/mobile-api/push/get-status-group-id.html
	Transfer struct {
		ID               int      `json:"id"`                           // id of the dispatch
		GroupId          int      `json:"group_id"`                     // id of the Group
		Status           string   `json:"status"`                       // Status of the dispatch. Can be failed, in_progress, pending, sent, see TransferStatus consts
		Errors           []string `json:"errors"`                       // Errors list
		Tag              string   `json:"tag"`                          // The dispatch Tag. A tag is an arbitrary string that labels every sending by the API. You can label an arbitrary number of sendings by one tag. The report displays tag at the second level.
		CreationDate     string   `json:"creation_date"`                // The Date the dispatch request was created.
		ClientTransferId *int64   `json:"client_transfer_id,omitempty"` // Sending ID specified by the user in the body of the SendBatch request. Would be nil, if Transfer is requested by TransferId, not ClientTransferId and GroupId
	}

	// PushBatchRequest The request to send a group of push notifications.
	// Documentation: https://appmetrica.yandex.com/docs/mobile-api/push/post-send-batch.html
	PushBatchRequest struct {
		GroupID          int      `json:"group_id"`           // GroupID. Required parameter.
		ClientTransferID int64    `json:"client_transfer_id"` // ID of the dispatch, specified by the user. It is used for checking the status of the sending.
		Tag              string   `json:"tag"`                // The dispatch tag. Required parameter.
		Batch            []*Batch `json:"batch"`              // Array of messages objects. It contains push messages with properties.
	}

	// Batch is an array of messages objects. It contains push messages with properties.
	Batch struct {
		Messages *Message  `json:"messages"` // Push message
		Devices  []*Device `json:"devices"`  // Devices to send push notifications to. Each dispatch of messages can contain up to 250,000 devices. Devices are grouped by id_type. One sending can have from 1 to 5 groups. All groups can contain up to 250,000 devices in total in a single HTTP request. For example, if the appmetrica_device_id group contains 100,000 devices, only 150,000 can be specified in the others.
	}

	// Message is a push message
	Message struct {
		Android *AndroidMessage `json:"android"` // AndroidMessage with platform-specific properties
		IOS     *IOSMessage     `json:"iOS"`     // IOSMessage with platform-specific properties
	}

	// AndroidMessage with platform-specific properties
	AndroidMessage struct {
		Silent     bool            `json:"silent"`                // A flag that indicates silent push sending. Possible values: true | false.
		Content    *AndroidContent `json:"content"`               // The content of the push message.
		OpenAction *AndroidAction  `json:"open_action,omitempty"` // The action to be taken when a user clicks on a push notification. If the field is empty, the user click opens the application.
	}

	// AndroidContent is the content of the push message.
	AndroidContent struct {
		Title            string `json:"title"`              // The title of the push message. The value is mandatory for non-silent push messages.
		Text             string `json:"text"`               // The text of the message. The value is mandatory for non-silent push messages.
		Icon             string `json:"icon"`               // The icon is shown in the notification bar. By default, the standard app icon is displayed. To change the icon, set the icon resource ID in the standard /res/drawable/ directory.
		IconBackground   string `json:"icon_background"`    // The color of the message icon. It is specified as a string in the format of the hex code #AARRGGBB. This field is available only for the Android platform.
		Image            string `json:"image"`              // The URL of the image which is displayed in the push message next to the notification text.
		Banner           string `json:"banner"`             // The URL of the image that is shown in the push message. This field is available only for the Android platform.
		Data             string `json:"data"`               // An arbitrary data string. You can pass any data you need as a string value. You can process the data string by using the appropriate AppMetrica Push SDK methods.
		ChannelID        string `json:"channel_id"`         // ID of the notification channel. If the ID is not specified, the default channel is used. Available for Android 8 or higher. For more information about channels, see Android documentation.
		Priority         int    `json:"priority"`           // Notification priority. Acceptable values are in the range of [-2; 2]. The platform determines the priority of messages and takes appropriate actions: interrupts the user (displays a message on the screen), or does not notify the user about the message. On different devices, priority is interpreted differently. This field is available only for the Android platform.
		CollapseKey      int    `json:"collapse_key"`       // Notification ID. The default value is 0. Ignored if there are no push notifications currently displayed for this application. If one or more notifications are displayed and the new message has the same notification ID, the content of this notification will be updated. If the ID is different, the new message is displayed. This field is available only for the Android platform.
		Vibration        []int  `json:"vibration"`          // Vibration pattern on message arrival. Format: [pause in ms, duration of vibration in ms, pause in ms, duration of vibration in ms, ...]. This field is available only for the Android platform.
		LedColor         string `json:"led_color"`          // LED color. It is specified as a string in the format of the hex code #RRGGBB. This field is available only for the Android platform.
		LedInterval      int    `json:"led_interval"`       // The duration of the glow of the LED indicator in ms. This field is available only for the Android platform.
		LedPauseInterval int    `json:"led_pause_interval"` // Pause time for the glow of the LED indicator in ms. This field is available only for the Android platform.
		TimeToLive       int    `json:"time_to_live"`       // The duration of the interval in seconds that FireBase will store the push message if the device is offline or out of range. This field is available only for the Android platform.
		Visibility       string `json:"visibility"`         // Displaying the push message on the lock screen. Ignored on Android 8 and higher (API level 26+), where it's set at the channel level. Acceptable values: secret, private, public. Not set by default. For more information about the visibility property, see the Android documentation (https://developer.android.com/reference/android/app/Notification#VISIBILITY_PRIVATE).
		Urgency          string `json:"urgency"`            // Urgency (priority) of push message delivery. Acceptable values: high, normal. The default value is high. Urgent push messages wake up the device, launch the app in background mode, and get access to the internet for a short time. Urgent push messages are delivered faster and more reliably. For more information about priority of FCM messages (https://firebase.google.com/docs/cloud-messaging/concept-options#setting-the-priority-of-a-message), see the Android documentation (https://developer.android.com/training/monitoring-device-state/doze-standby#using_fcm).
	}

	// AndroidAction is the action to be taken when a user clicks on a push notification. If the field is empty, the user click opens the application.
	AndroidAction struct {
		Deeplink string `json:"deeplink,omitempty"` // The deeplink with an application screen to take a user to after clicking on a push message.
	}

	// IOSMessage with platform-specific properties
	IOSMessage struct {
		Silent     bool        `json:"silent"`                // A flag that indicates silent push sending. Possible values: true | false.
		Content    *IOSContent `json:"content"`               // The content of the push message.
		OpenAction *IOSAction  `json:"open_action,omitempty"` // The action to be taken when a user clicks on a push notification. If the field is empty, the user click opens the application.
	}

	// IOSContent is the content of the push message.
	IOSContent struct {
		Title          string        `json:"title"`           // The title of the push message. The value is mandatory for non-silent push messages.
		Text           string        `json:"text"`            // The text of the message. The value is mandatory for non-silent push messages.
		Badge          int           `json:"badge"`           // Badge number to be displayed on the application icon on message arrival.
		Sound          string        `json:"sound"`           // The message sound. Possible values: default | disable
		ThreadID       string        `json:"thread_id"`       // ID for grouping push notifications. The value is specified in the threadIdentifier (https://developer.apple.com/documentation/usernotifications/unmutablenotificationcontent/1649872-threadidentifier?language=objc) property of the UNNotificationContent (https://developer.apple.com/documentation/usernotifications/unnotificationcontent) object.
		Category       string        `json:"category"`        // Push notifications category. The value is specified in the identifier property of the UNNotificationCategory object. More information about push actions and categories in the Apple documentation (https://developer.apple.com/documentation/usernotifications/declaring_your_actionable_notification_types?language=objc).
		MutableContent int           `json:"mutable_content"` // Indicates Notification Service Extension. If the value is 1, the push notification is processed by the extension. AppMetrica uses it to track the delivery of push notifications. To track the delivery, set up push notification statistics collection (https://appmetrica.yandex.com/docs/mobile-sdk-dg/push/ios-statistics-settings.html) and pass 1 as a field value. If the value is omitted, the number of delivered push notifications in reports is equal to the number of opened messages.
		Expiration     int           `json:"expiration"`      // The duration of time to continue trying to deliver the notification to the user's device. The value should be specified in seconds. If this time expires and the device is still unavailable (for example, it doesn't have internet access), the notification isn't delivered. By default, the time is unrestricted.
		Data           string        `json:"data"`            // The URL to go to when the push message is clicked.
		CollapseID     string        `json:"collapse_id"`     // Collapse ID (apns-collapse-id: https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingwithAPNs.html#//apple_ref/doc/uid/TP40008194-CH11-SW12). Multiple notifications with the same ID are displayed to the user as a single notification.
		Attachments    []*Attachment `json:"attachments"`     // An array of attachments to be added in a push message. Read more in the article "Step 6. (Optional) Configure uploading attached files." (https://appmetrica.yandex.com/docs/mobile-sdk-dg/push/ios-initialize.html#download-file) This field is only available for the iOS platform.
	}

	// IOSAction is the action to be taken when a user clicks on a push notification. If the field is empty, the user click opens the application.
	IOSAction struct {
		URL string `json:"url,omitempty"` // The URL to go to when the push message is clicked.
	}

	// Attachment to be added in a push message. Read more in the article "Step 6. (Optional) Configure uploading attached files." (https://appmetrica.yandex.com/docs/mobile-sdk-dg/push/ios-initialize.html#download-file) This field is only available for the iOS platform.
	Attachment struct {
		ID       string `json:"id"`        // ID of the push message contents. This field is only available for the iOS platform.
		FileURL  string `json:"file_url"`  // URL of the file from the push message. This field is only available for the iOS platform.
		FileType string `json:"file_type"` // Type of the attached file in the push message. For acceptable types, see File types in push messages (https://appmetrica.yandex.com/docs/mobile-api/push/file-type.html). This field is only available for the iOS platform.
	}

	// Device to send push notifications to
	Device struct {
		IDType   string   `json:"id_type"`   // The type of the ID. Acceptable values: appmetrica_device_id, ios_ifa, google_aid, android_push_token, ios_push_token, huawei_push_token, huawei_oaid.
		IDValues []string `json:"id_values"` // List of devices to send push messages to. The list can't be empty.
	}

	Error struct {
		ErrorType string `json:"error_type"`
		Message   string `json:"message"`
	}
)

const (
	host           = "https://push.api.appmetrica.yandex.net/push/v1"
	groupsEndpoint = "/management/groups"
	groupEndpoint  = "/management/group/"
	sendEndpoint   = "/send-batch"
	statusEndpoint = "/status/"
)
