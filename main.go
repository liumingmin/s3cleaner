package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//windows: go build -v --tags netgo -ldflags "-s -w -extldflags \"-static\"" -o s3cleaner.exe main.go
//win-linux: set GOOS=linux set GOARCH=amd64 go build -v --tags netgo -ldflags "-s -w -extldflags \"-static\"" -o s3cleaner main.go
//linux: GOOS=linux GOARCH=amd64 go build -v --tags netgo -ldflags '-s -w -extldflags "-static"' -o s3cleaner main.go

var (
	s3AK       = os.Getenv("S3_AK")
	s3SK       = os.Getenv("S3_SK")
	s3Endpoint = os.Getenv("S3_ENDPOINT")
	s3Bucket   = os.Getenv("S3_BUCKET")
	s3Region   = os.Getenv("S3_REGION")

	s3CopyBucket = os.Getenv("S3_COPY_BUCKET")
)

var (
	logger      *log.Logger
	s3client    *s3.S3
	timeOrig, _ = time.Parse("20060102", "20161212")
)

var mode = flag.Int("mode", 1, "mode1 scan, mode2 moveobject mode3 resotreobject")
var pageLen = flag.Int("page", 1, "scan max page len")
var expday = flag.Int("expday", 365*4, "day of expire")
var sample = flag.Int("sample", 1, "if 1 then sample")

func main() {
	flag.Parse()
	ctx := context.Background()

	fmt.Println(fmt.Sprintf("options: mode: %v, page: %v, expireday: %v, sample: %v", *mode, *pageLen, *expday, *sample))

	s3client = initS3Client()
	initSample()

	var matchFileCount, matchfilesize int64
	var total, totalSize int64

	if *mode == 1 {
		//statiscs
		total, totalSize = scan(ctx, s3Bucket, *sample, scanMatcher, func(obj *s3.Object) bool {
			matchFileCount++
			matchfilesize += *obj.Size
			return true
		})
	} else if *mode == 2 {
		//move to bucket tmp
		total, totalSize = scan(ctx, s3Bucket, *sample, scanMatcher, func(obj *s3.Object) bool {
			matchFileCount++
			matchfilesize += *obj.Size

			_, err := s3client.CopyObject(&s3.CopyObjectInput{
				CopySource: aws.String(fmt.Sprintf("%s/%s", s3Bucket, *obj.Key)),
				Bucket:     aws.String(s3CopyBucket),
				Key:        obj.Key,
			})
			if err != nil {
				logger.Fatalf("CopyObject err: %v\n", err)
			}

			_, err = s3client.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(s3Bucket),
				Key:    obj.Key,
			})
			if err != nil {
				logger.Fatalf("DeleteObject err: %v\n", err)
			}
			return true
		})
	} else if *mode == 3 {
		//restore
		total, totalSize = scan(ctx, s3CopyBucket, *sample, func(*s3.Object) bool {
			return true
		}, func(obj *s3.Object) bool {
			matchFileCount++
			matchfilesize += *obj.Size

			_, err := s3client.CopyObject(&s3.CopyObjectInput{
				CopySource: aws.String(fmt.Sprintf("%s/%s", s3CopyBucket, *obj.Key)),
				Bucket:     aws.String(s3Bucket),
				Key:        obj.Key,
			})
			if err != nil {
				logger.Fatalf("Restore err: %v\n", err)
			}

			_, err = s3client.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(s3CopyBucket),
				Key:    obj.Key,
			})
			if err != nil {
				logger.Fatalf("DeleteTmpObject err: %v\n", err)
			}
			return true
		})
	}

	reportStr := fmt.Sprintf("匹配文件: %v(%vGB), 总文件数: %v(%vGB)", matchFileCount, matchfilesize/1024/1024/1024,
		total, totalSize/1024/1024/1024)
	fmt.Println(reportStr)
}

func scanMatcher(obj *s3.Object) bool {
	keyName := *obj.Key

	if strings.HasSuffix(keyName, ".ts") {
		if len(keyName) > 73 {
			timestr := keyName[65:73]
			//fmt.Println(timestr)
			t, err := time.Parse("20060102", timestr)
			if err != nil {
				return false
			}
			if t.Before(time.Now().AddDate(0, 0, -*expday)) && t.After(timeOrig) {
				//fmt.Println(keyName)
				return true
			}
		}
	} else if strings.HasSuffix(keyName, ".mp4") {
		if len(keyName) > 43 {
			timestr := keyName[33:43] //13
			ts, err := strconv.ParseInt(timestr, 10, 64)
			if err != nil {
				return false
			}

			t := time.Unix(ts, 0)
			if t.Before(time.Now().AddDate(0, 0, -*expday)) && t.After(timeOrig) {
				//fmt.Println(keyName)
				return true
			}
		}
	}
	return false
}

func scan(ctx context.Context, srcBucket string, sample int, matcher, handler func(*s3.Object) bool) (int64, int64) {
	total := int64(0)
	totalSize := int64(0)

	i := 0
	params := &s3.ListObjectsInput{Bucket: aws.String(srcBucket)}
	err := s3client.ListObjectsPagesWithContext(ctx, params, func(output *s3.ListObjectsOutput, end bool) bool {
		for _, obj := range output.Contents {
			if matcher(obj) {
				handler(obj)

				if sample == 1 {
					sampleOutput(obj, total)
				}
			}

			total++
			totalSize += *obj.Size
		}
		i++

		if end {
			return false
		}
		return i < *pageLen
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ListObjectsPagesWithContext err: %v\n", err)
	}
	return total, totalSize
}

func initS3Client() *s3.S3 {
	sess, _ := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(s3AK, s3SK, ""),
		Region:           aws.String(s3Region),
		Endpoint:         aws.String(s3Endpoint),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	},
	)

	return s3.New(sess)
}

func sampleOutput(obj *s3.Object, index int64) {
	hashBase := uint32(*pageLen) / 10
	if hashBase < 2 {
		hashBase = 2
	}

	keyName := *obj.Key
	if uint32(index)%hashBase == 1 {
		logger.Println(keyName)
	}
}

func initSample() {
	if *sample == 1 {
		fPath := "sample_" + time.Now().Format("2006-01-02_15_04_05.txt")
		file, err := os.Create(fPath)
		if err != nil {
			log.Fatalln("fail to create log file", err)
			return
		}
		logger = log.New(file, "", log.LstdFlags)
	}
}
