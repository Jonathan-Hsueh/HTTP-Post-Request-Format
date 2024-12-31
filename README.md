This file just simply contains the contents of a couple potentially complicated topics into a nice place for students and people to learn from! 

A sample template for an HTTP Post request, this file has the basics to send a message in the form of a JSON to an HTTP link across the internet including Authenticated Hashing, all written in Golang. 
Here is a run down of the topics: 
  1. Marshalling: creating a JSON object of the data we wish to send, called a payload, and then marshaling the data into neat byte slices for the internet to read with clarity
  2. Generating TOTP: creating a Time-based one-time password (TOTP) authentication for a one-time use safe access passcode (using HMAC-SHA-512 Hash, a way to authenticate data)
  3. Passing and sending the Request: authenticate the data through headers and then send the created request object for varification

I hope this repo is helpful in understanding the mechanics of a post request!
