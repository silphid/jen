package model

type VarMap map[string]string

type Values struct {
	Variables    VarMap
	Placeholders VarMap
}
