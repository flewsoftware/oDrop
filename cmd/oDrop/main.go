package main

import (
	"flag"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net"
	"oDrop/core"
	"oDrop/utils"
	"oDrop/utils/speedwrap"
	"os"
	"strconv"
	"time"
)

func main() {
	var mode string
	var file string
	var n int

	var useLessCpuTimeExtractor = flag.Bool("lessCpuTime", false, "uses less cpu time whenever possible")
	var testingMode = flag.Bool("testingMode", false, "enables testing mode")
	mode, file, n = PromptUser()
	flag.Parse()

	if utils.ModeToSimple(mode) == "s" {

		// generates a random number to verify receiver
		r := utils.GetRandomNumber()
		fmt.Printf("Passcode is %s\n", r)

		// broadcasts the file
		fmt.Println("Waiting for connections")
		err := core.Send(core.SendDataCallback{
			SentCallback: func(c net.Conn) {
				fmt.Printf("Sent file to %v", c.RemoteAddr())
				os.Exit(0)
			},
			DataBroker: func(c net.Conn, reader io.Reader, size int64) {
				s := speedwrap.SW{}
				bar := progressbar.DefaultBytes(size, "Sending")
				bar.Describe("Sending")
				s.SetStartTime()
				io.Copy(io.MultiWriter(io.MultiWriter(c, bar), &s), reader)
				fmt.Printf("finished sending file with a avg speed of %d B/s", s.GetSpeedRound())
			},
		}, file, r)
		if err != nil {
			log.Fatalf("cant send file %v", err)
		}

	} else {
		var ip, port = "", ""
		if *testingMode {
			fmt.Println("testing mode")
			ip = "localhost"
			port = "6780"
		}
		err := core.Receive(file, strconv.Itoa(n), func(d io.Reader, f io.Writer, sizeD []byte) {
			size, err := strconv.Atoi(string(sizeD))
			if err != nil {
				log.Fatalf("cant get file size: %v\n", err)
			}

			bar := progressbar.DefaultBytes(
				int64(size),
				"downloading",
			)

			var downloadStartTime = time.Now()

			// copy contents of data to the file
			wb, err := io.Copy(io.MultiWriter(f, bar), d)
			if err != nil {
				log.Fatalln(err)
			}

			var endTime = time.Since(downloadStartTime)
			bar.Finish()

			if wb == 0 {
				fmt.Printf("got %d bytes passcode might be wrong", wb)
			} else {
				fmt.Printf("%d B written in %s (took %v to download)", wb, file, endTime)
			}
		}, ip, port, *useLessCpuTimeExtractor)
		if err != nil {
			log.Fatalf("cant receive file %v", err)
		}
	}

}

// this function return the mode and filename if the mode is send the filename is the name of the file to send
// else it is the name of the file to save as
func PromptUser() (string, string, int) {
	var (
		mode   string
		file   string
		number int
	)

	for {
		fmt.Print("Do you want to send/receive: ")
		_, _ = fmt.Scanln(&mode)
		mode = utils.RemoveWhitespace(mode)

		if mode == "send" || mode == "receive" || mode == "r" || mode == "s" {
			break
		} else if mode == "exit" {
			os.Exit(0)
		}
		fmt.Print("\n")

		fmt.Println("wrong input")
	}
	if utils.ModeToSimple(mode) == "s" {
		for {
			fmt.Print("enter the location of your file: ")
			_, _ = fmt.Scanln(&file)
			file = utils.RemoveWhitespace(file)

			if file != "" && utils.DoesFileExist(file) {
				break
			}
			fmt.Print("\n")
			fmt.Println("file doesnt exit")
		}
	} else {
		for {
			fmt.Print("enter the location to save file: ")
			_, _ = fmt.Scanln(&file)
			file = utils.RemoveWhitespace(file)

			if file != "" && utils.DoesFileExist(file) == false {
				break
			} else if file == "exit" {
				os.Exit(0)
			}
			fmt.Print("\n")

			fmt.Println("file exists")
		}
		for {
			fmt.Print("enter the pass code: ")
			_, err := fmt.Scanln(&number)
			fmt.Print("\n")
			if err != nil {
				log.Fatalln(err)
			} else {
				break
			}
		}
	}
	return mode, file, number
}
