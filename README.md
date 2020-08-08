# urldescribe

Go library to describe URLs. Essentially just returns whatever is in the <title> tag.

## Usage

```go
import (
    "fmt"

    "github.com/kari/urldescribe"
)

func main() {
    fmt.Println(resp := DescribeURL("https://github.com/kari/urldescribe"))
    // GitHub - kari/urldescribe: Go library to describe URLs
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[Apache License 2.0](https://choosealicense.com/licenses/apache-2.0/)
