package ipc

import (
	"reflect"
	"testing"
)

func BenchmarkMsgp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var expected = Msgp{
			Mtype: 9,
			Mtext: []byte("henrylee2cn"),
		}
		ptr, textSize := expected.marshal()
		if len(expected.Mtext) != textSize {
			b.Fatalf("expected:%v, actual:%v", len(expected.Mtext), textSize)
		}
		var actual Msgp
		err := actual.unmarshal(textSize, ptr)
		if err != nil {
			b.Fatal(err)
		}
		if !reflect.DeepEqual(expected, actual) {
			b.Fatalf("expected:%v, actual:%v", expected, actual)
		}
	}
}
