package models

type Report struct {
    ID         int    `json:"id"`
    UserID     int    `json:"user_id"`
    ReportName string `json:"report_name"`
    Content    string `json:"content"`
}
