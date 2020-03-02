variable "encryptionkey" {
}

provider "encrypted" {
  key = var.encryptionkey
}

data "encrypted_file" "example" {
  path         = "config/example.json"
  content_type = "json"
}

output "password" {
  value = data.encrypted_file.example.parsed["password"]
}
