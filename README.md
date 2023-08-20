# GymShark Packer

This application can be used as an HTTP API to calculate the number of packs needed to ship a given number of items. The API is flexible and can be easily extended to add or remove pack sizes.

To run the application, you can use the following commands:

`go run main.go`

or

`make run`.

This will start the application on port 3000. You can then use `curl` or `Postman` to make requests to the API.

For example, to get the number of packs needed to ship 1000 items, you would use the following curl command:

`curl -X POST http://localhost:3000/api/packer/pack -d '{"items": 1000}'`

The response will be the number of packs of certain capacity, which in this case is 1.
Response examples:

- items = 1:   {"packs": {"250": 1}}
- items = 250: {"packs": {"250": 1}}
- items = 251: {"packs": {?}}
- items = 500: {"packs": {"500": 1}}
- items = 501: {"packs": {"500": 1, "250": 1}}
- items = 12001: {"packs": {"5000": 2, "2000": 1 "250": 1}}
