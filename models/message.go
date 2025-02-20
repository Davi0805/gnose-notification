package models

import "encoding/json"

type Message struct {
    ID      string `json:"timestamp"`
    Content string `json:"content"`
    CompanyId string `json:companyId`
    UserId string `json:userId`
    Service string `json:service`
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (m *Message) UnmarshalBinary(data []byte) error {
    return json.Unmarshal(data, m)
}