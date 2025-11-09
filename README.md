# excel

Go package for importing Excel and CSV files into structs with validation and filtering support.

# Import File

Open the file and call the import method:

```go
file, err := os.Open("products.csv")
if err != nil {
    panic(err)
}
defer file.Close()

var products []Product
result, err := excel.Import(file, &products)
if err != nil {
    panic(err)
}

if len(result.ValidationErrors) > 0 {
    for _, err := range result.ValidationErrors {
        fmt.Printf("Row %d: %s\n", err.Row, err.Message)
    }
}
```

## Config struct using tags

Define your struct and attach mapping + validation tags:

```go
type Product struct {
    Name  string  `excel:"product_name" validate:"required|minLen:3"`
    SKU   string  `excel:"SKU" validate:"required|len:8|alphaNum" message:"required:SKU is required|len:SKU must be 8 chars"`
    Price float64 `excel:"Price"`
}
```

## Config using methods

If you prefer methods to tags or need dynamic config, implement:

```go
func (p *Product) Rules() map[string]string {
    return map[string]string{
        "Name": "required|minLen:3",
        "SKU":  "required|len:8|alphaNum",
    }
}

func (p *Product) Filters() map[string]string {
    return map[string]string{
        "Name": "trim|upperFirst",
        "SKU":  "trim|upper",
    }
}

func (p *Product) Messages() map[string]string {
    return map[string]string{
        "Name.required": "Product name is mandatory",
        "SKU.len":       "SKU must be exactly 8 characters",
    }
}

func (p *Product) Attributes() map[string]string {
    return map[string]string{
        "Name": "Product Name",
        "SKU":  "Stock Keeping Unit",
    }
}
```