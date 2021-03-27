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

	clients := []database.Client{}
	database.Connection.Find(&clients)

	var response []ClientResponse
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

	c.Status(http.StatusOK)
	return c.JSON(response)
}

func CreateClient(c *fiber.Ctx) error {
	var newClient = database.Client{}
	if err := c.BodyParser(newClient); err != nil {
		c.Status(http.StatusBadRequest)
		return nil
	}

	//TODO: Validate IP

	database.Connection.Create(newClient)
	c.Status(http.StatusOK)

	return c.JSON(&ClientResponse{
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
		c.Status(http.StatusBadRequest)
		return nil
	}

	var client = database.Client{}
	database.Connection.Find(&client, id)
	if client.PublicKey == "" {
		c.Status(http.StatusNotFound)
		return nil
	}

	var response = ClientResponse{
		Client:          client,
		LatestHandshake: "2021-01-09T21:15:37.189Z",
		Endpoint:        "77.18.62.145:15427",
		SentBytes:       32165498,
		ReceivedBytes:   132546,
	}

	c.Status(http.StatusOK)
	return c.JSON(response)
}

func DeleteClient(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusBadRequest)
		return nil
	}

	client := database.Client{}
	database.Connection.Find(&client, id)
	if client.PublicKey == "" {
		c.Status(http.StatusNotFound)
		return nil
	}

	database.Connection.Delete(&client)
	c.Status(http.StatusNoContent)
	return nil
}
