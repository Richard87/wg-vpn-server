package main

import (
	"encoding/binary"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"log"
)

type Client struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Ip        string `json:"ip"`
	PublicKey string `json:"publicKey"`
}
type allClients []Client

func ClientList() *allClients {
	var result = allClients{}

	_ = Db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("clients"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var client = &Client{}
			err := json.Unmarshal(v, client)
			if err != nil {
				log.Print(err)
				continue
			}

			result = append(result, *client)
		}

		return nil
	})

	return &result
}

func ClientFind(id int) *Client {
	var client *Client = nil

	_ = Db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("clients"))
		result := b.Get(itob(id))
		if result != nil {
			err := json.Unmarshal(result, &client)
			if err != nil {
				log.Print(err)
			}
		}

		return nil
	})

	return client
}

func ClientRemove(id int) {

	_ = Db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("clients"))
		err := b.Delete(itob(id))
		if err != nil {
			log.Print(err)
		}

		return nil
	})
}

func ClientCreate(newClient *Client) *Client {
	err := Db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte("clients"))
		id, _ := bucket.NextSequence()
		newClient.Id = int(id)

		// Marshal user data into bytes.
		jsonBuffer, err := json.Marshal(*newClient)
		if err != nil {
			log.Fatal(err)
		}

		return bucket.Put(itob(newClient.Id), jsonBuffer)
	})
	if err != nil {
		log.Fatal(err)
	}

	return newClient
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
