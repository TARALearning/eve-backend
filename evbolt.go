package eve

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

var (
	// EVBOLT_ROOT is the path where to store the *.db files
	EVBOLT_ROOT = "./tests/tmp/"
	// EVBOLT_KEY_TYPE_AUTO is key type for autoincrement/int bolt keys
	EVBOLT_KEY_TYPE_AUTO = "a"
	// EVBOLT_KEY_TYPE_CUSTOM is key type for custom/string (user given) bolt keys
	EVBOLT_KEY_TYPE_CUSTOM = "c"
)

// itob is required to convert a int value into a bolt specific key []byte value
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// boti is required to convert a bolt specific key []byte value into a int value
func boti(b []byte) int {
	if len(b) != 8 {
		return -1
	}
	i := int(b[7])
	return i
}

// boltPut puts a message with given custom/string id as a evschema html string into the given db bucket
func boltPut(id, message, db, bucket string) (string, error) {
	if DEBUG {
		log.Println("boltPut message::", message, "into", db, "::", bucket)
	}
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	err = cdb.Update(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		if DEBUG {
			log.Println("boltPut message to id::", id)
		}
		return sb.Put([]byte(id), []byte(message))
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

// cPut puts a given message to the given database, bucket with the given custom/string id
func EVBoltCustomPut(id, message, db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return id, err
	}
	return boltPut(id, message, db, bucket)
}

// aPut puts a given message as a evschema html string to the given database and bucket with a autoincrement/int key
func EVBoltAutoPut(message, db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if DEBUG {
		log.Println("aPut message::", message, "into", db, "::", bucket)
	}
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	var id uint64
	err = cdb.Update(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		id, err = sb.NextSequence()
		if err != nil {
			return err
		}
		if DEBUG {
			log.Println("aPut message to id::", id)
		}
		return sb.Put(itob(int(id)), []byte(message))
	})
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), nil
}

// cUpdate updates a given message as a evschema html string into the given database,bucket with a custom/string id
func EVBoltCustomUpdate(message, db, bucket, id string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return id, err
	}
	if DEBUG {
		log.Println("cUpdate message::", message, "into", db, "::", bucket, "with the given key::", id)
	}
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	err = cdb.Update(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		if DEBUG {
			log.Println("cUpdate message to id::", id)
		}
		return sb.Put([]byte(id), []byte(message))
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

// aUpdate updates a given message as a evschema html string into the given database,bucket with a autoincrement/int id
func EVBoltAutoUpdate(message, db, bucket, id string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if DEBUG {
		log.Println("aUpdate message::", message, "into", db, "::", bucket)
	}
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	nid, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	err = cdb.Update(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		if DEBUG {
			log.Println("aUpdate message to id::", id)
		}
		return sb.Put(itob(nid), []byte(message))
	})
	if err != nil {
		return "", err
	}
	return strconv.Itoa(nid), nil
}

// last returns the last entry of the given database bucket as a evschema html string
func EVBoltLast(db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	value := make([]byte, 0)
	rkey := make([]byte, 0)
	err = cdb.View(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		c := sb.Cursor()
		rkey, value = c.Last()
		return nil
	})
	if err != nil {
		return "", err
	}
	if string(rkey) == "" && string(value) == "" {
		return "", errors.New(NA)
	}
	if DEBUG {
		log.Println("last::value", string(value))
	}
	return string(value), nil
}

// first returns the first entry of the given database bucketas a evschema html string
func EVBoltFirst(db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	value := make([]byte, 0)
	rkey := make([]byte, 0)
	err = cdb.View(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		c := sb.Cursor()
		rkey, value = c.First()
		return nil
	})
	if err != nil {
		return "", err
	}
	if string(rkey) == "" && string(value) == "" {
		return "", errors.New(NA)
	}
	if DEBUG {
		log.Println("first::value", string(value))
	}
	return string(value), nil
}

// cGet returns the entry of the given database bucket with the given custom/string key as a evschema html string
func EVBoltCustomGet(db, bucket, key string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return key, err
	}
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	result := ""
	err = cdb.View(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		result = string(sb.Get([]byte(key)))
		return nil
	})
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New(NA)
	}
	if DEBUG {
		log.Println("cGet::value", result)
	}
	return result, nil
}

// aGet returns the entry of the given database bucket with the given autoincrement/int key as a evschema html string
func EVBoltAutoGet(db, bucket, key string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return key, err
	}
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	result := ""
	err = cdb.View(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		sb.ForEach(func(k, v []byte) error {
			if strconv.Itoa(boti(k)) == key {
				result = string(v)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New(NA)
	}
	if DEBUG {
		log.Println("aGet::value", result)
	}
	return result, nil
}

// boltAll returns all entries for the given database bucket as a evschema html string
func EVBoltAll(db, bucket string) (map[string]string, error) {
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	results := make(map[string]string, 0)
	err = cdb.View(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		sb.ForEach(func(rkey, v []byte) error {
			results[string(rkey)] = string(v)
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if DEBUG {
		log.Println("boltAll::values", results)
	}
	return results, nil
}

// all returns all entries for the given database bucket as a evschema html string
func EVBoltAllHtml(db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	res, err := EVBoltAll(db, bucket)
	allEntries := `<div itemscope="" itemtype="https://schema.org/ItemList">`
	counter := 0
	for k := range res {
		counter++
		allEntries += `<div itemprop="itemListElement" itemscope="" itemtype="http://schema.org/ListItem"><a itemprop="item" evboltkey="` + k + `" href="http://localhost:9092/0.0.1/eve/evbolt.html?database=` + db + `&bucket=` + bucket + `&key=` + k + `"><div itemprop="name">` + res[k] + `</div><meta itemprop="position" content="` + strconv.Itoa(counter) + `"/></a></div>`
	}
	return allEntries + `</div>`, nil
}

// all returns all entries for the given database bucket as a evschema html string
func EVBoltAllJson(db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	res, err := EVBoltAll(db, bucket)
	allEntries, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(allEntries), nil
}

// all returns all entries for the given database bucket as a evschema html string
func EVBoltAllString(db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	res, err := EVBoltAll(db, bucket)
	allEntries := ""
	for k := range res {
		allEntries += k + " " + res[k] + "\n"
	}
	return allEntries, nil
}

// bDelete deletes a entry of the given database bucket at the given autoincrement/int key as a evschema html string
func EVBoltAutoDelete(db, bucket, id string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if DEBUG {
		log.Println("aDelete message from", db, "::", bucket)
	}
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	nid, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	err = cdb.Update(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		if DEBUG {
			log.Println("aDelete message with id::", id)
		}
		return sb.Delete(itob(nid))
	})
	if err != nil {
		return "", err
	}
	return strconv.Itoa(nid), nil
}

// cDelete deletes a entry of the given database bucket at the given custom/string key as a evschema html string
func EVBoltCustomDelete(db, bucket, id string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if DEBUG {
		log.Println("cDelete message from", db, "::", bucket)
	}
	cdb, err := bolt.Open(EVBOLT_ROOT+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	err = cdb.Update(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		if DEBUG {
			log.Println("cDelete message with id::", id)
		}
		return sb.Delete([]byte(id))
	})
	if err != nil {
		return "", err
	}
	return id, nil
}
