# SEARCHTUBE

Search youtube, but not using the API.

_youtube's way smarter than me, so don't expect this to bypass it's rate limit_

## How to Use


Install it: `go get -u github.com/Pauloo27/searchtube`

```go
package main

import (
  "github.com/Pauloo27/searchtube"
  "fmt"
)

func main() {
  searchTerm := "Tutorial limpar casa"
  maxResults := 5
  results, err := searchtube.Search(searchTerm, maxResults)

  if err != nil {
    panic(err)
  }

  for i, result := range results {
    rank := i+1
    fmt.Printf(" #%d - %s by %s\n", rank, result.Title, result.Uploader)
  }
}
```

## License

This project is licensed under the [MIT license](./LICENSE).

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so.
