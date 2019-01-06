package gf_core

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)
//---------------------------------------------------
type Gf_s3_info struct {
	Client   *s3.S3
	Uploader *s3manager.Uploader
	Session  *session.Session
}
//---------------------------------------------------
func S3__init(p_runtime_sys *Runtime_sys) (*Gf_s3_info,*Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_s3.S3__init()")

	//--------------
	// DO NOT PUT credentials in code for production usage!
	// see https://www.socketloop.com/tutorials/golang-setting-up-configure-aws-credentials-with-official-aws-sdk-go
	// on setting creds from environment or loading from file

	//REMOVE!!
	aws_access_key_id     := "AKIAIHKJLFRT3S6BE2UA"
	aws_secret_access_key := "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"
	token                 := ""

	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
	_,err := creds.Get()

	//usr,_    := user.Current()
	//home_dir := usr.HomeDir
	//creds    := credentials.NewSharedCredentials(fmt.Sprintf("%s/.aws/credentials",home_dir),"default")
	//_, err := creds.Get()

	if err != nil {
		gf_err := Error__create("failed to acquire S3 static credentials - (credentials.NewStaticCredentials().Get())",
			"s3_credentials_error",nil,err,"gf_core",p_runtime_sys)
		return nil,gf_err
	}
	//--------------
	
	config := &aws.Config{
		Region          :aws.String("us-east-1"),
		Endpoint        :aws.String("s3.amazonaws.com"),
		S3ForcePathStyle:aws.Bool(true),      // <-- without these lines. All will fail! fork you aws!
		Credentials     :creds,
		//LogLevel        :0, // <-- feel free to crank it up 
	}
	sess := session.New(config)

	s3_uploader := s3manager.NewUploader(sess)
	s3_client   := s3.New(sess)

	s3_info := &Gf_s3_info{
		Client:  s3_client,
		Uploader:s3_uploader,
		Session: sess,
	}

	return s3_info,nil
}
//---------------------------------------------------
func S3__upload_file(p_target_file__local_path_str string,
			p_target_file__s3_path_str string,
			p_s3_bucket_name_str       string,
			p_s3_info                  *Gf_s3_info,
			p_runtime_sys              *Runtime_sys) (string,*Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_s3.S3__upload_file()")
	p_runtime_sys.Log_fun("INFO"     ,"p_s3_bucket_name_str       - "+p_s3_bucket_name_str)
	p_runtime_sys.Log_fun("INFO"     ,"p_target_file__s3_path_str - "+p_target_file__s3_path_str)

	//-----------------
	file, fs_err := os.Open(p_target_file__local_path_str)
	if fs_err != nil {
		gf_err := Error__create("failed to open a local file to upload it to S3",
			"file_open_error",
			&map[string]interface{}{
				"bucket_name_str":            p_s3_bucket_name_str,
				"target_file__local_path_str":p_target_file__local_path_str,
				"target_file__s3_path_str":   p_target_file__s3_path_str,
			},
			fs_err,"gf_core",p_runtime_sys)
		return "",gf_err
	}
	defer file.Close()
	//-----------------

	file_info,_   := file.Stat()
	var size int64 = file_info.Size()

	buffer := make([]byte, size)

	// read file content to buffer
	file.Read(buffer)

	file_bytes := bytes.NewReader(buffer) // convert to io.ReadSeeker type
	file_type  := http.DetectContentType(buffer)


	//Upload uploads an object to S3, intelligently buffering large files 
	//into smaller chunks and sending them in parallel across multiple goroutines.
	result,s3_err := p_s3_info.Uploader.Upload(&s3manager.UploadInput{
						    ACL        : aws.String("public-read"),
						    Bucket     : aws.String(p_s3_bucket_name_str),
						    Key        : aws.String(p_target_file__s3_path_str),
						    ContentType: aws.String(file_type),
						    Body       : file_bytes,
						})

	if s3_err != nil {
		gf_err := Error__create("failed to upload a file to an S3 bucket",
			"s3_file_upload_error",
			&map[string]interface{}{
				"bucket_name_str":            p_s3_bucket_name_str,
				"target_file__local_path_str":p_target_file__local_path_str,
				"target_file__s3_path_str":   p_target_file__s3_path_str,
			},
			s3_err,"gf_core",p_runtime_sys)
		return "",gf_err
	}

	r_str := fmt.Sprint(result)
	return r_str,nil
}
//---------------------------------------------------
func S3__copy_file(p_target_bucket_name_str string,
			p_source_bucket_and_file__s3_path_str string,
			p_target_file__s3_path_str            string,
			p_s3_info                             *Gf_s3_info,
			p_runtime_sys                         *Runtime_sys) *Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_s3.S3__copy_file()")

	fmt.Println("p_target_bucket_name_str              - "+p_target_bucket_name_str)
	fmt.Println("p_source_bucket_and_file__s3_path_str - "+p_source_bucket_and_file__s3_path_str)
	fmt.Println("p_target_file__s3_path_str            - "+p_target_file__s3_path_str)



	svc := s3.New(p_s3_info.Session)
	input := &s3.CopyObjectInput{
	    Bucket:     aws.String(p_target_bucket_name_str),
	    CopySource: aws.String(p_source_bucket_and_file__s3_path_str), //"/sourcebucket/HappyFacejpg"),
	    Key:        aws.String(p_target_file__s3_path_str),            //"HappyFaceCopyjpg"),
	}

	result, err := svc.CopyObject(input)
	if err != nil {
		gf_err := Error__create("failed to copy a file within S3",
			"s3_file_copy_error",
			&map[string]interface{}{
				"target_bucket_name_str":             p_target_bucket_name_str,
				"source_bucket_and_file__s3_path_str":p_source_bucket_and_file__s3_path_str,
				"target_file__s3_path_str":           p_target_file__s3_path_str,
			},
			err,"gf_core",p_runtime_sys)
		return gf_err
	}

	fmt.Println(result)

	return nil
}