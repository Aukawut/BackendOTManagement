package model

type WorkCell struct {
	ID_WORKGRP   int
	NAME_WORKGRP string
}

type WorkCellJoinFactory struct {
	ID_WORK_CELL  int
	NAME_WORKCELL string
	ID_FACTORY    int
	FACTORY_NAME  string
	ID_WORKGRP    int
}

type ReqWorkCellBody struct {
	ID_WORKGRP   int    `json:"workgroup"`
	ID_FACTORY   int    `json:"factory"`
	NAME_WORKCEL string `json:"name"`
}
