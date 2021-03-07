package main

type Client struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Ip        string `json:"ip"`
	PublicKey string `json:"publicKey"`
}
type allClients []Client

var clients = allClients{
	{
		Id:        0,
		Name:      "Richard-Workstation",
		Ip:        "10.0.0.1/32",
		PublicKey: "abcd1234",
	},
	{
		Id:        1,
		Name:      "Richard-Laptop",
		Ip:        "10.0.0.2/32",
		PublicKey: "helloWorld",
	},
}

func ClientFind(id int) *Client {
	for _, client := range clients {
		if client.Id == id {
			return &client
		}
	}

	return nil
}

func ClientCreate(newClient Client) Client {
	newClient.Id = len(clients)
	clients = append(clients, newClient)

	return newClient
}
