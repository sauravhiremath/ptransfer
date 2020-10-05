package ptransfer

// BUFFERSIZE is the minimum packet size for file transfers
const BUFFERSIZE = 1024

// NumberConnections represents number of max connections opened
const NumberConnections = 8

// IV is the key for AES-128
var IV = []byte("1234567812345678")

// ServerAddress represents IPv4 addrress of the sender
var ServerAddress string

// AESKey represents the password string for encryption
var AESKey string

// FileName represents the filename of the data being sent by client
var FileName string
