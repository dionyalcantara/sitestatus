package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	maxMonitorCycles = 10
	monitoringTime = 5
	statusOnline = "Online"
	statusOffline = "Offline"
)
	

func main() {
	for {
		option := introductions()

		switch option {
		case 1:
			fmt.Println("Monitorando os sites...")
			startMonitoring(readArchives())
		case 2:
			fmt.Println("Adicionando um site...")
			addSite()
		case 3:
			fmt.Println("Logs:")
			showLogs()
		case 4:
			fmt.Println("Sites cadastrados:")
			showSites()
		case 0:
			fmt.Println("Saindo...")
			os.Exit(0)
		default:
			fmt.Println("Comando Inválido")
			os.Exit(-1)
		}
	}

}

func introductions() int {
	fmt.Println("1 - Monitorar os sites")
	fmt.Println("2 - Adicionar um site")
	fmt.Println("3 - Exibir os logs")
	fmt.Println("4 - Exibir os sites cadastrados")
	fmt.Println("0 - Sair do programa")

	var option int
	fmt.Scan(&option)

	return option
}

func readArchives() []string {
	var archives []string
	file, err := os.OpenFile("sites.txt", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println("erro", err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		line = strings.TrimSpace(line)

		if err == io.EOF {
			break
		}

		archives = append(archives, line)
	}

	file.Close()

	return archives
}

func startMonitoring(sites []string) {
	if len(sites) == 0 {
		fmt.Println("Não há sites cadastrados")
		addSite()
		return
	}

	for i := 1; i < maxMonitorCycles; i++ {

		for _, site := range sites {
			res, err := http.Get(site)

			switch {
			case err != nil:
				fmt.Println(err.Error())
				registerLogs(site, statusOffline)
			case res.StatusCode != http.StatusOK:
				fmt.Println(site,"-", statusOffline)
				registerLogs(site, statusOffline)
			default:
				fmt.Println(site,"-", statusOnline)
				registerLogs(site, statusOnline)
			}
		}
		fmt.Println("")
		time.Sleep(monitoringTime * time.Second)

	}

}

func addSite() {
	fmt.Println("Digite o site que deseja adicionar:")
	var site string
	fmt.Scan(&site)

	sites, err := os.OpenFile("sites.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("erro", err)
	}

	sites.WriteString(site + "\n")

	sites.Close()
}

func showSites() {
	reader := readArchives()

	for _, site := range reader {
		fmt.Println(site)
	}
}

func registerLogs(site string, status string) {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("erro", err)
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " " + site + " - " + status + "\n")

	file.Close()
}

func showLogs() {
	arquivo, err := os.ReadFile("logs.txt")

	if err != nil {
		fmt.Println("erro", err)
	}

	fmt.Println(string(arquivo))
}
