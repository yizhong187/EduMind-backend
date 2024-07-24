

## Installation

As a library

```bash
go get github.com/lpernett/godotenv
```

or if you want to use it as a bin command

go >= 1.17
```shell
go install github.com/lpernett/godotenv/cmd/godotenv@latest
```

go < 1.17
```shell
go get github.com/lpernett/godotenv/cmd/godotenv
```

## Usage

Add your application configuration to your `.env` file in the root of your project:

```shell
S3_BUCKET=YOURS3BUCKET
SECRET_KEY=YOURSECRETKEYGOESHERE
```

Then in your Go app you can do something like

```go
package main

import (
    "log"
    "os"

    "github.com/lpernett/godotenv"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  s3Bucket := os.Getenv("S3_BUCKET")
  secretKey := os.Getenv("SECRET_KEY")

  // now do something with s3 or whatever
}
