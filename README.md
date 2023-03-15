# Whatsapp + ChatGPT integrations
This project was made using a GoLang. 
It is a serverless project using aws lambda to send ChatGPT questions using a 
twilio bot to receive the questions.
Serverless Framework is responsible for the infra.


> Obs: With you do not have knowledge with this framework, I suggest you see the [documentation](https://www.serverless.com/framework/docs) before usage.

To use this project you will need build the dockerfile image and push then on your aws account and change the repository name on serverless.yml file.

You can use this commands bellow, to make it:
```shell
-- to get your aws credentials, use this:
aws ecr get-login-password --region <ACCOUNT REGION> | docker login --username AWS --password-stdin <ACCOUNT ID>.dkr.ecr.<ACCOUNT REGION>.amazonaws.com

-- to build your dockerfile with a ECR repo name, use this: 
docker build -t <ACCOUNT ID>.dkr.ecr.<ACCOUNT REGION>.amazonaws.com/<IMAGE NAME>:<IMAGE VERSION> .

-- to send your image to ECR, use this:
docker push <ACCOUNT ID>.dkr.ecr.<ACCOUNT REGION>.amazonaws.com/<IMAGE NAME>:<IMAGE VERSION>
```

With your ECR configured and repository name is aponted on serverless.yaml on `functions.chat.image` you can deploy the application using de serverless deploy command.

```shell
serverless deploy
```

