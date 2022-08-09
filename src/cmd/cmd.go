package cmd

import (
	"log"
	"os"

	"github.com/jedzeins/image-board/src/batch"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seed",
	Short: "insert seed data (doggo pics) into s3",
	Long: `I built an image reverse proxy with nginx and s3 (localstack) and 
			it can be best demonstrated by having images in s3 on app start, which
			is why I'm writing a command so that it runs on init`,
	Run: func(cmd *cobra.Command, args []string) {

		localStackUrl := os.Getenv("S3_DOMAIN")
		BatchProgram, err := batch.NewBatch(os.Getenv("S3_BUCKET_NAME"), localStackUrl, &aws.Config{
			Region: aws.String(os.Getenv("S3_DEFAULT_REGION")),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("S3_USERNAME"),
				os.Getenv("S3_PASSWORD"), ""),
			S3ForcePathStyle: aws.Bool(true),
			Endpoint:         aws.String(localStackUrl),
		})
		if err != nil {
			log.Fatalf("error in batch init %v\n", err)
		}

		err = BatchProgram.CheckLocalStack()
		if err != nil {
			log.Fatalf("rootCmd.Run: %v\n", err)
		}

		files, err := BatchProgram.GetFilePaths("./images")
		if err != nil {
			log.Fatalf("rootCmd.Run: %v\n", err)
		}

		err = BatchProgram.PreloadImages(files)
		if err != nil {
			log.Fatalf("rootCmd.Run: %v\n", err)
		}

		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Execute: %v\n", err)
	}
}
