package ptransfer

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const encryptedFile = "temp_encypt_tbd"

// RunServer does
func RunServer() {
	logger := log.WithFields(log.Fields{
		"function": "main",
	})
	logger.Info("Initializing")
	var wg sync.WaitGroup
	wg.Add(NumberConnections)
	for id := 0; id < NumberConnections; id++ {
		go listenerThread(id, &wg)
	}
	wg.Wait()
}

func listenerThread(id int, wg *sync.WaitGroup) {
	logger := log.WithFields(log.Fields{
		"function": "listenerThread@" + ServerAddress + ":" + strconv.Itoa(27000+id),
	})

	defer wg.Done()

	err := listener(id)
	if err != nil {
		logger.Error(err)
	}
}

func listener(id int) (err error) {
	port := strconv.Itoa(27001 + id)
	logger := log.WithFields(log.Fields{
		"function": "listener@" + ServerAddress + ":" + port,
	})
	server, err := net.Listen("tcp", ServerAddress+":"+port)
	if err != nil {
		return errors.Wrap(err, "Error listening on "+ServerAddress+":"+port)
	}
	defer server.Close()
	logger.Info("waiting for connections")
	//Spawn a new goroutine whenever a client connects
	for {
		connection, err := server.Accept()
		if err != nil {
			return errors.Wrap(err, "problem accepting connection")
		}
		logger.Info("Client connected")
		go sendFileToClient(id, connection)
	}
}

//This function is to 'fill'
func fillString(retunString string, toLength int) string {
	for {
		lengthString := len(retunString)
		if lengthString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func createCipher() cipher.Block {
	c, err := aes.NewCipher([]byte(AESKey))
	if err != nil {
		log.Fatalf("Failed to create the AES cipher: %s", err)
	}
	return c
}

func encryption(plainText string) {
	bytes := []byte(plainText)
	blockCipher := createCipher()
	stream := cipher.NewCTR(blockCipher, IV)
	stream.XORKeyStream(bytes, bytes)
	err := ioutil.WriteFile(fmt.Sprintf(encryptedFile), bytes, 0644)
	if err != nil {
		log.Fatalf("Writing encryption file: %s", err)
	} else {
		fmt.Printf("Message encrypted in file: %s\n\n", encryptedFile)
	}
}

func sendFileToClient(id int, connection net.Conn) {
	logger := log.WithFields(log.Fields{
		"function": "sendFileToClient #" + strconv.Itoa(id),
	})
	defer connection.Close()

	// Open File
	var plainText string
	var plainChunk = make([]byte, 16)
	fileTemp, err := os.Open(FileName)
	for {
		_, err = fileTemp.Read(plainChunk)
		if err == io.EOF {
			//End of file reached, break out of for loop
			logger.Info("EOF")
			break
		}
		plainText += string(plainChunk)
		logger.Info("1st Step done\n")
	}

	// AES Encrptptppt
	encryption(plainText)

	//Open the file that needs to be send to the client
	file, err := os.Open(encryptedFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	//Get the Filename and filesize
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	numChunks := math.Ceil(float64(fileInfo.Size()) / float64(BUFFERSIZE))
	chunksPerWorker := int(math.Ceil(numChunks / float64(NumberConnections)))

	bytesPerConnection := int64(chunksPerWorker * BUFFERSIZE)
	if id+1 == NumberConnections {
		bytesPerConnection = fileInfo.Size() - (NumberConnections-1)*bytesPerConnection
	}
	fileSize := fillString(strconv.FormatInt(int64(bytesPerConnection), 10), 10)

	FileName := fillString(fileInfo.Name(), 64)

	if id == 0 || id == NumberConnections-1 {
		logger.Infof("numChunks: %v", numChunks)
		logger.Infof("chunksPerWorker: %v", chunksPerWorker)
		logger.Infof("bytesPerConnection: %v", bytesPerConnection)
		logger.Infof("FileName: %v", fileInfo.Name())
	}

	logger.Info("sending")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(FileName))
	sendBuffer := make([]byte, BUFFERSIZE)

	chunkI := 0
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			//End of file reached, break out of for loop
			logger.Info("EOF")
			break
		}
		if (chunkI >= chunksPerWorker*id && chunkI < chunksPerWorker*id+chunksPerWorker) || (id == NumberConnections-1 && chunkI >= chunksPerWorker*id) {
			connection.Write(sendBuffer)
		}
		chunkI++
	}
	fmt.Println("File has been sent, closing connection!")
	return
}
