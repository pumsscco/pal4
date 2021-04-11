package main

import (
    "fmt"
)
type Upgrade struct {
    Level,Experience,MaxHP,MaxMP int //等级、经验值、精、神
    Physical,Toughness,Speed,Lucky,Will  int   //武、防、速、运、灵
    FendOff  float32  //格挡率
    FendOffPer  string  //格挡率百分比形式
}
type Upgrades struct {
    Role string         //角色名称
    UpgradeList []Upgrade
}
func getUpgradeRole(role string) (upgrades Upgrades)  {
    roleId:=getId("Role",role)
    upgradeSql:=`
        select level,experience,max_hp,max_mp,physical,toughness,speed,lucky,will,fend_off from UpgradeData where role_id=?
    `
    upgradeList:=[]Upgrade{}
    rows,_ := Db.Query(upgradeSql,roleId)
    for rows.Next() {
        upgrade := Upgrade{}
        rows.Scan(
            &upgrade.Level,&upgrade.Experience,&upgrade.MaxHP,&upgrade.MaxMP,&upgrade.Physical,&upgrade.Toughness,&upgrade.Speed,
            &upgrade.Lucky,&upgrade.Will,&upgrade.FendOff,
        )
        upgrade.FendOffPer=fmt.Sprintf("%s%%",perDisp(float32(upgrade.FendOff*100)))
        upgradeList = append(upgradeList, upgrade)
    }
    rows.Close()
    upgrades.UpgradeList=upgradeList
    upgrades.Role=role
    return
}
