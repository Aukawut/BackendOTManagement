package model

type Factory struct {
	ID_FACTORY    int
	FACTORY_NAME  string
	ID_GROUP_DEPT int
	CREATED_AT    interface{}
	NAME_GROUP    string
}

type BodyInsertFactory struct {
	GroupDept   int    `json:"groupDept"`
	FactoryName string `json:"factoryName"`
}
