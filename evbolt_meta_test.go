package eve

import (
	"os"
	"testing"
)

func Test_DbBucketId(t *testing.T) {
	if dbBucketID("testdb", "testbucket") != "e28eeb9e6b5288e6ddb883637d00c4bc" {
		t.Error("the dbCucketId does not work as expected!")
	}
}

func Test_EVBoltMetaDbBucketID(t *testing.T) {
	id, err := EVBoltMetaDbBucketID("testdb", "testbucket")
	if err != nil {
		t.Error(err)
	}
	if id != "e28eeb9e6b5288e6ddb883637d00c4bc" {
		t.Error("the EVBoltMetaDbBucketID does not work as expected!")
	}
}

func Test_EVBoltMetaLogDbBucket(t *testing.T) {
	id, err := EVBoltMetaLogDbBucket("testdb", "testbucket")
	if err != nil {
		t.Error(err)
	}
	if id != "e28eeb9e6b5288e6ddb883637d00c4bc" {
		t.Error("the EVBoltMetaDbBucketID does not work as expected!")
	}
	os.Remove(EVBOLT_ROOT + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltMetaAll(t *testing.T) {
	all, err := EVBoltMetaAll()
	if err != nil {
		t.Error(err)
	}
	t.Log(all)
	os.Remove(EVBOLT_ROOT + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltMetaAllBucketsForDb(t *testing.T) {
	all, err := EVBoltMetaAllBucketsForDb("testdb")
	if err != nil {
		t.Error(err)
	}
	t.Log(all)
	os.Remove(EVBOLT_ROOT + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltMetaAllDbsForBucket(t *testing.T) {
	all, err := EVBoltMetaAllDbsForBucket("testbucket")
	if err != nil {
		t.Error(err)
	}
	t.Log(all)
	os.Remove(EVBOLT_ROOT + string(os.PathSeparator) + "evbolt.meta.db")
}
