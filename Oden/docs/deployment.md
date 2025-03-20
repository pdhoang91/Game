# Oden - Deployment Guide

This document describes how to deploy the Oden Idle RPG game to production environments.

## Client Deployment (Unity iOS)

### Prerequisites

- Unity 2021.3 LTS or newer
- Xcode 13 or newer
- Apple Developer Account
- MacOS computer for iOS builds

### Building the iOS App

1. Open the Unity project from the `client` directory.

2. Configure the project for iOS:
   - Go to `File > Build Settings`
   - Select `iOS` platform
   - Click `Switch Platform` if not already selected

3. Configure player settings:
   - Click `Player Settings`
   - Set `Product Name` to "Oden"
   - Set `Bundle Identifier` to "com.yourcompany.oden"
   - Set minimum iOS version to iOS 13.0
   - Configure other settings as needed (icons, splash screen, etc.)

4. Set the API endpoint:
   - Open `Assets/Scripts/Config/GameConfig.cs`
   - Set `ApiEndpoint` to your production API URL

5. Build the project:
   - In Build Settings, click `Build`
   - Choose a location to save the Xcode project
   - Wait for Unity to generate the Xcode project

6. Open the Xcode project:
   - Open the generated `.xcodeproj` file in Xcode
   - Select your Development Team in `Signing & Capabilities`
   - Configure any additional iOS settings

7. Build and Archive:
   - Connect an iOS device or use a simulator
   - Select `Product > Archive` to create a distributable build
   - Follow App Store Connect instructions to upload the build

## Server Deployment (Golang on AWS)

### Prerequisites

- AWS Account
- AWS CLI configured
- Terraform (optional, for infrastructure as code)
- Domain name (optional, for production)

### AWS Infrastructure Setup

1. Create a VPC:
   - Log in to AWS Console
   - Go to VPC service
   - Create a new VPC with at least 2 public subnets across different AZs
   - Create an Internet Gateway and attach it to the VPC
   - Configure route tables to allow internet access

2. Set up RDS Database:
   - Go to RDS service
   - Create a new MySQL database instance
   - Select MySQL 8.0 or compatible version
   - Choose `dev/test` or `production` template based on your needs
   - Configure instance size (t3.micro for MVP is sufficient)
   - Set up storage (20GB is usually enough for MVP)
   - Configure VPC, subnets, and security groups
   - Create a database named `oden`
   - Note the endpoint, username, and password

3. Create S3 Bucket for Assets:
   - Go to S3 service
   - Create a new bucket named `oden-assets`
   - Configure appropriate permissions (public read for assets)
   - Set up CORS configuration if needed

4. Set up EC2 Instance:
   - Go to EC2 service
   - Launch a new t3.micro instance (suitable for MVP)
   - Select Amazon Linux 2 or Ubuntu Server as the AMI
   - Configure VPC, subnet, and security group
   - Add storage (8GB is sufficient for MVP)
   - Create or select a key pair for SSH access
   - Launch the instance and note its public IP

5. Configure Security Groups:
   - API Server: Allow inbound traffic on port 80 (HTTP), 443 (HTTPS), and 22 (SSH)
   - Database: Allow inbound connections from the API server security group

### Database Initialization

1. Connect to your RDS instance from your local machine or EC2 instance:
   ```bash
   mysql -h <rds-endpoint> -P 3306 -u <username> -p
   ```

2. Create the necessary database tables:
   - Use the SQL scripts in `server/internal/db/migrations` to set up the schema
   - For example:
     ```bash
     mysql -h <rds-endpoint> -P 3306 -u <username> -p oden < server/internal/db/migrations/001_initial_schema.sql
     ```

### Server Deployment

1. SSH into your EC2 instance:
   ```bash
   ssh -i <your-key.pem> ec2-user@<ec2-public-ip>
   ```

2. Install required packages:
   ```bash
   # On Amazon Linux
   sudo yum update -y
   sudo yum install -y golang git
   
   # On Ubuntu
   sudo apt update
   sudo apt install -y golang git
   ```

3. Clone the repository:
   ```bash
   git clone https://your-repository-url/oden.git
   cd oden/server
   ```

4. Configure the server:
   - Create a configuration file:
     ```bash
     cp internal/config/config.example.json internal/config/config.json
     ```
   - Edit the configuration with your production values:
     ```bash
     vim internal/config/config.json
     ```
   - Set the database connection, JWT secret, and other parameters

5. Build the application:
   ```bash
   go build -o oden-server ./cmd/api
   ```

6. Set up a systemd service for automatic restart:
   ```bash
   sudo vim /etc/systemd/system/oden.service
   ```
   
   With contents:
   ```
   [Unit]
   Description=Oden Game Server
   After=network.target
   
   [Service]
   User=ec2-user
   WorkingDirectory=/home/ec2-user/oden/server
   ExecStart=/home/ec2-user/oden/server/oden-server
   Restart=always
   
   [Install]
   WantedBy=multi-user.target
   ```

7. Start the service:
   ```bash
   sudo systemctl enable oden
   sudo systemctl start oden
   ```

8. (Optional) Set up Nginx as reverse proxy:
   ```bash
   sudo yum install -y nginx
   sudo vim /etc/nginx/conf.d/oden.conf
   ```
   
   With contents:
   ```
   server {
       listen 80;
       server_name api.oden-game.com;
   
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

9. (Optional) Set up SSL with Let's Encrypt:
   ```bash
   sudo amazon-linux-extras install epel
   sudo yum install -y certbot python2-certbot-nginx
   sudo certbot --nginx -d api.oden-game.com
   ```

## Continuous Deployment (Optional)

For more advanced deployment, consider setting up:

1. GitHub Actions or AWS CodePipeline for CI/CD
2. AWS Elastic Beanstalk for easier application deployment
3. AWS ECS/EKS for containerized deployment
4. AWS CloudFormation or Terraform for infrastructure as code

## Monitoring and Logging

1. Set up CloudWatch for monitoring:
   - Create alarms for server CPU, memory, and disk usage
   - Set up log groups for application logs

2. Configure application logging:
   - Update the Golang application to send logs to CloudWatch
   - Set appropriate log levels for production

## Troubleshooting

- Check server logs:
  ```bash
  sudo journalctl -u oden
  ```

- Check Nginx logs:
  ```bash
  sudo tail -f /var/log/nginx/error.log
  ```

- Test database connection:
  ```bash
  mysql -h <rds-endpoint> -P 3306 -u <username> -p oden
  ```

- Check server status:
  ```bash
  sudo systemctl status oden
  ``` 