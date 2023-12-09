package tree

import (
	"fmt"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func TestTree(t *testing.T) {
	sss := ptrace.NewSpanSlice()

	currTime := time.Now()
	duration := time.Minute
	latency := time.Hour
	for i := 0; i < 5; i++ {
		sss.AppendEmpty()
		s := sss.At(i)
		s.SetName(fmt.Sprintf("span-%v", i))
		// id := pcommon.NewSpanIDEmpty()
		id := [8]byte{0, 0, 0, 0, 0, 0, 0, uint8(i + 1)}
		s.SetSpanID(id)
		if i == 0 {
			s.SetParentSpanID(pcommon.NewSpanIDEmpty())
		} else {
			s.SetParentSpanID([8]byte{0, 0, 0, 0, 0, 0, 0, uint8(i)})
		}
		s.SetStartTimestamp(pcommon.NewTimestampFromTime(currTime))
		currTime = currTime.Add(duration)
		s.SetEndTimestamp(pcommon.NewTimestampFromTime(currTime))
		currTime = currTime.Add(latency)
	}
	sss.AppendEmpty()
	s := sss.At(5)
	s.SetName("span-y")
	s.SetParentSpanID([8]byte{0, 0, 0, 0, 0, 0, 0, uint8(2)})
	// for i := 0; i < 5; i++ {
	// 	s := sss.At(i)
	// 	t.Log(s.Name(), s.ParentSpanID(), s.SpanID(), s.StartTimestamp(), s.EndTimestamp())
	// }

	for i := 6; i < 10; i++ {
		sss.AppendEmpty()
		s := sss.At(i)
		s.SetName(fmt.Sprintf("span-%v", i))
		// id := pcommon.NewSpanIDEmpty()
		id := [8]byte{0, 0, 0, 0, 0, 0, uint8(i + 1), 0}
		s.SetSpanID(id)
		if i == 6 {
			s.SetParentSpanID(pcommon.NewSpanIDEmpty())
		} else {
			s.SetParentSpanID([8]byte{0, 0, 0, 0, 0, 0, uint8(i), 0})
		}
		s.SetStartTimestamp(pcommon.NewTimestampFromTime(currTime))
		currTime = currTime.Add(duration)
		s.SetEndTimestamp(pcommon.NewTimestampFromTime(currTime))
		currTime = currTime.Add(latency)
	}

	forest := NewForest()
	forest.AddSpans(sss)
	forest.Print()
}
