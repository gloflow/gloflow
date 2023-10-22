/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_aws

import (
	"os"
	"bytes"
	"fmt"
	"net/http"
	"mime"
	"time"
	"path/filepath"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

type GFs3Info struct {
	Client   *s3.S3
	Uploader *s3manager.Uploader
	Session  *session.Session
}

//---------------------------------------------------

func S3getFile(pTargetFileS3pathStr string,
	pTargetFileLocalPathStr string,
	pS3bucketNameStr        string,
	pS3info                 *GFs3Info,
	pRuntimeSys             *gf_core.RuntimeSys) *gf_core.GFerror {
	
	fmt.Printf("target_file_s3_path - %s\n", pTargetFileS3pathStr)
	fmt.Printf("s3_bucket_name      - %s\n", pS3bucketNameStr)
	
	downloader := s3manager.NewDownloader(pS3info.Session)

	// create a local host FS file to store the downloaded image into
	file, err := os.Create(pTargetFileLocalPathStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create local file on host FS, to save a downloaded S3 file to.",
			"file_create_error", 
			map[string]interface{}{
				"target_file__s3_path_str":    pTargetFileS3pathStr,
				"target_file__local_path_str": pTargetFileLocalPathStr,
				"s3_bucket_name_str":          pS3bucketNameStr,
			}, err, "gf_core", pRuntimeSys)
		return gfErr
	}

	// write downloaded S3 file contents to the local FS file
	bytesDownloadedInt, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(pS3bucketNameStr),
		Key:    aws.String(pTargetFileS3pathStr),
	})

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to download an image from S3 bucket",
			"s3_file_download_error", nil, err, "gf_core", pRuntimeSys)
		return gfErr
	}
	fmt.Printf("file downloaded, %d bytes\n", bytesDownloadedInt)


	return nil
}

//---------------------------------------------------
// S3_INIT

func S3init(pRuntimeSys *gf_core.RuntimeSys) (*GFs3Info, *gf_core.GFerror) {
	
	config := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("s3.amazonaws.com"),
		S3ForcePathStyle: aws.Bool(true),      // <-- without these lines. All will fail! fork you aws!
		// Credentials:      creds,
		// LogLevel:         0, // <-- feel free to crank it up 
	}

	//--------------
	/*
	// STATIC_CREDENTIALS - they're non-empty and should be constructed. otherwise AWS creds are acquired
	//                      by the AWS client from the environment.
	if pAccessKeyIDstr != "" {

		creds  := credentials.NewStaticCredentials(pAccessKeyIDstr, pSecretAccessKeyStr, pTokenStr)
		_, err := creds.Get()

		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to acquire S3 static credentials - (credentials.NewStaticCredentials().Get())",
				"s3_credentials_error", nil, err, "gf_core", pRuntimeSys)
			return nil, gfErr
		}

		config.Credentials = creds
	}
	*/

	//--------------

	sess := session.New(config)

	s3_uploader := s3manager.NewUploader(sess)
	s3_client   := s3.New(sess)

	s3_info := &GFs3Info{
		Client:   s3_client,
		Uploader: s3_uploader,
		Session:  sess,
	}

	return s3_info, nil
}

//---------------------------------------------------
// S3__GENERATE_PRESIGNED_URL

func S3generatePresignedUploadURL(pTargetFileS3pathStr string,
	pS3bucketNameStr string,
	pS3info          *GFs3Info,
	pRuntimeSys      *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	// INPUT
	fileEXTstr     := filepath.Ext(pTargetFileS3pathStr)
	contentTypeStr := mime.TypeByExtension(fileEXTstr)

	s3putObjectParams := &s3.PutObjectInput{
		ACL:         aws.String("public-read"),
		Bucket:      aws.String(pS3bucketNameStr),
		Key:         aws.String(pTargetFileS3pathStr),
		ContentType: aws.String(contentTypeStr),
	}

	req, _ := pS3info.Client.PutObjectRequest(s3putObjectParams)

	// PRESIGN
	presignedURLstr, err := req.Presign(time.Minute * 1)
	if err != nil { // resp is now filled
		gfErr := gf_core.ErrorCreate("failed to generate pre-signed S3 putObject URL",
			"s3_file_upload_url_presign_error", nil, err, "gf_core", pRuntimeSys)
		return "", gfErr
	}

	return presignedURLstr, nil
}

//---------------------------------------------------
// S3__UPLOAD_FILE

func S3putFile(pTargetFileLocalPathStr string,
	pTargetFileS3pathStr string,
	pS3bucketNameStr     string,
	pS3info              *GFs3Info,
	pRuntimeSys          *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	//-----------------
	file, fsErr := os.Open(pTargetFileLocalPathStr)
	if fsErr != nil {
		gfErr := gf_core.ErrorCreate("failed to open a local file to upload it to S3",
			"file_open_error",
			map[string]interface{}{
				"bucket_name_str":             pS3bucketNameStr,
				"target_file__local_path_str": pTargetFileLocalPathStr,
				"target_file__s3_path_str":    pTargetFileS3pathStr,
			},
			fsErr, "gf_core", pRuntimeSys)
		return "", gfErr
	}
	defer file.Close()
	
	//-----------------

	fileInfo, _  := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)

	// read file content to buffer
	file.Read(buffer)

	fileBytes := bytes.NewReader(buffer) // convert to io.ReadSeeker type
	fileType  := http.DetectContentType(buffer)

	// Upload uploads an object to S3, intelligently buffering large files 
	// into smaller chunks and sending them in parallel across multiple goroutines.
	result, s3err := pS3info.Uploader.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Bucket:      aws.String(pS3bucketNameStr),
		Key:         aws.String(pTargetFileS3pathStr),
		ContentType: aws.String(fileType),
		Body:        fileBytes,
	})

	if s3err != nil {
		gfErr := gf_core.ErrorCreate("failed to upload a file to an S3 bucket",
			"s3_file_upload_error",
			map[string]interface{}{
				"bucket_name_str":             pS3bucketNameStr,
				"target_file__local_path_str": pTargetFileLocalPathStr,
				"target_file__s3_path_str":    pTargetFileS3pathStr,
			},
			s3err, "gf_core", pRuntimeSys)
		return "", gfErr
	}

	rStr := fmt.Sprint(result)
	return rStr, nil
}

//---------------------------------------------------
// S3__COPY_FILE
 
func S3copyFile(pSourceBucketStr string,
	pSourceFileS3pathStr string,
	pTargetBucketNameStr string,
	pTargetFileS3pathStr string,
	pS3info              *GFs3Info,
	pRuntimeSys          *gf_core.RuntimeSys) *gf_core.GFerror {

	fmt.Printf("source_bucket        - %s\n", pSourceBucketStr)
	fmt.Printf("source_file__s3_path - %s\n", pSourceFileS3pathStr)
	fmt.Printf("target_bucket_name   - %s\n", pTargetBucketNameStr)
	fmt.Printf("target_file__s3_path - %s\n", pTargetFileS3pathStr)

	source_bucket_and_file__s3_path_str := filepath.Clean(fmt.Sprintf("/%s/%s", pSourceBucketStr, pSourceFileS3pathStr))

	svc   := s3.New(pS3info.Session)
	input := &s3.CopyObjectInput{
		CopySource: aws.String(source_bucket_and_file__s3_path_str),
	    Bucket:     aws.String(pTargetBucketNameStr),
	    Key:        aws.String(pTargetFileS3pathStr),
	}

	result, err := svc.CopyObject(input)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to copy a file within S3",
			"s3_file_copy_error",
			map[string]interface{}{
				"source_bucket_and_file__s3_path_str": source_bucket_and_file__s3_path_str,
				"target_bucket_name_str":              pTargetBucketNameStr,
				"target_file__s3_path_str":            pTargetFileS3pathStr,
			},
			err, "gf_core", pRuntimeSys)
		return gfErr
	}

	fmt.Println(result)

	return nil
}