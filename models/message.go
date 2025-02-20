package models

import "encoding/json"

type Message struct {
    ID      string `json:"id"`
    Content string `json:"content"`
    CompanyId string `json:companyId`
    UserId string `json:userId`
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (m *Message) UnmarshalBinary(data []byte) error {
    return json.Unmarshal(data, m)
}