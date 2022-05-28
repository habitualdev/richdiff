# RichDiff 

### Easy to use Library for extracting and parsing PE Rich Signatures

Performs easy to use operations on PE Rich Signatures. Supports outputing the signature in a human readable table, as 
as a JSON object (with the option to save it to a file).

TODO: 
 - Add built in support for diffing two or more signatures
 - Ability to render rich signature as a picture for visual comparisons

Example usage:


```
package main

import (
    "github.com/roaldi/richdiff"
    "fmt"
)

func main(){

# Read from filesystem
results, _ := richdiff.RichFileExtraction("sample.exe")

# richdiff.RichExtraction() alternatively uses []byte as an input

# Sorts the results by product ID
results.Sort()

# Saves json to file
results.WriteToFile("sample.json")

# Prints json to the console
fmt.Println(results.String())


# prints the results in a table
results.RichTable()

}
```