package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	
	scanner := bufio.NewScanner(os.Stdin)  //создание нового сканера
	fmt.Printf("domain,hasMX,hasSPF,sprRecord,hasDMARC,dmarcRecord\n")  //вывод порядкак ответа о проверках

	for scanner.Scan() {
		checkDomain(scanner.Text())  //вызываем метод проверяющий домен передовая в него последний сгенерированый методом Scan() токен в формате string
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string){

	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecord, err := net.LookupMX(domain)  //получаем DNS MX запись для domain

	if err != nil {
		log.Fatal("Error: %v\n", err)
	}
	if len(mxRecord) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)  //получаем DNS TXT запись для domain

	if err != nil {
		log.Fatal("Error:%v\n", err)
	}

	for _, record := range txtRecords {           //проверяем наличие хотя-бы 1 записи в txtRecords
		if strings.HasPrefix(record, "v=spf1") {  //с префиксом "v=spf1"
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)  //получаем DNS TXT запись для _dmarc.domain
	if err != nil {
		log.Printf("ErrorL%v\n", err)
	}

	for _, record := range dmarcRecords {           //проверяем наличие хотя-бы 1 записи в dmarcRecords
		if strings.HasPrefix(record, "v=DMARC1") {  //с префиксом "v=DMARC1"
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%v %v %v %v %v %v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)  //выводим результаты проверок в одну строку
}