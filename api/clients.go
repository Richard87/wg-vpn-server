package api

import (
	"github.com/Richard87/wg-vpn-server/database"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type ClientResponse struct {
	database.Client
	LatestHandshake string `json:"latestHandshake"`
	Endpoint        string `json:"endpoint"`
	SentBytes       int64  `json:"sentBytes"`
	ReceivedBytes   int64  `json:"receivedBytes"`
}

func GetClients(c *fiber.Ctx) error {

	clients := make([]database.Client, 100)
	database.Connection.Find(&clients)

	response := make([]ClientResponse, 100)
	for _, c := range clients {
		res := ClientResponse{
			Client:          c,
			LatestHandshake: "2021-01-09T21:15:37.189Z",
			Endpoint:        "77.18.62.145:15427",
			SentBytes:       32165498,
			ReceivedBytes:   132546,
		}
		response = append(response, res)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func CreateClient(c *fiber.Ctx) error {
	var newClient = database.Client{}
	if err := c.BodyParser(newClient); err != nil {
		return c.Status(http.StatusBadRequest).Format("Bad request")
	}

	//TODO: Validate IP

	database.Connection.Create(newClient)

	return c.Status(http.StatusOK).JSON(&ClientResponse{
		Client:          newClient,
		LatestHandshake: "2021-01-09T21:15:37.189Z",
		Endpoint:        "77.18.62.145:15427",
		SentBytes:       32165498,
		ReceivedBytes:   132546,
	})
}

func GetClient(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).Format("Bad request")
	}

	client := new(database.Client)
	database.Connection.Find(client, id)
	if client.PublicKey == "" {
		return c.Status(http.StatusNotFound).Format("Not found")
	}

	var response = ClientResponse{
		Client:          *client,
		LatestHandshake: "2021-01-09T21:15:37.189Z",
		Endpoint:        "77.18.62.145:15427",
		SentBytes:       32165498,
		ReceivedBytes:   132546,
	}

	return c.Status(http.StatusOK).JSON(response)
}

func DeleteClient(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).Format("Bad request")
	}

	client := database.Client{}
	database.Connection.Find(&client, id)
	if client.PublicKey == "" {
		return c.Status(http.StatusNotFound).Format("Not found")
	}

	database.Connection.Delete(&client)
	return c.Status(http.StatusNoContent).JSON(nil)
}
