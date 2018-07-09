package model

type Account struct {
	Type     string   `json:"type,omitempty"`
	Pid      string   `json:"pid,omitempty"`
	Email    string   `json:"email,omitempty"`
	Password string   `json:"password,omitempty"`
	Courses  []Course `json:"courses"`
}

type Credentials struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Profile struct {
	Type      string `json:"type,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}

type Session struct {
	Type string `json:"type,omitempty"`
	Pid  string `json:"pid,omitempty"`
}
