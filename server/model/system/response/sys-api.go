package response

import "time"

type SysApiResp struct {
	Id          uint      `json:"id"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	Group       string    `json:"group,omitzero"`
	Method      string    `json:"method"`
	Disabled    bool      `json:"disabled"`
	Sort        uint      `json:"sort"`
	CreatedAt   time.Time `json:"created_at,omitzero"`
	UpdatedAt   time.Time `json:"updated_at,omitzero"`
}

type SysApiGroupResp struct {
	Group string       `json:"group"`
	Items []SysApiResp `json:"items"`
}
