package data

type Command struct {
    Id string `json:"id"`
    ProcessId string `json:"processId"`
    Name string `json:"name"`
    Body string `json:"body"`
}
