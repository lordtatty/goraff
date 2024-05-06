package nodeactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type LLMClient interface {
	Chat(systemMsg, userMsg string, stream chan string) (string, error)
}

type LLM struct {
	SystemMsg      string
	UserMsg        string
	Client         LLMClient
	IncludeOutputs []string
	StateVals      map[string]string
}

func (l *LLM) Do(s *goraff.NodeState, r *goraff.StateReadOnly, triggeringNodeID string) error {
	fmt.Println("Running LLM Node" + l.StateVals["name"])
	for key, val := range l.StateVals {
		s.Set(key, val)
	}
	msg := l.buildIncludes(r)
	msg = msg + "\n\n" + l.UserMsg
	streamCh := make(chan string)
	var err error
	go func() {
		_, e := l.Client.Chat(l.SystemMsg, msg, streamCh)
		err = e
		close(streamCh)
	}()
	result := ""
	for r := range streamCh {
		result += r
		fmt.Println("sending message: " + result)
		s.Set("result", result)
	}
	if err != nil {
		return fmt.Errorf("failed to chat: %w", err)
	}
	return nil
}

func (l *LLM) buildIncludes(r *goraff.StateReadOnly) string {
	result := ""
	for _, output := range l.IncludeOutputs {
		wants := r.NodeState(output).Get("result")
		name := r.NodeState(output).Get("name")
		wantStr := fmt.Sprintf("NAME: %s", name)
		resultStr := fmt.Sprintf("RESULT: %s", wants)
		result += fmt.Sprintf("### OUTPUT BLOCK START###\n%s\n%s\n### OUTPUT BLOCK END###\n", wantStr, resultStr)
	}
	return result
}
