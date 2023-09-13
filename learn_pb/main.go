package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"time"

	pb "github.com/iziyang/learn_pb/tutorialpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// 创建一个 AddressBook 对象
	addressBook := &pb.AddressBook{
		People: []*pb.Person{
			{
				Name:  "Alice",
				Id:    1,
				Email: "alice@example.com",
				Phones: []*pb.Person_PhoneNumber{
					{
						Number: "555-1234",
						Type:   pb.Person_HOME,
					},
				},
				LastUpdated: timestamppb.New(time.Now()),
			},
			{
				Name:  "Bob",
				Id:    2,
				Email: "bob@example.com",
				Phones: []*pb.Person_PhoneNumber{
					{
						Number: "555-5678",
						Type:   pb.Person_WORK,
					},
				},
				LastUpdated: timestamppb.New(time.Now()),
			},
		},
	}

	// 序列化 AddressBook 对象
	data, err := proto.Marshal(addressBook)
	if err != nil {
		fmt.Println("Error during serialization:", err)
		return
	}

	// 反序列化 AddressBook 对象
	addressBook2 := &pb.AddressBook{}
	err = proto.Unmarshal(data, addressBook2)
	if err != nil {
		fmt.Println("Error during deserialization:", err)
		return
	}

	fmt.Println("Original AddressBook:", addressBook)
	fmt.Println("Deserialized AddressBook:", addressBook2)
}
