package twiliogo

import (
	"encoding/json"
	"net/url"
)

// IPMessage is a IP Messaging Message resource.
type IPMessage struct {
	Sid         string `json:"sid"`
	AccountSid  string `json:"account_sid"`
	ServiceSid  string `json:"service_sid"`
	To          string `json:"to"` // channel sid
	Attributes  string `json:"attributes"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	WasEdited   bool   `json:"was_edited"`
	From        string `json:"from"` // identity
	Body        string `json:"body"`
	URL         string `json:"url"`
}

// IPMessageList gives the results for querying the set of messages. Returns the first page
// by default.
type IPMessageList struct {
	Client   Client
	Messages []IPMessage `json:"messages"`
	Meta     Meta        `json:"meta"`
}

// SendIPMessageToChannel sends a message to a channel.
func SendIPMessageToChannel(client *TwilioIPMessagingClient, serviceSid, channelSid, from, body, attributes string) (*IPMessage, error) {
	var message *IPMessage

	params := url.Values{}
	params.Set("Body", body)
	if from != "" {
		params.Set("From", from)
	}
	params.Set("Attributes", attributes)

	res, err := client.post(params, "/Services/"+serviceSid+"/Channels/"+channelSid+"/Messages")

	if err != nil {
		return message, err
	}

	message = new(IPMessage)
	err = json.Unmarshal(res, message)

	return message, err
}

// UpdateIPMessage updates ane existing IP Message.
func UpdateIPMessage(client *TwilioIPMessagingClient, serviceSid, channelSid, messageSid, body, attributes string) (*IPMessage, error) {
	var message *IPMessage

	params := url.Values{}
	params.Set("Body", body)
	params.Set("Attributes", attributes)

	res, err := client.post(params, "/Services/"+serviceSid+"/Channels/"+channelSid+"/Messages/"+messageSid)

	if err != nil {
		return message, err
	}

	message = new(IPMessage)
	err = json.Unmarshal(res, message)

	return message, err
}

// GetIPChannelMessage returns the specified IP Message in the channel.
func GetIPChannelMessage(client *TwilioIPMessagingClient, serviceSid, channelSid, sid string) (*IPMessage, error) {
	var message *IPMessage

	res, err := client.get(url.Values{}, "/Services/"+serviceSid+"/Channels/"+channelSid+"/Messages/"+sid)

	if err != nil {
		return nil, err
	}

	message = new(IPMessage)
	err = json.Unmarshal(res, message)

	return message, err
}

// ListIPMessages returns the first page of messages for a channel.
func ListIPMessages(client *TwilioIPMessagingClient, serviceSid, channelSid string) (*IPMessageList, error) {
	var messageList *IPMessageList

	body, err := client.get(nil, "/Services/"+serviceSid+"/Channels/"+channelSid+"/Messages")

	if err != nil {
		return messageList, err
	}

	messageList = new(IPMessageList)
	messageList.Client = client
	err = json.Unmarshal(body, messageList)

	return messageList, err
}

// GetMessages recturns the current page of messages.
func (c *IPMessageList) GetMessages() []IPMessage {
	return c.Messages
}

// GetAllMessages returns all of the messages from all of the pages (from here forward).
func (c *IPMessageList) GetAllMessages() ([]IPMessage, error) {
	messages := c.Messages
	t := c

	for t.HasNextPage() {
		var err error
		t, err = t.NextPage()
		if err != nil {
			return nil, err
		}
		messages = append(messages, t.Messages...)
	}
	return messages, nil
}

// HasNextPage returns whether or not there is a next page of messages.
func (c *IPMessageList) HasNextPage() bool {
	return c.Meta.NextPageUrl != ""
}

// NextPage returns the next page of messages.
func (c *IPMessageList) NextPage() (*IPMessageList, error) {
	if !c.HasNextPage() {
		return nil, Error{"No next page"}
	}

	return c.getPage(c.Meta.NextPageUrl)
}

// HasPreviousPage indicates whether or not there is a previous page of results.
func (c *IPMessageList) HasPreviousPage() bool {
	return c.Meta.PreviousPageUrl != ""
}

// PreviousPage returns the previous page of messages.
func (c *IPMessageList) PreviousPage() (*IPMessageList, error) {
	if !c.HasPreviousPage() {
		return nil, Error{"No previous page"}
	}

	return c.getPage(c.Meta.NextPageUrl)
}

// FirstPage returns the first page of messages.
func (c *IPMessageList) FirstPage() (*IPMessageList, error) {
	return c.getPage(c.Meta.FirstPageUrl)
}

func (c *IPMessageList) getPage(uri string) (*IPMessageList, error) {
	var messageList *IPMessageList

	client := c.Client

	body, err := client.get(nil, uri)

	if err != nil {
		return messageList, err
	}

	messageList = new(IPMessageList)
	messageList.Client = client
	err = json.Unmarshal(body, messageList)

	return messageList, err
}
