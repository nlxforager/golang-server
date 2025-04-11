package helloworldpb

import (
	"log"
	"os"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestWritePersonWritesPerson(t *testing.T) {
	p := &HelloRequest{
		Name: "HELLO",
	}

	buf, err := proto.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile("test.data", buf, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}

	f, err := os.ReadFile("test.data")
	if err != nil {
		log.Fatalln("Failed to read address book:", err)
	}

	var fp HelloRequest

	err = proto.Unmarshal(f, &fp)
	if err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	if p.Name != fp.Name {
		log.Fatalln("Name mismatch")
	}

}
