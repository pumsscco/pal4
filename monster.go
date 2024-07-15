package main

import (
	"fmt"
	"reflect"
	"strings"
)

// monster表分五组enemy结构展示
type EnemyCommon struct {
	Id                                       int    //敌人ID
	Name, Icon, Model, Description           string //名称、图标、模型、描述
	IsBoss, Level, Experience, Target, Range int    //   种族、是否BOSS、经验、等级、一般攻击目标、一般攻击方式
	Count                                    int    //存盘统计
	StatsCount                               string //存盘统计，是或否
	Boss                                     string //种族中文名、是否BOSS，是或否
	TargetName, RangeName                    string //一般攻击目标、一般攻击方式查询出原数据
}
type EnemyCommons struct {
	Part            string
	EnemyCommonList []EnemyCommon
}

func getEnemyCommon() (enemyCommons EnemyCommons) {
	sql := `
        select id,name,icon,model,description,is_boss,level,experience,physical_atk_target,physical_atk_range,count from Monster order by level
    `
	enemyCommonList := []EnemyCommon{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		enemyCommon := EnemyCommon{}
		rows.Scan(
			&enemyCommon.Id, &enemyCommon.Name, &enemyCommon.Icon, &enemyCommon.Model,
			&enemyCommon.Description, &enemyCommon.IsBoss, &enemyCommon.Level, &enemyCommon.Experience,
			&enemyCommon.Target, &enemyCommon.Range, &enemyCommon.Count,
		)
		enemyCommon.Boss = getIsBoss(enemyCommon.IsBoss)
		enemyCommon.TargetName = getName("PhysicalAttackTarget", enemyCommon.Target)
		enemyCommon.RangeName = getName("PhysicalAttackType", enemyCommon.Range)
		if enemyCommon.Count == 1 {
			enemyCommon.StatsCount = "是"
		} else {
			enemyCommon.StatsCount = "否"
		}
		enemyCommonList = append(enemyCommonList, enemyCommon)
	}
	rows.Close()
	enemyCommons.EnemyCommonList = enemyCommonList
	enemyCommons.Part = "日常属性"
	return
}

type EnemyBasic struct {
	Id, HP, Rage, MP                        int    //敌人ID、精、气、神
	Name, Wuling                            string //名称、五灵合并
	Physical, Toughness, Speed, Lucky, Will int    //武防速运灵
	Water, Fire, Thunder, Air, Earth        int    //五灵属性
}
type EnemyBasics struct {
	Part           string //本部分的名称
	EnemyBasicList []EnemyBasic
}

func getEnemyBasic() (enemyBasics EnemyBasics) {
	sql := `
        select id,name,max_hp,rage,max_mp,physical,toughness,speed,lucky,will,water,fire,thunder,air,earth from Monster order by max_hp
    `
	enemyBasicList := []EnemyBasic{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		enemyBasic := EnemyBasic{}
		rows.Scan(
			&enemyBasic.Id, &enemyBasic.Name, &enemyBasic.HP, &enemyBasic.Rage, &enemyBasic.MP, &enemyBasic.Physical,
			&enemyBasic.Toughness, &enemyBasic.Speed, &enemyBasic.Lucky, &enemyBasic.Will, &enemyBasic.Water,
			&enemyBasic.Fire, &enemyBasic.Thunder, &enemyBasic.Air, &enemyBasic.Earth,
		)
		wulingAttribute := [][]string{
			{"Water", "水"}, {"Fire", "火"}, {"Thunder", "雷"},
			{"Air", "风"}, {"Earth", "土"},
		}
		v := reflect.ValueOf(enemyBasic)
		for _, f := range wulingAttribute {
			if fv := v.FieldByName(f[0]); fv.Int() > 0 {
				enemyBasic.Wuling += fmt.Sprintf("%s:%d ", f[1], fv.Int())
			}
		}
		enemyBasic.Wuling = strings.TrimSuffix(enemyBasic.Wuling, " ")
		enemyBasicList = append(enemyBasicList, enemyBasic)
	}
	rows.Close()
	enemyBasics.EnemyBasicList = enemyBasicList
	enemyBasics.Part = "八基本属性与五灵"
	return
}

type EnemyResistance struct {
	Id                                                  int
	Name                                                string
	PhysicalExtract, WaterExtract, FireExtract          float32 //物理、水、火吸收
	ThunderExtract, AirExtract, EarthExtract            float32 //雷风土吸收
	PhysicalExtractPer, WaterExtractPer, FireExtractPer string  //物理、水、火吸收
	ThunderExtractPer, AirExtractPer, EarthExtractPer   string  //雷风土吸收
	PhysicalReact, WaterReact, FireReact                float32 //物理、水、火反弹
	ThunderReact, AirReact, EarthReact                  float32 //雷风土反弹
	React                                               string  //反弹合并
	SoundWounded1, SoundWounded2, SoundWounded3         string  //受伤音效1～3
}
type EnemyResistances struct {
	Part                string
	EnemyResistanceList []EnemyResistance
}

func getEnemyResistance() (enemyResistances EnemyResistances) {
	sql := `
        select id,name,physical_extract,water_extract,fire_extract,thunder_extract,air_extract,earth_extract,
        physical_react,water_react,fire_react,thunder_react,air_react,earth_react,
        sound_wounded1,sound_wounded2,sound_wounded3 from Monster
    `
	enemyResistanceList := []EnemyResistance{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		enemyResistance := EnemyResistance{}
		rows.Scan(
			&enemyResistance.Id, &enemyResistance.Name, &enemyResistance.PhysicalExtract, &enemyResistance.WaterExtract,
			&enemyResistance.FireExtract, &enemyResistance.ThunderExtract, &enemyResistance.AirExtract,
			&enemyResistance.EarthExtract, &enemyResistance.PhysicalReact, &enemyResistance.WaterReact,
			&enemyResistance.FireReact, &enemyResistance.ThunderReact, &enemyResistance.AirReact,
			&enemyResistance.EarthReact, &enemyResistance.SoundWounded1, &enemyResistance.SoundWounded2,
			&enemyResistance.SoundWounded3,
		)
		reactAttribute := [][]string{
			{"PhysicalReact", "物理"}, {"WaterReact", "水"}, {"FireReact", "火"},
			{"ThunderReact", "雷"}, {"AirReact", "风"}, {"EarthReact", "土"},
		}
		v := reflect.ValueOf(&enemyResistance).Elem()
		for _, f := range reactAttribute {
			if fv := v.FieldByName(f[0]); fv.Float() > 0 {
				enemyResistance.React += fmt.Sprintf("%s:%s%% ", f[1], perDisp(float32(fv.Float()*100)))
			}
		}
		extractAttribute := []string{
			"PhysicalExtract", "WaterExtract", "FireExtract",
			"ThunderExtract", "AirExtract", "EarthExtract",
		}
		for _, f := range extractAttribute {
			fv := v.FieldByName(f)
			fvp := v.FieldByName(f + "Per")
			if fv.Float() > 0 {
				fvp.SetString(fmt.Sprintf("%s%%", perDisp(float32(fv.Float()*100))))
			} else {
				fvp.SetString("0")
			}
		}
		enemyResistance.React = strings.TrimSuffix(enemyResistance.React, " ")
		enemyResistanceList = append(enemyResistanceList, enemyResistance)
	}
	rows.Close()
	enemyResistances.EnemyResistanceList = enemyResistanceList
	enemyResistances.Part = "抗性、反弹、受伤音效"
	return
}

type EnemySkill struct {
	Id                                                               int     //敌人ID
	Name                                                             string  //名称
	PhysicalAdditional                                               int     //物理伤害追加
	AdditionalRate                                                   string  //爆击、格挡、命中、反击追加合并
	AdditionalCritical, FendOff, AdditionalHitting, CounterPunchRate float32 //爆击、格挡、命中、反击追加
	Skill1, Skill2, Skill3, Skill4, Skill5                           int     //技能1～5的ID，对应Magic和Stunt两张表中的ID
	Skills                                                           string  //技能1～5，格式为1:特技5833:剑啸九天 2:仙术5742:雨恨云愁
}
type EnemySkills struct {
	Part           string
	EnemySkillList []EnemySkill
}

func getEnemySkill() (enemySkills EnemySkills) {
	sql := `
        select id,name,physical_additional,additional_critical,fend_off,additional_hitting,counterpunch_rate,
        skill1,skill2,skill3,skill4,skill5 from Monster
`
	enemySkillList := []EnemySkill{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		enemySkill := EnemySkill{}
		rows.Scan(
			&enemySkill.Id, &enemySkill.Name, &enemySkill.PhysicalAdditional, &enemySkill.AdditionalCritical,
			&enemySkill.FendOff, &enemySkill.AdditionalHitting, &enemySkill.CounterPunchRate, &enemySkill.Skill1,
			&enemySkill.Skill2, &enemySkill.Skill3, &enemySkill.Skill4, &enemySkill.Skill5,
		)
		additionalAttribute := [][]string{
			{"AdditionalCritical", "暴击"}, {"FendOff", "格挡"}, {"AdditionalHitting", "命中"}, {"CounterPunchRate", "反击"},
		}
		v := reflect.ValueOf(&enemySkill).Elem()
		for _, f := range additionalAttribute {
			if fv := v.FieldByName(f[0]); fv.Float() > 0 {
				enemySkill.AdditionalRate += fmt.Sprintf("%s:%s%% ", f[1], perDisp(float32(fv.Float()*100)))
			}
		}
		enemySkill.Skills = getSkills([]int{enemySkill.Skill1, enemySkill.Skill2, enemySkill.Skill3, enemySkill.Skill4, enemySkill.Skill5})
		enemySkillList = append(enemySkillList, enemySkill)
	}
	rows.Close()
	enemySkills.EnemySkillList = enemySkillList
	enemySkills.Part = "物理追加与技能"
	return
}

type EnemyDrop struct {
	Id                                         int     //敌人ID
	Name                                       string  //名称
	StolenProperty, StolenNumber, StolenMoney  int     //可偷窃物品ID、可偷窃物品数量、可偷窃金钱数量
	Drop1ID, Drop2ID, Drop3ID, Drop4ID         int     //掉落物品ID
	Drop1, Drop2, Drop3, Drop4                 string  //掉落物品名称
	Drop1Rate, Drop2Rate, Drop3Rate, Drop4Rate float32 //对应的掉落机率
	Drop1Per, Drop2Per, Drop3Per, Drop4Per     string  //对应的掉落百分率
	MaxDropMoney, MinDropMoney                 int     //掉落金钱数量范围
	Stolen                                     string  //偷窃信息合并，物品及数量或金钱数量如止血草:1或钱:1000
	DropMoney                                  string  //掉落金钱数量范围，如100~200，如上下限相同，则直接显示该数值
}
type EnemyDrops struct {
	Part          string
	EnemyDropList []EnemyDrop
}

func getEnemyDrop() (enemyDrops EnemyDrops) {
	sql := `
        select id,name,stolen_property,stolen_number,stolen_money,drop1id,drop1rate,drop2id,drop2rate,
        drop3id,drop3rate,drop4id,drop4rate,max_drop_money,min_drop_money from Monster
    `
	enemyDropList := []EnemyDrop{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		enemyDrop := EnemyDrop{}
		rows.Scan(
			&enemyDrop.Id, &enemyDrop.Name, &enemyDrop.StolenProperty, &enemyDrop.StolenNumber,
			&enemyDrop.StolenMoney, &enemyDrop.Drop1ID, &enemyDrop.Drop1Rate, &enemyDrop.Drop2ID,
			&enemyDrop.Drop2Rate, &enemyDrop.Drop3ID, &enemyDrop.Drop3Rate, &enemyDrop.Drop4ID,
			&enemyDrop.Drop4Rate, &enemyDrop.MaxDropMoney, &enemyDrop.MinDropMoney,
		)
		if enemyDrop.StolenProperty != 0 {
			paName := getName("Property", enemyDrop.StolenProperty)
			if enemyDrop.StolenNumber > 1 {
				enemyDrop.Stolen = fmt.Sprintf("%s*%d", paName, enemyDrop.StolenNumber)
			} else {
				enemyDrop.Stolen = fmt.Sprintf("%s", paName)
			}
		} else if enemyDrop.StolenMoney != 0 {
			enemyDrop.Stolen = fmt.Sprintf("%s*%d", "金钱", enemyDrop.StolenMoney)
		}
		enemyDrop.Drop1 = getName("Property", enemyDrop.Drop1ID)
		enemyDrop.Drop2 = getName("Property", enemyDrop.Drop2ID)
		enemyDrop.Drop3 = getName("Property", enemyDrop.Drop3ID)
		enemyDrop.Drop4 = getName("Property", enemyDrop.Drop4ID)
		if enemyDrop.Drop1 != "" {
			enemyDrop.Drop1Per = fmt.Sprintf("%s%%", perDisp(float32(enemyDrop.Drop1Rate*100)))
		}
		if enemyDrop.Drop2 != "" {
			enemyDrop.Drop2Per = fmt.Sprintf("%s%%", perDisp(float32(enemyDrop.Drop2Rate*100)))
		}
		if enemyDrop.Drop3 != "" {
			enemyDrop.Drop3Per = fmt.Sprintf("%s%%", perDisp(float32(enemyDrop.Drop3Rate*100)))
		}
		if enemyDrop.Drop4 != "" {
			enemyDrop.Drop4Per = fmt.Sprintf("%s%%", perDisp(float32(enemyDrop.Drop4Rate*100)))
		}
		if enemyDrop.MaxDropMoney != enemyDrop.MinDropMoney {
			enemyDrop.DropMoney = fmt.Sprintf("%d~%d", enemyDrop.MinDropMoney, enemyDrop.MaxDropMoney)
		} else if enemyDrop.MaxDropMoney != 0 {
			enemyDrop.DropMoney = fmt.Sprintf("%d", enemyDrop.MaxDropMoney)
		}
		enemyDropList = append(enemyDropList, enemyDrop)
	}
	rows.Close()
	enemyDrops.EnemyDropList = enemyDropList
	enemyDrops.Part = "偷窃与掉落"
	return
}

type EnemySummary struct {
	Id                                                               int     //敌人ID
	Name                                                             string  //名称
	Level, Experience 								int    //   等级、经验
	HP                                                               int     //精
	//五灵类
	Water, Fire, Thunder, Air, Earth                                 int     //五灵属性
	Wuling                                                           string  //五灵属性合并
	//吸收类
	PhysicalExtract, WaterExtract, FireExtract          float32 //物理、水、火吸收
	ThunderExtract, AirExtract, EarthExtract            float32 //雷风土吸收
	Immune 												string  //仅免疫的情况
	Absorb												string  //达到吸收的水平了
	Extract                                               string  //吸收合并
	//反弹类
	PhysicalReact, WaterReact, FireReact                float32 //物理、水、火反弹
	ThunderReact, AirReact, EarthReact                  float32 //雷风土反弹
	React                                               string  //反弹合并
	Property 								string  //五灵、吸收、反弹合并
	//机率追加类
	AdditionalCritical, FendOff, AdditionalHitting, CounterPunchRate float32 //爆击、格挡、命中、反击追加
	AdditionalRate                                                   string  //爆击、格挡、命中、反击追加合并
	//技能类
	Skill1, Skill2, Skill3, Skill4, Skill5                           int     //技能1～5的ID，对应Magic和Stunt两张表中的ID
	Skills                                                           string  //技能合并
	//偷窃类
	StolenProperty, StolenNumber, StolenMoney                        int     //可偷窃物品ID、可偷窃物品数量、可偷窃金钱数量
	//掉落类
	Drop1ID, Drop2ID, Drop3ID, Drop4ID                               int     //掉落物品ID
	Drop1, Drop2, Drop3, Drop4                                       string  //掉落物品名称
	Drop1Rate, Drop2Rate, Drop3Rate, Drop4Rate                       float32 //对应的掉落机率
	Drop1Per, Drop2Per, Drop3Per, Drop4Per                           string  //对应的掉落百分率
	Stolen                                                           string  //偷窃信息合并，物品及数量或金钱数量如止血草*9或钱*1000
	Drop                                                             string  // 掉落信息合并
}
type EnemySummarys struct {
	Part             string
	EnemySummaryList []EnemySummary
}

/*
最终格式类似下面这样
敌人名称	等级	经验	精		属性		机率追加	技能		可偷物品	掉落物品

	琼浆玉露25%

沉霜		36		700		2560	水属性		暴击10% 	冰刺		四合寒香	护心丹41.67%

	仙术免疫	格挡30%		云之幻境				戒冥玉41.67%
													百花醉仙露33.33%

切记：属性不要通过描述来获得，不准的
*/
func getEnemySummary() (enemySummarys EnemySummarys) {
	//第一步，先把大量原始数据从数据库中取出
	sql := `
        select id,name,level,experience,max_hp,
		water,fire,thunder,air,earth,
		physical_extract,water_extract,fire_extract,thunder_extract,air_extract,earth_extract,
		physical_react,water_react,fire_react,thunder_react,air_react,earth_react,
		additional_critical,fend_off,additional_hitting,counterpunch_rate,
		skill1,skill2,skill3,skill4,skill5,
		stolen_property,stolen_number,stolen_money,
		drop1id,drop1rate,drop2id,drop2rate,drop3id,drop3rate,drop4id,drop4rate from Monster
    `
	enemySummaryList := []EnemySummary{}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		enemySummary := EnemySummary{}
		rows.Scan(
			&enemySummary.Id, &enemySummary.Name, &enemySummary.Level, &enemySummary.Experience, &enemySummary.HP, 
			&enemySummary.Water,&enemySummary.Fire,&enemySummary.Thunder,&enemySummary.Air,&enemySummary.Earth,
			&enemySummary.PhysicalExtract,&enemySummary.WaterExtract,&enemySummary.FireExtract,
			&enemySummary.ThunderExtract,&enemySummary.AirExtract,&enemySummary.EarthExtract,
			&enemySummary.PhysicalReact,&enemySummary.WaterReact,&enemySummary.FireReact,
			&enemySummary.ThunderReact,&enemySummary.AirReact,&enemySummary.EarthReact,
			&enemySummary.AdditionalCritical,&enemySummary.FendOff,&enemySummary.AdditionalHitting,&enemySummary.CounterPunchRate,
			&enemySummary.Skill1,&enemySummary.Skill2,&enemySummary.Skill3,&enemySummary.Skill4,&enemySummary.Skill5,
			&enemySummary.StolenProperty, &enemySummary.StolenNumber,&enemySummary.StolenMoney, 
			&enemySummary.Drop1ID, &enemySummary.Drop1Rate, &enemySummary.Drop2ID,&enemySummary.Drop2Rate, 
			&enemySummary.Drop3ID, &enemySummary.Drop3Rate, &enemySummary.Drop4ID,&enemySummary.Drop4Rate,
		)
		//然后开始进行深加工，首选加五灵属性
		wulingAttribute := [][]string{
			{"Water", "水"}, {"Fire", "火"}, {"Thunder", "雷"},{"Air", "风"}, {"Earth", "土"},
		}
		v := reflect.ValueOf(enemySummary)
		for _, f := range wulingAttribute {
			if fv := v.FieldByName(f[0]); fv.Int() > 0 {
				enemySummary.Wuling += fmt.Sprintf("%s、", f[1])
			}
		}
		if enemySummary.Wuling != "" {
			enemySummary.Wuling = strings.TrimSuffix(enemySummary.Wuling, "、")
			enemySummary.Wuling += "属性"	
		} else {
			enemySummary.Wuling = "无属性"	
		}
		//然后加工吸收类
		extractAttribute := [][]string{
			{"PhysicalExtract", "物理"}, {"WaterExtract", "水"}, {"FireExtract", "火"},
			{"ThunderExtract", "雷"}, {"AirExtract", "风"}, {"EarthExtract", "土"},
		}
		//v := reflect.ValueOf(&enemySummary).Elem()
		ve := v.Elem()
		var flagImmune,flagAbsorb int
		for _, f := range extractAttribute {
			fv := ve.FieldByName(f[0])
			if fv.Float() == 1.0 {
				enemySummary.Immune += fmt.Sprintf("%s、", f[1])
				if f[1] != "物理" {
					flagImmune+=1
				}
			} else if fv.Float() > 1.0 {
				enemySummary.Absorb += fmt.Sprintf("%s、", f[1])
				if f[1] != "物理" {
					flagAbsorb+=1
				}
			}
		}
		if flagImmune == 5 {
			enemySummary.Extract="仙术免疫"
		} else if flagAbsorb == 5 {
			enemySummary.Extract="仙术吸收"
		} else if enemySummary.Immune != "" || enemySummary.Absorb != ""{
			enemySummary.Immune = strings.TrimSuffix(enemySummary.Immune, "、")
			enemySummary.Extract += fmt.Sprintf("%s免疫 ",enemySummary.Immune)
			enemySummary.Absorb = strings.TrimSuffix(enemySummary.Absorb, "、")
			enemySummary.Extract += fmt.Sprintf("%s吸收",enemySummary.Absorb)
		}
		enemySummary.Extract = strings.TrimSuffix(enemySummary.Extract, " ")
		//然后加工反弹类
		reactAttribute := [][]string{
			{"PhysicalReact", "物理"}, {"WaterReact", "水"}, {"FireReact", "火"},
			{"ThunderReact", "雷"}, {"AirReact", "风"}, {"EarthReact", "土"},
		}
		//v := reflect.ValueOf(&enemySummary).Elem()
		flagReact := 0
		for _, f := range reactAttribute {
			if fv := ve.FieldByName(f[0]); fv.Float() > 0 {
				enemySummary.React += fmt.Sprintf("%s、", f[1])
				if f[1] != "物理" {
					flagReact+=1
				}
			}
		}
		if flagReact == 5 {
			enemySummary.React = "仙术反弹"
		} else {
			enemySummary.React = strings.TrimSuffix(enemySummary.React, "、")
			enemySummary.React += fmt.Sprintf("%s反弹",enemySummary.React)
		}
		enemySummary.Property = fmt.Sprintf(
			"%s/%s/%s",enemySummary.Wuling,enemySummary.Extract,enemySummary.React,
		)
		//加工机率追加类
		additionalAttribute := [][]string{
			{"AdditionalCritical", "暴击"}, {"FendOff", "格挡"}, {"AdditionalHitting", "命中"}, {"CounterPunchRate", "反击"},
		}
		//v := reflect.ValueOf(&enemySummary).Elem()
		for _, f := range additionalAttribute {
			if fv := ve.FieldByName(f[0]); fv.Float() > 0 {
				enemySummary.AdditionalRate += fmt.Sprintf("%s%s%%/", f[1], perDisp(float32(fv.Float()*100)))
			}
		}
		enemySummary.AdditionalRate = strings.TrimSuffix(enemySummary.AdditionalRate, "/")
		//加工技能类
		enemySummary.Skills = getSkillCombo([]int{enemySummary.Skill1, enemySummary.Skill2, enemySummary.Skill3, enemySummary.Skill4, enemySummary.Skill5})
		
		//加工偷窃类
		if enemySummary.StolenProperty != 0 {
			paName := getName("Property", enemySummary.StolenProperty)
			if enemySummary.StolenNumber > 1 {
				enemySummary.Stolen = fmt.Sprintf("%s*%d", paName, enemySummary.StolenNumber)
			} else {
				enemySummary.Stolen = fmt.Sprintf("%s", paName)
			}
		} else if enemySummary.StolenMoney != 0 {
			enemySummary.Stolen = fmt.Sprintf("%s*%d", "金钱", enemySummary.StolenMoney)
		}
		//加工掉落类
		enemySummary.Drop1 = getName("Property", enemySummary.Drop1ID)
		enemySummary.Drop2 = getName("Property", enemySummary.Drop2ID)
		enemySummary.Drop3 = getName("Property", enemySummary.Drop3ID)
		enemySummary.Drop4 = getName("Property", enemySummary.Drop4ID)
		if enemySummary.Drop1 != "" {
			enemySummary.Drop1Per = fmt.Sprintf("%s%%", perDisp(float32(enemySummary.Drop1Rate*100)))
			enemySummary.Drop += fmt.Sprintf("%s%s%%/",enemySummary.Drop1,enemySummary.Drop1Per)
		}
		if enemySummary.Drop2 != "" {
			enemySummary.Drop2Per = fmt.Sprintf("%s%%", perDisp(float32(enemySummary.Drop2Rate*100)))
			enemySummary.Drop += fmt.Sprintf("%s%s%%/",enemySummary.Drop2,enemySummary.Drop2Per)
		}
		if enemySummary.Drop3 != "" {
			enemySummary.Drop3Per = fmt.Sprintf("%s%%", perDisp(float32(enemySummary.Drop3Rate*100)))
			enemySummary.Drop += fmt.Sprintf("%s%s%%/",enemySummary.Drop3,enemySummary.Drop3Per)
		}
		if enemySummary.Drop4 != "" {
			enemySummary.Drop4Per = fmt.Sprintf("%s%%", perDisp(float32(enemySummary.Drop4Rate*100)))
			enemySummary.Drop += fmt.Sprintf("%s%s%%",enemySummary.Drop4,enemySummary.Drop4Per)
		}
		enemySummaryList = append(enemySummaryList, enemySummary)
	}
	rows.Close()
	enemySummarys.EnemySummaryList = enemySummaryList
	enemySummarys.Part = "主要属性汇总"
	return
}
