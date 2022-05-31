# Terraform Provider Encrypted

## Encrypt Your Terraform Secrets Easily

Terraform-provider-encrypted lets you manage encrypted files to store your Terraform secrets. (A bit like [Chef's databags](https://docs.chef.io/data_bags/#encrypt-a-data-bag-item)).

## Usage

### Installation

```
go install github.com/Scalingo/terraform-provider-encrypted
```

### Generate an Encryption Key

To generate a new encryption key. You can use OpenSSL CLI:

```
openssl rand -hex 32 | tr '[a-z]' '[A-Z]'
```

The key must be 16, 24 or 32 bytes stored in hex format.

### Configuration

To configure the provider you must give it an encryption key or use the ENCRYPTION_KEY environment variable.

```
provider "encrypted" {
  key = var.encryptionkey
}
```

### Datasource

This provider exposes a single resource: `encrypted_file`.

The resource can be used like this:


```
data "encrypted_file" "my_secret" {
  path         = "secrets/aws.json"
  content_type = "json"
}
```

This resource accepts the following parameters:

* path: Path to the encrypted file
* content_type: Format of the specified file (For now only "`json`" and "`yaml`" is supported)
* data_path: path to the root element

And has the following outputs:

* parsed: if the content is a map, the decrypted content will be there
* array: if the content is an array, the decrypted content will be there
* value: the decrypted file content

### Utility

To create an encrypted file you can use our tool present in the cmd/encrypt package.

#### Installation

```
go install github.com/Scalingo/terraform-provider-encrypted/cmd/encrypt
```

Configuration:

The CLI is configured using environment variables:

* ENCRYPTION_KEY (required): Value of the encryption key
* EDITOR (optional, default: 'nvim') program used to edit the encrypted files
* EDITOR_ARGUMENTS (optional, default: '') program arguments used to edit the encrypted files

Usage:

```
encrypt path/to/my/file.json
```

### Examples

A simple example can be found in the examples/simple directory.
