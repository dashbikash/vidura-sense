package natsio

import (
	"time"

	"github.com/dashbikash/vidura-sense/internal/system"
	"github.com/nats-io/nats.go"
)

func init() {
	nc, _ := nats.Connect(nats.DefaultURL)

	defer nc.Drain()
	js, _ := nc.JetStream()
	if _, err := js.CreateKeyValue(&nats.KeyValueConfig{Bucket: system.Config.Data.NatsIO.KvBuckets.RobotsTxt, TTL: time.Hour * 24 * 7}); err != nil {
		system.Log.Debug("Key Value Bucket Already Exist")
	}

}

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
