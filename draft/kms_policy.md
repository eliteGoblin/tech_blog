


It gives the AWS account that owns the KMS key full access to the KMS key.

Unlike other AWS resource policies, a AWS KMS key policy does not automatically give permission to the account or any of its users. To give permission to account administrators, the key policy must include an explicit statement that provides this permission, like this one.

It allows the account to use IAM policies to allow access to the KMS key, in addition to the key policy.

Without this permission, IAM policies that allow access to the key are ineffective, although IAM policies that deny access to the key are still effective.

It reduces the risk of the key becoming unmanageable by giving access control permission to the account administrators, including the account root user, which cannot be deleted.


[Key policies in AWS KMS](https://docs.aws.amazon.com/kms/latest/developerguide/key-policies.html)
A key policy is an resource policy for an AWS KMS key. Key policies are the primary way to control access to KMS keys. Every KMS key must have exactly one key policy. The statements in the key policy determine who has permission to use the KMS key and how they can use it. You can also use IAM policies and grants to control access to the KMS key, but every KMS key must have a key policy.


Unless the key policy explicitly allows it, you cannot use IAM policies to allow access to a KMS key. Without permission from the key policy, IAM policies that allow permissions have no effect. 

The default key policy enables IAM policies

Unlike IAM policies, which are global, key policies are Regional. A key policy controls access only to a KMS key in the same Region. It has no effect on KMS keys in other Regions.

All KMS keys must have a key policy. IAM policies are optional. To use an IAM policy to control access to a KMS key, the key policy for the KMS key must give the account permission to use IAM policies. Specifically, the key policy must include the policy statement that enables IAM policies.

并不是自己理解的: 需要both: 

*  和其他一样, 只要IAM policy 和 resource based policy 有一样， 就行

It gives the AWS account that owns the KMS key full access to the KMS key.

Unlike other AWS resource policies, a AWS KMS key policy does not automatically give permission to the account or any of its users. To give permission to account administrators, the key policy must include an explicit statement that provides this permission, like this one.

```
{
  "Sid": "Enable IAM policies",
  "Effect": "Allow",
  "Principal": {
    "AWS": "arn:aws:iam::111122223333:root"
   },
  "Action": "kms:*",
  "Resource": "*"
}
```

# Test



Create test IAM role: `kms-exports-role`, with empty permission.

Create KMS key with alias: `kms-exports-test`

Test encrypt: 
```
aws kms sign \
--region us-east-1 \
--key-id daba4d8c-85fb-4c1f-902d-a21993e69839 \
--message 'hello world' \
--message-type RAW \
--signing-algorithm RSASSA_PKCS1_V1_5_SHA_256
```

加入AWS profile: 

```
[profile exports-test]
output = json
role_arn = arn:aws:iam::961173933985:role/kms-exports-role
source_profile = dev
region = us-east-1
```

got error: 

```
An error occurred (AccessDeniedException) when calling the Sign operation: User: arn:aws:sts::961173933985:assumed-role/kms-exports-role/botocore-session-1649832014 is not authorized to perform: kms:Sign on resource: arn:aws:kms:us-east-1:961173933985:key/daba4d8c-85fb-4c1f-902d-a21993e69839 because no identity-based policy allows the kms:Sign action
```

## Key Policy test

加入KEY policy: 

```yaml
{
    "Sid": "Allow use of the key",
    "Effect": "Allow",
    "Principal": {
        "AWS": "arn:aws:iam::961173933985:role/kms-exports-role"
    },
    "Action": [
        "kms:DescribeKey",
        "kms:GetPublicKey",
        "kms:Sign",
        "kms:Verify"
    ],
    "Resource": "*"
}
```

成功. remove, 失败

## IAM Identity policy

移除resource based policy:

```
An error occurred (AccessDeniedException) when calling the Sign operation: User: arn:aws:sts::961173933985:assumed-role/kms-exports-role/botocore-session-1649847451 is not authorized to perform: kms:Sign on resource: arn:aws:kms:us-east-1:961173933985:key/daba4d8c-85fb-4c1f-902d-a21993e69839 because no identity-based policy allows the kms:Sign action
```

添加Policy: 

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "kms:DescribeKey",
                "kms:GetPublicKey",
                "kms:Sign",
                "kms:Verify"
            ],
            "Resource": "arn:aws:kms:us-east-1:961173933985:key/daba4d8c-85fb-4c1f-902d-a21993e69839"
        }
    ]
}
```

成功: 
```
{
    "KeyId": "arn:aws:kms:us-east-1:961173933985:key/daba4d8c-85fb-4c1f-902d-a21993e69839",
    "Signature": "LvyVoIcUetOqhM06ilrJnJ+UmJw/e9DALS3zUpNvAzejSlSLPh1tlJMRsXob4ed58eWcgT1fMFUtEJ9UDLiKSoqU6GeiuLvtvxSgXt/MdG4hYQr1b3yf9G6riYGHdqZyeh5aBFFnfBeEBvXYINZz2TgQ91s/x80MoUrSNVlERW+864JqAzUrJaq+c88wiW3gg1PpWSpPxuSau+rGxFqvUfjVI1u1oTntxLWMGS14Hh32Y2DAh4s0VgltFEQzIVFl1yBlWFE2PPv8nT9s0DzIv2bVoAH1uf1zx2nLq8AQGnT7Mps4kSVgrPI2IWi/Hz9396Xsp5K2IlhbCkAQPjI3Mg==",
    "SigningAlgorithm": "RSASSA_PKCS1_V1_5_SHA_256"
}
```