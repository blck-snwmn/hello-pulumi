package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var accountID = os.Getenv("CF_ACCOUNT_ID")

func main() {
	ctx := context.Background()
	w, err := auto.NewLocalWorkspace(ctx)
	if err != nil {
		fmt.Printf("Failed to setup and run http server: %v\n", err)
		os.Exit(1)
	}
	err = w.InstallPlugin(ctx, "aws", "v5.39.0")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s\n", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		executeHandler(w, r, true)
	case http.MethodPut:
		executeHandler(w, r, false)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ws, err := auto.NewLocalWorkspace(ctx,
		auto.WorkDir("."),
	)
	if err != nil {
		log.Printf("Failed to setup and run http server: %v\n", err)
		http.Error(w, "Failed to setup and run http server", http.StatusInternalServerError)
		return
	}
	stacks, err := ws.ListStacks(ctx)
	if err != nil {
		log.Printf("Failed to list stacks: %v\n", err)
		http.Error(w, "Failed to list stacks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stacks)
}

func executeHandler(w http.ResponseWriter, r *http.Request, isNew bool) {
	bucketName := r.URL.Query().Get("bn")
	if bucketName == "" {
		http.Error(w, "bucketName is required", http.StatusBadRequest)
		return
	}
	preview := r.URL.Query().Get("p")

	ctx := r.Context()
	program := generateRunFunc(accountID, bucketName)

	var (
		s   auto.Stack
		err error
	)
	if isNew {
		s, err = auto.NewStackInlineSource(ctx, "dev-server", "pulumi-server", program)
	} else {
		s, err = auto.SelectStackInlineSource(ctx, "dev-server", "pulumi-server", program)
	}
	if err != nil {
		log.Printf("Failed to create stack: %v\n", err)
		http.Error(w, "Failed to create stack", http.StatusInternalServerError)
		return
	}
	if preview == "true" {
		res, err := s.Preview(ctx, optpreview.ProgressStreams(os.Stdout))
		if err != nil {
			log.Printf("Failed to update stack: %v\n", err)
			http.Error(w, "Failed to update stack", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(res.ChangeSummary)
		return
	}
	res, err := s.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		log.Printf("Failed to update stack: %v\n", err)
		http.Error(w, "Failed to update stack", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res.Summary)
	return
}

func generateRunFunc(accountID, bucketName string) pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		p, err := aws.NewProvider(ctx, "aws.cloudflare_r2", &aws.ProviderArgs{
			Profile:                   pulumi.String("pulumir2"), // your profile name
			Region:                    pulumi.String("auto"),
			SkipCredentialsValidation: pulumi.Bool(true),
			SkipRegionValidation:      pulumi.Bool(true),
			SkipRequestingAccountId:   pulumi.Bool(true),
			SkipMetadataApiCheck:      pulumi.Bool(true),
			Endpoints: aws.ProviderEndpointArray{
				aws.ProviderEndpointArgs{
					S3: pulumi.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = s3.NewBucketV2(ctx, bucketName, nil, pulumi.Provider(p))
		if err != nil {
			return err
		}
		return nil
	}
}
