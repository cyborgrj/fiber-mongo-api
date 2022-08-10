package models

import (
	"encoding/json"
	"fiber-mongo-api/custom_errors"
	"io/ioutil"
	"net/http"
)

type Address struct {
	Logradouro  string `json:"logradouro"`
	Cep         string `json:"cep" validate:"required"`
	Bairro      string `json:"bairro"`
	Cidade      string `json:"localidade"`
	Uf          string `json:"uf"`
	Complemento string `json:"complemento"`
}

func GetData(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	jsonByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return jsonByte, nil

}

// Para fazer: função isValid para o CEP e se retorna algo, cep que não retorna endereço algum etc.
/*/ func isCepValid (cep string) bool {

}

*/
func ToAdress(cep string) (*Address, error) {
	endereco := &Address{}

	if cep == "" {
		return nil, custom_errors.ErrCEPnaoInformado
	}

	data, errData := GetData("https://viacep.com.br/ws/" + cep + "/json/")
	if errData != nil {
		return nil, custom_errors.ErrCEPnaoRecuperado
	}

	err := json.Unmarshal(data, &endereco)
	if err != nil {
		return nil, err
	}

	return endereco, nil
}
