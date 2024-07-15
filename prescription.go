package main

import (
    "fmt"
    "strings"
)
type Prescription struct {
    Id              int                 //物品ID
    Name,Description,Attribute  string  //名称、描述说明、属性说明（对应PRA_PROPERTY）
    ProductId,NeedPotential,NeedLing int    //合成物品ID、潜力需求量、灵容量需求量
    Skill,Price int //技能ID、价格
    Effect1,Effect2,Effect3,Effect,BuyScene string //粒子特效剑、双剑、琴、粒子特效组合、购买场景
    EaUrl string
}
type Prescriptions struct {
    Type string         //配方子类型，剑、身体防具、锻冶等
    IsRZ,IsNotRZ,IsDY,IsZL bool //配方大类的标志
    PrescriptionList []Prescription
}
func getPrescriptionType(prescriptionType []string) (prescriptions Prescriptions)  {
    switch prescriptionType[0] {
    case "熔铸图谱":
        prescriptions.IsRZ=true
    case "锻造图谱":
        prescriptions.IsDY,prescriptions.IsNotRZ=true,true
    case "注灵图谱":
        prescriptions.IsZL,prescriptions.IsNotRZ=true,true
    }
    prescriptionList:=[]Prescription{}
    prescriptionSql:=""
    typeId:=getId("PrescriptionType",prescriptionType[0])
    if prescriptionType[0]=="熔铸图谱" {
        equipTypeId:=getId("EquipType",prescriptionType[1])
        prescriptionSql=fmt.Sprintf(`
            select pra.id,pra.name,pra.description,attribute,product_id,need_potential,need_ling,pra.skill_id,pra.price,pra.ef2,pra.ef3,ef4 
            from Prescription pra join Equip ea on pra.product_id=ea.id where pra.type=? and ea.type=%d
            order by pra.price,ea.id
        `,equipTypeId)
    } else {
        prescriptionSql=`
            select id,name,description,attribute,product_id,need_potential,need_ling,skill_id,price,ef2,ef3,ef4 
            from Prescription where type=?
            order by price,need_potential,need_ling
        `
    }
    rows,_ := Db.Query(prescriptionSql,typeId)
    for rows.Next() {
        prescription := Prescription{}
        rows.Scan(
            &prescription.Id,&prescription.Name,&prescription.Description,&prescription.Attribute,
            &prescription.ProductId,&prescription.NeedPotential,&prescription.NeedLing,&prescription.Skill,
            &prescription.Price,&prescription.Effect1,&prescription.Effect2,&prescription.Effect3,
        )
        for _,eff:=range []string{prescription.Effect1,prescription.Effect2,prescription.Effect3} {
            if eff!="" {
                prescription.Effect+=fmt.Sprintf("%s ",eff)
            }
        }
        prescription.Effect=strings.TrimSuffix(prescription.Effect," ")
        prescription.BuyScene=getBuyScene(prescription.Id)
        for k,v:=range routeMapEquip {
            if v==prescriptionType[1] {
                prescription.EaUrl=k
                break
            }
        }
        prescriptionList = append(prescriptionList, prescription)
    }
    rows.Close()
    prescriptions.PrescriptionList=prescriptionList
    prescriptions.Type=prescriptionType[1]
    return
}
