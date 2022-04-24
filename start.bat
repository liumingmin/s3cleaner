set S3_AK=
set S3_SK=
set S3_ENDPOINT=http://localhost:9005
set S3_BUCKET=test
set S3_REGION=zh-south-1
set S3_COPY_BUCKET=testtmp

./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=0  > output0.txt
./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=1  > output1.txt
./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=2  > output2.txt
./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=3  > output3.txt
./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=4  > output4.txt
./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=5  > output5.txt
./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=6  > output6.txt
./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=7  > output7.txt
./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=8  > output8.txt
./s3cleaner --mode=2 --page=100000 --hashn=10 --hashi=9  > output9.txt