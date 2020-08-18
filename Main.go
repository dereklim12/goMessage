package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var headers []string
var messages []string
var allMsgs []string
var collatedMsgs [][]string
var output string

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func myError(text string) error {
	//fmt.Println(&errorString{text})
	return &errorString{text}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	dat, err := ioutil.ReadFile(path + "/src/github.com/dereklim7777777/decode/message.txt")
	check(err)
	//@todo make sure integers only 1 & 0 by checking ASCII values is 48 or 49 if integer
	/*for k, v := range dat {
		//check if value is line feed ASCII code 10
		if v == 10 {
			v =
		}
	}*/
	content := string(dat)
	//regular expression to remove all carriage return and new lines
	var re = regexp.MustCompile(`\r?\n`)
	content = re.ReplaceAllString(content, "")
	//fmt.Printf("%T %v", content, content)
	return content
}

func convertBinary(message string) int {
	var res int
	switch message {
	case "000":
		res = 1
	case "001":
		res = 1
	case "010":
		res = 2
	case "011":
		res = 3
	case "100":
		res = 4
	case "101":
		res = 5
	case "110":
		res = 6
	case "111":
		res = 7
	}
	return res
}

//method to convert string to byte
func strToBytes(message string) []byte {
	b := []byte(message)
	return b
}

//method to convert byte to string
func bytesToStr(bytes []byte) string {
	s := string(bytes)
	return s
}

//trim 3 characters key from message
func trimKeys(message string) string {
	message = message[3:]
	return message
}

//trim characters added to msgs from message
func trimMsg(message string, d int) string {
	message = message[d:]
	return message
}

//generate ending segment
func genEndSeg(keyLen int) string {
	endSeg := ""
	for i := 0; i < keyLen; i++ {
		endSeg += "1"
	}
	return endSeg
}

func isEndSeg(str string, endSeg string, message string, d int) (bool, string) {
	isEndSeg := false
	trimmedMessage := ""
	if str == endSeg {
		isEndSeg = true
		trimmedMessage = trimMsg(message, d)
		//fmt.Println("str", str, "-", endSeg, "Original", message, "d:", d, "trimmedMessage", trimmedMessage, "Line 124")
	}
	return isEndSeg, trimmedMessage
}

func terminate() {
	//fmt.Println(allMsgs, "Terminated --- line 133")
	log.Fatal("Terminated --- line 134")
	return
}

//method to process a single message from getMessages method
func getMessages(messages []string) []string {
	var vMsgs []string
	for v := range messages {
		msgsArr := processNewMessage(messages[v])
		for n := range msgsArr {
			vMsgs = append(vMsgs, msgsArr[n])
		}
		return vMsgs
	}
	return vMsgs
}

func getSequence(length int) []string {
	s := []string{"0", "00", "01", "10", "000", "001", "010", "011", "100",
		"101", "110", "0000", "0001", "0010", "0011", "0100", "0101", "0110",
		"0111", "1000", "1001", "1010", "1011", "1100", "1110", "00000"}
	return s
}

func getNewHeaders(vMsgs []string) string {
	message := ""
	sequence := getSequence(len(headers[0]))
	var newHeader = make(map[string]string)
	for v := range vMsgs {
		//fmt.Println(string(headers[0]), vMsgs[v], v, string(headers[0][v]), sequence[v], "Line 151")
		newHeader[string(sequence[v])] =
			string(headers[0][v])
	}
	for v := range vMsgs {
		//fmt.Println(newHeader[vMsgs[v]], "Line 192")
		message += newHeader[vMsgs[v]]
	}
	fmt.Println(message, "Line 194")
	return message
}

func processNewMessage(message string) []string {
	var vMsgs []string
	var str string
	//fmt.Println(message, "---154")
	str = message[:3]
	//trim keys which is made up of first 3 elements
	//fmt.Println(message, str, "---157")
	message = trimKeys(message)
	//@todo check str is valid binary format
	msgLen := convertBinary(str)
	//fmt.Println(msgLen, "---161")
	//get message segments based on key
	vMsgs = getNewMessage(message, msgLen)
	collatedMsgs = append(collatedMsgs, vMsgs)
	//fmt.Println(vMsgs, collatedMsgs, "Line 162")
	//vHeaders := getNewHeaders(vMsgs)
	output += getNewHeaders(vMsgs)
	terminate()
	return vMsgs
}

func getNewMessage(message string, msgLen int) []string {
	var vMsgs []string
	//fmt.Println(message, "--- Line 170")
	str := ""
	mb := strToBytes(message)
	endSeg := genEndSeg(msgLen)
	d, n := 0, 0
	for i := range mb {
		//fmt.Println(string(mb[i]), "--- Line 171")
		if n < msgLen {
			str += string(mb[i])
			d++
		} else {
			if len(str)%2 == 0 {
				d = d - 1 //@todo need to check this for other permutations
			}
			isEndSeg, trimmedMessage := isEndSeg(str, endSeg, message, d)
			//fmt.Println(isEndSeg, str, endSeg, vMsgs, trimmedMessage, message, d, "--- Line 178")
			if isEndSeg == true {
				//fmt.Println(trimmedMessage, "ln 187")
				trimmedMessage = trimmedMessage[len(endSeg):]
				//trimmedMessage = trimmedMessage[:len(endSeg)]
				for v := range vMsgs {
					allMsgs = append(allMsgs, vMsgs[v])
				}
				//fmt.Println(vMsgs, trimmedMessage, "--- Line 187")
				processNewMessage(trimmedMessage)
				d = 0
				n = 0
			} else {
				vMsgs = append(vMsgs, str)
				str = string(mb[i])
				d = 1
				n = 0
			}
		}
		n++
		//fmt.Println(d, n, vMsgs, "--- Line 198")
	}
	//fmt.Println(vMsgs, allMsgs, "--- Line 201")
	return allMsgs
}

func main() {
	input := readFile()
	var i, n = 0, 0 //i for headers index, n for messages index
	var header string
	var message string
	for k, v := range input {
		//fmt.Printf("%T %v %v", v, string(v), v)
		//fmt.Println(" Key:", k, "indexH", i, "indexM", n, " Line 55")
		if v != 48 && v != 49 {
			header += string(v)
			if k+1 < len(input) {
				if input[k+1] == 48 || input[k+1] == 49 {
					i++
					headers = append(headers, header)
					header = ""
				}
			}
		}
		if v == 48 || v == 49 {
			message += string(v)
			if k+1 < len(input) {
				if input[k+1] != 48 && input[k+1] != 49 {
					n++
					messages = append(messages, message)
					message = ""
				}
			}
		}
	}
	if header != "" {
		headers = append(headers, header)
		header = ""
	}
	if message != "" {
		messages = append(messages, message)
		message = ""
	}
	//fmt.Println("Headers:", headers)
	//fmt.Println(string(headers[0][0]))
	//fmt.Println("Messages:", messages)
	//fmt.Println(string(messages[0][0]))
	/*for k, v := range headers {
		fmt.Println(k, v)
	}*/
	getMessages(messages)
	//fmt.Println("Msgs:", msgs, "Line 290")
	//fmt.Println("Output:", output, "Line 291")
}
