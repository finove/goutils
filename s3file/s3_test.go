package s3file

import (
	"log"
	"testing"
)

func TestAll(t *testing.T) {
	testStore(t)
}

func testStore(t *testing.T) {
	var err error
	var ss S3Store
	err = ss.Connect("key11", "secretkey222", "ap-northeast-1", "http://s3addr:9000")
	if err != nil {
		t.Fatalf("setup s3 store fail:%v", err)
	}
	bkts, _ := ss.ListBuckets()
	log.Println(bkts)
	for _, bkt := range bkts {
		ss.ListBucketFiles(bkt, "")
		ss.ListBucketFiles(bkt, "202012/")
	}
	log.Println(ss.GetPresignedURL("bucket", "202010/key111222333.tar.gz"))
}
