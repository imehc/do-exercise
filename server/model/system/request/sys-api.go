package request

type UpdateSysApiReq struct {
	Description string `json:"description"`
	Group       string `json:"group"`
	Disabled    bool   `json:"disabled"`
	Sort        uint   `json:"sort"`
}
