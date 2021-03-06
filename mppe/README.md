# Ministério Público de Pernambuco

Este crawler tem como objetivo a recuperação de informações sobre folhas de pagamentos dos funcionários do Ministério Público de Pernambuco. O site com as informações pode ser acessado [aqui](https://transparencia.mppe.mp.br/contracheque).

O crawler está estruturado como uma CLI. Você passa dois argumentos (mês e ano) e serão baixadas oito planilhas no formato XLSX, 

## Como usar

### Executando com Docker

- Inicialmente é preciso instalar o [Docker](https://docs.docker.com/install/). 

- Construção da imagem:

```sh
docker build --build-arg GIT_COMMIT=$(git rev-list -1 HEAD) -t mppe .
```

- Execução:
	- Para executar é necessário passar o .env e um volume, caso deseje persistência dos dados. No .env, o campo ```OUTPUT_PATH``` indica o path relativo dentro do container. 
	- OBS: O path dos arquivos retornado no [CrawlingResult](https://github.com/dadosjusbr/storage/blob/master/agency.go) será relativo ao container. Montar o volume de dados no mesmo path para diversos containers pode ser boa prática.
	- Um arquivo .env.example na pasta raíz indica as variáveis de ambiente que precisam ser passadas para o coletor.


```sh

docker run \
--mount type=bind,source="$(pwd)"/FILES_DIR,target=/OUTPUT_DIR \
--env-file=.env \
mppe --mes=${MES} --ano=${ANO}
```

- No comando de run:
	- ```--mount type=bind,source="$(pwd)"/FILES_DIR,target=/OUTPUT_DIR``` faz o bind de um diretório existente em sua máquina chamado FILES_DIR e um dentro container chamado OUTPUT_DIR, o mesmo passado dentro do .env.
	- ```--env-file=.env``` especifica o path para o env-file.
	- ```mppe --mes=${MES} --ano=${ANO}``` é o nome do container que queremos executar e os argumentos que serão passados para a função de entrada.

  
### Executando sem uso do docker:

- É preciso ter o compilador de Go instalado em sua máquina. Mais informações [aqui](https://golang.org/dl/).

- Um arquivo .env.example na pasta raíz indica as variáveis de ambiente que precisam ser passadas para o coletor.
- O resultado do coletor, [CrawlingResult](https://github.com/dadosjusbr/storage/blob/master/agency.go), possui um campo que indica o commit do git usado para dar o build. Para que ele seja setado adequadamente, é precisso passar o commit como argumento do build.
 

```sh
go get
go build
./mppe --mes=${MES} --ano=${ANO}
```


## Dificuldades para libertação dos dados

- Não há API;
- Utilização de formato próprietário: as planilhas são baixadas no formato XLSX;
- Números sem função aparente nas URLs: no meio das URLs existem números que muitas vezes parecem não ter função, pois os links funcionam com ou sem eles. Exemplo de url para baixar os membros ativos de 02/2019: se você tirar o `445-`, o download ainda acontece https://transparencia.mppe.mp.br/contracheque/category/445-remuneracao-de-todos-os-membros-ativos-2019?download=4936:membros-ativos-02-2019;
- Falta de padrão nos nomes dos meses nas URLs: alguns estão com números, outros com o nome dos meses;
- Inconsistências nas URLs: por exemplo os relatórios de 2014 a partir de fevereiro estarem com 2015 na URL, e não 2014;
- Dados disponíveis somente a partir de 2011;
- Alguns arquivos possuem uma variação absurdamente grande uns com outros de tal forma que torna muito dificil a criação coletores, isso hora criando arquivos com multiplas folhas sem se saber qual a que de fato possui os dados hora fornecendo planilhas vazias; 
