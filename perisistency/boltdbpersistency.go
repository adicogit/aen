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
		db.bolt, err = bolt.Open(db.dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Log.Error("Unable to open DB", "dbPath", db.dbPath, "error", err)
			log.Log.Debug("Exiting OpenDB")
			return err
		}
	}
	log.Log.Debug("Exiting OpenDB")
	return nil
}

func (db *boltDBPersistency) CloseDB() {
	log.Log.Debug("Entering CloseDB")
	if db.bolt != nil {
		db.bolt.Close()
		db.bolt = nil
	}
	log.Log.Debug("Exiting 	CloseDB")
}

func (db *boltDBPersistency) getKeyAndBucketUsingTX(tx *bolt.Tx, key string, create bool) (string, *bolt.Bucket, error) {
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
			if create {
				bucketToUse, err = tx.CreateBucketIfNotExists([]byte(bucket))
			} else {
				bucketToUse = tx.Bucket([]byte(bucket))
			}
		} else {
			if create {
				bucketToUse, err = bucketToUse.CreateBucketIfNotExists([]byte(bucket))
			} else {
				bucketToUse = bucketToUse.Bucket([]byte(bucket))
			}
		}
		if bucketToUse == nil && err == nil {
			if create {
				err = errors.New("Bucket not found")
			} else {
				// For read operations, just return nil bucket to indicate not found
				log.Log.Debug("Bucket not found during read", "key", key)
				return keyName, nil, nil
			}
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

// Removed getKeyAndBucket as it incorrectly returned bucket pointers across transactions.

func (db *boltDBPersistency) ReadData(key string, data interface{}) error {
	log.Log.Debug("Entering ReadData")

	// If DB is not opened,  open it
	err := db.OpenDB()
	if err != nil {
		return err
	}

	if err := db.bolt.View(func(tx *bolt.Tx) error {
		keyName, bucketToUse, err := db.getKeyAndBucketUsingTX(tx, key, false)
		if err != nil {
			return err
		}
		if bucketToUse == nil {
			// Bucket not found, nothing to read
			return nil
		}
		// get specified ley value
		dataRead := bucketToUse.Get([]byte(keyName))
		if dataRead == nil {
			// Key not found
			return nil
		}
		err = convertFromByteArray(dataRead, data)
		if err != nil {
			log.Log.Error("Failed to deserialize object for specified key", "key", keyName, "data read", dataRead, "error", err)
			return err
		}
		return nil
	}); err != nil {
		log.Log.Error("Failed to read object for specified key", "key", key, "error", err)
		return err
	}

	log.Log.Debug("Exiting ReadData")
	return nil
}

func (db *boltDBPersistency) WriteData(key string, data interface{}) error {
	log.Log.Debug("Entering WriteData")
	// If DB is not opened,  open it
	err := db.OpenDB()
	if err != nil {
		return err
	}

	byteArray, err := convertToByteArray(data)
	if err != nil {
		log.Log.Error("Failed to serialize specified object", "data", data, "error", err)
		return err
	}

	if err := db.bolt.Update(func(tx *bolt.Tx) error {
		keyName, bucketToUse, err := db.getKeyAndBucketUsingTX(tx, key, true)

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
	err := enc.Encode(data)
	if err != nil {
		log.Log.Error("Unable to convert data into byte array", "data", data, "error", err)
	}
	return buffer.Bytes(), err
}

func convertFromByteArray(dataToRead []byte, data interface{}) error {
	buffer := bytes.NewBuffer(dataToRead)

	dec := gob.NewDecoder(buffer)
	err := dec.Decode(data)
	if err != nil {
		log.Log.Error("Unable to convert data into object", "data", dataToRead, "error", err)
	}
	return err
}
