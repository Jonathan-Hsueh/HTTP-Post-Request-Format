package main // Package that allows for main and standard programs

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Define the string variables that contain my information
	email := ""
	info1 := ""
	info2 := ""

	// Create a Json Payload which is a string representation of data in a JSON format and will be the body of the Post Request
	jsonPayload := map[string]string{ // Create a map object to represent Key-Value Pairs (which means a type where both keys and values are strings (kind of a fancy struct)) that will easily translate to JSON format
		"email":      email,
		"json_info1": info1,
		"json_info2": info2,
	}

	// json.marshal converts a map into json format 'byte slice' (a single line of binary code representing the data) : variable payloadBytes
	bytesPayload, err := json.Marshal(jsonPayload)
	if err != nil { // Nil is the 0 value or null value for many types
		fmt.Println("Error marshalling Payload to JSON")
		return
	}

	// Generate the TOTP Password (use an email object)
	secretPass := email + ""
	totpPass, err := makeTOTP(secretPass)
	if err != nil {
		fmt.Println("Error generating TOTP Password")
		return
	}

	// Pass the JSON payload (data transmitted over the internet) into the HTTP Post request
	url := "https://your.url.to.send"
	httpReq, err := http.NewRequest("POST (or) GET", url, bytes.NewBuffer(bytesPayload)) // bytes.NewBuffer converts the payload into data that can be sent across the HTTP network
	if err != nil {
		fmt.Println("Error creating POST request")
		return
	}

	// Add the Headers
	httpReq.Header.Set("Content-Type", "application/json") // Change the Content-Type header
	httpReq.SetBasicAuth(email, totpPass)

	// Send the Request
	client := &http.Client{}
	resp, err := client.Do(httpReq) // client.Do executes the post request
	if err != nil {
		fmt.Println("Error sending POST request")
		return
	}

	// Handling the Request Response
	defer resp.Body.Close() // closes the response body after reading
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Submission passed")
	} else {
		fmt.Println("Submission didn't pass, and here is the status code: %d\n", resp.StatusCode)
	}
}

// Function that generates the TOTP value, functins into can return two values (in this case a result and an error message)
func makeTOTP(secretPass string) (string, error) {
	// Get the current unix time in whatever intervals you like! example is 60 second intervals
	interval := time.Now().Unix() / 60

	// Convert the time interval into typically 8 byte slices
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, uint64(interval)) // big-endian format is most significant digit first

	// Generate something called a HMAC-SHA-512 Hash which is used to authenticate data and secure communications
	// It creates a hash based message authentication code
	key := []byte(secretPass)      // create an array of bytes using secretPass to convert into a byte slice
	h := hmac.New(sha512.New, key) // creates a new hmac object with the sha512 function using our byte slice
	h.Write(buffer)                // hash the time interval
	hash := h.Sum(nil)             // finalize the hash

	// Extract the Dynamic Offset (Which is how many values we want to extract) and is determined by the last byte of the hash
	offset := hash[len(hash)-1] & 0x0F

	// Extract the 4 byte binary code which is taken from the hash starting at the offsets
	code := (int64(hash[offset])&0x7F)<<24 |
		(int64(hash[offset+1])&0xFF)<<16 |
		(int64(hash[offset+2])&0xFF)<<8 |
		(int64(hash[offset+3]) & 0xFF)

	totp := code % 10000000000 // ensures it is a ten digit number

	// Returns the totp formatted as a zero-padded string
	return fmt.Sprintf("%010d", totp), nil
}
