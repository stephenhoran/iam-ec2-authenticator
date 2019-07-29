# **IAM-EC2-Authenticator**
IAM-EC2-Authenticator is a simple tool used to create and manage ssh-keys on 
EC2 (currently only Linux is support) Servers. No IAM credentials are stored on
the server itself. Once a group is configured in IAM, iam-ec2-authenticator will
IAM periodically to determine if any new users have been added to that group. If 
a new user it present it will simple create a new user an add their ssh public
key automatically. 