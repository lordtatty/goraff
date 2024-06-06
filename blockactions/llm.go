package blockactions

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
}

func (l *LLM) Do(s *goraff.Node, r *goraff.ReadableGraph, triggeringNode *goraff.ReadableNode) error {
	fmt.Println("Running LLM Node")
	msg, err := l.buildIncludes(r)
	if err != nil {
		return fmt.Errorf("error building includes: %w", err)
	}
	msg = msg + "\n\n" + l.UserMsg
	streamCh := make(chan string)
	go func() {
		_, e := l.Client.Chat(l.SystemMsg, msg, streamCh)
		err = e
		close(streamCh)
	}()
	result := ""
	for r := range streamCh {
		result += r
		s.SetStr("result", result)
	}
	if err != nil {
		return fmt.Errorf("failed to chat: %w", err)
	}
	return nil
}

func (l *LLM) buildIncludes(r *goraff.ReadableGraph) (string, error) {
	result := ""
	for _, output := range l.IncludeOutputs {
		n, err := r.FirstNodeByName(output)
		if err != nil {
			return "", fmt.Errorf("error getting node state: %w", err)
		}
		wants := n.FirstStr("result")
		n, err = r.FirstNodeByName(output)
		if err != nil {
			return "", fmt.Errorf("error getting node state: %w", err)
		}
		name := n.FirstStr("name")
		wantStr := fmt.Sprintf("NAME: %s", name)
		resultStr := fmt.Sprintf("RESULT: %s", wants)
		result += fmt.Sprintf("### OUTPUT BLOCK START###\n%s\n%s\n### OUTPUT BLOCK END###\n", wantStr, resultStr)
	}
	return result, nil
}
