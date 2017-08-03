
Create a cluster
http://docs.aws.amazon.com/AmazonECS/latest/developerguide/create_cluster.html


aws logs create-log-group --log-group-name awslogs-wordpress --region ap-northeast-1

Managing secrets
https://aws.amazon.com/blogs/security/how-to-manage-secrets-for-amazon-ec2-container-service-based-applications-by-using-amazon-s3-and-docker/


http://docs.aws.amazon.com/AmazonECS/latest/developerguide/cmd-ecs-cli-compose.html

ecs-cli compose [--verbose] [--file compose-file] [--project-name project-name] [--task-role-arn role_value] [--cluster cluster_name] [--region region] [subcommand] [arguments] 

ecs-cli compose --verbose --file docker-compose.yml  --project-name dev-lvl-api --task-role-arn arn:aws:iam::365769272576:role/ecsInstanceRole --cluster dev-leveledup-api --region us-west-2 service --profile leveledup