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
	// evBoltRoot is the path where to store the *.db files
	evBoltRoot = "./tests/tmp/"
)

const (
	// EvBoltKeyTypeAuto is key type for autoincrement/int bolt keys
	EvBoltKeyTypeAuto = "a"
	// EvBoltKeyTypeCustom is key type for custom/string (user given) bolt keys
	EvBoltKeyTypeCustom = "c"
)

// SetEvBoltRoot sets the root path where the databases should be saved
func SetEvBoltRoot(rootPath string) {
	evBoltRoot = rootPath
}

// EVBoltItob is required to convert a int value into a bolt specific key []byte value
func EVBoltItob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// EVBoltBoti is required to convert a bolt specific key []byte value into a int value
func EVBoltBoti(b []byte) int {
	if len(b) != 8 {
		return -1
	}
	i := int(b[7])
	return i
}

// EVBoltPut puts a message with given custom/string id as a evschema html string into the given db bucket
func EVBoltPut(id, message, db, bucket string) (string, error) {
	if debug {
		log.Println("EVBoltPut message::", message, "into", db, "::", bucket)
	}
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	err = cdb.Update(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		if debug {
			log.Println("EVBoltPut message to id::", id)
		}
		return sb.Put([]byte(id), []byte(message))
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

// EVBoltCustomPut puts a given message to the given database, bucket with the given custom/string id
func EVBoltCustomPut(id, message, db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return id, err
	}
	return EVBoltPut(id, message, db, bucket)
}

// EVBoltAutoPut puts a given message as a evschema html string to the given database and bucket with a autoincrement/int key
func EVBoltAutoPut(message, db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if debug {
		log.Println("EVBoltAutoPut message::", message, "into", db, "::", bucket)
	}
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
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
		if debug {
			log.Println("EVBoltAutoPut message to id::", id)
		}
		return sb.Put(EVBoltItob(int(id)), []byte(message))
	})
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), nil
}

// EVBoltCustomUpdate updates a given message as a evschema html string into the given database,bucket with a custom/string id
func EVBoltCustomUpdate(message, db, bucket, id string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return id, err
	}
	if debug {
		log.Println("EVBoltCustomUpdate message::", message, "into", db, "::", bucket, "with the given key::", id)
	}
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	err = cdb.Update(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		if debug {
			log.Println("EVBoltCustomUpdate message to id::", id)
		}
		return sb.Put([]byte(id), []byte(message))
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

// EVBoltAutoUpdate updates a given message as a evschema html string into the given database,bucket with a autoincrement/int id
func EVBoltAutoUpdate(message, db, bucket, id string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if debug {
		log.Println("EVBoltAutoUpdate message::", message, "into", db, "::", bucket)
	}
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
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
		if debug {
			log.Println("EVBoltAutoUpdate message to id::", id)
		}
		return sb.Put(EVBoltItob(nid), []byte(message))
	})
	if err != nil {
		return "", err
	}
	return strconv.Itoa(nid), nil
}

// EVBoltLast returns the last entry of the given database bucket as a evschema html string
func EVBoltLast(db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
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
	if debug {
		log.Println("EVBoltLast::value", string(value))
	}
	return string(value), nil
}

// EVBoltFirst returns the first entry of the given database bucketas a evschema html string
func EVBoltFirst(db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
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
	if debug {
		log.Println("EVBoltFirst::value", string(value))
	}
	return string(value), nil
}

// EVBoltCustomGet returns the entry of the given database bucket with the given custom/string key as a evschema html string
func EVBoltCustomGet(db, bucket, key string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return key, err
	}
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
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
	if debug {
		log.Println("EVBoltCustomGet::value", result)
	}
	return result, nil
}

// EVBoltAutoGet returns the entry of the given database bucket with the given autoincrement/int key as a evschema html string
func EVBoltAutoGet(db, bucket, key string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return key, err
	}
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	result := ""
	err = cdb.View(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		sb.ForEach(func(k, v []byte) error {
			if strconv.Itoa(EVBoltBoti(k)) == key {
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
	if debug {
		log.Println("EVBoltAutoGet::value", result)
	}
	return result, nil
}

// EVBoltAll returns all entries for the given database bucket as a evschema html string
func EVBoltAll(db, bucket string) (map[string]string, error) {
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
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
	if debug {
		log.Println("EVBoltAll::values", results)
	}
	return results, nil
}

// EVBoltAllHTML returns all entries for the given database bucket as a evschema html string
func EVBoltAllHTML(db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	res, err := EVBoltAll(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	allEntries := `<div itemscope="" itemtype="https://schema.org/ItemList">`
	counter := 0
	for k := range res {
		counter++
		allEntries += `<div itemprop="itemListElement" itemscope="" itemtype="http://schema.org/ListItem"><a itemprop="item" evboltkey="` + k + `" href="http://127.0.0.1:9092/` + VERSION + `eve/evbolt.html?database=` + db + `&bucket=` + bucket + `&key=` + k + `"><div itemprop="name">` + res[k] + `</div><meta itemprop="position" content="` + strconv.Itoa(counter) + `"/></a></div>`
	}
	return allEntries + `</div>`, nil
}

// EVBoltAllJSON returns all entries for the given database bucket as a evschema html string
func EVBoltAllJSON(db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	res, err := EVBoltAll(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	allEntries, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(allEntries), nil
}

// EVBoltAllString returns all entries for the given database bucket as a evschema html string
func EVBoltAllString(db, bucket string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	res, err := EVBoltAll(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	allEntries := ""
	for k := range res {
		allEntries += k + " " + res[k] + "\n"
	}
	return allEntries, nil
}

// EVBoltAutoDelete deletes a entry of the given database bucket at the given autoincrement/int key as a evschema html string
func EVBoltAutoDelete(db, bucket, id string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if debug {
		log.Println("EVBoltAutoDelete message from", db, "::", bucket)
	}
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
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
		if debug {
			log.Println("EVBoltAutoDelete message with id::", id)
		}
		return sb.Delete(EVBoltItob(nid))
	})
	if err != nil {
		return "", err
	}
	return strconv.Itoa(nid), nil
}

// EVBoltCustomDelete deletes a entry of the given database bucket at the given custom/string key as a evschema html string
func EVBoltCustomDelete(db, bucket, id string) (string, error) {
	_, err := EVBoltMetaLogDbBucket(db, bucket)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if debug {
		log.Println("EVBoltCustomDelete message from", db, "::", bucket)
	}
	cdb, err := bolt.Open(evBoltRoot+string(os.PathSeparator)+db, 0777, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer cdb.Close()
	cdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	err = cdb.Update(func(tx *bolt.Tx) error {
		sb := tx.Bucket([]byte(bucket))
		if debug {
			log.Println("EVBoltCustomDelete message with id::", id)
		}
		return sb.Delete([]byte(id))
	})
	if err != nil {
		return "", err
	}
	return id, nil
}
