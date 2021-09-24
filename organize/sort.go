package organize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var FileLists struct {
	empresa         []string
	estabelecimento []string
	socio           []string
	simples         []string
	cnae            []string
	situacao        []string
	municipio       []string
	natureza        []string
	pais            []string
	socios          []string
}

func GetFiles() {

	// TODO: Tirar os dois pootos qnd rodar normal.
	dataFolder := "../data"

	organizedFiles := FileLists

	// Passa por todos os arquivos da pasta de dados
	err := filepath.Walk(dataFolder, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(strings.ToUpper(path), "EMPRECSV") {
			organizedFiles.empresa = append(organizedFiles.empresa, path)
		} else if strings.Contains(strings.ToUpper(path), "ESTABELE") {
			organizedFiles.estabelecimento = append(organizedFiles.estabelecimento, path)
		}
		// files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Empresas: ", organizedFiles.empresa)
	fmt.Println()
	fmt.Println("Estabelecimento: ", organizedFiles.estabelecimento)

}

// go test -v github.com/cuducos/minha-receita/organize -run TestGetFiles
