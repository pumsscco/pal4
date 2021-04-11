package main

import (
    "fmt"
)
type Magic struct {
    Id int //仙术ID
    Name,Description,Attribute  string   //仙术名称、仙术说明、效果说明
    AiCmdType,Target,Wuling,AttachedSkill,ConsumedMP int  //Ai指令类型、目标、五灵属性、附加技能、耗神值
    AiCmdName,TargetName,WulingName string  //以上各项的实际名称
    Animation,TargetEffect,TargetBind,Color string //施放者动作ID、目标方特效ID、目标方特效挂载点、五灵颜色
}
type Magics struct {
    Type string         //仙术类型
    IsSelf bool //是否我方仙术
    MagicList []Magic
}
func getMagicType(magicType string) (magics Magics)  {
    propSelect:=""
    switch magicType {
    case "我方":
        propSelect=" attribute!=''"
        magics.IsSelf=true
    case "敌方":
        propSelect=" attribute=''"
    }
    magicSql:=fmt.Sprintf(`
        select id,name,description,attribute,ai_cmd_type,target,wuling,attached_skill,consumed_mp,
        animation,target_ef,target_bind from Magic where %s
    `,propSelect)
    magicList:=[]Magic{}
    rows,_ := Db.Query(magicSql)
    for rows.Next() {
        magic := Magic{}
        rows.Scan(
            &magic.Id,&magic.Name,&magic.Description,&magic.Attribute,&magic.AiCmdType,&magic.Target,
            &magic.Wuling,&magic.AttachedSkill,&magic.ConsumedMP,&magic.Animation,&magic.TargetEffect,
            &magic.TargetBind,
        )
        magic.AiCmdName=getName("AiCommandType",magic.AiCmdType)
        magic.TargetName=getName("SkillTarget",magic.Target)
        magic.WulingName=getName("WuLing",magic.Wuling)
        switch magic.WulingName {
        case "水":
            magic.Color="water"
        case "火":
            magic.Color="fire"
        case "雷":
            magic.Color="thunder"
        case "风":
            magic.Color="air"
        case "土":
            magic.Color="earth"
        }
        magicList = append(magicList, magic)
    }
    rows.Close()
    magics.MagicList=magicList
    magics.Type=magicType
    return
}
