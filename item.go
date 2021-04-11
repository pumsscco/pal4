package main

import (
    "fmt"
    "reflect"
    "strings"
)
type Scenes struct {
    Type,Path   string  //场景类型以及其url路径
    SceneList []string   //场景名列表
}
func getScenes(sceneType string) (scenes Scenes) {
    sql:=""
    if sceneType=="迷宫" {
        sql=`select name from SceneName where id regexp "^M"`
        scenes.Path="m"
    } else if sceneType=="城镇" {
        sql=`select name from SceneName where id regexp "^Q"`
        scenes.Path="c"
    }
    rows,_ := Db.Query(sql)
    sceneList:=[]string{}
    for rows.Next() {
        var scene string
        rows.Scan(&scene)
        sceneList = append(sceneList, scene)
    }
    scenes.Type=sceneType
    scenes.SceneList=sceneList
    return
}
type Item struct {
    Section,ItemId,Model,Texture  string //区块编号、物品ID、模型、贴图
    SectionName,Apperance,Outer string //外观，依据模型与贴图分析面得，含室外区域
    CoorX,CoorY,CoorZ float32 //东西坐标、上下坐标、南北坐标
    EquipId,PropertyId,ItemNum int  //装备ID、道具ID、前两项之一的数量
    EquipName,PropertyName string //上项ID所对应名称
    Money   int //钱数量
    Item1Id,Item1Num,Item2Id,Item2Num,Item3Id,Item3Num int 
    Item4Id,Item4Num,Item5Id,Item5Num,Item6Id,Item6Num int
    ItemMoney int
    ItemAll string //把物品全部组合在一起
}
type Items struct {
    Scene string  //场景中文名
    HasOuter bool
    ItemList []Item
}
//依据场景的中文名，找出此场景可拾取的全部物品
func getPickup(scene string) (items Items)  {
    hasOuter:=false
    sceneId:=getSceneId(scene)
    sql:=`
        select section,item_id,model,texture,coor_x,coor_y,coor_z,equip_id,property_id,item_num,money,item1id,
        item1num,item2id,item2num,item3id,item3num,item4id,item4num,item5id,item5num,item6id,item6num,item_money 
        from SceneItem where scene=?
    `
    itemList:=[]Item{}
    rows,_ := Db.Query(sql,sceneId)
    for rows.Next() {
        item:=Item{}
        rows.Scan(
            &item.Section,&item.ItemId,&item.Model,&item.Texture,
            &item.CoorX,&item.CoorY,&item.CoorZ,&item.EquipId,&item.PropertyId,
            &item.ItemNum,&item.Money,&item.Item1Id,&item.Item1Num,
            &item.Item2Id,&item.Item2Num,&item.Item3Id,&item.Item3Num,
            &item.Item4Id,&item.Item4Num,&item.Item5Id,&item.Item5Num,
            &item.Item6Id,&item.Item6Num,&item.ItemMoney,
        )
        item.SectionName=getSectionName(sceneId,item.Section)
        item.Apperance=getApperance(item.Model)
        switch {
        case item.EquipId!=0:
                item.EquipName=getName("Equip",item.EquipId)
                if item.ItemNum>1 {
                    item.ItemAll+=fmt.Sprintf("%s*%d ",item.EquipName,item.ItemNum)
                } else if item.ItemNum==1 {
                    item.ItemAll+=fmt.Sprintf("%s ",item.EquipName)
                }
        case item.PropertyId!=0:
                item.PropertyName=getName("Property",item.PropertyId)
                if item.ItemNum>1 {
                    item.ItemAll+=fmt.Sprintf("%s*%d ",item.PropertyName,item.ItemNum)
                } else if item.ItemNum==1 {
                    item.ItemAll+=fmt.Sprintf("%s ",item.PropertyName)
                }
        case item.Money!=0:
                item.ItemAll+=fmt.Sprintf("金钱*%d ",item.Money)
        }
        //只有在非单独的装备、道具或是钱时，并且物品数量不为0时，才可能是宝箱类
        if item.EquipId==0 && item.PropertyId==0 && item.Money==0 && item.ItemNum==0 {
            v:=reflect.ValueOf(item)
            for i:=1;i<=6;i++ {
                f:=fmt.Sprintf("Item%dId",i)
                if fv:=v.FieldByName(f); fv.Int()>0 {
                    //先查道具，再查装备
                    itemName:=getName("Property",int(fv.Int()))
                    if itemName=="" {
                        itemName=getName("Equip",int(fv.Int()))
                    }
                    if itemName!="" {
                        fn:=fmt.Sprintf("Item%dNum",i)
                        fvn:=v.FieldByName(fn).Int()
                        if fvn>1 {
                            item.ItemAll+=fmt.Sprintf("%s*%d ",itemName,fvn)
                        } else if fvn==1 {
                            item.ItemAll+=fmt.Sprintf("%s ",itemName)
                        }
                    }
                }
            }
        }
        item.Outer=getOuter(sceneId,item.Section)  
        if item.Outer!="" {
            hasOuter=true
        }
        if item.SectionName!="" {
            itemList = append(itemList,item)
        }
    }
    rows.Close()
    if strings.HasPrefix(sceneId,"M") {
        items.HasOuter=false
    } else {
        items.HasOuter=hasOuter
    }
    items.ItemList=itemList
    items.Scene=scene
    return
}
