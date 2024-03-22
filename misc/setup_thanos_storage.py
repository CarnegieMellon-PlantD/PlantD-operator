import json
import base64

from kubernetes import client, config
import boto3
import yaml
import argparse


# 1. Create S3 bucket for Thanos
def create_s3_bucket(bucket_name):
    s3_client = boto3.client('s3')
    s3_client.create_bucket(Bucket=bucket_name)

# 2. Create IAM policy to access the bucket
def create_policy(policy_name):
    # Create the IAM policy document
    policy_document = {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Action": [
                    "s3:PutObject",
                    "s3:GetObject",
                    "s3:AbortMultipartUpload",
                    "s3:ListBucket",
                    "s3:DeleteObject",
                    "s3:ListMultipartUploadParts"
                ],
                "Resource": [
                    "arn:aws:s3:::plantd-thanos",
                    "arn:aws:s3:::plantd-thanos/*"
                ]
            }
        ]
    }

    iam = boto3.client('iam')

    try:
        response = iam.create_policy(
            PolicyName=policy_name,
            PolicyDocument=json.dumps(policy_document)
        )
        return response['Policy']['Arn']
    except iam.exceptions.EntityAlreadyExistsException:
        print("Policy {} already exists".format(policy_name))

# 3. Attach iam policy to user
def attach_policy_to_user(policy_arn, username):
    print("Attaching policy {} to user {}".format(policy_arn, username))
    iam = boto3.client('iam')
    iam.attach_user_policy(
        UserName=username,
        PolicyArn=policy_arn
    )
    print(f"Policy attached to user '{username}' successfully.")

# 4. Create Kubernetes secret for Thanos
def create_k8s_secret(access_key, secret_key, bucket_name, namespace):
    thanos_config = {
        'type': 'S3',
        'config': {
            'bucket': bucket_name,
            'endpoint': 's3.us-east-1.amazonaws.com',
            'access_key': access_key,
            'secret_key': secret_key
        }
    }

    thanos_yaml = yaml.dump(thanos_config)
    encoded_secret = base64.b64encode(thanos_yaml.encode('utf-8')).decode()

    config.load_kube_config()
    api = client.CoreV1Api()

    secret = client.V1Secret(
        metadata=client.V1ObjectMeta(name='thanos-objstore-config', namespace=namespace),
        data={'thanos.yaml': encoded_secret},
        type='Opaque'
    )

    # Create the secret in the Kubernetes cluster
    api.create_namespaced_secret(namespace=namespace, body=secret)
    print("Kubernetes secret 'thanos-objstore-config' created successfully.")

# Example usage
if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Create S3 bucket, IAM role, and Kubernetes secret for Thanos")
    parser.add_argument("bucket_name", help="Name of the S3 bucket for Thanos")
    parser.add_argument("access_key", help="AWS access key")
    parser.add_argument("secret_key", help="AWS secret key")
    args = parser.parse_args()

    bucket_name = args.bucket_name
    access_key = args.access_key
    secret_key = args.secret_key

    boto3.setup_default_session(
        aws_access_key_id=access_key,
        aws_secret_access_key=secret_key
    )

    namespace = "plantd-operator-system"
    identity = boto3.client('sts').get_caller_identity()
    username = identity['Arn'].split('/')[-1]
    account_id = identity['Account']
    print(username, account_id)

    create_s3_bucket(bucket_name)
    policy_name = 'plantd-thanos-s3-1'
    policy_arn = create_policy(policy_name)
    if policy_arn is None:
        policy_arn = 'arn:aws:iam::{}:policy/{}'.format(account_id, policy_name)
    attach_policy_to_user(policy_arn, username)
    create_k8s_secret(access_key, secret_key, bucket_name, namespace)
