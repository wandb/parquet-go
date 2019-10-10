package main

import (
	"log"
	"time"
	"encoding/json"

	"github.com/wandb/parquet-go-source/local"
	"github.com/wandb/parquet-go/reader"
	"github.com/wandb/parquet-go/writer"
	"github.com/wandb/parquet-go/parquet"
)

type Student struct {
	Name    string  `parquet:"name=name, type=UTF8, encoding=PLAIN_DICTIONARY"`
	Age     int32   `parquet:"name=age, type=INT32"`
	Id      int64   `parquet:"name=id, type=INT64"`
	Weight  float32 `parquet:"name=weight, type=FLOAT"`
	Sex     bool    `parquet:"name=sex, type=BOOLEAN"`
	Day     int32   `parquet:"name=day, type=DATE"`
	Scores	map[string]int32 `parquet:"name=scores, type=MAP, keytype=UTF8, valuetype=INT32"`
	Ignored int32   //without parquet tag and won't write
}

func main() {
	var err error
	fw, err := local.NewLocalFileWriter("to_json.parquet")
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}

	//write
	pw, err := writer.NewParquetWriter(fw, new(Student), 4)
	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}

	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.CompressionType = parquet.CompressionCodec_SNAPPY
	num := 10
	for i := 0; i < num; i++ {
		stu := Student{
			Name:   "StudentName",
			Age:    int32(20 + i%5),
			Id:     int64(i),
			Weight: float32(50.0 + float32(i)*0.1),
			Sex:    bool(i%2 == 0),
			Day:    int32(time.Now().Unix() / 3600 / 24),
			Scores: map[string]int32{
				"math": int32(90 + i%5),
				"physics": int32(90 + i%3),
				"computer": int32(80 + i%10),
			},
		}
		if err = pw.Write(stu); err != nil {
			log.Println("Write error", err)
		}
	}
	if err = pw.WriteStop(); err != nil {
		log.Println("WriteStop error", err)
		return
	}
	log.Println("Write Finished")
	fw.Close()

	///read
	fr, err := local.NewLocalFileReader("to_json.parquet")
	if err != nil {
		log.Println("Can't open file")
		return
	}

	pr, err := reader.NewParquetReader(fr, nil, 4)
	if err != nil {
		log.Println("Can't create parquet reader", err)
		return
	}
	
	num = int(pr.GetNumRows())
	res, err := pr.ReadByNumber(num)
	if err != nil {
		log.Println("Can't read", err)
		return
	}

	jsonBs, err := json.Marshal(res)
	if err != nil {
		log.Println("Can't to json", err)
		return
	}

	log.Println(string(jsonBs))

	pr.ReadStop()
	fr.Close()

}
