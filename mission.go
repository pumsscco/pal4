package main

import (
    "fmt"
)
type Mission struct {
    Trunk,QuestId int //任务变量、任务编号
    Name,Picture,Description  string   //任务名称、任务图片、说明文字
    StoryPer float32  //剧情完成度  
    StoryShow int //任务变灰标识
    StoryShowName string //任务变灰是或否
}
type Missions struct {
    Type string         //任务类型
    IsNotMain,IsNotDelegate bool //主线及委托的反向标志位
    MissionList []Mission
}
func getMissionType(missionType string) (missions Missions)  {
    pattern:=""
    switch missionType {
    case "主线":
        pattern=" depended_id<200"
        missions.IsNotDelegate=true
    case "委托":
        pattern=" depended_id between 200 and 299"
        missions.IsNotMain=true
    case "支线":
        pattern=" depended_id>=300"
        missions.IsNotDelegate=true
        missions.IsNotMain=true
    }
    missionSql:=fmt.Sprintf(`
        select trunk,quest_id,name,picture,description,story_per,story_show from Mission where %s
    `,pattern)
    missionList:=[]Mission{}
    rows,_ := Db.Query(missionSql)
    for rows.Next() {
        mission := Mission{}
        rows.Scan(
            &mission.Trunk,&mission.QuestId,&mission.Name,&mission.Picture,
            &mission.Description,&mission.StoryPer,&mission.StoryShow,
        )
        if mission.StoryShow==1 {
            mission.StoryShowName="是"
        } else {
            mission.StoryShowName="否"
        }
        missionList = append(missionList, mission)
    }
    rows.Close()
    missions.MissionList=missionList
    missions.Type=missionType
    return
}