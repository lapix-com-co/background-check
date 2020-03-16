# Background Check

Gives information about any colombian person and supports the following providers:

- **Siri** [Contraloria][siri]

## Siri

The basic person properties as name, middle name, last name are needed because *Siri* has some primitive catpcha.

#### Usage

```go
package main

import (
    "background-check-co/siri"
    "encoding/json"
    "time"
    "log"
    "fmt"
)

func main() {
    handler := siri.NewQuery(siri.Delay(time.Second * 4))
    annotations, err := handler.Pull(&siri.PullInput{
        FirstName:      "John",
        MiddleName:     "Alexander",
        LastName:       "Doe",
        SecondSurname:  "Bell",
        DocumentType:   siri.DNI,
        DocumentNumber: "52644444",
    })

    if err != nil {
        log.Fatal(err)
    }

    b, _ := json.Marshal(annotations)
    fmt.Print(string(b))
}
```

#### TODO

- [ ] Add an easy way to check get the user's name from the document number. 
This will make easier the SIRI validation. This resource need to be lightweight because will be used
before every source check if needed.
- [ ] Add RUT [validation][rut]
- [ ] Add [Police background check][police].
- [ ] Add [Sibor][sibor]
- [ ] What about Panama papers?

[rut]: https://muisca.dian.gov.co/WebRutMuisca/DefConsultaEstadoRUT.faces
[siri]: https://www.procuraduria.gov.co/CertWEB/Certificado.aspx
[police]: https://antecedentes.policia.gov.co:7005/WebJudicial/antecedentes.xhtml
[sibor]: https://www.contraloria.gov.co/control-fiscal/responsabilidad-fiscal/control-fiscal/responsabilidad-fiscal/certificado-de-antecedentes-fiscales/persona-natural