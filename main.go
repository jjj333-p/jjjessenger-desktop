package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type messageObj struct {
	content, name                                     *string
	renderedContent, renderedName, renderedMessageObj fyne.CanvasObject
	timestamp                                         time.Time
}

// returns msg object configured with the input data
func NewMessageObject(content string, authorname string) messageObj {

	//create rendered text
	rn := widget.NewLabel(authorname)
	rc := widget.NewLabel(content)

	//export all the configured data feilds
	return messageObj{
		content:            &content,
		name:               &authorname,
		renderedContent:    rc,
		renderedName:       rn,
		timestamp:          time.Now().UTC(),
		renderedMessageObj: container.New(layout.NewBorderLayout(rn, nil, nil, nil), rn, rc),
	}

}

// struct to pair chatUI to its id
type chatUI struct {
	id             *int
	ui, sendbutton fyne.CanvasObject
	renderedmsgs   *fyne.Container
	msgbox         *container.Scroll
	msgs           []messageObj
	mainentry      *widget.Entry
}

// handling new messages
func (chat chatUI) append(content string, authorname string) {

	//make object from what is currently in the text entry box
	msg := NewMessageObject(content, authorname)

	//clear the text input box
	chat.mainentry.SetText("")

	//add msg object to the array of msgs
	chat.msgs = append(chat.msgs, msg)

	//add msg to rendered obj
	chat.renderedmsgs.Add(msg.renderedMessageObj)

	//scroll box to bottom to see new msg
	chat.msgbox.ScrollToBottom()

	// Create a map to hold the request data
	requestData := map[string]string{
		"authorname": authorname,
		"content":    content,
	}

	// Marshal the map into a JSON object
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		panic(err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/post", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// if resp.Body == nil {
	// 	fmt.Println("No body response")
	// 	return
	// }

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Print the response body
	fmt.Println(string(body))

}

// creat the chat box for a particular chat
func NewChatUI(id int) chatUI {

	//create object and add the id to it.
	chat := chatUI{
		id: &id,
	}

	//create typing box
	chat.mainentry = widget.NewMultiLineEntry()

	//that background text thingy before you type
	chat.mainentry.SetPlaceHolder("eeee...")

	//create object for the rendered messages
	chat.renderedmsgs = container.New(layout.NewVBoxLayout())

	//add scroll to the rendered text
	chat.msgbox = container.NewVScroll(chat.renderedmsgs)

	//array to store msgs
	chat.msgs = []messageObj{}

	//button to send message in mainentry, and add message to chat on send button
	chat.sendbutton = widget.NewButton("send ^", func() { chat.append(chat.mainentry.Text, "you") })

	//put together entry and send button
	input := container.NewBorder(nil, nil, nil, chat.sendbutton, chat.mainentry)

	//attaching entry and button to the msgbox
	chat.ui = container.NewBorder(nil, input, nil, nil, chat.msgbox)

	//return out object to use.
	return chat

}

func main() {

	//make objects for the app
	a := app.New()
	w := a.NewWindow("pain")

	divider := container.NewHSplit(widget.NewTextGridFromString("aaaaaaa"), NewChatUI(0).ui)

	w.SetContent(divider)

	//default size of the window
	w.Resize(fyne.NewSize(1000, 500))

	//just posting
	w.ShowAndRun()

}
