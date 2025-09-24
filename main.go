package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type ImageAction struct {
	Detail struct {
		Result         string `json:"result"`
		RepositoryName string `json:"repository-name"`
		ImageDigest    string `json:"image-digest"`
		ActionType     string `json:"action-type"`
		ImageTag       string `json:"image-tag"`
	} `json:"detail"`
}

func checkScriptPath(scriptPath string) bool {
	info, err := os.Stat(scriptPath)
	if err != nil {
		fmt.Println("Error unable to stat script file: ", err)
		return false
	}
	if (info.Mode() & 0111) == 0 {
		fmt.Printf("Script at '%s' is not executable\n", scriptPath)
		return false
	}
	if !filepath.IsAbs(scriptPath) && !strings.HasPrefix(scriptPath, "./") {
		fmt.Printf("%s is not a valid script path\n", scriptPath)
		return false
	}
	return true
}

func main() {
	// Parse command line arguments
	queueName := flag.String("queue-name", "", "The name of the SQS queue")
	region := flag.String("region", "us-east-1", "The AWS region to use")
	flag.Parse()

	// Get the non-flag arguments
	args := flag.Args()

	// Check that an external script path is provided
	if len(args) < 1 {
		fmt.Println("External script path not provided")
		os.Exit(1)
	}

	// Get the external script path
	scriptPath := args[0]

	// Ensure that the queue name is provided
	if *queueName == "" {
		fmt.Println("Error: Queue name is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !checkScriptPath(scriptPath) {
		os.Exit(1)
	}

	// Run main process queue function
	if !processQueue(region, queueName, scriptPath) {
		os.Exit(1)
	}
}

func processQueue(region *string, queueName *string, scriptPath string) bool {
	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(*region), // Change this to your desired region
	})

	if err != nil {
		fmt.Println("Error creating session: ", err)
		return false
	}

	// Create a new SQS client
	svc := sqs.New(sess)

	// Retrieve the queue URL using the queue name
	queueURLResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queueName,
	})
	if err != nil {
		fmt.Println("Error getting queue URL: ", err)
		return false
	}

	queueURL := queueURLResult.QueueUrl

	fmt.Println("Using queue URL: ", *queueURL)

	// Create an input for receiving messages
	input := &sqs.ReceiveMessageInput{
		QueueUrl: queueURL,
		MaxNumberOfMessages: aws.Int64(
			10,
		), // Change this to the maximum number of messages you want to receive
		WaitTimeSeconds: aws.Int64(
			20,
		), // Change this to the maximum amount of time to wait for a message
		VisibilityTimeout: aws.Int64(5), // Set the visibility timeout to 5 seconds
	}

	for {
		// Receive messages from the queue
		result, err := svc.ReceiveMessage(input)
		if err != nil {
			fmt.Println("Error receiving message: ", err)
			return false
		}

		// Process each message received
		for _, message := range result.Messages {
			fmt.Println("Received message: ", *message.Body)

			var action ImageAction
			err := json.Unmarshal([]byte(*message.Body), &action)
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				return false
			}

			repositoryName := action.Detail.RepositoryName
			imageTag := action.Detail.ImageTag

			fmt.Printf("Repository Name: %s\n", repositoryName)
			fmt.Printf("Image Tag: %s\n", imageTag)

			// Call the external script for each message received
			cmd := exec.Command(scriptPath, repositoryName, imageTag) //nolint:all
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error executing script:", err)
			} else if cmd.ProcessState.ExitCode() == 0 {
				// Delete the message from the queue if the script exits with a zero exit code
				_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      queueURL,
					ReceiptHandle: message.ReceiptHandle,
				})
				if err != nil {
					fmt.Println("Error deleting message:", err)
				}
			}

			output := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
			for _, line := range output {
				fmt.Printf("Script output: %s\n", line)
			}
		}
	}
}
