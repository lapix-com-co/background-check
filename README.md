# Annotations Handler

Gives information about any colombian person and supports the following providers:

- **Siri** Contraloria

## Siri

The basic person properties as name, middle name, last name needed because *Siri* has some primitive catpcha.

#### Usage

```go
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
```

#### TODO
- Add first name question
