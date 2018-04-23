package main

//Ingore - rule for ignore
type Ingore struct {
	Mask     string `json:"mask"`
	Line     int    `json:"line"`
	Ruleid   string `json:"ruleId"`
	Ingoreid string `json:"ingoreId"`
}
