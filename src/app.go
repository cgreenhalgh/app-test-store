package main

import (
	//"encoding/json"
	//"errors"
	"log"
	//"net/http"
	//"net/url"
	"os"
	//"strings"
	//"strconv"
	//"sync"
	"time"
	
	//"github.com/gorilla/mux"
	//databox "github.com/me-box/lib-go-databox"
	databox "github.com/cgreenhalgh/lib-go-databox"
)

var dataStoreHref = os.Getenv("DATABOX_STORE_JSON_ENDPOINT")
var dataStoreHref2 = os.Getenv("DATABOX_STORE_TIMESERIES_ENDPOINT")

// Note: must match manifest!
const STORE_TYPE = "store-json"

type TsTest func(ts databox.TimeSeries_0_2_0, arg1 string, arg2 int) error

func doTest(name string, test TsTest, arg1 string, arg2 int, ts databox.TimeSeries_0_2_0, reps int) {
	start := time.Now()
	for i:=0; i<reps; i++ {
		test(ts, arg1, arg2)
	}
	stop := time.Now()
	elapsed := stop.Sub(start)
	log.Printf("test %s, args len %d, %d took %f average", name, len(arg1), arg2, elapsed.Seconds()/float64(reps))
}

var start = time.Now()
var next = start
var written = 0

func insert(ts databox.TimeSeries_0_2_0, arg1 string, arg2 int) error {
	//log.Printf("test!")
	ts.WriteRawValueAt(arg1, next)
	written = written+1
	next = next.Add(1*time.Second)
	return nil
}

func readLatest(ts databox.TimeSeries_0_2_0, arg string, arg2 int) error {
	//log.Printf("test!")
	_,err := ts.ReadLatest()
	return err
}

func readRange1(ts databox.TimeSeries_0_2_0, arg string, arg2 int) error {
	//log.Printf("test!")
	from := next.Add(time.Second*time.Duration(-arg2))
	to := next.Add(1*time.Second)
	_,err := ts.ReadRange(from, to)
	return err
}

func insertUntilN(ts databox.TimeSeries_0_2_0, arg string, n int) {
	for written < n {
		insert(ts, arg, 0)
	}
}

// Note: must match manifest!
const STORE_TYPE2 = "store-json"

type TsTest2 func(ts string, arg1 string, arg2 int) error

func doTest2(name string, test TsTest2, arg1 string, arg2 int, ts string, reps int) {
	start := time.Now()
	for i:=0; i<reps; i++ {
		test(ts, arg1, arg2)
	}
	stop := time.Now()
	elapsed := stop.Sub(start)
	log.Printf("test %s, args len %d, %d took %f average", name, len(arg1), arg2, elapsed.Seconds()/float64(reps))
}

var start2 = time.Now()
var next2 = start
var written2 = 0

func insert2(ts string, arg1 string, arg2 int) error {
	//log.Printf("test!")
	// TODO at time...
	_,err := databox.StoreTsWrite(ts, arg1) // next
	if err != nil {
		log.Printf("Error insert: %s", err.Error())
	}
	written2 = written2+1
	next2 = next2.Add(1*time.Second)
	return nil
}

func readLatest2(ts string, arg string, arg2 int) error {
	//log.Printf("test!")
	_,err := databox.StoreTsLatest(ts)
	return err
}

func readRange2(ts string, arg string, arg2 int) error {
	//log.Printf("test!")
	from := next2.Add(time.Second*time.Duration(-arg2))
	to := next2.Add(1*time.Second)
	// TODO proper range; times
	_,err := databox.StoreTsRange(ts, int(from.Unix()), int(to.Unix()))
	return err
}

func insertUntilN2(ts string, arg string, n int) {
	for written2 < n {
		insert2(ts, arg, 0)
	}
}

func main() {
	//Wait for my store to become active
	databox.WaitForStoreStatus(dataStoreHref)
	databox.WaitForStoreStatus(dataStoreHref2)
	log.Printf("Stores ready")
	
	COUNTS := []int { 0, 100, 1000, 10000 } //, 1000000 }

	var metadata = databox.StoreMetadata{
			Description:    "test",
			ContentType:    "application/json",
			Vendor:         "test",
			DataSourceType: "test",
			DataSourceID:   "test",
			StoreType:      STORE_TYPE,
			IsActuator:     false,
			Unit:           "",
			Location:       "",
		}
	
	ts,err := databox.MakeStoreTimeSeries_0_2_0(dataStoreHref, metadata.DataSourceID, metadata.StoreType)
	if err != nil {
		log.Printf("Error making timeseries api: %s", err.Error())
		return
	}

	// first...
	insert(ts, "\"\"", 0)
	for i:=0; i<len(COUNTS); i++ {
		insertUntilN(ts, "\"\"", COUNTS[i])
		log.Printf("Written %d", written)
		doTest("read latest", readLatest, "", 0, ts, 10)
		for j:=0; j<=i; j++ {
			log.Printf("range -%d", COUNTS[j])
			doTest("read range", readRange1, "", COUNTS[j], ts, 10)
		}
		doTest("insert", insert, "\"\"", 0, ts, 10)
	}

	// second...
	log.Printf("==========================")
	log.Printf("timeseries...")
	ts2 := dataStoreHref2+"/test"
	log.Printf("timeseries URL: %s", ts2)
	// doesn't work with string
	insert2(ts2, "{}", 0)
	for i:=0; i<len(COUNTS); i++ {
		insertUntilN2(ts2, "{}", COUNTS[i])
		log.Printf("Written %d", written2)
		doTest2("read latest", readLatest2, "", 0, ts2, 10)
		for j:=0; j<=i; j++ {
			log.Printf("range -%d", COUNTS[j])
			doTest2("read range", readRange2, "", COUNTS[j], ts2, 10)
		}
		doTest2("insert", insert2, "{}", 0, ts2, 10)
	}
	
	log.Printf("Bye")
}
