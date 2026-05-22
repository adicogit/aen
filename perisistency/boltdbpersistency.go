package persistency

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strings"
	"time"

	"aen.it/poolmanager/config"
	"aen.it/poolmanager/log"
	"github.com/boltdb/bolt"
)

const (
	BUCKET_SEPARATOR = ":"
)

type boltDBPersistency struct {
	dbPath string
	bolt   *bolt.DB
}

var BoltDBPersistency *boltDBPersistency

func init() {
	log.Log.Debug("Entering BoltDBPersistency init")
	log.Log.Info("Create BoltDBPersistency")
	BoltDBPersistency = &boltDBPersistency{}
	BoltDBPersistency.loadFromConfig()
	log.Log.Debug("Exiting BoltDBPersistency init")
}

func (db *boltDBPersistency) loadFromConfig() {
	log.Log.Debug("Entering loadFromConfig")
	db.dbPath = config.BoltDBConfig.DBpath
	log.Log.Debug("Exiting loadFromConfig")
}

func (db *boltDBPersistency) OpenDB() error {
	log.Log.Debug("Entering OpenDB")
	var err error
	err = nil
	if db.bolt == nil {
		db.bolt, err = bolt.Open(config.BoltDBConfig.DBpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	}
	log.Log.Debug("Exiting 	OpenDB")
	return err
}

func (db *boltDBPersistency) CloseDB() {
	log.Log.Debug("Entering CloseDB")
	if db.bolt != nil {
		db.bolt.Close()
		db.bolt = nil
	}
	log.Log.Debug("Exiting 	CloseDB")
}

func (db *boltDBPersistency) getKeyAndBucketUsingTX(tx *bolt.Tx, key string) (string, *bolt.Bucket, error) {
	log.Log.Debug("Entering getKeyAndBucketUsingTX", "key", key)

	log.Log.Info("Build bucket chanin and read key name")
	buckets := strings.Split(key, BUCKET_SEPARATOR)
	if len(buckets) <= 1 {
		err := errors.New("Key name must contain at least one bucket name. It must be in the format <bucket name>:...:<bucket name>:<key name>")
		log.Log.Error("Wrong key format fo persisting layer", "error", err)
		log.Log.Debug("Exiting getKeyAndBucket", "key", key, "bucket list", buckets)
		return "", nil, err
	}
	keyName := buckets[len(buckets)-1]
	buckets = buckets[:len(buckets)-1]

	// Get bucket to be used. It is taken from key that must be in the format bucket:...:bucket:key
	var bucketToUse *bolt.Bucket

	for _, bucket := range buckets {
		var err error
		if bucketToUse == nil {
			bucketToUse, err = tx.CreateBucketIfNotExists([]byte(bucket))
		} else {
			bucketToUse, err = bucketToUse.CreateBucketIfNotExists([]byte(bucket))
		}
		if err != nil {
			log.Log.Error("Unable to reach target bucket for provided key", "key", key, "error", err)
			log.Log.Debug("Exiting getKeyAndBucketUsingTX", "key", key, "error", err)
			return keyName, nil, err
		}
	}

	log.Log.Debug("Exiting getKeyAndBucketUsingTX", "key", key)
	return keyName, bucketToUse, nil
}

func (db *boltDBPersistency) getKeyAndBucket(key string) (string, *bolt.Bucket, error) {
	log.Log.Debug("Entering getKeyAndBucket", "key", key)

	var keyName string
	var bucketToUse *bolt.Bucket
	var err error
	if err = db.bolt.Update(func(tx *bolt.Tx) error {
		keyName, bucketToUse, err = db.getKeyAndBucketUsingTX(tx, key)
		return err
	}); err != nil {
		log.Log.Error("Failed to create needed chain for specified key", "key", key, "error", err)
		log.Log.Debug("Exiting getKeyAndBucket", "final key name", keyName, "error", err)
		return keyName, nil, err
	}

	log.Log.Debug("Exiting getKeyAndBucket", "final key name", keyName)
	return keyName, bucketToUse, nil
}

func (db *boltDBPersistency) ReadData(key string, data interface{}) error {
	log.Log.Debug("Entering ReadData")

	// If DB is not opened,  open it
	db.OpenDB()

	keyName, bucketToUse, err := db.getKeyAndBucket(key)
	if err != nil {
		log.Log.Error("Failed to create needed chain for specified key", "key", key, "error", err)
		return err
	}
	if err := db.bolt.View(func(tx *bolt.Tx) error {
		// get specified ley value
		dataRead := bucketToUse.Get([]byte(keyName))
		err = convertFromByteArray(dataRead, data)
		if err != nil {
			log.Log.Error("Failed to deserialize object for specified key", "key", keyName, "data read", dataRead, "error", err)
			return err
		}
		return nil
	}); err != nil {
		log.Log.Error("Failed to read object for specified key", "key", keyName, "error", err)
		return err
	}

	log.Log.Debug("Exiting ReadData")
	return nil
}

func (db *boltDBPersistency) WriteData(key string, data interface{}) error {
	log.Log.Debug("Entering WriteData")
	// If DB is not opened,  open it
	db.OpenDB()

	byteArray, err := convertToByteArray(data)
	if err != nil {
		log.Log.Error("Failed to serialize specified object", "data", data, "error", err)
		return err
	}

	if err := db.bolt.Update(func(tx *bolt.Tx) error {
		keyName, bucketToUse, err := db.getKeyAndBucketUsingTX(tx, key)

		// get value fot specified key
		err = bucketToUse.Put([]byte(keyName), byteArray)

		return err
	}); err != nil {
		log.Log.Error("Failed to read object for specified key", "key", key, "error", err)
		return err
	}

	log.Log.Debug("Exiting WriteData")
	return nil
}

func convertToByteArray(data interface{}) ([]byte, error) {
	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(&data)
	if err != nil {
		log.Log.Error("Unable to convert data into byte array", "data", data, "error", err)
	}
	return buffer.Bytes(), err
}

func convertFromByteArray(dataToRead []byte, data interface{}) error {
	buffer := bytes.NewBuffer(dataToRead)

	dec := gob.NewDecoder(buffer)
	err := dec.Decode(&data)
	if err != nil {
		log.Log.Error("Unable to convert data into object", "data", dataToRead, "error", err)
	}
	return err
}
