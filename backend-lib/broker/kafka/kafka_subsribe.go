package libkafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

func (b *broker) Subscribe(ctx context.Context, input SubInput) (output SubOutput, err error) {
	reader := kafka.NewReader(input.Config)

	readerWrapper := Reader{
		R:            reader,
		subTracer:    b.subTracer,
		commitTracer: b.commitTracer,
		groupID:      input.Config.GroupID,
		unmarshal:    input.Unmarshal,
	}
	if input.Unmarshal == nil {
		readerWrapper.unmarshal = json.Unmarshal
	}

	output = SubOutput{
		Reader: readerWrapper,
	}
	return
}
