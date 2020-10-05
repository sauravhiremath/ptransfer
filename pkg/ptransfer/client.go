package ptransfer

import (
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
	log "github.com/sirupsen/logrus"
)

func decryption() []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(FileName))
	if err != nil {
		log.Fatalf("Reading encrypted file: %s", err)
	}
	blockCipher := createCipher()
	stream := cipher.NewCTR(blockCipher, IV)
	stream.XORKeyStream(bytes, bytes)
	return bytes
}

//RunClient does
func RunClient() {
	uiprogress.Start()
	var wg sync.WaitGroup
	wg.Add(NumberConnections)
	bars := make([]*uiprogress.Bar, NumberConnections)
	for id := 0; id < NumberConnections; id++ {
		go func(id int) {
			defer wg.Done()
			port := strconv.Itoa(27001 + id)
			connection, err := net.Dial("tcp", ServerAddress+":"+port)
			if err != nil {
				panic(err)
			}
			defer connection.Close()

			bufferFileName := make([]byte, 64)
			bufferFileSize := make([]byte, 10)

			connection.Read(bufferFileSize)
			fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
			bars[id] = uiprogress.AddBar(int(fileSize+1028) / 1024).AppendCompleted().PrependElapsed()

			connection.Read(bufferFileName)
			FileName = strings.Trim(string(bufferFileName), ":")
			os.Remove(FileName + "." + strconv.Itoa(id))
			newFile, err := os.Create(FileName + "." + strconv.Itoa(id))
			if err != nil {
				panic(err)
			}
			defer newFile.Close()

			var receivedBytes int64
			for {
				if (fileSize - receivedBytes) < BUFFERSIZE {
					io.CopyN(newFile, connection, (fileSize - receivedBytes))
					// Empty the reng bytes that we don't need from the network buffer
					connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
					break
				}
				io.CopyN(newFile, connection, BUFFERSIZE)
				//Increment the counter
				receivedBytes += BUFFERSIZE
				bars[id].Incr()
			}
		}(id)
	}
	wg.Wait()

	// cat the file
	os.Remove(FileName)
	finished, err := os.Create(FileName)
	defer finished.Close()
	if err != nil {
		log.Fatal(err)
	}
	for id := 0; id < NumberConnections; id++ {
		fh, err := os.Open(FileName + "." + strconv.Itoa(id))
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(finished, fh)
		if err != nil {
			log.Fatal(err)
		}
		fh.Close()
		os.Remove(FileName + "." + strconv.Itoa(id))
	}

	// decryption
	ioutil.WriteFile(FileName+"Decrypted", decryption(), 0644)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("\n\n\nDownloaded " + FileName + "!")
	time.Sleep(1 * time.Second)
}
