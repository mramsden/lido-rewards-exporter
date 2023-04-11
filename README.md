# Deployment

## Code

After packaging the code using `make package` then run the following to deploy:

```
aws s3 cp build/main.zip s3://lido-rewards-exporter-source/
aws lambda update-function-code --function-name lido-rewards-exporter --s3-bucket lido-rewards-exporter-source --s3-key main.zip
```

## Infrastructure

Run the following to update the infrastructure in AWS:

```
aws cloudformation deploy --stack-name lido-rewards-exporter --template-file infrastructure/app.yml --capabilities CAPABILITY_NAMED_IAM
```
