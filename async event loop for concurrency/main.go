package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Node struct {
	NodeID    string
	NodeIDs   []string
	NextMsgID int
	mu        sync.Mutex
	outMu     sync.Mutex

	inbox chan Message
	wg    sync.WaitGroup
}

type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

func (n *Node) Send(dest string, body map[string]interface{}) {
	n.mu.Lock()
	msgID := n.NextMsgID
	n.NextMsgID++
	src := n.NodeID
	n.mu.Unlock()

	body["msg_id"] = msgID

	msg := Message{Src: src, Dest: dest, Body: body}
	output, _ := json.Marshal(msg)

	n.outMu.Lock()
	fmt.Println(string(output))
	n.outMu.Unlock()
}

func (n *Node) Reply(request Message, body map[string]interface{}) {
	if msgID, ok := request.Body["msg_id"].(float64); ok {
		body["in_reply_to"] = int(msgID)
	}
	n.Send(request.Src, body)
}

func (n *Node) worker() {
	defer n.wg.Done()

	for msg := range n.inbox {
		n.HandleMessage(msg)
	}
}

func (n *Node) HandleMessage(msg Message) {
	body := msg.Body

	msgType, ok := body["type"].(string)
	if !ok {
		return
	}

	switch msgType {

	case "init":
		n.mu.Lock()
		n.NodeID = body["node_id"].(string)

		rawIDs := body["node_ids"].([]interface{})
		n.NodeIDs = []string{} // reset to avoid duplicates
		for _, id := range rawIDs {
			n.NodeIDs = append(n.NodeIDs, id.(string))
		}
		n.mu.Unlock()

		n.Reply(msg, map[string]interface{}{
			"type": "init_ok",
		})

	case "echo":
		n.Reply(msg, map[string]interface{}{
			"type": "echo_ok",
			"echo": body["echo"],
		})
	}
}

func main() {
	node := &Node{
		inbox: make(chan Message, 100),
	}

	numWorkers := 1
	node.wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go node.worker()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			continue
		}
		node.inbox <- msg
	}

	close(node.inbox)
	node.wg.Wait()
}
