package SurveyMonkey_GetResponse

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"fmt"
	"strconv"
	"io/ioutil"
	"net/http"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error)  {

	// do eval
	fmt.Println("Starting the application...")
		//accessToken := "z8UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ"
		//surveyName := "FLG_2_QA_Variety"

		accessToken := context.GetInput("Access_Token").(string)
		surveyName := context.GetInput("Survey_Name").(string)
		
		jsonstr := ""
		jsonSR := ""
		activityOutput := `{ "survey" : { "questions" : [ { "title" : "", "id" : "", "validation" : "", "position": number, "subtype" : "", "family" : "", "type" : "", "visible" : boolean, "answers" : {"rows": [],"other": [],"choices": []}, "responses" : [] } ] } }`

		request, _ := http.NewRequest("GET", "https://api.surveymonkey.com/v3/surveys?title="+surveyName, nil)
		request.Header.Set("Authorization", "bearer "+accessToken)
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res_surveyID, err_surveyID := client.Do(request)
		surveyID := ""
		if err_surveyID != nil {
			fmt.Printf("The HTTP request for getting SurveyID failed with error %s\n", err_surveyID)
		} else {
			res_surveyID, _ := ioutil.ReadAll(res_surveyID.Body)
			surveyID = gjson.Get(string(res_surveyID), "data.0.id").String()
			//fmt.Println("res_surveyID", surveyID)
		}

		link_surveyDetails := "https://api.surveymonkey.com/v3/surveys/"+surveyID+"/details"
		request, _ = http.NewRequest("GET", link_surveyDetails, nil)
		request.Header.Set("Authorization", "bearer "+accessToken)
		request.Header.Set("Content-Type", "application/json")
		client = &http.Client{}
		res_surveyDetails, err_surveyDetails := client.Do(request)
		if err_surveyDetails != nil {
			fmt.Printf("The HTTP request for getting SurveyDetails failed with error %s\n", err_surveyDetails)
		} else {
			surveyDetails, _ := ioutil.ReadAll(res_surveyDetails.Body)
			//fmt.Println(string(surveyDetails))
			//set surveyDetails
			jsonstr = 	`{ "surveydetails" : `+string(surveyDetails)+`}`
		}
		//jsonSR = `{ "surveyresponses" : `+string(jsonIPSR)+`}`

		link_surveyResponse := "https://api.surveymonkey.com/v3/surveys/"+surveyID+"/responses/bulk"
		request, _ = http.NewRequest("GET", link_surveyResponse, nil)
		request.Header.Set("Authorization", "bearer "+accessToken)
		request.Header.Set("Content-Type", "application/json")
		client = &http.Client{}
		res_surveyResponse, err_surveyResponse := client.Do(request)
		if err_surveyResponse != nil {
			fmt.Printf("The HTTP request for getting SurveyDetails failed with error %s\n", err_surveyResponse)
		} else {
			surveyResponse, _ := ioutil.ReadAll(res_surveyResponse.Body)
			//fmt.Println(string(surveyResponse))
			//set surveyresponses
			jsonSR = `{ "surveyresponses" : `+string(surveyResponse)+`}`
		}
		//fmt.Println(jsonstr)
		//fmt.Println(jsonSR)
/*-----------------------------------------------------------------------------------------------------------*/

		activityOutput = setSurveyDetails(jsonstr, jsonSR, activityOutput)
		context.SetOutput("Response_Json", activityOutput)
		
	return true, nil
}

func setSurveyDetails(jsonstr string, jsonSR string, activityOutput string) string {
		//set metadata
		questions := gjson.Get(jsonstr, "surveydetails.pages.0.questions")
		for _, que := range questions.Array() {
			//fmt.Println("que= ",que.String())
			//queIndex := "survey.questions."+gjson.Get(que.String(), "survey.questions.#").String()+".title"
			queIndex := gjson.Get(activityOutput, "survey.questions.#").String()
			//fmt.Println("queIndex= "+queIndex)
			//set heading
			activityOutput_tmp, _ := sjson.Set(activityOutput, "survey.questions."+queIndex+".title", gjson.Get(que.String(), "headings.0.heading").String())
			//set question id
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".id", gjson.Get(que.String(), "id").String())
			//set
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".validation", gjson.Get(que.String(), "validation").String())
			//set
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".position", gjson.Get(que.String(), "position").String())
			//set
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".subtype", gjson.Get(que.String(), "subtype").String())
			//set
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".family", gjson.Get(que.String(), "family").String())
			//set
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".type", gjson.Get(que.String(), "required.type").String())
			//set
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".visible", gjson.Get(que.String(), "visible").String())
			//set answers-rows
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.rows.0.id", "")
			rows := gjson.Get(que.String(), "answers.rows")
			for r, row := range rows.Array() {
				// tmp := strings.Join([]string{"survey.questions."+queIndex+".answers.rows."+strconv.Itoa(r)+".visible"},"")
				// tmp := "survey.questions."+queIndex+".answers.rows."+strconv.Itoa(r)+".visible"
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.rows."+strconv.Itoa(r)+".visible", gjson.Get(row.String(), "visible").String())
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.rows."+strconv.Itoa(r)+".text", gjson.Get(row.String(), "text").String())
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.rows."+strconv.Itoa(r)+".position", gjson.Get(row.String(), "position").String())
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.rows."+strconv.Itoa(r)+".id", gjson.Get(row.String(), "id").String())
			}
			//set answers-other
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.other.id", gjson.Get(que.String(), "answers.other.id").String())
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.other.visible", gjson.Get(que.String(), "answers.other.visible").String())
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.other.is_answer_choice", gjson.Get(que.String(), "answers.other.is_answer_choice").String())
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.other.apply_all_rows", gjson.Get(que.String(), "answers.other.apply_all_rows").String())
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.other.text", gjson.Get(que.String(), "answers.other.text").String())
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.other.position", gjson.Get(que.String(), "answers.other.position").String())
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.other.num_chars", gjson.Get(que.String(), "answers.other.num_chars").String())
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.other.error_text", gjson.Get(que.String(), "answers.other.error_text").String())
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.other.num_lines", gjson.Get(que.String(), "answers.other.num_lines").String())
			//set answer-choices
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.choices.0.id", "")
			choice := gjson.Get(que.String(), "answers.choices")
			for c, ch := range choice.Array() {
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.choices."+strconv.Itoa(c)+".visible", gjson.Get(ch.String(), "visible").String())
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.choices."+strconv.Itoa(c)+".text", gjson.Get(ch.String(), "text").String())
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.choices."+strconv.Itoa(c)+".position", gjson.Get(ch.String(), "position").String())
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.choices."+strconv.Itoa(c)+".is_na", gjson.Get(ch.String(), "is_na").String())
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.choices."+strconv.Itoa(c)+".weight", gjson.Get(ch.String(), "weight").String())
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.choices."+strconv.Itoa(c)+".description", gjson.Get(ch.String(), "description").String())
					activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".answers.choices."+strconv.Itoa(c)+".id", gjson.Get(ch.String(), "id").String())
			}
		////////////////////////////////////////////////////////////////////////////
			activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".responses.0.id", "")
			responses := gjson.Get(jsonSR, "surveyresponses.data")
			for rs, res := range responses.Array() {
						activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".responses."+strconv.Itoa(rs)+".id", gjson.Get(res.String(), "id").String())
						//set responses-answers
						curr_que := `pages.0.questions.#[id="`+gjson.Get(que.String(), "id").String()+`"].answers`
						answers := gjson.Get(res.String(), curr_que)
						for a, ans := range answers.Array() {
								activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".responses."+strconv.Itoa(rs)+".answers."+strconv.Itoa(a)+".text", gjson.Get(ans.String(), "text").String())
								activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".responses."+strconv.Itoa(rs)+".answers."+strconv.Itoa(a)+".choice_id", gjson.Get(ans.String(), "choice_id").String())
								activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".responses."+strconv.Itoa(rs)+".answers."+strconv.Itoa(a)+".row_id", gjson.Get(ans.String(), "row_id").String())
								//set answer title from surveydetails
								if gjson.Get(activityOutput_tmp, "survey.questions."+queIndex+".responses."+strconv.Itoa(rs)+".answers."+strconv.Itoa(rs)+".text").String()=="" {
									activityOutput_tmp, _ = sjson.Set(activityOutput_tmp, "survey.questions."+queIndex+".responses."+strconv.Itoa(rs)+".answers."+strconv.Itoa(rs)+".text", gjson.Get(que.String(), `answers.choices.#[id="`+gjson.Get(ans.String(), "choice_id").String()+`"].text`).String())
								}
								//fmt.Println("------------>",gjson.Get(ans.String(), "choice_id").String())
						}
			}
		////////////////////////////////////////////////////////////////////////////
			//Update actual output var
			activityOutput = activityOutput_tmp
		}//end of outer loop
		// v1, _ := sjson.Set(activityOutput, "survey.questions.0.title", gjson.Get(jsonstr, "surveydetails.pages.0.questions.0.headings.0.heading").String())
		// fmt.Println("Output= ",v2)
		//fmt.Println("activityOutput= ",activityOutput)
		return activityOutput
}
