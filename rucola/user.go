package main

import (
  "fmt"
  "crypto/rand"
  "encoding/hex"
)

func userExists(client Client) bool {
  _, err := getUser(client.Name)

  if err != nil {
    fmt.Println("false")
    return false
  }
  fmt.Println("true")
  return true
}

func generateId() string {
  randomBytes := make([]byte, 32)
  rand.Read(randomBytes)

  return string(hex.EncodeToString(randomBytes)[:16])
}

func authUser(client Client) (bool, error) {
  user, _ := getUser(client.Name)

  fmt.Println("id: '"+client.Id +"' vs '" + user.Id + "'")

  if user.Id == client.Id {
    return true, nil
  } else {
    return false, nil
  }
}
