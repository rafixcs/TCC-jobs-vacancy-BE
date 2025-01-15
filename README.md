# TCC-jobs-vacancy-BE
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues-pr/rafixcs/TCC-jobs-vacancy-BE)


# JobsVacancy Server

Backend para o trabalho de conclusao de curso de pos graduacao desenvolvimento full stack. O objetivo desse desenvolvimento eh construir uma aplicacao para a divulgacao de vagas de emprego

**Índice**
- [Princípios](#principios)
- [Documentaçao](#documentacao)
- [Requisitos](#requisitos)
- [Estrutura](#estrutura)
- [Executando](#executando)

## Principios
Os principios seguidos no desenvolvimento desta aplicacao incluem:
- **Independencia**: para cada responsabilidade dentro do sistema eh devidamente separado
- **Escalabilidade**: a arquitetura foi pensada para que possa suportar o crescimento rapido e manter facilidade na manutecao
- **Seguranca**: foram utilizados praticas de seguranca para garantir a seguranca das funcionalidades da aplicacao

## Documentaçao

Para documentar as APIs criadas para o servidor, foi utilizado a biblioteca *swaggo* para gerar a documentacao para o formato swagger.

A documentacao eh gerada a partir da escrita de comentarios sobre os metodos responsaveis pelas APIs.

```// Auth godoc
// @Summary authenticate user
// @Description authenticate user
// @Tags Auth
// @Param authrequest body AuthRequest true "Change password"
// @Success 200 {object} AuthResponse "Success"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/auth [post]
func Auth(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	}

	authDomain := authfactory.CreateAuthDomain()

	token, roleId, err := authDomain.UserAuth(authRequest.Email, authRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	authResponse := AuthResponse{
		Token:  token,
		RoleId: roleId,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&authResponse)
}
```

Com isso podemos gerar a documentacao com o comando:
`swag init --output ./docs --dir ./src/ --pd`

Quando o servidor estiver rodando, ao acessarmos o endereco: `/swagger/index.html#/Auth/post_api_v1_auth`

<p align="center">
  <img src=".git/images/swagger_example.png" />
</p>


## Requisitos

## Executando