Sample of creating a Cloudflare R2 bucket using pulumi

# Setup
1. [Generate an Access Key](https://developers.cloudflare.com/r2/api/s3/tokens/)
2. Configure your AWS credentials. 
    Here we are creating a profile called `pulumir2`

```
$ aws configure --profile pulumir2
AWS Access Key ID [None]: <Your R2 Access Key ID>
AWS Secret Access Key [None]: <Your R2 Secret Access Key>
Default region name [None]: auto
Default output format [None]: 
```

# Deploy
```
$ CF_ACCOUNT_ID=<YOUR_ACCOUNT_ID> pulumi up -y 
```
