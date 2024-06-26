package blockactions_test

import (
	"strings"
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/blockactions"
	"github.com/lordtatty/goraff/mocks"
	"github.com/lordtatty/goraff/notifiers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLLM_Do(t *testing.T) {
	assert := assert.New(t)
	mClient := mocks.NewLLMClient(t)

	// Define the expected input arguments for the mocked method call
	expectedSystemMsg := "system message"
	expectedUserMsg := "user message"
	expectedMessages := []string{"result1", "result2", "result3"} // Define the expected messages to be received
	expectedResult := "result1result2result3"
	expectedError := error(nil)

	mClient.EXPECT().
		Chat(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(systemMsg, userMsg string, stCh chan string) (string, error) {
			// Use a buffered channel to ensure all messages are sent before returning
			done := make(chan struct{})

			go func() {
				defer close(done) // Close the channel when all messages are sent
				for _, msg := range expectedMessages {
					stCh <- msg
				}
			}()

			// Wait for the 'done' channel to be closed
			<-done
			return expectedResult, expectedError
		})

	// Execute the method under test
	sut := &blockactions.LLM{
		SystemMsg: expectedSystemMsg,
		UserMsg:   expectedUserMsg,
		Client:    mClient,
	}
	msgIdx := 0
	ntfy := &notifiers.GraphNotifier{}
	s := &goraff.Graph{Notifier: ntfy}
	r := goraff.NewReadableGraph(s)
	ntfy.Listen(func(ntfy goraff.GraphChangeNotification) {
		msgIdx++
		want := strings.Join(expectedMessages[:msgIdx], "")
		n, err := r.FirstNodeByName("node1")
		assert.NoError(err)
		assert.Equal(want, n.FirstStr("result"))
	})
	n := s.NewNode("node1", nil)

	err := sut.Do(n, r, nil)
	assert.NoError(err)
	assert.Equal(msgIdx, len(expectedMessages))
	assert.Equal(expectedResult, n.Get().FirstStr("result"))
}
