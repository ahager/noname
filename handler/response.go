package handler

type JsonResponse struct {
    FlagName     string `json:"name"`
    FlagStatus   string `json:"status"`
    FlagSticky   string `json:"sticky"`
    FlagRatio    int    `json:"ratio"`

    ClientId   string   `json:"clientId"`
}
