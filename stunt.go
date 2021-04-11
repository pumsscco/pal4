package main

import (
    "fmt"
)
type Stunt struct {
    Id int //特技ID
    Name,Description,Attribute  string   //特技名称、特技说明、效果说明(对应PROP)
    AiCmdType,Target,RequireLevel,AttachedSkill,ActType int  
    AiCmdName,TargetName,ActTypeName string  //以上各项的实际名称
    Animation,TargetEffect,TargetBind string //施放者动作ID、目标方特效ID、目标方特效挂载点
    TianheLvl,LingshaLvl,MengliLvl,ZiyingLvl int //各人的修习等级要求
    Role string //修习角色
}
type Stunts struct {
    Type string         //特技类型
    IsSelf bool //是否我方特技
    StuntList []Stunt
}
func getStuntType(stuntType string) (stunts Stunts)  {
    propSelect:=""
    switch stuntType {
    case "我方":
        propSelect=" attribute!=''"
        stunts.IsSelf=true
    case "敌方":
        propSelect=" attribute=''"
    }
    stuntSql:=fmt.Sprintf(`
        select id,name,description,attribute,ai_cmd_type,target,tianhe_lv_lmt,lingsha_lv_lmt,mengli_lv_lmt,ziying_lv_lmt lv_lmt,
        attached_skill,act_type,animation,target_ef,target_bind from Stunt where %s
    `,propSelect)
    stuntList:=[]Stunt{}
    rows,_ := Db.Query(stuntSql)
    for rows.Next() {
        stunt := Stunt{}
        rows.Scan(
            &stunt.Id,&stunt.Name,&stunt.Description,&stunt.Attribute,&stunt.AiCmdType,&stunt.Target,&stunt.TianheLvl,
            &stunt.LingshaLvl,&stunt.MengliLvl,&stunt.ZiyingLvl,&stunt.AttachedSkill,&stunt.ActType,&stunt.Animation,
            &stunt.TargetEffect,&stunt.TargetBind,
        )
        switch {
        case stunt.LingshaLvl!=0:
            stunt.Role="韩菱纱"
        case stunt.MengliLvl!=0:
            stunt.Role="柳梦璃"
        case stunt.ZiyingLvl!=0:
            stunt.Role="慕容紫英"
        default:
            stunt.Role="云天河"
        }
        stunt.RequireLevel=stunt.TianheLvl|stunt.LingshaLvl|stunt.MengliLvl|stunt.ZiyingLvl
        stunt.AiCmdName=getName("AiCommandType",stunt.AiCmdType)
        stunt.TargetName=getName("SkillTarget",stunt.Target)
        stunt.ActTypeName=getName("ActType",stunt.ActType)
        stuntList = append(stuntList, stunt)
    }
    rows.Close()
    stunts.StuntList=stuntList
    stunts.Type=stuntType
    return
}
