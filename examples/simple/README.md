## Run This Example


### Edit the config/example.json
```
$ EDITOR=nvim ENCRYPTION_KEY=BB11898935FC019FFD0AC161E1D8CDB2F570E49538DEB55D5EA049BDC5CAE53B encrypt config/example.json
```

## See the output secret:

```
$ terraform refresh
var.encryptionkey
  Enter a value: BB11898935FC019FFD0AC161E1D8CDB2F570E49538DEB55D5EA049BDC5CAE53B

data.encrypted_file.example: Refreshing state...

Outputs:

password = 5up3R53cr3t
```
