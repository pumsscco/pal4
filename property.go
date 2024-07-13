package main

type Property struct {
    Id              int                 //物品ID
    Name,Description,Attribute  string  //名称、描述说明、属性说明（对应attribute）
    Model,Texture     string   //模型、贴图
    Level int    //道具等级
    LevelName string //等级名称
    Price,AttachedSkill          int    //价格、挂载技能
    BuyScene string //购买场景
}
type Properties struct {
    Type string         //道具子类型，食物、香料、矿石、尸块等
    IsDH bool //是否恢复类
    IsNotCSJQ bool //既不是尸块，也不是剧情类，可购买
    IsNotCJQ bool //不是香料、材料类、剧情类三类之一，可挂载技能
    IsNotJQ bool //不是剧情类，有价格
    PropertyList []Property
}
//依据装备类型的中文名，获得该类物品的全部属性
func getPropertyType(propertyType []string) (properties Properties)  {
    properties.IsNotCSJQ=true
    properties.IsNotCJQ=true
    properties.IsNotJQ=true
    if propertyType[0]=="恢复" && propertyType[1]=="其它恢复类" {
        properties.IsDH=true
    }
    if propertyType[1]=="尸块" || propertyType[0]=="剧情" {
        properties.IsNotCSJQ=false
    }
    if propertyType[1]=="香料" || propertyType[0]=="材料" || propertyType[0]=="剧情" {
        properties.IsNotCJQ=false
    }
    if propertyType[0]=="剧情" {
        properties.IsNotJQ=false
    }
    //先依名称取类型ID
    typeId:=getId("PropertyClass",propertyType[0])
    //利用类型ID取原始数据
    propertyList:=[]Property{}
    propertySql:=`
        select id,name,description,attribute,model,texture,property_level,price,attached_skill 
        from Property where type=? 
    `
    switch propertyType[1] {
    case "食物":
        propertySql+=` and model regexp "^SW"`
    case "其它恢复类":
        propertySql+=` and model not regexp "^SW"`
    case "香料":
        propertySql+=` and model regexp "^CX"`
    case "其它辅助类":
        propertySql+=` and model not regexp "^CX"`
    case "矿石":
        propertySql+=` and model regexp "^CK" and attribute="熔铸、锻冶的材料"`
    case "尸块":
        propertySql+=` and attribute="注灵的材料"`
    case "其它材料":
        propertySql+=` and (model regexp "^CQ" or attribute="")`
    }
    propertySql+=` order by price`
    rows,_ := Db.Query(propertySql,typeId)
    for rows.Next() {
        property := Property{}
        rows.Scan(
            &property.Id,&property.Name,&property.Description,&property.Attribute,&property.Model,
            &property.Texture,&property.Level,&property.Price,&property.AttachedSkill,
        )
        if properties.IsDH {
            property.LevelName=getName("PropertyLevel",property.Level)
        }
        if properties.IsNotCSJQ {
            property.BuyScene=getBuyScene(property.Id)
        }
        propertyList = append(propertyList, property)
    }
    rows.Close()
    properties.PropertyList=propertyList
    properties.Type=propertyType[1]
    return
}
