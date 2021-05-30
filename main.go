package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/klauspost/reedsolomon"
)

var (
	decode = flag.Bool("decode", false, "decode or not")
	fileName = flag.String("file", "", "file to encode/decode")
	encodePath = flag.String("encode_path", "", "folder to save encode files")
)

// usage
// ./main --file=./file//testfile.txt -encode_path=./encode
// ./main --file=./file//testfile_recover.txt -encode_path=./encode --decode


func main() {
	flag.Parse()

	if fileName == nil || encodePath == nil{
		panic("need filename and encode path")
	}

	encoder, err := reedsolomon.New(5, 3)
	if err != nil{
		panic(err)
	}

	if *decode{
		shards := make([][]byte, 8)
		var missingShards []int
		for i:= 0; i< 8; i ++{
			encodeFile := path.Join(*encodePath, fmt.Sprintf("encode_%d", i))
			data, err := ioutil.ReadFile(encodeFile)
			if err == nil {
				shards[i] = data
			} else if os.IsNotExist(err){
				missingShards = append(missingShards, i)
				continue
			}else {
				panic(err)
			}

		}
		err = encoder.Reconstruct(shards)
		if err != nil{
			panic(err)
		}
		for _, index := range missingShards{
			encodeFile := path.Join(*encodePath, fmt.Sprintf("encode_%d", index))
			err := ioutil.WriteFile(encodeFile, shards[index], 0644)
			if err != nil{
				panic(err)
			}
		}

		fmt.Printf("reconstruct data done\n")
		f, err := os.OpenFile(*fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil{
			panic(err)
		}
		dataSize := 0
		for i := 0; i < 5; i++{
			dataSize += len(shards[i])
		}
		err = encoder.Join(f, shards, dataSize)
		if err != nil{
			panic(err)
		}
		fmt.Printf("recover file success")
	}else{
		data, err := ioutil.ReadFile(*fileName)
		if err != nil{
			panic(err)
		}
		shards, err := encoder.Split(data)
		if err != nil{
			panic(err)
		}
		fmt.Printf("split data into 5+3=%d shards success.\n", len(shards))
		err = encoder.Encode(shards)
		if err != nil{
			panic(err)
		}
		fmt.Printf("encode data success.")
		err = os.MkdirAll(*encodePath, 0777)
		if err != nil{
			panic(err)
		}

		for i, s := range shards{
			err := ioutil.WriteFile(path.Join(*encodePath, fmt.Sprintf("encode_%d", i)), s, 0644)
			if err != nil{
				panic(err)
			}
		}
	}
}
