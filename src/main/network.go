package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

func Write(endpoint *string, nodeID *ua.NodeID, value *ua.Variant, ctx context.Context, c *opcua.Client) {
	req := &ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{
			{
				NodeID:      nodeID,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					EncodingMask: ua.DataValueValue,
					Value:        value,
				},
			},
		},
	}

	resp, err := c.WriteWithContext(ctx, req)
	if err != nil {
		log.Fatalf("Read failed: %s", err)
	}
	log.Printf("%v", resp.Results[0])
}

// Use this function to read the value of a specific node based on the consumer
func Read(endpoint *string, nodeID *ua.NodeID, ctx context.Context, c *opcua.Client) {
	req := &ua.ReadRequest{
		MaxAge: 2000,
		NodesToRead: []*ua.ReadValueID{
			{NodeID: nodeID},
		},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	resp, err := c.ReadWithContext(ctx, req)
	if err != nil {
		log.Fatalf("Read failed: %s", err)
	}
	if resp.Results[0].Status != ua.StatusOK {
		log.Fatalf("Status not OK: %v", resp.Results[0].Status)
	}
	log.Printf("%#v", resp.Results[0].Value.Value())
	fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())

}

// Use this function to connect to a server
func Connect(endpoint string, ctx context.Context) *opcua.Client {
	// Connect to the server
	c := opcua.NewClient(endpoint)
	if err := c.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	return c
}
