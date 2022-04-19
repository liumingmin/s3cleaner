#!/usr/bin/env bash

export S3_AK=Vg6p9p/WM55ZbiZkE8Vyzw==
export S3_SK=r0yRc7Yxc0fB7yWRoaWJrvLlC3hShtqBFfqj13PKTLo=
export S3_ENDPOINT=http://localhost:9005
export S3_BUCKET=test
export S3_REGION=zh-south-1
export S3_COPY_BUCKET=testtmp

nohup ./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=0  > output0.txt &
nohup ./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=1  > output1.txt &
nohup ./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=2  > output2.txt &
nohup ./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=3  > output3.txt &
nohup ./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=4  > output4.txt &
nohup ./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=5  > output5.txt &
nohup ./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=6  > output6.txt &
nohup ./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=7  > output7.txt &
nohup ./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=8  > output8.txt &
nohup ./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=9  > output9.txt &