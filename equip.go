package main

import (
    "fmt"
    "strings"
    "reflect"
    "html/template"
    "time"
    "encoding/json"
)
//将类别提取出后的结构组合，适用于自动依据类别填充模板
type Equip struct {
    Id              int                 //物品ID   id
    Name,Description,Attribute     string      //名称  name、描述  description、属性综合 attribute
    Model,Texture     string   //模型 model、贴图 texture
    LvLmt,Price,LingCap          int         //等级要求 level_limit、价格 price、灵容量 ling_capacity
    Potential,MaxHP,AdditionalRage int //潜力 forge_potential、精上限提升 max_hp、怒气增加 additional_rage
    MaxMP,Physical,Toughness,Speed,Lucky,Will int //神上限提升 max_mp、武防速运灵 physical toughness speed lucky will
    Water,Fire,Thunder,Air,Earth int         //水火雷风土五灵属性 water fire thunder air earth
    WaterAdditional,FireAdditional,ThunderAdditional int  //水火雷伤害追加
    AirAdditional,EarthAdditional int //风土伤害追加
    PhysicalExtract,WaterExtract,FireExtract float32   //物理、水、火吸收
    ThunderExtract,AirExtract,EarthExtract float32 //雷风土吸收
    PhysicalReact,WaterReact,FireReact float32   //物理、水、火反弹
    ThunderReact,AirReact,EarthReact float32 //雷风土反弹
    AdditionalCritical,FendOff,AdditionalHitting float32 //暴击、格挡、命中追加
    Effect1 string //刀光特效、购买场景
    BuyScene template.HTML
}
type Equips struct {
    Type string         //装备类型，剑、双剑等
    IsDswordOrSword bool //剑与双剑类型的判定
    EquipList []Equip
}
//依据装备类型的中文名，获得该类物品的全部属性
func getEquipType(equipType string) (equips Equips)  {
    //经验证，在gin框架中实现的，用回这里，也是行的，不过有个问题，在redis的key中，字符串含有中文，验证起来不方便，也不太规范，必须只含英文、数字，或是下划线才合适
    k:=fmt.Sprintf("pal4:equip:%s",equipType)
    val, err := client.Get(k).Result()
	if err == nil {
        json.Unmarshal([]byte(val),&equips)
		return
	} else {
        logger.Println("get equip from redis error: ",err)
    }
    equipAttribute:=[][]string{
        {"MaxHP","精上限"},{"AdditionalRage","气"},{"MaxMP","神上限"},
        {"Physical","武"},{"Toughness","防"},{"Speed","速"},{"Lucky","运"},{"Will","灵"},
        {"Water","水"},{"Fire","火"},{"Thunder","雷"},{"Air","风"},{"Earth","土"},
        {"WaterAdditional","水伤害"},{"FireAdditional","火伤害"},{"ThunderAdditional","雷伤害"},
        {"AirAdditional","风伤害"},{"EarthAdditional","土伤害"},{"PhysicalExtract","物理吸收"},
        {"WaterExtract","水吸收"},{"FireExtract","火吸收"},{"ThunderExtract","雷吸收"},
        {"AirExtract","风吸收"},{"EarthExtract","土吸收"},{"PhysicalReact","物理反弹"},
        {"WaterReact","水反弹"},{"FireReact","火反弹"},{"ThunderReact","雷反弹"},
        {"AirReact","风反弹"},{"EarthReact","土反弹"},
        {"AdditionalCritical","暴击"},{"FendOff","格挡"},{"AdditionalHitting","命中"},
        {"LingCap","灵蕴"},{"Potential","潜力"},
    }
    if equipType=="双剑" || equipType=="剑" {
        equips.IsDswordOrSword=true
    }
    typeId:=getId("EquipType",equipType)
    //利用类型ID取原始数据
    equipList:=[]Equip{}
    equipSql:=`
        select id,name,description,model,texture,tianhe_lv_lmt|lingsha_lv_lmt|mengli_lv_lmt|ziying_lv_lmt lvl_lmt as lvl_lmt,price,
        ling_capacity,forge_potential,max_hp,additional_rage,max_mp,physical,toughness,speed,lucky,will,water,fire,
        thunder,air,earth,water_additional,fire_additional,thunder_additional,air_additional,earth_additional,
        physical_extract,water_extract,fire_extract,thunder_extract,air_extract,earth_extract,physical_react,
        water_react,fire_react,thunder_react,air_react,earth_react,additional_critical,fend_off,additional_hitting,ef1 
        from Equip where type=? order by lvl_lmt, price
    `
    //rows,_ := Db.Query(equipSql,typeId)
    rows,_ := Db.Query(equipSql,typeId)
    for rows.Next() {
        equip := Equip{}
        rows.Scan(
            &equip.Id,&equip.Name,&equip.Description,&equip.Model,&equip.Texture,&equip.LvLmt,&equip.Price,
            &equip.LingCap,&equip.Potential,&equip.MaxHP,&equip.AdditionalRage,&equip.MaxMP,&equip.Physical,
            &equip.Toughness,&equip.Speed,&equip.Lucky,&equip.Will,&equip.Water,&equip.Fire,&equip.Thunder,
            &equip.Air,&equip.Earth,&equip.WaterAdditional,&equip.FireAdditional,&equip.ThunderAdditional,
            &equip.AirAdditional,&equip.EarthAdditional,&equip.PhysicalExtract,&equip.WaterExtract,&equip.FireExtract,
            &equip.ThunderExtract,&equip.AirExtract,&equip.EarthExtract,&equip.PhysicalReact,&equip.WaterReact,
            &equip.FireReact,&equip.ThunderReact,&equip.AirReact,&equip.EarthReact,&equip.AdditionalCritical,
            &equip.FendOff,&equip.AdditionalHitting,&equip.Effect1,
        )
        //利用反射合并属性
        v:=reflect.ValueOf(equip)
        for _,f:= range equipAttribute {
            fv:=v.FieldByName(f[0])
            t:=fv.Type()
            if t.String()=="int" && fv.Int()>0 {
                if f[1]=="灵蕴" || f[1]=="潜力" {
                    equip.Attribute+=fmt.Sprintf("%s:%d ",f[1],fv)
                } else {
                    equip.Attribute+=fmt.Sprintf("%s+%d ",f[1],fv)
                }
            } else if t.String()=="float32" && fv.Float()>0 {
                equip.Attribute+=fmt.Sprintf("%s+%s%% ",f[1],perDisp(float32(fv.Float()*100)))
            }
        }
        equip.Attribute=strings.TrimSuffix(equip.Attribute," ")
        buyScene:=getBuyScene(equip.Id)
        if buyScene=="" {
            for k,v:= range routeMapEquip {
                if v==equipType && k!="p" {
                    equip.BuyScene=template.HTML(fmt.Sprintf(`<a href="/%s/%s\">%s配方</a>`,"pra",k,equipType))
                    break
                }
            }
        } else {
            equip.BuyScene=template.HTML(buyScene)
        }        
        equipList = append(equipList, equip)
    }
    rows.Close()
    equips.EquipList=equipList
    equips.Type=equipType
    
    s,err:=json.Marshal(equips)
    if err!=nil {
        logger.Println("equip serialize error: ",err)
    } else {
        client.Set(k, string(s), 36*time.Hour)
    }
    return
}
