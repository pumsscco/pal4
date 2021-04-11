package main

import (
    "fmt"
    "html/template"
//    "strings"
)
type Things struct {
    Cat,Type   string  //物品的大类与细类
    ThingList []string   //物品列表
}
func getThings(cat,thingType string) (things Things) {
    sql:=""
    switch cat {
    case "ea":
        eaType:=getId("EquipType",routeMapEquip[thingType])
        sql=fmt.Sprintf("select name from Equip where type=%d",eaType)
    case "pa":
        paType:=getId("PropertyClass",routeMapProperty[thingType][0])
        sql=fmt.Sprintf("select name from Property where type=%d",paType)
        switch routeMapProperty[thingType][1] {
        case "食物":
            sql+=` and model regexp "^SW"`
        case "其它恢复类":
            sql+=` and model not regexp "^SW"`
        case "香料":
            sql+=` and model regexp "^CX"`
        case "其它辅助类":
            sql+=` and model not regexp "^CX"`
        case "矿石":
            sql+=` and model regexp "^CK" and attribute="熔铸、锻冶的材料"`
        case "尸块":
            sql+=` and attribute="注灵的材料"`
        case "其它材料":
            sql+=` and (model regexp "^CQ" or attribute="")`
        }
    case "pra":
        praType:=getId("PrescriptionType",routeMapPrescription[thingType][0])
        switch routeMapPrescription[thingType][0] {
        case "熔铸图谱":
            eaType:=getId("EquipType",routeMapPrescription[thingType][1])
            sql=fmt.Sprintf(`
                select pra.name from Prescription pra join Equip ea on pra.product_id=ea.id where pra.type=%d and ea.type=%d
            `,praType,eaType)
        default:
            sql=fmt.Sprintf("select name from Prescription where type=%d",praType)
        }
    }
    rows,_ := Db.Query(sql)
    for rows.Next() {
        var name string
        rows.Scan(&name)
        things.ThingList = append(things.ThingList, name)
    }
    things.Type=thingType
    things.Cat=cat
    return
}
type Method struct {
    Name,StolenEnemy,DropEnemy  string 
    BuyScene template.HTML
    HasNonePickup,HasPickup bool 
    PickList []Pick
}
//依据大类，子类与名称，获得各种得到该物品的渠道
func findMethod(cat,t,thing string) (method Method)  {
    switch cat {
    case "ea":
        id:=getId("Equip",thing)
        //购买
        method.BuyScene=template.HTML(getBuyScene(id))
        if method.BuyScene=="" && t!="p" {
            tName:=routeMapEquip[t]
            method.BuyScene=template.HTML(fmt.Sprintf(`<a href="/%s/%s\">%s配方</a>`,"pra",t,tName))
        }
        //拾取
        method.PickList=pickItem(id)
    case "pa":
        id:=getId("Property",thing)
        method.BuyScene=template.HTML(getBuyScene(id))
        method.StolenEnemy=getStolen(id)
        method.DropEnemy=getDrop(id)
        method.PickList=pickItem(id)
    case "pra":
        id:=getId("Prescription",thing)
        method.BuyScene=template.HTML(getBuyScene(id))
    }
    if method.BuyScene!="" || method.StolenEnemy!="" || method.DropEnemy!="" {
        method.HasNonePickup=true
    }
    if method.PickList!=nil {
        method.HasPickup=true
    }
    method.Name=thing
    return
}