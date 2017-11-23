package eve

import (
	"crypto/md5"
	"encoding/hex"
)

var (
	// MetaDB contains the meta data which buckets are in what database
	MetaDB = "evbolt.meta.db"
	// MetaBUCKET the name of the bucket in the meta db
	MetaBUCKET = "meta"
	// MetaName the name of the meta result set
	MetaName = "EVBolt Meta"
)

// dbBucketID returns the md5 hashed {db}/{bucket} id as a string
func dbBucketID(db, bucket string) string {
	msg := db + "/" + bucket
	hasher := md5.New()
	hasher.Write([]byte(msg))
	return hex.EncodeToString(hasher.Sum(nil))
}

// EVBoltMetaDbBucketID returns the md5 hashed {db}/{bucket} id as a evschema html string
func EVBoltMetaDbBucketID(db, bucket string) (string, error) {
	return dbBucketID(db, bucket), nil
}

// EVBoltMetaLogDbBucket writes the db bucket (as a evschema html string) into the meta storage
func EVBoltMetaLogDbBucket(db, bucket string) (string, error) {
	msg := db + "/" + bucket
	id := dbBucketID(db, bucket)
	return EVBoltPut(id, msg, MetaDB, MetaBUCKET)
}

// EVBoltMetaAll returns all entries in the meta storage as a evschema html string
func EVBoltMetaAll() (string, error) {
	res, err := EVBoltAll(MetaDB, MetaBUCKET)
	if err != nil {
		return "", err
	}
	all := ""
	for k := range res {
		all += k + " " + res[k] + "\n"
	}
	return all, nil
}

// EVBoltMetaAllBucketsForDb returns all buckets for a given database from the meta storage as a evschema html string
// todo: finish implementation
func EVBoltMetaAllBucketsForDb(db string) (string, error) {
	res, err := EVBoltAll(MetaDB, MetaBUCKET)
	if err != nil {
		return "", err
	}
	all := ""
	for k := range res {
		all += k + " " + res[k] + "\n"
	}
	return all, nil
}

// EVBoltMetaAllDbsForBucket returns all databases for the given bucket from the meta storage as a evschema html string
// todo: finish implementation
func EVBoltMetaAllDbsForBucket(bucket string) (string, error) {
	res, err := EVBoltAll(MetaDB, MetaBUCKET)
	if err != nil {
		return "", err
	}
	all := ""
	for k := range res {
		all += k + " " + res[k] + "\n"
	}
	return all, nil
}
