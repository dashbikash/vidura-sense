package natsio

import (
	"github.com/dashbikash/vidura-sense/internal/system"
	"github.com/nats-io/nats.go"
)

func KVPut(bucket string, key string, val string) {
	nc, _ := nats.Connect(nats.DefaultURL)

	defer nc.Drain()

	js, _ := nc.JetStream()
	kv, _ := js.KeyValue(bucket)
	kv.Put(key, []byte(val))

}

func KVGet(bucket string, key string) string {
	nc, _ := nats.Connect(nats.DefaultURL)

	defer nc.Drain()

	js, _ := nc.JetStream()
	kv, _ := js.KeyValue(bucket)
	entry, err := kv.Get(key)
	if err != nil {
		system.Log.Error(err.Error())
		return ""
	}
	return string(entry.Value())

}
