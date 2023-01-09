package main

import (
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	// Define a flag "file" with a default value of "input.txt"
	inputFileName := flag.String("file", "input.txt", "input file name")

	// Define a flag "decrypt"
	decrypt := flag.Bool("decrypt", false, "decrypt file")
	flag.Parse()

	inputFile, err := os.Open(*inputFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer inputFile.Close()

	// Create the output file with the same name as the input file, but with the .encrypted extension added
	outputFileName := *inputFileName
	if !*decrypt {
		outputFileName += ".encrypted"
	} else {
		outputFileName = outputFileName[:len(outputFileName)-len(".encrypted")]
	}
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Read in the key
	key := []byte("aaaa567812345678") // <-- Change this key

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stream := cipher.NewCTR(block, []byte("aaaa567812345678")) // <-- Change this key

	// Encrypt or decrypt based on the value of the "decrypt" flag
	if !*decrypt {
		// Encrypt the input file and write it to the output file
		io.Copy(outputFile, cipher.StreamReader{S: stream, R: inputFile})
	} else {
		// Read the contents of the input file into a byte slice
		inputData, err := ioutil.ReadAll(inputFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Decrypt the input data
		stream.XORKeyStream(inputData, inputData)

		// Write the decrypted data to the output file
		_, err = outputFile.Write(inputData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Close the input file
	inputFile.Close()

	// Remove the input file if the "decrypt" flag is not set
	if !*decrypt {
		err = os.Remove(*inputFileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}