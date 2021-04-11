package main

type Question struct {
    Problem,Answer1,Answer2,Answer3 string //问题、答案1、答案2、答案3
    RightAnswer    int   //正确答案
}
type Questions struct {
    Type string         //问题类型
    QuestionList []Question
}
func getQuestionType(dbname string) (questions Questions)  {
    db:=0
    switch dbname {
    case "仙剑历史":
        db=1
    case "仙剑故事":
        db=2
    case "仙剑世界":
        db=3
    }
    questionList:=[]Question{}
    questionSql:=`
        select question,answer1,answer2,answer3,right_answer from GameQuestion where db=?
    `
    rows,_ := Db.Query(questionSql,db)
    for rows.Next() {
        question := Question{}
        rows.Scan(
            &question.Problem,&question.Answer1,&question.Answer2,&question.Answer3,&question.RightAnswer,
        )
        questionList = append(questionList, question)
    }
    rows.Close()
    questions.QuestionList=questionList
    questions.Type=dbname
    return
}