package types

import "time"

type ForwardablePacketData struct {
	Data []byte `json:"data"`
	Memo string `json:"memo"`
}

type IBCEndpoint struct {
	Port    string `json:"port,omitempty"`
	Channel string `json:"channel,omitempty"`
}
type PacketMetadata struct {
	Forward *ForwardMetadata `json:"cargo"`
}

type ForwardMetadata struct {
	Endpoint IBCEndpoint   `json:"endpoint,omitempty"`
	Timeout  time.Duration `json:"timeout,omitempty"`
	Retries  *uint8        `json:"retries,omitempty"`
	Next     *string       `json:"next,omitempty"`
}
