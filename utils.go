package main

import (
	"fmt"
	"strings"
)

//获取购买场景
func getBuyScene(itemId int) (buyScene string) {
	//先拿商店列表
	shopSql := "select shop_id,open_condition,open_num from Good where goods_id=?"
	rows, _ := Db.Query(shopSql, itemId)
	for rows.Next() {
		//获取此店的全部信息
		shop := Shop{}
		rows.Scan(&shop.Id, &shop.OpenCondition, &shop.OpenNum)
		if shop.Id > "shop036" {
			shop.Period = "后期"
		}
		shop.Scene = getShopScene(shop.Id)
		if shop.OpenCondition > 0 {
			shop.Commission = getCommission(shop.OpenCondition, shop.OpenNum)
		}
		//利用店铺信息构造购买场景信息列表
		buyScene += shop.Scene
		if shop.Period != "" {
			buyScene += fmt.Sprintf("(%s)", shop.Period)
		}
		if shop.Commission != "" {
			buyScene += fmt.Sprintf("-委托:%s", shop.Commission)
		}
		buyScene += " "
	}
	rows.Close()
	buyScene = strings.TrimSuffix(buyScene, " ")
	return
}

//最简的唯一敌人结构
type EnemyUnique struct {
	Id   int
	Name string
}

//获得物品的偷窃途径
func getStolen(itemId int) (stolenEnemy string) {
	sql := "select id,name from Monster where stolen_property=?"
	rows, _ := Db.Query(sql, itemId)
	for rows.Next() {
		enemyUnique := EnemyUnique{}
		rows.Scan(&enemyUnique.Id, &enemyUnique.Name)
		stolenEnemy += fmt.Sprintf("%d:%s ", enemyUnique.Id, enemyUnique.Name)
	}
	rows.Close()
	stolenEnemy = strings.TrimSuffix(stolenEnemy, " ")
	return
}

//获得物品的掉落途径
func getDrop(itemId int) (dropEnemy string) {
	sql := "select id,name from Monster where drop1id=? or drop2id=? or drop3id=? or drop4id=?"
	rows, _ := Db.Query(sql, itemId, itemId, itemId, itemId)
	for rows.Next() {
		enemyUnique := EnemyUnique{}
		rows.Scan(&enemyUnique.Id, &enemyUnique.Name)
		dropEnemy += fmt.Sprintf("%d:%s ", enemyUnique.Id, enemyUnique.Name)
	}
	rows.Close()
	dropEnemy = strings.TrimSuffix(dropEnemy, " ")
	return
}

type Pick struct {
	Scene, Section, ItemId, Model, Texture string //场景与区块编号、物品ID、模型、贴图
	SceneName, SectionName, Apperance      string
	CoorX, CoorY, CoorZ                    float32 //东西坐标、上下坐标、南北坐标
}

//获得物品的拾取途径
func pickItem(itemId int) (picks []Pick) {
	sql := fmt.Sprintf(`select scene,section,item_id,model,texture,coor_x,coor_y,coor_z 
    from SceneItem where (equip_id=%d and item_num!=0) or (property_id=%d and item_num!=0)
    `, itemId, itemId)
	for i := 1; i <= 6; i++ {
		sql += fmt.Sprintf(" or (item%did=%d and item%dnum!=0)", i, itemId, i)
	}
	rows, _ := Db.Query(sql)
	for rows.Next() {
		pick := Pick{}
		rows.Scan(
			&pick.Scene, &pick.Section, &pick.ItemId, &pick.Model, &pick.Texture,
			&pick.CoorX, &pick.CoorY, &pick.CoorZ,
		)
		pick.SceneName = getSceneName(pick.Scene)
		pick.SectionName = getSectionName(pick.Scene, pick.Section)
		pick.Apperance = getApperance(pick.Model)
		if pick.SectionName != "" {
			picks = append(picks, pick)
		}
	}
	rows.Close()
	return
}

//以更好看的方式，显示全部的百分比数值
func perDisp(f float32) (fs string) {
	fs = fmt.Sprintf("%.2f", f)
	for {
		hasDot, TrailZero := strings.Contains(fs, "."), strings.HasSuffix(fs, "0")
		if !TrailZero || !hasDot {
			break
		} else {
			fs = strings.TrimSuffix(fs, "0")
			fs = strings.TrimSuffix(fs, ".")
		}
	}
	return
}

//依据场景与区块的编码，获得区块的中文名，要注意
func getSectionName(scene, section string) (name string) {
	sql := "select name from Scene where scene=? and section=?"
	Db.QueryRow(sql, scene, section).Scan(&name)
	if strings.HasSuffix(section, "Y") {
		name += "(夜)"
	}
	return
}

//获取室外区域
func getOuter(scene, section string) (name string) {
	//分析区块编码，看其是否属于子区块
	outerSect := ""
	if strings.HasPrefix(section, "N") || strings.HasPrefix(section, scene) {
		return
	} else {
		outerSect = scene + string(section[0])
	}
	if outerSect != "" {
		sql := "select name from Scene where scene=? and section=?"
		Db.QueryRow(sql, scene, outerSect).Scan(&name)
	}
	return
}

//依据物品的模型，得到物品的外观类别
func getApperance(model string) (apperance string) {
	switch {
	case model == "OM06":
		apperance = "特殊矿石"
	case model == "OM07":
		apperance = "小宝箱"
	case model == "OM08":
		apperance = "大宝箱"
	case model == "OM09":
		apperance = "隐藏宝箱"
	case model == "OM10":
		apperance = "上锁宝箱"
	case model == "OQ20":
		apperance = "金钱"
	case strings.HasPrefix(model, "CK"):
		apperance = "矿石"
	case strings.HasPrefix(model, "CQ"):
		apperance = "其它材料"
	case strings.HasPrefix(model, "DF"):
		apperance = "辅助道具"
	case strings.HasPrefix(model, "SW"):
		apperance = "食物"
	case strings.HasPrefix(model, "DH"):
		apperance = "其它恢复道具"
	case strings.HasPrefix(model, "DG"):
		apperance = "攻击道具"
	case strings.HasPrefix(model, "JQ"):
		apperance = "剧情道具"
	case strings.HasPrefix(model, "WT"):
		apperance = "剑"
	case strings.HasPrefix(model, "WL"):
		apperance = "双剑"
	case strings.HasPrefix(model, "WM"):
		apperance = "琴"
	case strings.HasPrefix(model, "P"):
		apperance = "配饰"
	case strings.HasPrefix(model, "X"):
		apperance = "足部防具"
	case strings.HasPrefix(model, "M"):
		apperance = "头部防具"
	case strings.HasPrefix(model, "Y"):
		apperance = "身体防具"
	default:
		apperance = "其它"
	}
	return
}

//两个函数，一个通过ID查名称，另一个反过来
func getName(table string, id int) (name string) {
	sql := fmt.Sprintf("select name from %s where id=%d", table, id)
	Db.QueryRow(sql).Scan(&name)
	return
}
func getId(table, name string) (id int) {
	sql := fmt.Sprintf(`select id from %s where name="%s"`, table, name)
	Db.QueryRow(sql).Scan(&id)
	return
}
func getSceneName(id string) (name string) {
	sql := "select name from SceneName where id=?"
	Db.QueryRow(sql, id).Scan(&name)
	return
}
func getSceneId(name string) (id string) {
	sql := "select id from SceneName where name=?"
	Db.QueryRow(sql, name).Scan(&id)
	return
}
func getIsBoss(id int) (name string) {
	sql := "select name from BOOL where id=?"
	Db.QueryRow(sql, id).Scan(&name)
	if name == "TRUE" {
		name = "是"
	} else if name == "FALSE" {
		name = "否"
	}
	return
}

//依据技能ID获得技能种类及其名称
func getSkills(skillIds []int) (skills string) {
	//只要不是0,就查仙术和特技两张表
	for i, s := range skillIds {
		if s != 0 {
			if magic := getName("Magic", s); magic != "" {
				skills += fmt.Sprintf("%d:%s%d:%s ", i+1, "仙术", s, magic)
			} else if stunt := getName("Stunt", s); stunt != "" {
				skills += fmt.Sprintf("%d:%s%d:%s ", i+1, "特技", s, stunt)
			}
		}
	}
	skills = strings.TrimSuffix(skills, " ")
	return
}

type Shop struct {
	Id, Period             string //店铺ID,出现时期
	OpenCondition, OpenNum int    //开放变量及其值
	Commission, Scene      string //委托任务，与场景中文名
}

//依据商店ID获得场景名称
func getShopScene(shopId string) (scene string) {
	sceneIdSql := "select position from ShopProperty where id=?"
	var sceneId string
	Db.QueryRow(sceneIdSql, shopId).Scan(&sceneId)
	scene = getSceneName(sceneId)
	return
}

//依据开放变量及值，获得委托任务
func getCommission(openCondition, openNum int) (commission string) {
	commissionSql := "select name from Mission where trunk=? and quest_id=?"
	Db.QueryRow(commissionSql, openCondition, openNum).Scan(&commission)
	return
}
