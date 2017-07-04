package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 2
const delay = 3

func main() {

	for {
		exibeIntroducao()
		exibeMenu()

		comando := leComando()
		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:

			imprimeLogs()
		case 0:
			fmt.Println("Saindo do Programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço esse comando")
			os.Exit(-1)
		}
	}

}

func exibeIntroducao() {
	fmt.Println("Olá!")
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("")
	return comandoLido
}

func exibeMenu() {
	fmt.Println("1 - Iniciar monitoramento")
	fmt.Println("2 - Exibir os Logs")
	fmt.Println("0 - Sair")
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")
	sites := leSitesDoArquivo()
	for i := 0; i < monitoramentos; i++ {
		fmt.Println("Rodada", i+1, "de monitoramento!")
		for i, site := range sites {
			fmt.Println("Testando site", i+1, ":", site)
			testaSite(site)
			fmt.Println("")
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}

}

func testaSite(site string) {
	resp, error := http.Get(site)

	if error != nil {
		fmt.Println("Ocorreu um erro:", error)
	}
	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "carregado com sucesso!")
		registraLog(site, true)
	} else {
		fmt.Println("Ocorreu algum erro ao carregar o site:", site, "! Status Code:", resp.StatusCode)
		registraLog(site, false)
	}
}

func leSitesDoArquivo() []string {
	var sites []string
	arquivo, error := os.Open("sites.txt")

	if error != nil {
		fmt.Println("Ocorreu um erro:", error)
	}

	leitor := bufio.NewReader(arquivo)
	for {
		linha, error := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if error == io.EOF {
			break
		}
	}

	arquivo.Close()

	return sites
}

func registraLog(site string, status bool) {
	arquivo, error := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if error != nil {
		fmt.Println("Ocorreu um erro:", error)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")
	arquivo.Close()
}

func imprimeLogs() {
	fmt.Println("Exibindo Logs...")
	arquivo, erro := ioutil.ReadFile("logs.txt")
	if erro != nil {
		fmt.Println(erro)
	}

	fmt.Println(string(arquivo))
}
