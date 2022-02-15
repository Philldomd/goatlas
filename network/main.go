package network

import (
  "io/ioutil"
  "log"
  "net/http"
)

type Network struct {
  Header http.Header
  Body []byte
}

func NewNetwork() *Network {
  resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
  if err != nil {
    log.Fatalln(err)
  }

  b, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatalln(err)
  }

  net := Network{
    Header: resp.Header,
    Body: b,
  }
  return &net
}
