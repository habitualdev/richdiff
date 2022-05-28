# RichDiff 

### Easy to use Library for extracting and parsing PE Rich Signatures

Performs easy to use operations on PE Rich Signatures. Supports outputing the signature in a human readable table, as 
as a JSON object (with the option to save it to a file).

Example image output:

<img src="resources/img.png" alt="Example image output" width="100" height="100">


Example usage:


```
package main

import (
    "github.com/roaldi/richdiff"
    "fmt"
)

func main(){

// Read from filesystem
results, err := richdiff.RichFileExtraction("sample.exe")

// richdiff.RichExtraction() alternatively uses []byte as an input

// Sorts the results by product ID
results.Sort()

// Saves json to file
results.WriteToFile("sample.json")

// Prints json to the console
fmt.Println(results.String())

// prints the results in a table
results.RichTable()

// create a png from the DECRYPTED rich signature
img, err := richdiff.RichFileToImage("sample.exe")

// Diff the results with another richdiff results object
changelog, numberOfDiffs, err := results.Diff(richdiff.RichResults{})

fmt.Println(err.Error())

}
```
